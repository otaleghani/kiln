---
title: Stats Command
description: The stats command analyzes your vault to provide a quantitative summary of your knowledge base, including file counts and word metrics.
---

# Stats Command

The `stats` command provides a quick, quantitative overview of your knowledge base.

It scans your input directory and calculates metrics regarding the size and complexity of your vault. This is useful for tracking your writing progress or getting a sense of the scale of your documentation.

## Usage

```bash
./kiln stats
```

**Example Output:**
```
METRIC        VALUE
------        -----
Total Notes   38
Total Words   6533
Longest Note  Client Side Navigation.md (412 words)
```

## Flags

| Flag       | Short version | Default value | Description                                  |
| ---------- | ------------- | ------------- | -------------------------------------------- |
| `--input`  | `-i`          | `"vault"`     | The path to the directory containing your Markdown notes to be analyzed. |