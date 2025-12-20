---
title: Search
description: Kiln features a real-time client-side search that instantly filters your sidebar, allowing you to quickly locate notes and folders by name.
---
# Search

Kiln includes a built-in **Quick Find** feature designed to help you navigate large vaults with ease.

Located at the top of the sidebar, the search input allows you to filter your file tree in real-time. As you type, the Explorer immediately updates to show only the items that match your query.

## Functionality

The search engine runs entirely on the **Client Side**, meaning it does not need to communicate with a server to return results. This ensures zero latency—results appear instantly as you type.

### Smart Expansion
To ensure you can always find what you are looking for, the search logic automatically manages your folder structure:

* **Filtering:** Files and folders that do not match your query are hidden.
* **Auto-Expand:** If a matching file is nested deep within closed folders, Kiln will automatically expand the necessary parent directories to reveal the match.

*Note: Currently, this feature searches file and folder names (titles). Full-text content search is planned for future updates.*