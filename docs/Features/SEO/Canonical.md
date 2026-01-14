---
title: Canonical tag
description: Learn how to optimize your site's SEO by adding canonical tags to your website pages.
---
# Canonical tag
Without a canonical tag, if your site is accessible via both `http` and `https`, or `www` and `non-www`, or with trailing slashes, Google might view them as duplicates.

To add a correct canonical tag to you kiln website you just need to specify in the [[generate]] command the `--url` flag. Remember: the URL has to be the final URL of your index page, like this:
```bash
kiln generate --url "https://www.example.com" # No slashes at the end!
```