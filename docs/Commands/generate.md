---
title: generate command
description: Test
---
# `generate`

The `generate` command is used to build the static website from your vault. Is the primary command that you'll be using.

## Flags
There are some flags that you can use to customize the output:

| Flag       | Short version | Default value | Description                                                                                                       |
| ---------- | ------------- | ------------- | ----------------------------------------------------------------------------------------------------------------- |
| `--theme`  | `-t`          | `"default"`   | Color theme for the website. View [[Themes]] for more infomation.                                                 |
| `--font`   | `-f`          | `"inter"`     | Font family for the website. View [[Fonts]] for more information.                                                 |
| `--url`    | `-u`          | `""`          | The base URL of the website. Used for navigation setup, [[Sitemap.xml]] generation and [[Robots.txt]] generation. |
| `--name`   | `-n`          | `"My Notes"`  | The name of the website. Used in [[Meta Tags]] and as the website title.                                          |
| `--input`  | `-i`          | `"vault"`     | Name of the directory containing your vault.                                                                      |
| `--output` | `-o`          | `"public"`    | Name of the directory where your static website will be generated.                                                |
