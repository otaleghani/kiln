(function() {
    // Cache the data promise to prevent re-fetching on every HTMX navigation
    let graphDataPromise = null;

    // Function to coordinate initialization
    function setupGraph() {
        const globalContainer = document.getElementById('global-graph-container');
        const localContainer = document.getElementById('local-graph-container');

        // Only load if a graph container is present
        if (!globalContainer && !localContainer) return;

        // Load D3 from CDN if not already present
        if (typeof d3 === 'undefined') {
            // Check if we are already loading it to prevent duplicates
            if (!document.getElementById('d3-script')) {
                const script = document.createElement('script');
                script.id = 'd3-script';
                script.src = 'https://d3js.org/d3.v7.min.js';
                script.onload = () => initGraph(globalContainer, localContainer);
                script.onerror = () => console.error("Failed to load D3.js. Please check your internet connection.");
                document.head.appendChild(script);
            }
        } else {
            // D3 already loaded, run immediately
            initGraph(globalContainer, localContainer);
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
        mutations.forEach(m => {
            if (m.type === 'attributes' && (m.attributeName === 'data-theme' || m.attributeName === 'class')) {
                shouldUpdate = true;
            }
        });
        if (shouldUpdate) {
            const global = document.getElementById('global-graph-container');
            const local = document.getElementById('local-graph-container');
            if (global || local) {
                initGraph(global, local);
            }
        }
    });
    
    themeObserver.observe(document.documentElement, { attributes: true, attributeFilter: ['data-theme', 'class'] });

    if (window.matchMedia) {
        window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', () => {
             const global = document.getElementById('global-graph-container');
            const local = document.getElementById('local-graph-container');
            if (global || local) {
                initGraph(global, local);
            }
        });
    }

    function initGraph(globalContainer, localContainer) {
        // Initialize cache if empty
        if (!graphDataPromise) {
            graphDataPromise = fetch('/graph.json')
                .then(res => {
                    if (!res.ok) throw new Error(`HTTP error! status: ${res.status}`);
                    return res.json();
                })
                .then(data => {
                    // 1. Pre-process: Filter out links pointing to non-existent nodes
                    const nodeIds = new Set(data.nodes.map(n => n.id));
                    data.links = data.links.filter(l => nodeIds.has(l.source) && nodeIds.has(l.target));

                    // 2. Calculate Degrees (Connections per node)
                    const degreeMap = {};
                    data.links.forEach(l => {
                        degreeMap[l.source] = (degreeMap[l.source] || 0) + 1;
                        degreeMap[l.target] = (degreeMap[l.target] || 0) + 1;
                    });
                    
                    data.nodes.forEach(n => {
                        n.degree = degreeMap[n.id] || 0;
                    });

                    return data;
                })
                .catch(err => {
                    console.error("Graph loading failed:", err);
                    graphDataPromise = null; 
                    if (window.location.protocol === 'file:') {
                        console.warn("Fetch API blocked by file:// protocol.");
                    }
                });
        }

        graphDataPromise.then(data => {
            if (!data) return;

            if (globalContainer) {
                renderGraph(globalContainer, JSON.parse(JSON.stringify(data)), false);
            }
            if (localContainer) {
                const pageTitleEl = document.getElementById('page-title-data');
                if (pageTitleEl) {
                    const currentTitle = pageTitleEl.dataset.title;
                    const localData = filterLocalData(JSON.parse(JSON.stringify(data)), currentTitle);
                    if (localData.nodes.length > 0) {
                        localContainer.style.display = ''; 
                        renderGraph(localContainer, localData, true);
                    } else {
                        localContainer.style.display = 'none';
                    }
                }
            }
        });
    }

    function filterLocalData(data, currentId) {
        const linkedIds = new Set();
        linkedIds.add(currentId);
        
        // 1. Identify neighbors
        const validLinks = data.links.filter(l => {
            const isSource = l.source === currentId;
            const isTarget = l.target === currentId;
            if (isSource) linkedIds.add(l.target);
            if (isTarget) linkedIds.add(l.source);
            return isSource || isTarget;
        });

        // 2. Filter nodes based on neighbors
        const validNodes = data.nodes.filter(n => linkedIds.has(n.id));

        return { nodes: validNodes, links: validLinks };
    }

    function renderGraph(container, data, isLocal) {
        const rect = container.getBoundingClientRect();
        const width = rect.width || 800;
        const height = rect.height || 600;

        // --- THEME INTEGRATION ---
        const style = getComputedStyle(document.documentElement);
        
        const accentColor = style.getPropertyValue('--accent-color').trim() || '#7e6df7';
        const textColor = style.getPropertyValue('--text-color').trim() || '#ccc';
        const neutralNodeColor = style.getPropertyValue('--color-comment').trim() || '#888';
        const neutralLinkColor = style.getPropertyValue('--sidebar-border').trim() || '#999';

        container.innerHTML = '';

        const zoom = d3.zoom()
            .scaleExtent([0.1, 4])
            .on("zoom", (event) => {
                svgGroup.attr("transform", event.transform);
            });

        const svg = d3.select(container).append("svg")
            .attr("width", "100%")
            .attr("height", "100%")
            .attr("viewBox", [0, 0, width, height])
            .call(zoom)
            .on("dblclick.zoom", null);

        const svgGroup = svg.append("g");

        // Helper to calculate radius based on connections
        // Base size 4 + square root of degree * 2 (prevents huge nodes)
        const getNodeRadius = (d) => 4 + (Math.sqrt(d.degree || 0) * 2);

        // Force Simulation
        const simulation = d3.forceSimulation(data.nodes)
            .force("link", d3.forceLink(data.links).id(d => d.id).distance(isLocal ? 100 : 50))
            .force("charge", d3.forceManyBody().strength(isLocal ? -300 : -100))
            .force("center", d3.forceCenter(width / 2, height / 2))
            .force("collide", d3.forceCollide().radius(d => getNodeRadius(d) + 2)); // Dynamic collision radius

        const link = svgGroup.append("g")
            .attr("stroke", neutralLinkColor)
            .attr("stroke-opacity", 0.6)
            .selectAll("line")
            .data(data.links)
            .join("line")
            .attr("stroke-width", 1);

        const node = svgGroup.append("g")
            .selectAll("circle")
            .data(data.nodes)
            .join("circle")
            .attr("r", d => getNodeRadius(d)) 
            .attr("fill", neutralNodeColor) // Always start neutral
            .call(drag(simulation));

        const label = svgGroup.append("g")
            .selectAll("text")
            .data(data.nodes)
            .join("text")
            .attr("dx", d => getNodeRadius(d) + 4)
            .attr("dy", ".35em")
            .text(d => d.label)
            .style("fill", textColor)
            .style("font-size", "10px")
            .style("pointer-events", "none")
            .style("opacity", isLocal ? 1 : 0.7);

        // Hover & Interaction
        node.on("mouseover", function(event, d) {
            // Highlight Node (Scale up slightly)
            const currentR = getNodeRadius(d);
            d3.select(this)
                .attr("r", currentR * 1.2)
                .attr("fill", accentColor);
            
            // Highlight Connected Links
            link.style('stroke', l => (l.source === d || l.target === d) ? accentColor : neutralLinkColor);
            link.style('stroke-opacity', l => (l.source === d || l.target === d) ? 1 : 0.2);
            link.attr('stroke-width', l => (l.source === d || l.target === d) ? 2 : 1);

            // Highlight Label
            label.filter(l => l === d).style("opacity", 1).style("font-weight", "bold");
        })
        .on("mouseout", function(event, d) {
            // Reset Node
            const currentR = getNodeRadius(d);
            d3.select(this)
                .attr("r", currentR)
                .attr("fill", neutralNodeColor);
            
            // Reset Links
            link.style('stroke', neutralLinkColor);
            link.style('stroke-opacity', 0.6);
            link.attr('stroke-width', 1);

            // Reset Label
            label.filter(l => l === d).style("opacity", isLocal ? 1 : 0.7).style("font-weight", "normal");
        })
        .on("click", function(event, d) {
            window.location.href = d.url;
        });

        simulation.on("tick", () => {
            link
                .attr("x1", d => d.source.x)
                .attr("y1", d => d.source.y)
                .attr("x2", d => d.target.x)
                .attr("y2", d => d.target.y);

            node
                .attr("cx", d => d.x)
                .attr("cy", d => d.y);
            
            label
                .attr("x", d => d.x)
                .attr("y", d => d.y);
        });

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

            return d3.drag()
                .on("start", dragstarted)
                .on("drag", dragged)
                .on("end", dragended);
        }
    }
})();
