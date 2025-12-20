---
title: Doctor Command
description: The doctor command scans your vault for integrity issues, such as broken Wikilinks, ensuring your site is error-free before deployment.
---

# Doctor Command

The `doctor` command serves as a diagnostic tool for your vault. It scans your content to identify structural issues that could lead to a poor user experience, such as broken links or missing references.

It is best practice to run this command before deploying your site or after a major refactor of your notes.

## Functionality

Currently, the Doctor performs the following check:

* **Broken Link Detection:** It scans every WikiLink (`[[Note Name]]`) in your vault and verifies that the target file actually exists. If you renamed a file but forgot to update the links pointing to it, the Doctor will flag it.

## Usage

```bash
./kiln doctor
```

**Example Output:**
```
kiln: 2025/12/20 11:11:16 Diagnosing vault...
kiln: 2025/12/20 11:11:16 No broken links found
```

## Flags

| Flag       | Short version | Default value | Description                                  |
| ---------- | ------------- | ------------- | -------------------------------------------- |
| `--input`  | `-i`          | `"vault"`     | Name of the directory containing your vault. |