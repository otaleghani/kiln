---
title: Styles
description: Explore styling options in Kiln. From raw CSS to a complete guide on setting up Tailwind CSS, learn how to implement your own design engine with Kiln's unopinionated "Bring Your Own Engine" philosophy.
---
# Styling Your Site
Kiln is unopinionated when it comes to styling. We follow a **"Bring Your Own Engine"** philosophy.

Because Kiln simply renders HTML and copies static assets, you can use any styling methodology you prefer—from raw CSS to Sass, Tailwind, or Bootstrap.

## Raw CSS
If you want to keep things simple, just write standard CSS. Kiln automatically copies any static files found in your project to the output directory.

### 1. Create your CSS file
Create a folder named `assets` (or any name you prefer) and add your stylesheet.

```Plaintext
my-project/
├── assets/
│   └── style.css
├── posts/
└── layout.html
```

### 2. Link it in your Layout
In your `layout.html`, link to the file. Note that the path should be relative to the root of the built site.

```html
<head>
    <title>{{ .Page | get "Title" }}</title>
    <link rel="stylesheet" href="/assets/style.css">
</head>
```

When you run `kiln generate`, the `assets` folder is copied to the build folder, and your site is styled.

## Tailwind CSS
Tailwind CSS is a popular utility-first framework. Since Kiln is a static site generator, you can easily use the **Tailwind CLI** to scan your HTML templates and generate the CSS.

We recommend using the **Standalone CLI** for Tailwind so you don't have to deal with `node_modules`, keeping your project clean just like Kiln.

### 1. Get the Tailwind CLI
Download the executable for your platform from the [Tailwind Releases page](https://github.com/tailwindlabs/tailwindcss/releases/) and place it in your project root (rename it to `tailwindcss` for ease of use).

```bash
# Example for Mac/Linux
chmod +x tailwindcss
```

If you have Node.js available you could use `npm` instead, like so:
```bash
npm install tailwindcss @tailwindcss/cli
```

### 2. Create Input CSS
Create a source CSS file (e.g., `src/input.css`) and add the Tailwind directive:
```css
@import "tailwindcss";
```

### 3. Run the Watcher
Now, run Tailwind in "watch" mode. We tell it to take `src/input.css` and output it to `assets/style.css`.

```bash
./tailwindcss -i ./src/input.css -o ./assets/style.css --watch
```

### 4. Link and Serve
In your `layout.html`, link to the **output** file:

```html
<link rel="stylesheet" href="/assets/style.css">
```

### 5. Start using Tailwind in your HTML
You are all set! Now you just need to start using Tailwind.

```html
<!doctype html>
<html>
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<link href="./assets/style.css" rel="stylesheet">
	</head>
	<body>
		<h1 class="text-3xl font-bold underline">Hello world!</h1>
	</body>
</html>
```

### 6. View the result
After some changes, you can view the website using the [[serve]] command, like so:

```bash
kiln serve
```

**Workflow:**

1. You edit your `layout.html` and add a class like `text-red-500`.
2. The `tailwindcss` process sees the change and regenerates `assets/style.css`.
3. You refresh your browser (served by Kiln) and see the changes.