// Client-side JS for default layout: sidebar, search, theme toggle, and local graph. @feature:layouts
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
  const iframe = document.querySelector("iframe.giscus-frame");
  if (!iframe) return;

  const current = document.documentElement.getAttribute("data-theme");
  const sysDark = window.matchMedia("(prefers-color-scheme: dark)").matches;

  // If 'current' exists, use it.
  // If 'current' is missing, check system preference (sysDark).
  let target;
  if (current) {
    target = current; // e.g., if data-theme="dark", target is "dark"
  } else {
    target = sysDark ? "dark" : "light";
  }

  // Now 'target' holds the correct theme name ("dark" or "light")
  const themeUrl =
    target === "dark"
      ? "{{.BaseURL}}/giscus-theme-dark.css"
      : "{{.BaseURL}}/giscus-theme-light.css";

  iframe.contentWindow.postMessage(
    {
      giscus: {
        setConfig: {
          theme: themeUrl,
        },
      },
    },
    "https://giscus.app",
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

// Sidebar toggle via event delegation (survives htmx swaps without re-binding)
window.setupSidebarDelegation = function () {
  if (window._sidebarDelegationBound) return;
  window._sidebarDelegationBound = true;

  document.body.addEventListener("click", (e) => {
    const leftBtn = e.target.closest(".sidebar-toggle.left-toggle");
    const rightBtn = e.target.closest(".sidebar-toggle.right-toggle");

    if (leftBtn) {
      e.stopPropagation();
      window.toggleSidebar("left-sidebar", "right-sidebar");
      return;
    }
    if (rightBtn) {
      e.stopPropagation();
      window.toggleSidebar("right-sidebar", "left-sidebar");
      return;
    }

    // Mobile auto-close: close open sidebars when clicking a link or graph node
    if (window.innerWidth >= 1280) return;
    const link = e.target.closest("a");
    const isGraphNode =
      ["circle", "text"].includes(e.target.tagName.toLowerCase()) &&
      (e.target.closest("#global-graph-container") ||
        e.target.closest("#local-graph-container"));
    if (!link && !isGraphNode) return;

    window.closeSidebar("left-sidebar");
    window.closeSidebar("right-sidebar");
  });
};

window.toggleSidebar = function (id, otherId) {
  const sidebar = document.getElementById(id);
  if (!sidebar) return;

  // On mobile, close the other sidebar first
  if (window.innerWidth < 1280) {
    window.closeSidebar(otherId);
  }

  sidebar.classList.toggle("collapsed");

  if (window.innerWidth >= 1280) {
    localStorage.setItem(id, sidebar.classList.contains("collapsed"));
  }
};

window.closeSidebar = function (id) {
  const sidebar = document.getElementById(id);
  if (!sidebar) return;
  // On mobile, "collapsed" = visible, so removing it hides the sidebar.
  // On desktop, "collapsed" = hidden, so removing it shows the sidebar.
  // We only auto-close on mobile, where removing "collapsed" is correct.
  if (sidebar.classList.contains("collapsed")) {
    sidebar.classList.remove("collapsed");
  }
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

window.initLightbox = function () {
  // Create overlay once
  if (!document.getElementById('img-lightbox')) {
    const overlay = document.createElement('div');
    overlay.id = 'img-lightbox';
    overlay.className = 'hidden';
    overlay.innerHTML = '<button id="img-lightbox-close" type="button" aria-label="Close">&times;</button><img src="" alt="">';
    document.body.appendChild(overlay);

    const close = () => overlay.classList.add('hidden');
    overlay.addEventListener('click', (e) => {
      if (e.target !== overlay.querySelector('img')) close();
    });
    document.addEventListener('keydown', (e) => {
      if (e.key === 'Escape' && !overlay.classList.contains('hidden')) close();
    });
  }

  // Bind expand buttons (idempotent via data attribute)
  document.querySelectorAll('.img-expand-btn').forEach((btn) => {
    if (btn.dataset.lightboxBound) return;
    btn.dataset.lightboxBound = '1';
    btn.addEventListener('click', (e) => {
      e.preventDefault();
      e.stopPropagation();
      const figure = btn.closest('.img-figure');
      const img = figure.querySelector('img');
      if (!img) return;
      const overlay = document.getElementById('img-lightbox');
      // Use original src (highest resolution available)
      overlay.querySelector('img').src = img.src;
      overlay.querySelector('img').alt = img.alt;
      overlay.classList.remove('hidden');
    });
  });
};

// Calls every init function
window.initAll = function () {
  window.initThemeToggle();
  if (window.initFullSearch) {
    window.initFullSearch();
  }
  window.setupSidebarDelegation();
  window.highlightSidebarLink();
  window.addCopyButtons();
  window.initLightbox();
  document.addEventListener("keydown", function(e) {
    if ((e.ctrlKey || e.metaKey) && e.key === "k") {
      e.preventDefault();
      if (window.openSearchModal) window.openSearchModal();
    }
  });

  Promise.all([window.initMathJax(), window.initMermaid()]);
};

document.addEventListener("DOMContentLoaded", () => {
  window.initAll();
});

document.addEventListener("htmx:afterSwap", () => {
  window.initAll();

  // MathJax specific re-render if it was already loaded
  if (window.MathJax && window.MathJax.typesetPromise) {
    window.MathJax.typesetPromise();
  }
});

window.addEventListener("message", function (event) {
  // Security check: Only allow messages from Giscus
  if (event.origin !== "https://giscus.app") return;

  // Check if the message is specifically about Giscus data
  // (Giscus sends a distinctive message structure)
  if (!(typeof event.data === "object" && event.data.giscus)) return;

  // Fire your theme function
  // We double-check that the frame exists just to be safe
  const iframe = document.querySelector("iframe.giscus-frame");
  if (iframe) {
    window.changeGiscusTheme();
  }
});
