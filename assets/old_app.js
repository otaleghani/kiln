/**
 * Kiln Logic - app.js
 */

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

window.initMathJax = async function () {
  // Basic heuristic: check for delimiters in text content
  // We search the main content only to be efficient
  const content = document.querySelector(".markdown-body");
  if (!content) return;

  const text = content.innerText;
  // Check for $$ or \[ or \(
  if (!text.includes("$$") && !text.includes("\\(") && !text.includes("\\["))
    return;

  // Lazy Load MathJax
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
      pageReady: () => {
        return window.MathJax.startup.defaultPageReady();
      },
    },
  };

  await window.loadScript(
    "https://cdn.jsdelivr.net/npm/mathjax@4implement-trie-prefix-tree/es5/tex-svg.js",
    "mathjax-script",
  );
};

window.initLucide = async function () {
  // Only load if icons exist
  if (!document.querySelector("[data-lucide]")) return;

  await window.loadScript("https://unpkg.com/lucide@latest", "lucide-script");
  if (window.lucide) window.lucide.createIcons();
};

// Sidebar
window.setupSidebarInteraction = function () {
  const toggleBtns = document.querySelectorAll(".sidebar-toggle.left-toggle");
  const sidebar = document.getElementById("sidebar");
  if (!toggleBtns.length || !sidebar) return;

  toggleBtns.forEach((btn) => {
    const newBtn = btn.cloneNode(true);
    btn.parentNode.replaceChild(newBtn, btn);
    newBtn.addEventListener("click", () => {
      sidebar.classList.toggle("collapsed");
      localStorage.setItem(
        "sidebar-collapsed",
        sidebar.classList.contains("collapsed"),
      );
    });
  });
};

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
        "right-sidebar-collapsed",
        sidebar.classList.contains("collapsed"),
      );
    });
  });
};

window.setupMobileAutoClose = function () {
  document.body.addEventListener("click", (e) => {
    if (window.innerWidth > 768) return;

    const link = e.target.closest("a");
    const isGraphNode =
      ["circle", "text"].includes(e.target.tagName.toLowerCase()) &&
      (e.target.closest("#global-graph-container") ||
        e.target.closest("#local-graph-container"));

    if (!link && !isGraphNode) return;

    // Close Sidebars
    ["sidebar", "right-sidebar"].forEach((id) => {
      const el = document.getElementById(id);
      if (el && el.classList.contains("collapsed")) {
        el.classList.remove("collapsed");
        localStorage.setItem(
          id === "sidebar" ? "sidebar-collapsed" : "right-sidebar-collapsed",
          "true",
        );
      }
    });
  });
};

window.initSidebarSearch = function () {
  const searchInput = document.getElementById("sidebar-search");
  if (!searchInput) return;

  // Remove old listeners by cloning (optional but safer in SPAs)
  const newInput = searchInput.cloneNode(true);
  searchInput.parentNode.replaceChild(newInput, searchInput);

  newInput.addEventListener("input", (e) => {
    const term = e.target.value.toLowerCase().trim();
    const items = document.querySelectorAll(".sidebar.left-sidebar li");

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

// --- 3. Graph & Canvas Logic ---
window.toggleGraphExpand = function () {
  const wrapper = document.getElementById("local-graph-wrapper");
  if (!wrapper) return;
  wrapper.classList.toggle("expanded");

  const icon = wrapper.querySelector('button[title="Expand"] i');
  if (icon) {
    icon.setAttribute(
      "data-lucide",
      wrapper.classList.contains("expanded") ? "minimize-2" : "maximize-2",
    );
    if (window.lucide) window.lucide.createIcons();
  }
  setTimeout(() => {
    window.dispatchEvent(new Event("resize"));
  }, 50);
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

// --- 4. Utilities ---
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

window.highlightSidebarLink = function () {
  document
    .querySelectorAll(".sidebar a")
    .forEach((el) => el.classList.remove("active"));

  const normalize = (p) => {
    if (!p) return "";
    try {
      p = decodeURIComponent(p);
    } catch (e) {}
    return p.length > 1 && p.endsWith("/") ? p.slice(0, -1) : p;
  };

  const currentPath = normalize(window.location.pathname);
  const links = document.querySelectorAll(".sidebar a");

  for (const link of links) {
    const linkPath = normalize(link.getAttribute("href"));
    if (linkPath === currentPath) {
      link.classList.add("active");
      let parent = link.parentElement;
      while (parent) {
        if (parent.tagName === "DETAILS") parent.open = true;
        parent = parent.parentElement;
      }
      break;
    }
  }
};

// This should be in your initAll or DOMContentLoaded
window.initSidebarState = function () {
  const sidebar = document.getElementById("sidebar");
  const rightSidebar = document.getElementById("right-sidebar");

  // Double requestAnimationFrame ensures the DOM has fully painted
  // the collapsed state before we turn on animations.
  requestAnimationFrame(() => {
    requestAnimationFrame(() => {
      if (sidebar) sidebar.classList.add("animate-ready");
      if (rightSidebar) rightSidebar.classList.add("animate-ready");
    });
  });
};

// window.initAll = function() {
//     window.initSidebarState();
//     window.initSidebarSearch();
//     window.setupSidebarInteraction();
//     window.setupRightSidebarInteraction();
//     window.initThemeToggle();
//     window.addCopyButtons();
//     window.highlightSidebarLink();
//
//     // Animate sidebars in only after JS loads
//     const sb = document.getElementById('sidebar');
//     const rsb = document.getElementById('right-sidebar');
//     requestAnimationFrame(() => {
//         if (sb) sb.classList.add('animate-ready');
//         if (rsb) rsb.classList.add('animate-ready');
//     });
// };

window.initAll = function () {
  // UI Setup (Instant)
  window.setupSidebarInteraction();
  window.setupRightSidebarInteraction();
  window.initThemeToggle();
  window.addCopyButtons();
  window.highlightSidebarLink();
  window.initSidebarSearch();

  // Lazy Load Heavy Libs (Conditional)
  // We run these concurrently
  Promise.all([
    window.initMermaid(),
    window.initMathJax(),
    window.initLucide(),
  ]);
};

// --- 5. Event Listeners ---
document.addEventListener("DOMContentLoaded", () => {
  window.initAll();
  window.setupMobileAutoClose();
  // if (window.lucide) window.lucide.createIcons();
});

document.addEventListener("htmx:afterSwap", (evt) => {
  // if (window.lucide) window.lucide.createIcons();
  // if (window.mermaid) window.mermaid.run({ querySelector: '.mermaid' });
  // if (window.MathJax && window.MathJax.typesetPromise) window.MathJax.typesetPromise();
  //
  // window.initAll(); // Re-attach listeners to new DOM elements
});

document.addEventListener("htmx:afterSwap", (evt) => {
  // Re-run checks on new content
  window.initAll();

  // MathJax specific re-render if it was already loaded
  if (window.MathJax && window.MathJax.typesetPromise) {
    window.MathJax.typesetPromise();
  }
});
