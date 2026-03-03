// @feature:search Client-side full-text search with inverted index and modal overlay results.
(function () {
  var BASE_URL = "{{.BaseURL}}";
  var MAX_RESULTS = 15;
  var SNIPPET_LEN = 120;

  var invertedIndex = null;
  var indexEntries = null;

  function fetchIndex() {
    if (window._searchIndex) {
      return Promise.resolve(window._searchIndex);
    }
    return fetch(BASE_URL + "/search-index.json")
      .then(function (res) { return res.json(); })
      .then(function (data) {
        window._searchIndex = data;
        return data;
      });
  }

  function buildInvertedIndex(entries) {
    var idx = {};
    for (var i = 0; i < entries.length; i++) {
      var text = ((entries[i].title || "") + " " + (entries[i].content || "")).toLowerCase();
      var words = text.split(/\s+/);
      for (var w = 0; w < words.length; w++) {
        var word = words[w].replace(/[^a-z0-9]/g, "");
        if (!word) continue;
        if (!idx[word]) idx[word] = new Set();
        idx[word].add(i);
      }
    }
    return idx;
  }

  function searchEntries(query) {
    if (!indexEntries) return [];
    if (!query || query.trim().length === 0) {
      return indexEntries.slice(0, MAX_RESULTS);
    }

    if (!invertedIndex) return [];
    var tokens = query.toLowerCase().split(/\s+/).filter(function (t) {
      return t.replace(/[^a-z0-9]/g, "").length > 0;
    });
    if (tokens.length === 0) return indexEntries.slice(0, MAX_RESULTS);

    var resultSet = null;
    var indexKeys = Object.keys(invertedIndex);

    for (var t = 0; t < tokens.length; t++) {
      var token = tokens[t].replace(/[^a-z0-9]/g, "");
      if (!token) continue;
      var matching = new Set();
      for (var k = 0; k < indexKeys.length; k++) {
        if (indexKeys[k].indexOf(token) === 0) {
          var entries = invertedIndex[indexKeys[k]];
          entries.forEach(function (idx) { matching.add(idx); });
        }
      }
      if (resultSet === null) {
        resultSet = matching;
      } else {
        var intersection = new Set();
        resultSet.forEach(function (idx) {
          if (matching.has(idx)) intersection.add(idx);
        });
        resultSet = intersection;
      }
    }

    if (!resultSet) return [];
    var results = [];
    resultSet.forEach(function (idx) { results.push(indexEntries[idx]); });
    return results.slice(0, MAX_RESULTS);
  }

  function snippet(content, query) {
    if (!content) return "";
    if (!query || query.trim().length === 0) {
      var s = content.substring(0, SNIPPET_LEN);
      if (SNIPPET_LEN < content.length) s = s + "...";
      return s;
    }
    var lower = content.toLowerCase();
    var tokens = query.toLowerCase().split(/\s+/);
    var pos = -1;
    for (var i = 0; i < tokens.length; i++) {
      pos = lower.indexOf(tokens[i]);
      if (pos >= 0) break;
    }
    if (pos < 0) {
      var s = content.substring(0, SNIPPET_LEN);
      if (SNIPPET_LEN < content.length) s = s + "...";
      return s;
    }
    var start = Math.max(0, pos - 30);
    var s = content.substring(start, start + SNIPPET_LEN);
    if (start > 0) s = "..." + s;
    if (start + SNIPPET_LEN < content.length) s = s + "...";
    return s;
  }

  function createOverlay() {
    var overlay = document.createElement("div");
    overlay.id = "search-overlay";
    overlay.className = "hidden";

    var modal = document.createElement("div");
    modal.id = "search-modal";

    var input = document.createElement("input");
    input.id = "search-modal-input";
    input.type = "text";
    input.placeholder = "Search...";
    input.autocomplete = "off";

    var hint = document.createElement("div");
    hint.className = "search-hint";
    hint.textContent = navigator.platform.indexOf("Mac") > -1 ? "⌘K" : "Ctrl+K";

    var results = document.createElement("div");
    results.id = "search-modal-results";

    modal.appendChild(input);
    modal.appendChild(hint);
    modal.appendChild(results);
    overlay.appendChild(modal);
    document.body.appendChild(overlay);

    overlay.addEventListener("click", function (e) {
      if (e.target === overlay) hideOverlay();
    });

    return overlay;
  }

  function getOverlay() {
    return document.getElementById("search-overlay") || createOverlay();
  }

  function showOverlay() {
    var overlay = getOverlay();
    overlay.classList.remove("hidden");
    var input = document.getElementById("search-modal-input");
    if (input) {
      input.value = "";
      input.focus();
    }
    showResults(searchEntries(""), "");
  }

  function hideOverlay() {
    var overlay = document.getElementById("search-overlay");
    if (overlay) overlay.classList.add("hidden");
  }

  function showResults(results, query) {
    var container = document.getElementById("search-modal-results");
    if (!container) return;
    container.innerHTML = "";
    if (results.length === 0) {
      var msg = document.createElement("div");
      msg.className = "search-empty";
      msg.textContent = "No results found";
      container.appendChild(msg);
      return;
    }

    for (var i = 0; i < results.length; i++) {
      var entry = results[i];
      var item = document.createElement("a");
      item.href = entry.url;
      item.className = "search-result-item";
      item.setAttribute("data-index", i);

      var title = document.createElement("div");
      title.className = "search-result-title";
      title.textContent = entry.title || "";
      item.appendChild(title);

      if (entry.folder) {
        var folder = document.createElement("div");
        folder.className = "search-result-folder";
        folder.textContent = entry.folder;
        item.appendChild(folder);
      }

      var snip = document.createElement("div");
      snip.className = "search-result-snippet";
      snip.textContent = snippet(entry.content, query);
      item.appendChild(snip);

      item.addEventListener("click", function () {
        hideOverlay();
      });
      container.appendChild(item);
    }
  }

  function getHighlightedIndex(container) {
    var items = container.querySelectorAll(".search-result-item");
    for (var i = 0; i < items.length; i++) {
      if (items[i].classList.contains("highlighted")) return i;
    }
    return -1;
  }

  function setHighlighted(container, index) {
    var items = container.querySelectorAll(".search-result-item");
    for (var i = 0; i < items.length; i++) {
      items[i].classList.remove("highlighted");
    }
    if (index >= 0 && index < items.length) {
      items[index].classList.add("highlighted");
      items[index].scrollIntoView({ block: "nearest" });
    }
  }

  window.initFullSearch = function () {
    fetchIndex().then(function (data) {
      indexEntries = data;
      invertedIndex = buildInvertedIndex(data);
    });

    var overlay = getOverlay();
    var modalInput = document.getElementById("search-modal-input");
    var resultsContainer = document.getElementById("search-modal-results");

    modalInput.addEventListener("input", function (e) {
      var term = e.target.value.trim();
      var results = searchEntries(term);
      showResults(results, term);
    });

    modalInput.addEventListener("keydown", function (e) {
      var items = resultsContainer.querySelectorAll(".search-result-item");
      if (items.length === 0 && e.key !== "Escape") return;

      var current = getHighlightedIndex(resultsContainer);

      if (e.key === "ArrowDown") {
        e.preventDefault();
        setHighlighted(resultsContainer, Math.min(current + 1, items.length - 1));
      } else if (e.key === "ArrowUp") {
        e.preventDefault();
        setHighlighted(resultsContainer, Math.max(current - 1, 0));
      } else if (e.key === "Enter") {
        e.preventDefault();
        if (current >= 0 && items[current]) {
          window.location.href = items[current].href;
        }
      } else if (e.key === "Escape") {
        hideOverlay();
      }
    });
  };

  window.openSearchModal = function () {
    showOverlay();
  };
})();
