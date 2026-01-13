/**
 * Canvas Renderer
 * Handles the rendering of JSON Canvas nodes and edges.
 * Updated to support injected HTML content and Images with floating headers.
 */
window.JsonCanvasRenderer = class JsonCanvasRenderer {
  constructor() {
    this.viewport = document.getElementById("viewport");
    this.world = document.getElementById("canvas-world");
    this.edgesLayer = document.getElementById("edges-layer");
    if (this.edgesLayer) {
      this.edgesLayer.style.zIndex = "50";
      this.edgesLayer.style.pointerEvents = "none";
    }
    this.state = {
      scale: 1,
      panX: 0,
      panY: 0,
      isDragging: false,
      lastMouseX: 0,
      lastMouseY: 0,
    };
    this.data = null;
    this.nodeElements = new Map();

    // Bind methods for removal later
    this.handleWheel = this.handleWheel.bind(this);
    this.handleMouseDown = this.handleMouseDown.bind(this);
    this.handleMouseMove = this.handleMouseMove.bind(this);
    this.handleMouseUp = this.handleMouseUp.bind(this);

    // Bind touch methods
    this.handleTouchStart = this.handleTouchStart.bind(this);
    this.handleTouchMove = this.handleTouchMove.bind(this);
    this.handleTouchEnd = this.handleTouchEnd.bind(this);

    this.setupInteractions();
  }

  cleanup() {
    if (this.viewport) {
      this.viewport.removeEventListener("wheel", this.handleWheel);
      this.viewport.removeEventListener("mousedown", this.handleMouseDown);
      // Remove touch listeners
      this.viewport.removeEventListener("touchstart", this.handleTouchStart);
      this.viewport.removeEventListener("touchmove", this.handleTouchMove);
      this.viewport.removeEventListener("touchend", this.handleTouchEnd);
    }
    window.removeEventListener("mousemove", this.handleMouseMove);
    window.removeEventListener("mouseup", this.handleMouseUp);
  }

  load(data) {
    this.data = data;
    this.render();
    this.centerCanvas();
  }

  render() {
    this.world.querySelectorAll(".node").forEach((el) => el.remove());
    this.edgesLayer.innerHTML = "";
    this.setupSVGDefs();
    this.nodeElements.clear();
    if (!this.data) return;

    // Render groups first so they appear behind content
    const groups = this.data.nodes.filter((n) => n.type === "group");
    groups.forEach((node) => this.createNode(node));

    const nodes = this.data.nodes.filter((n) => n.type !== "group");
    nodes.forEach((node) => this.createNode(node));

    if (this.data.edges)
      this.data.edges.forEach((edge) => this.createEdge(edge));
  }

  setupSVGDefs() {
    // Create defs and marker for the arrowhead
    const defs = document.createElementNS("http://www.w3.org/2000/svg", "defs");
    const marker = document.createElementNS(
      "http://www.w3.org/2000/svg",
      "marker",
    );
    marker.setAttribute("id", "arrow-head");
    marker.setAttribute("viewBox", "0 0 10 10");
    marker.setAttribute("refX", "6"); // Tip of the arrow
    marker.setAttribute("refY", "5"); // Center y
    marker.setAttribute("markerWidth", "4");
    marker.setAttribute("markerHeight", "4");
    marker.setAttribute("orient", "auto");

    const path = document.createElementNS("http://www.w3.org/2000/svg", "path");
    path.setAttribute("d", "M 0 0 L 10 5 L 0 10 z"); // Triangle shape
    //path.setAttribute('d', 'M 0 2.5 L 10 5 L 0 7.5 z');
    path.setAttribute("fill", "#999"); // Default gray color matching edges

    marker.appendChild(path);
    defs.appendChild(marker);
    this.edgesLayer.appendChild(defs);
  }

  slugify(text) {
    return text
      .toString()
      .toLowerCase()
      .trim()
      .replace(/\s+/g, "-") // Replace spaces with -
      .replace(/[^\w\-]+/g, "") // Remove all non-word chars
      .replace(/\-\-+/g, "-"); // Replace multiple - with single -
  }

  createNode(nodeData) {
    const el = document.createElement("div");
    el.id = `node-${nodeData.id}`;
    el.className = `node ${nodeData.type}`;
    el.style.left = `${nodeData.x}px`;
    el.style.top = `${nodeData.y}px`;
    el.style.width = `${nodeData.width}px`;
    el.style.height = `${nodeData.height}px`;

    // Apply color if present
    if (nodeData.color) {
      el.classList.add(`color-${nodeData.color}`);
      el.style.setProperty("--node-color", nodeData.color);
    }

    switch (nodeData.type) {
      case "text":
        el.innerHTML = `<div class="node-note">
                    ${marked.parse(nodeData.text || "")}
                </div>`;
        break;

      case "file":
        // Check if it is an image
        if (
          nodeData.isImage ||
          /\.(png|jpg|jpeg|gif|svg|webp)$/i.test(nodeData.file)
        ) {
          // IMAGE NODE
          const imgSrc = nodeData.src || nodeData.file;
          el.innerHTML = `<div class="canvas-node-image">
                        <img src="${imgSrc}" draggable="false" alt="${nodeData.file}" />
                    </div>`;
          el.style.border = "none";
          el.style.background = "transparent";
          el.style.boxShadow = "none";
        } else {
          // MARKDOWN NOTE NODE
          // Configure the Wrapper (el) to be visible and allow overflow (resets the default node)
          el.style.overflow = "visible";
          el.style.background = "transparent";
          el.style.boxShadow = "none";
          el.style.border = "none";

          // Create node header with link
          const headerHeight = 24;
          const headerEl = document.createElement("div");
          headerEl.className = "canvas-node-header";
          const cleanName = nodeData.file.replace(/\.md$/i, "");
          const parts = cleanName.split("/").map((p) => this.slugify(p));
          const url = "/" + parts.join("/");
          headerEl.innerHTML = `<a href="${url}" target="_blank">${cleanName}</a>`;
          el.appendChild(headerEl);

          // Create the actual note content and handles clipping
          const boxEl = document.createElement("div");
          boxEl.className = "canvas-node-content-box";

          if (nodeData.htmlContent) {
            // Hydrated content (Injected by Generator)
            boxEl.innerHTML = `<div class="node-note">
                            ${nodeData.htmlContent}
                        </div>`;
            boxEl.classList.add("is-note");
          } else if (/\.md$/i.test(nodeData.file)) {
            // Content missing: Show loading state and fetch client-side
            boxEl.innerHTML = `<div class="canvas-node-loading-note">
                            <span>Loading note...</span>
                        </div>`;
            boxEl.classList.add("is-note");
            this.fetchAndRenderNote(nodeData, boxEl);
          } else {
            // Fallback for unknown file types
            boxEl.innerHTML = `<div class="canvas-node-fallback">
                            <span>${nodeData.file}</span>
                        </div>`;
          }

          el.appendChild(boxEl);
        }
        break;

      case "link":
        // UPDATED LINK RENDERING
        el.style.overflow = "visible";
        el.style.background = "transparent";
        el.style.boxShadow = "none";
        el.style.border = "none";

        // Header (Shows URL and allows opening in new tab)
        const linkHeader = document.createElement("div");
        linkHeader.className = "canvas-node-header";
        linkHeader.innerHTML = `<a href="${nodeData.url}" target="_blank">${nodeData.url}</a>`;
        el.appendChild(linkHeader);

        // Content Box (The Iframe)
        const linkContent = document.createElement("div");
        linkContent.className = "canvas-node-content-box";
        // Remove padding for iframes so they fill the box
        linkContent.style.padding = "0";
        linkContent.style.overflow = "hidden";

        // Iframe
        // Note: Some sites (like Google/GitHub) block embedding via X-Frame-Options headers.
        const iframe = document.createElement("iframe");
        iframe.src = nodeData.url;
        iframe.style.width = "100%";
        iframe.style.height = "100%";
        iframe.style.border = "none";
        // Sandbox attributes for security
        iframe.setAttribute(
          "sandbox",
          "allow-scripts allow-same-origin allow-popups allow-forms",
        );

        linkContent.appendChild(iframe);
        el.appendChild(linkContent);
        break;
      // el.innerHTML = `<a href="${nodeData.url}" target="_blank" class="canvas-node-link"><span>${nodeData.url}</span></a>`;
      // break;

      case "group":
        if (nodeData.label) {
          const label = document.createElement("div");
          label.className = "canvas-group-label";
          label.textContent = nodeData.label;
          el.appendChild(label);
        }
        el.style.backgroundColor = "rgba(0,0,0,0.02)";
        el.style.border = "2px dashed rgba(0,0,0,0.1)";
        el.style.pointerEvents = "none"; // Let clicks pass through to nodes below/inside
        break;
    }

    this.world.appendChild(el);
    this.nodeElements.set(nodeData.id, el);
  }

  /**
   * Client-side fetcher for when the generator hasn't pre-injected content.
   */
  async fetchAndRenderNote(nodeData, container) {
    try {
      const cleanPath = nodeData.file.replace(/\.md$/i, "");
      const parts = cleanPath.split("/").map((p) => this.slugify(p));
      const slugPath = parts.join("/");

      // Try standard index.html path first
      const url = `/${slugPath}/index.html`;

      const response = await fetch(url);

      if (!response.ok) {
        // Try alternate URL format
        const altResponse = await fetch(`/${slugPath}`);
        if (!altResponse.ok) throw new Error("Note not found");
        return this.processResponse(altResponse, container);
      }

      await this.processResponse(response, container);
    } catch (err) {
      console.warn("Canvas: Could not fetch embedded note", nodeData.file, err);
      container.innerHTML = `<div class="flex flex-col items-center justify-center h-full text-gray-500 p-4">
                 <i class="ph ph-warning-circle text-2xl mb-2 text-red-400"></i>
                 <span class="text-xs text-center">Unable to load content</span>
            </div>`;
    }
  }

  async processResponse(response, container) {
    const html = await response.text();
    const parser = new DOMParser();
    const doc = parser.parseFromString(html, "text/html");

    // Find content
    const content = doc.querySelector("#content");

    if (content) {
      // Strip out TOC if it exists inside the content to save space
      const toc = content.querySelector(".toc");
      if (toc) toc.remove();

      // Render content only (Header is already handled in createNode)
      container.innerHTML = `<div class="content">${content.innerHTML}</div>`;
    }
  }

  createEdge(edgeData) {
    const fromEl = this.nodeElements.get(edgeData.fromNode);
    const toEl = this.nodeElements.get(edgeData.toNode);
    if (!fromEl || !toEl) return;

    const fromNode = this.data.nodes.find((n) => n.id === edgeData.fromNode);
    const toNode = this.data.nodes.find((n) => n.id === edgeData.toNode);
    const start = this.getAnchorPoint(fromNode, edgeData.fromSide);
    const end = this.getAnchorPoint(toNode, edgeData.toSide);

    const path = document.createElementNS("http://www.w3.org/2000/svg", "path");
    path.classList.add("edge");

    if (edgeData.color) {
      path.setAttribute("stroke", edgeData.color);
    } else {
      path.setAttribute("stroke", "#999"); // Default stroke color if none provided
    }

    // Add arrow marker to the end of the line (default for most edges)
    path.setAttribute("marker-end", "url(#arrow-head)");

    const cp1 = this.getControlPoint(start, edgeData.fromSide);
    const cp2 = this.getControlPoint(end, edgeData.toSide);
    const d = `M ${start.x} ${start.y} C ${cp1.x} ${cp1.y}, ${cp2.x} ${cp2.y}, ${end.x} ${end.y}`;
    path.setAttribute("d", d);
    this.edgesLayer.appendChild(path);
  }

  getAnchorPoint(node, side) {
    let x = node.x + node.width / 2;
    let y = node.y + node.height / 2;
    if (!side) return { x, y };
    switch (side) {
      case "top":
        y = node.y;
        break;
      case "bottom":
        y = node.y + node.height;
        break;
      case "left":
        x = node.x;
        break;
      case "right":
        x = node.x + node.width;
        break;
    }
    return { x, y };
  }

  getControlPoint(point, side) {
    const curvature = 100;
    let x = point.x,
      y = point.y;
    if (!side) return { x, y };
    switch (side) {
      case "top":
        y -= curvature;
        break;
      case "bottom":
        y += curvature;
        break;
      case "left":
        x -= curvature;
        break;
      case "right":
        x += curvature;
        break;
    }
    return { x, y };
  }

  handleWheel(e) {
    e.preventDefault();
    const delta = -e.deltaY * 0.001;
    this.state.scale = Math.min(Math.max(0.1, this.state.scale + delta), 5);
    this.updateTransform();
  }

  handleMouseDown(e) {
    if (e.target.closest(".md-content a")) return;

    this.state.isDragging = true;
    this.state.lastMouseX = e.clientX;
    this.state.lastMouseY = e.clientY;
    this.viewport.style.cursor = "grabbing";
  }

  handleMouseMove(e) {
    if (!this.state.isDragging) return;
    e.preventDefault();
    const dx = e.clientX - this.state.lastMouseX;
    const dy = e.clientY - this.state.lastMouseY;
    this.state.panX += dx;
    this.state.panY += dy;
    this.state.lastMouseX = e.clientX;
    this.state.lastMouseY = e.clientY;
    this.updateTransform();
  }

  handleMouseUp() {
    this.state.isDragging = false;
    if (this.viewport) this.viewport.style.cursor = "default";
  }

  handleTouchStart(e) {
    // Allow standard clicks on links/buttons
    if (e.target.closest("a, button, input, .canvas-node-header")) return;

    if (e.touches.length === 1) {
      // Single touch = Drag
      this.state.isDragging = true;
      this.state.lastMouseX = e.touches[0].clientX;
      this.state.lastMouseY = e.touches[0].clientY;
    } else if (e.touches.length === 2) {
      // Two fingers = Start Pinch Zoom
      this.state.isDragging = false;
      this.state.initialPinchDistance = this.getPinchDistance(e);
      this.state.initialScale = this.state.scale;
    }
  }

  handleTouchMove(e) {
    // 1. Handle Panning (1 finger)
    if (this.state.isDragging && e.touches.length === 1) {
      e.preventDefault(); // Stop the browser from scrolling the page

      const dx = e.touches[0].clientX - this.state.lastMouseX;
      const dy = e.touches[0].clientY - this.state.lastMouseY;

      this.state.panX += dx;
      this.state.panY += dy;

      this.state.lastMouseX = e.touches[0].clientX;
      this.state.lastMouseY = e.touches[0].clientY;

      this.updateTransform();
    }
    // 2. Handle Zooming (2 fingers)
    else if (e.touches.length === 2 && this.state.initialPinchDistance) {
      e.preventDefault();

      const currentDist = this.getPinchDistance(e);
      const zoomFactor = currentDist / this.state.initialPinchDistance;

      // Calculate new scale based on initial scale * zoom factor
      let newScale = this.state.initialScale * zoomFactor;
      newScale = Math.min(Math.max(0.1, newScale), 5); // Clamp limits

      this.state.scale = newScale;
      this.updateTransform();
    }
  }

  handleTouchEnd(e) {
    // If all fingers lifted, stop dragging
    if (e.touches.length === 0) {
      this.state.isDragging = false;
      this.state.initialPinchDistance = null;
    }
  }

  // Helper to calculate distance between two fingers
  getPinchDistance(e) {
    return Math.hypot(
      e.touches[0].clientX - e.touches[1].clientX,
      e.touches[0].clientY - e.touches[1].clientY,
    );
  }

  updateTransform() {
    this.world.style.transform = `translate(${this.state.panX}px, ${this.state.panY}px) scale(${this.state.scale})`;
  }

  setupInteractions() {
    this.viewport.addEventListener("wheel", this.handleWheel, {
      passive: false,
    });
    this.viewport.addEventListener("mousedown", this.handleMouseDown);
    window.addEventListener("mousemove", this.handleMouseMove);
    window.addEventListener("mouseup", this.handleMouseUp);

    // Passive: false is required to use preventDefault() to stop page scrolling
    this.viewport.addEventListener("touchstart", this.handleTouchStart, {
      passive: false,
    });
    window.addEventListener("touchmove", this.handleTouchMove, {
      passive: false,
    });
    window.addEventListener("touchend", this.handleTouchEnd);
  }

  updateTransform() {
    this.world.style.transform = `translate(${this.state.panX}px, ${this.state.panY}px) scale(${this.state.scale})`;
  }

  zoomIn() {
    this.state.scale = Math.min(this.state.scale + 0.2, 5);
    this.updateTransform();
  }
  zoomOut() {
    this.state.scale = Math.max(0.1, this.state.scale - 0.2);
    this.updateTransform();
  }
  resetView() {
    this.centerCanvas();
  }

  centerCanvas() {
    if (!this.data || !this.data.nodes || this.data.nodes.length === 0) return;
    let minX = Infinity,
      minY = Infinity,
      maxX = -Infinity,
      maxY = -Infinity;
    this.data.nodes.forEach((n) => {
      minX = Math.min(minX, n.x);
      minY = Math.min(minY, n.y);
      maxX = Math.max(maxX, n.x + n.width);
      maxY = Math.max(maxY, n.y + n.height);
    });
    const contentWidth = maxX - minX;
    const contentHeight = maxY - minY;
    const viewWidth = this.viewport.clientWidth;
    const viewHeight = this.viewport.clientHeight;
    const padding = 100;
    const scaleX = (viewWidth - padding * 2) / contentWidth;
    const scaleY = (viewHeight - padding * 2) / contentHeight;
    this.state.scale = Math.min(Math.min(scaleX, scaleY), 1);
    this.state.panX =
      (viewWidth - contentWidth * this.state.scale) / 2 -
      minX * this.state.scale;
    this.state.panY =
      (viewHeight - contentHeight * this.state.scale) / 2 -
      minY * this.state.scale;
    this.updateTransform();
  }
};
