---
title: Changelog
description: A list of all the addition to kiln.
---
# Changelog

## v0.3.13
- FIX: mobile sidebar was by default open on mobile

## v0.3.12
- FIX: sitemap not generating

## v0.3.11
- FIX: folder pages overrides

## v0.3.10
- FIX: giscus theme was inverted (again...)

## v0.3.9
- FIX: giscus theme was inverted

## v0.3.8
- FIX: giscus change theme not being loaded when giscus loads

## v0.3.7
- ADD: Comments support 
- FIX: Links to javascript redirects

## v0.3.6
- FIX: Navbar had links to redirects for folders

## v0.3.5
- ADD: Handle for `_redirects` file (for [[Cloudflare Pages]])
- FIX: Canonical tags in `flat-urls` mode
- FIX: `![new]` callout had no style
- FIX: Canvas notes not being displayed

## v0.3.4
- ADD: Canonical tags

## v0.3.3
- FIX: Navbar was displaying notes with the same name
- ADD: Documentation about overlapping files
- FIX: Various small UI fixes

## v0.3.2
- FIX: Command `version` was not updating correctly

## v0.3.1
- FIX: `simple_layout.html` and `default_layout.html` had issues with scrolling
- FIX: `canvas.js` didn't have mobile events

## v0.3.0
- ADD: [[Bases]] page generation
- ADD: [[Folders]] page generation
- ADD: [[Tags]] page generation
- ADD: Simple layout
- ADD: New [[Themes]] (`tokionight`, `rose-pine`, `gruvbox`, `everforest`, `cyberdream`)
- ADD: New [[Fonts]] (`lora`, | `libre-baskerville`, `noto-serif`, `ibm-plex-sans`, `google-sans`, `roboto`
- CHORE: Default layout refactor
	- REMOVED: Lucide as a dependency
	- ADD: TailwindCSS instead of raw CSS


## v0.2.6
- FIX: Static data, like images, where not copied over

## v0.2.5
- FIX: Flat-urls where not working properly
- FIX: `deploy.sh` and flat-urls flag

## v0.2.4
- FIX: Wikilinks labels where not displaying

## v0.2.3
- FIX: Anchor links where not slugified
- FIX: CNAME and favicon.ico file where not copied over

## v0.2.2
- FIX: Embed empty files

## v0.2.1
- ADD: Better logs
- CHORE: Refactored most of the markdown parsing
- CHORE: Refactored default generation

## v0.2.0
- ADD: [[What is Custom Mode|Custom mode]]

## v0.1.5
- FIX: Graph not showing links
- FIX: Default output directory not being set

## v0.1.4
- WIP: Custom mode
- ADD: Text embedding

## v0.1.3
- FIX: Missing site name in configuration
- FIX: [[Doctor]] command not working properly
- FIX: Local graph links not working
- FIX: URL generation
- FIX: UI mobile experience

## v0.1.2
- ADD: CNAME file handling
- FIX: Layout footer
- FIX: UI 

## v0.1.1
- FIX: Deployment crash

## v0.1.0
- Initial release