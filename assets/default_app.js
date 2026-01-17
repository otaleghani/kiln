// Script loading helper
window.loadScript = function (src, id) {
  return new Promise((resolve, reject) => {
    if (document.getElementById(id)) return resolve(); // Already loaded
    const s = document.createElement("script");
    s.src = src;
    s.id = id;
    s.defer = true;
    s.onload = resolve;
    s.onerror = reject;
    document.head.appendChild(s);
  });
};

window.initMathJax = async function () {
  const content = document.querySelector("#content");
  if (!content) return;

  const text = content.innerText;
  // Check for $$ or \( or \[
  // Note: We check for delimiters to avoid loading heavy scripts unnecessarily
  if (!text.includes("$$") && !text.includes("\\(") && !text.includes("\\[")) {
    return;
  }

  // Check if already loaded - If typesetPromise exists, the library is active. Just tell it to render.
  if (window.MathJax && window.MathJax.typesetPromise) {
    await window.MathJax.typesetPromise();
    return;
  }

  // Configure it if not loaded
  if (!window.MathJax) {
    window.MathJax = {
      tex: {
        inlineMath: [
          ["$", "$"],
          ["\\(", "\\)"],
        ],
        displayMath: [
          ["$$", "$$"],
          ["\\[", "\\]"],
        ],
        processEscapes: true,
      },
      svg: { fontCache: "global" },
      startup: {
        // This handles the VERY FIRST render when the script loads
        pageReady: () => {
          return window.MathJax.startup.defaultPageReady();
        },
      },
    };
  }

  // Load script
  await window.loadScript(
    "https://cdn.jsdelivr.net/npm/mathjax@4/tex-svg.js",
    "mathjax-script",
  );
};

// Lazyload mermaid
window.initMermaid = async function () {
  const graphs = document.querySelectorAll(".mermaid");
  if (graphs.length === 0) return; // Stop if no diagrams

  // Save original content for theme switching re-renders
  graphs.forEach((g) => {
    if (!g.getAttribute("data-original"))
      g.setAttribute("data-original", g.innerHTML);
  });

  // Lazy Load Mermaid
  try {
    await import("https://cdn.jsdelivr.net/npm/mermaid@10/dist/mermaid.esm.min.mjs").then(
      (m) => {
        window.mermaid = m.default;
        const theme =
          document.documentElement.getAttribute("data-theme") === "dark"
            ? "dark"
            : "default";
        window.mermaid.initialize({ startOnLoad: false, theme: theme });
        window.mermaid.run({ querySelector: ".mermaid" });
      },
    );
  } catch (e) {
    console.warn("Mermaid failed to load", e);
  }
};

// Changes the giscus theme based on the current theme
window.changeGiscusTheme = function () {
  // Get the Giscus iframe
  const iframe = document.querySelector("iframe.giscus-frame");

  if (!iframe) return;

  // Define the URL based on the mode ('light' or 'dark')
  // REPLACE these URLs with your actual file paths
  const current = document.documentElement.getAttribute("data-theme");
  const sysDark = window.matchMedia("(prefers-color-scheme: dark)").matches;
  let target = !current
    ? sysDark
      ? "light"
      : "dark"
    : current === "dark"
      ? "light"
      : "dark";
  const themeUrl =
    target === "dark"
      ? "{{.BaseURL}}/giscus-theme-dark.css"
      : "{{.BaseURL}}/giscus-theme-light.css";

  // Send the message to Giscus
  iframe.contentWindow.postMessage(
    {
      giscus: {
        setConfig: {
          theme: themeUrl,
        },
      },
    },
    "https://giscus.app", // This origin is required for security
  );
};

// Theme toggle logic
window.initThemeToggle = function () {
  const btn = document.getElementById("theme-toggle");
  if (!btn) return;

  const newBtn = btn.cloneNode(true);
  btn.parentNode.replaceChild(newBtn, btn);

  newBtn.addEventListener("click", async () => {
    try {
      const current = document.documentElement.getAttribute("data-theme");
      const sysDark = window.matchMedia("(prefers-color-scheme: dark)").matches;
      let target = !current
        ? sysDark
          ? "light"
          : "dark"
        : current === "dark"
          ? "light"
          : "dark";

      document.documentElement.setAttribute("data-theme", target);
      localStorage.setItem("theme", target);
      window.changeGiscusTheme();

      // Re-render Mermaid if loaded
      if (window.mermaid) {
        const mermaidTheme = target === "dark" ? "dark" : "default";
        window.mermaid.initialize({ startOnLoad: false, theme: mermaidTheme });
        const graphs = document.querySelectorAll(".mermaid");
        graphs.forEach((graph) => {
          const original = graph.getAttribute("data-original");
          if (original) {
            graph.removeAttribute("data-processed");
            graph.innerHTML = original;
          }
        });
        await window.mermaid.run({ querySelector: ".mermaid" });
      }
    } catch (e) {
      console.error(e);
    }
  });
};

// Left sidebar logic
window.setupLeftSidebarInteraction = function () {
  const toggleBtns = document.querySelectorAll(".sidebar-toggle.left-toggle");
  const sidebar = document.getElementById("left-sidebar");
  if (!toggleBtns.length || !sidebar) return;

  toggleBtns.forEach((btn) => {
    const newBtn = btn.cloneNode(true);
    btn.parentNode.replaceChild(newBtn, btn);
    newBtn.addEventListener("click", () => {
      sidebar.classList.toggle("collapsed");
      localStorage.setItem(
        "left-sidebar",
        sidebar.classList.contains("collapsed"),
      );
    });
  });
};

// Right sidebar logic
window.setupRightSidebarInteraction = function () {
  const toggleBtns = document.querySelectorAll(".sidebar-toggle.right-toggle");
  const sidebar = document.getElementById("right-sidebar");
  if (!toggleBtns.length || !sidebar) return;

  toggleBtns.forEach((btn) => {
    const newBtn = btn.cloneNode(true);
    btn.parentNode.replaceChild(newBtn, btn);
    newBtn.addEventListener("click", () => {
      sidebar.classList.toggle("collapsed");
      localStorage.setItem(
        "right-sidebar",
        sidebar.classList.contains("collapsed"),
      );
    });
  });
};

// Sidebars autoclose on mobile logic
window.setupMobileAutoClose = function () {
  document.body.addEventListener("click", (e) => {
    if (window.innerWidth > 1280) return;

    const link = e.target.closest("a");
    const isGraphNode =
      ["circle", "text"].includes(e.target.tagName.toLowerCase()) &&
      (e.target.closest("#global-graph-container") ||
        e.target.closest("#local-graph-container"));

    if (!link && !isGraphNode) return;

    // Close Sidebars
    ["left-sidebar", "right-sidebar"].forEach((id) => {
      const el = document.getElementById(id);
      if (el && !el.classList.contains("collapsed")) {
        el.classList.add("collapsed");
        localStorage.setItem(
          id === "left-sidebar" ? "left-sidebar" : "right-sidebar",
          "true",
        );
      }
    });
  });
};

// Navbar search
window.initNavbarSearch = function () {
  const searchInput = document.getElementById("navbar-search");
  if (!searchInput) return;

  // Remove old listeners by cloning (optional but safer in SPAs)
  const newInput = searchInput.cloneNode(true);
  searchInput.parentNode.replaceChild(newInput, searchInput);

  newInput.addEventListener("input", (e) => {
    const term = e.target.value.toLowerCase().trim();
    const items = document.querySelectorAll("#left-sidebar li");

    items.forEach((item) => {
      const text = item.textContent.toLowerCase();
      const matches = text.includes(term);
      item.style.display = matches ? "" : "none";
      if (matches && term) {
        const details = item.querySelector("details");
        if (details) details.open = true;
        let parent = item.parentElement;
        while (parent && parent.closest(".sidebar")) {
          if (parent.tagName === "DETAILS") parent.open = true;
          parent = parent.parentElement;
        }
      }
    });
  });
  // Restore focus if needed after swap
  if (document.activeElement !== newInput) newInput.focus();
};

// Called from HTML to init canvas with Go data
window.initCanvasMode = function (canvasData) {
  if (window.renderer) window.renderer.cleanup();
  window.CANVAS_DATA = canvasData;

  const tryInitCanvas = () => {
    if (window.JsonCanvasRenderer) {
      window.renderer = new window.JsonCanvasRenderer();
      window.renderer.load(window.CANVAS_DATA);
      if (window.lucide) window.lucide.createIcons();
    } else {
      setTimeout(tryInitCanvas, 50);
    }
  };
  tryInitCanvas();
};

// Expands the graph
window.toggleGraphExpand = function () {
  const wrapper = document.getElementById("local-graph-wrapper");
  if (!wrapper) return;
  wrapper.classList.toggle("expanded");
};

// Highlights the sidebar link
window.highlightSidebarLink = function () {
  document
    .querySelectorAll("#left-sidebar a")
    .forEach((el) => el.classList.remove("text-accent"));

  const normalize = (p) => {
    if (!p) return "";
    try {
      p = decodeURIComponent(p);
    } catch (e) {}
    return p.length > 1 && p.endsWith("/") ? p.slice(0, -1) : p;
  };

  const currentPath = normalize(window.location.pathname);
  const links = document.querySelectorAll("#left-sidebar a");

  for (const link of links) {
    const linkPath = normalize(link.getAttribute("href"));
    if (linkPath === currentPath) {
      link.classList.add("text-accent");
      let parent = link.parentElement;
      while (parent) {
        if (parent.tagName === "DETAILS") parent.open = true;
        parent = parent.parentElement;
      }
      break;
    }
  }
};

window.addCopyButtons = function () {
  document.querySelectorAll(".chroma").forEach((block) => {
    if (block.querySelector(".copy-code-btn")) return;
    const btn = document.createElement("button");
    btn.className = "copy-code-btn";
    btn.textContent = "Copy";
    btn.addEventListener("click", () => {
      const code = block.querySelector("code").innerText;
      navigator.clipboard
        .writeText(code)
        .then(() => {
          btn.textContent = "Copied!";
          setTimeout(() => {
            btn.textContent = "Copy";
          }, 2000);
        })
        .catch((err) => {});
    });
    block.appendChild(btn);
  });
};

// Calls every init function
window.initAll = function () {
  window.changeGiscusTheme();
  window.initThemeToggle();
  window.initNavbarSearch();
  window.setupRightSidebarInteraction();
  window.setupLeftSidebarInteraction();
  window.highlightSidebarLink();
  window.addCopyButtons();

  Promise.all([window.initMathJax(), window.initMermaid()]);
};

document.addEventListener("DOMContentLoaded", () => {
  window.initAll();
  window.setupMobileAutoClose();
});

document.addEventListener("htmx:afterSwap", () => {
  window.initAll();

  // MathJax specific re-render if it was already loaded
  if (window.MathJax && window.MathJax.typesetPromise) {
    window.MathJax.typesetPromise();
  }
});
