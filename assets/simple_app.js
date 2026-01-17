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

window.initToggles = function () {
  const menuBtn = document.getElementById("menu-button");
  const menuIcon = document.getElementById("menu-icon");
  const menu = document.getElementById("menu-wrapper");

  window.toggleMenu = () => {
    window.closeAllExcept("menu");
    if (menu.classList.contains("hidden")) {
      menu.classList.remove("hidden");
      menuIcon.classList.add("text-accent");
      menuIcon.classList.remove("text-foreground");
    } else {
      menu.classList.add("hidden");
      menuIcon.classList.remove("text-accent");
      menuIcon.classList.add("text-foreground");
    }
  };

  const localGraphBtn = document.getElementById("local-graph-button");
  const localGraphWrapper = document.getElementById("local-graph-wrapper");
  const localGraphIcon = document.getElementById("local-graph-icon");
  window.toggleLocalGraph = () => {
    window.closeAllExcept("local-graph");
    if (localGraphWrapper.classList.contains("hidden")) {
      localGraphWrapper.classList.remove("hidden");
      localGraphIcon.classList.add("text-accent");
      localGraphIcon.classList.remove("text-foreground");
    } else {
      localGraphWrapper.classList.add("hidden");
      localGraphIcon.classList.remove("text-accent");
      localGraphIcon.classList.add("text-foreground");
    }
  };

  const tocBtn = document.getElementById("toc-button");
  const tocWrapper = document.getElementById("toc-wrapper");
  const tocIcon = document.getElementById("toc-icon");
  window.toggleTOC = () => {
    window.closeAllExcept("toc");
    if (tocWrapper.classList.contains("hidden")) {
      tocWrapper.classList.remove("hidden");
      tocIcon.classList.add("text-accent");
      tocIcon.classList.remove("text-foreground");
    } else {
      tocWrapper.classList.add("hidden");
      tocIcon.classList.remove("text-accent");
      tocIcon.classList.add("text-foreground");
    }
  };

  window.closeAllExcept = (item) => {
    if (item !== "menu" && menu && menuIcon) {
      menu.classList.add("hidden");
      menuIcon.classList.remove("text-accent");
      menuIcon.classList.add("text-foreground");
    }

    if (item !== "local-graph" && localGraphWrapper && localGraphIcon) {
      localGraphWrapper.classList.add("hidden");
      localGraphIcon.classList.remove("text-accent");
      localGraphIcon.classList.add("text-foreground");
    }

    if (item !== "toc" && tocWrapper && tocIcon) {
      tocWrapper.classList.add("hidden");
      tocIcon.classList.remove("text-accent");
      tocIcon.classList.add("text-foreground");
    }
  };

  if (localGraphBtn) {
    localGraphBtn.addEventListener("click", window.toggleLocalGraph);
  }
  if (menuBtn) {
    menuBtn.addEventListener("click", window.toggleMenu);
  }
  if (tocBtn) {
    tocBtn.addEventListener("click", window.toggleTOC);
  }
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
  window.changeGiscusTheme();
  window.initThemeToggle();
  window.initToggles();
  window.addCopyButtons();
  Promise.all([window.initMathJax(), window.initMermaid()]);
};

document.addEventListener("DOMContentLoaded", () => {
  localStorage.setItem("menu", "close");
  window.initAll();
});

document.addEventListener("htmx:afterSwap", () => {
  localStorage.setItem("menu", "close");
  window.initAll();
});
