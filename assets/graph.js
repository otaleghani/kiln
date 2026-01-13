(function () {
  // 1. Define Base URL globally for the script (Removing trailing slash if present)
  const BASE_URL = "{{.BaseURL}}".replace(/\/$/, "");

  // Cache the data promise to prevent re-fetching on every HTMX navigation
  let graphDataPromise = null;

  // Function to coordinate initialization
  function setupGraph() {
    const globalContainer = document.getElementById("global-graph-container");
    const localContainer = document.getElementById("local-graph-container");
    const localContainerWrapper = document.getElementById(
      "local-graph-wrapper",
    );

    // Only load if a graph container is present
    if (!globalContainer && !localContainer) return;

    // Load D3 from CDN if not already present
    if (typeof d3 === "undefined") {
      // Check if we are already loading it to prevent duplicates
      if (!document.getElementById("d3-script")) {
        const script = document.createElement("script");
        script.id = "d3-script";
        script.src = "https://d3js.org/d3.v7.min.js";
        script.onload = () =>
          initGraph(globalContainer, localContainer, localContainerWrapper);
        script.onerror = () =>
          console.error(
            "Failed to load D3.js. Please check your internet connection.",
          );
        document.head.appendChild(script);
      }
    } else {
      // D3 already loaded, run immediately
      initGraph(globalContainer, localContainer, localContainerWrapper);
    }
  }

  // Listen for initial load
  document.addEventListener("DOMContentLoaded", setupGraph);

  // Listen for HTMX content swaps
  document.addEventListener("htmx:afterSwap", setupGraph);
  document.addEventListener("htmx:historyRestore", setupGraph);

  // --- THEME CHANGE LISTENER ---
  const themeObserver = new MutationObserver((mutations) => {
    let shouldUpdate = false;
    mutations.forEach((m) => {
      if (
        m.type === "attributes" &&
        (m.attributeName === "data-theme" || m.attributeName === "class")
      ) {
        shouldUpdate = true;
      }
    });
    if (shouldUpdate) {
      const global = document.getElementById("global-graph-container");
      const local = document.getElementById("local-graph-container");
      const localWrapper = document.getElementById("local-graph-wrapper");
      if (global || local) {
        initGraph(global, local, localWrapper);
      }
    }
  });

  themeObserver.observe(document.documentElement, {
    attributes: true,
    attributeFilter: ["data-theme", "class"],
  });

  if (window.matchMedia) {
    window
      .matchMedia("(prefers-color-scheme: dark)")
      .addEventListener("change", () => {
        const global = document.getElementById("global-graph-container");
        const local = document.getElementById("local-graph-container");
        const localWrapper = document.getElementById("local-graph-wrapper");
        if (global || local) {
          initGraph(global, local, localWrapper);
        }
      });
  }

  function initGraph(globalContainer, localContainer, localContainerWrapper) {
    // Initialize cache if empty
    if (!graphDataPromise) {
      // Use the injected BaseURL for fetching
      graphDataPromise = fetch(BASE_URL + "/graph.json")
        .then((res) => {
          if (!res.ok) throw new Error(`HTTP error! status: ${res.status}`);
          return res.json();
        })
        .then((data) => {
          // Pre-process: Filter out links pointing to non-existent nodes
          const nodeIds = new Set(data.nodes.map((n) => n.id));
          data.links = data.links.filter(
            (l) => nodeIds.has(l.source) && nodeIds.has(l.target),
          );

          // Calculate Degrees (Connections per node)
          const degreeMap = {};
          data.links.forEach((l) => {
            degreeMap[l.source] = (degreeMap[l.source] || 0) + 1;
            degreeMap[l.target] = (degreeMap[l.target] || 0) + 1;
          });

          data.nodes.forEach((n) => {
            n.degree = degreeMap[n.id] || 0;
          });

          return data;
        })
        .catch((err) => {
          console.error("Graph loading failed:", err);
          graphDataPromise = null;
          if (window.location.protocol === "file:") {
            console.warn("Fetch API blocked by file:// protocol.");
          }
        });
    }

    graphDataPromise.then((data) => {
      if (!data) return;

      if (globalContainer) {
        renderGraph(globalContainer, JSON.parse(JSON.stringify(data)), false);
      }
      if (localContainer) {
        const pageTitleEl = document.getElementById("page-title-data");
        if (pageTitleEl) {
          const currentTitle = pageTitleEl.dataset.title;
          const localData = filterLocalData(
            JSON.parse(JSON.stringify(data)),
            currentTitle,
          );
          if (localData.nodes.length > 0) {
            localContainerWrapper.style.display = "";
            renderGraph(localContainer, localData, true);
          } else {
            localContainerWrapper.style.display = "none";
          }
        }
      }
    });
  }

  function filterLocalData(data, currentId) {
    const linkedIds = new Set();
    linkedIds.add(currentId);

    // Identify neighbors
    const validLinks = data.links.filter((l) => {
      const isSource = l.source === currentId;
      const isTarget = l.target === currentId;
      if (isSource) linkedIds.add(l.target);
      if (isTarget) linkedIds.add(l.source);
      return isSource || isTarget;
    });

    // Filter nodes based on neighbors
    const validNodes = data.nodes.filter((n) => linkedIds.has(n.id));

    return { nodes: validNodes, links: validLinks };
  }

  function renderGraph(container, data, isLocal) {
    const rect = container.getBoundingClientRect();
    const width = rect.width || 800;
    const height = rect.height || 600;

    // --- THEME INTEGRATION ---
    const style = getComputedStyle(document.documentElement);
    const accentColor =
      style.getPropertyValue("--accent-color").trim() || "#7e6df7";
    const textColor = style.getPropertyValue("--text-color").trim() || "#ccc";
    const neutralNodeColor =
      style.getPropertyValue("--color-comment").trim() || "#888";
    const neutralLinkColor =
      style.getPropertyValue("--sidebar-border").trim() || "#999";
    const folderNodeColor =
      style.getPropertyValue("--color-yellow").trim() || "#FFD700";
    const assetNodeColor =
      style.getPropertyValue("--color-red").trim() || "#FF0000";

    container.innerHTML = "";

    const zoom = d3
      .zoom()
      .scaleExtent([0.1, 4])
      .on("zoom", (event) => {
        svgGroup.attr("transform", event.transform);
      });

    const svg = d3
      .select(container)
      .append("svg")
      .attr("width", "100%")
      .attr("height", "100%")
      .attr("viewBox", [0, 0, width, height])
      .call(zoom)
      .on("dblclick.zoom", null);

    const svgGroup = svg.append("g");

    // Helper to calculate radius
    const getNodeRadius = (d) => 4 + Math.sqrt(d.degree || 0) * 2;

    // Force Simulation
    const simulation = d3
      .forceSimulation(data.nodes)
      .force(
        "link",
        d3
          .forceLink(data.links)
          .id((d) => d.id)
          .distance(isLocal ? 100 : 50),
      )
      .force("charge", d3.forceManyBody().strength(isLocal ? -300 : -100))
      .force("center", d3.forceCenter(width / 2, height / 2))
      .force(
        "collide",
        d3.forceCollide().radius((d) => getNodeRadius(d) + 2),
      );

    const link = svgGroup
      .append("g")
      .attr("stroke", neutralLinkColor)
      .attr("stroke-opacity", 0.6)
      .selectAll("line")
      .data(data.links)
      .join("line")
      .attr("stroke-width", 1);

    // Create Anchors (<a>) first
    const nodeGroup = svgGroup
      .append("g")
      .selectAll("a")
      .data(data.nodes)
      .join("a")
      .attr("cursor", "pointer")
      // 1. Calculate the Href
      .attr("href", (d) => {
        const targetPath = d.url.startsWith("/") ? d.url : "/" + d.url;
        if (BASE_URL && targetPath.startsWith(BASE_URL)) {
          return targetPath;
        }
        return BASE_URL + targetPath;
      });
    // 2. Apply HTMX Attributes directly to the link
    // .attr("hx-boost", "true")
    // .attr("hx-target", "#kiln-main")
    // .attr("hx-swap", "outerHTML")
    // .attr("hx-select", "#kiln-main, #right-sidebar-content")
    // .attr("hx-push-url", "true"); // Explicitly ensure URL is pushed

    // Append Circles to the Anchors
    // We assign this to 'node' so the tick function and hover effects still work on the circle
    const node = nodeGroup
      .append("circle")
      .attr("r", (d) => getNodeRadius(d))
      .attr("fill", (d) =>
        d.type == "folder"
          ? folderNodeColor
          : d.type == ".md" || d.type == ".canvas" || d.type == ".base"
            ? neutralNodeColor
            : assetNodeColor,
      )
      .call(drag(simulation));

    const label = svgGroup
      .append("g")
      .selectAll("text")
      .data(data.nodes)
      .join("text")
      .attr("dx", (d) => getNodeRadius(d) + 4)
      .attr("dy", ".35em")
      .text((d) => d.label)
      .style("fill", textColor)
      .style("font-size", "10px")
      .style("pointer-events", "none")
      .style("opacity", isLocal ? 1 : 0.7);

    // Hover Effects (Unchanged)
    node
      .on("mouseover", function (event, d) {
        const currentR = getNodeRadius(d);
        d3.select(this)
          .attr("r", currentR * 1.2)
          .attr("fill", accentColor);

        link.style("stroke", (l) =>
          l.source === d || l.target === d ? accentColor : neutralLinkColor,
        );
        link.style("stroke-opacity", (l) =>
          l.source === d || l.target === d ? 1 : 0.2,
        );
        link.style("stroke-width", (l) =>
          l.source === d || l.target === d ? 2 : 1,
        );

        label
          .filter((l) => l === d)
          .style("opacity", 1)
          .style("font-weight", "bold");
      })
      .on("mouseout", function (event, d) {
        const currentR = getNodeRadius(d);
        d3.select(this)
          .attr("r", currentR)
          .attr(
            "fill",
            d.type == "folder"
              ? folderNodeColor
              : d.type == ".md" || d.type == ".canvas" || d.type == ".base"
                ? neutralNodeColor
                : assetNodeColor,
          );

        link.style("stroke", neutralLinkColor);
        link.style("stroke-opacity", 0.6);
        link.attr("stroke-width", 1);

        label
          .filter((l) => l === d)
          .style("opacity", isLocal ? 1 : 0.7)
          .style("font-weight", "normal");
      });

    // REMOVED: The manual .on("click") handler is gone.
    // The <a> tag handles it natively now.

    simulation.on("tick", () => {
      link
        .attr("x1", (d) => d.source.x)
        .attr("y1", (d) => d.source.y)
        .attr("x2", (d) => d.target.x)
        .attr("y2", (d) => d.target.y);

      node.attr("cx", (d) => d.x).attr("cy", (d) => d.y);

      label.attr("x", (d) => d.x).attr("y", (d) => d.y);
    });

    // --- NEW: Tell HTMX to process the new DOM elements ---
    if (typeof htmx !== "undefined") {
      htmx.process(container);
    }

    function drag(simulation) {
      function dragstarted(event) {
        if (!event.active) simulation.alphaTarget(0.3).restart();
        event.subject.fx = event.subject.x;
        event.subject.fy = event.subject.y;
      }

      function dragged(event) {
        event.subject.fx = event.x;
        event.subject.fy = event.y;
      }

      function dragended(event) {
        if (!event.active) simulation.alphaTarget(0);
        event.subject.fx = null;
        event.subject.fy = null;
      }

      return d3
        .drag()
        .on("start", dragstarted)
        .on("drag", dragged)
        .on("end", dragended);
    }

    if (typeof htmx !== "undefined") {
      htmx.process(container);
    }
  }
})();
