---
title: Hidden files and folders
description: Learn how to hide certain files and folders in Kiln. This allows you to hide for example folders that are only important in the vault, like the template folder.
---
# Hidden files and folders

Right now Kiln hides by default every `dotfile`, so everything that starts with a `.`, being that a file or a folder, is hidden from the output. This feature cannot be turned off right now.

## `_hidden_` prefix

If you wish to hide certain files and folders, just use the `_hidden_` prefix. An example usage could be hiding your template folder. If you add this prefix to the template folder they will not be processed by kiln.

> [!info]
> This current solution will eventually be changed again into a more stable version that works for every situation.

