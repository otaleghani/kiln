// Client-side JS for simple layout: sidebar, search, theme toggle, and local graph. @feature:layouts
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

// Panel definitions: each maps a button ID to its wrapper and icon IDs.
window._panels = [
  { btn: "menu-button", wrapper: "menu-wrapper", icon: "menu-icon" },
  {
    btn: "local-graph-button",
    wrapper: "local-graph-wrapper",
    icon: "local-graph-icon",
  },
  { btn: "toc-button", wrapper: "toc-wrapper", icon: "toc-icon" },
];

// Close a single panel by looking up its current DOM elements.
window._closePanel = function (panel) {
  const wrapper = document.getElementById(panel.wrapper);
  const icon = document.getElementById(panel.icon);
  if (wrapper) wrapper.classList.add("hidden");
  if (icon) {
    icon.classList.remove("text-accent");
    icon.classList.add("text-foreground");
  }
};

// Toggle a panel: close all others, then toggle the target.
window._togglePanel = function (panel) {
  window._panels.forEach((p) => {
    if (p.wrapper !== panel.wrapper) window._closePanel(p);
  });
  const wrapper = document.getElementById(panel.wrapper);
  const icon = document.getElementById(panel.icon);
  if (!wrapper) return;
  const opening = wrapper.classList.contains("hidden");
  wrapper.classList.toggle("hidden");
  if (icon) {
    icon.classList.toggle("text-accent", opening);
    icon.classList.toggle("text-foreground", !opening);
  }
};

// Named toggle functions referenced by inline scripts (e.g. TOC link clicks).
window.toggleMenu = () =>
  window._togglePanel(window._panels[0]);
window.toggleLocalGraph = () =>
  window._togglePanel(window._panels[1]);
window.toggleTOC = () =>
  window._togglePanel(window._panels[2]);

// Event delegation: one listener that survives htmx swaps.
window.initToggles = function () {
  if (window._togglesDelegationBound) return;
  window._togglesDelegationBound = true;

  document.body.addEventListener("click", (e) => {
    for (const panel of window._panels) {
      if (e.target.closest("#" + panel.btn)) {
        window._togglePanel(panel);
        return;
      }
    }
  });
};

// Lazyload MathJax
window.initMathJax = async function () {
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

// Calls every init function
window.initAll = function () {
  window.initThemeToggle();
  window.initToggles();
  window.addCopyButtons();
  Promise.all([window.initMathJax(), window.initMermaid()]);
};

document.addEventListener("DOMContentLoaded", () => {
  window.initAll();
});

document.addEventListener("htmx:afterSwap", () => {
  // Close all panels on navigation
  window._panels.forEach((p) => window._closePanel(p));
  window.initAll();
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
