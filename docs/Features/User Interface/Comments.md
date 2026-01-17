---
title: Comments
description: "Want to add comments for your website? Not a problem: you could use Kiln supports Giscus and it's perfect to add UGC into your vault."
---
# Comments

To add comments into your website we recommend to use [giscus](https://giscus.app/), a comments system powered by GitHub Discussions. This is overall the easiest and fastest way to add comments to your website, being that you are most likely already using GitHub to host your notes. 

## `giscus` guide

To enable comments you'll need to make sure that
- Your repository is public. You can change the visibility on the **Settings** page of your repo. Find out more about this [here](https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/managing-repository-settings/setting-repository-visibility#making-a-repository-public). 
- You have the `Discussions` feature turned **on**. Here's the [official GitHub tutorial on how to enable Discussions](docs.github.com/en/github/administering-a-repository/managing-repository-settings/enabling-or-disabling-github-discussions-for-a-repository).
- You have the `giscus` app installed for your repository. To do this, just follow this [link](https://github.com/apps/giscus) and install the app for your account giving it visibility of your repository.

After that hope into [the official giscus website](https://giscus.app/) and follow the steps under the `Configuration` heading to configure the script that you'll add into your notes. Just follow the instructions and setup the comments how you want them. Here's an example of the final script, in this case is the script used for this page:
```html
<script src="https://giscus.app/client.js"
    data-repo="otaleghani/kiln"
    data-repo-id="R_kgDOQpz0LQ"
    data-category="Comments"
    data-category-id="DIC_kwDOQpz0Lc4C1ENW"
    data-mapping="url"
	data-strict="1"
    data-reactions-enabled="1"
    data-emit-metadata="0"
    data-input-position="top"
    data-theme="https://giscus.app/themes/custom_example.css"
    data-lang="en"
    data-loading="lazy"
    crossorigin="anonymous"
    async>
</script>
```

## Themes

Kiln automatically tries to sync the theme of `giscus` with the current active theme. You can find a list of all the available [[Themes]] here. This happens on page load and whenever the user changes theme.

## Example

<script src="https://giscus.app/client.js"
        data-repo="otaleghani/kiln"
        data-repo-id="R_kgDOQpz0LQ"
        data-category="Comments"
        data-category-id="DIC_kwDOQpz0Lc4C1ENW"
        data-mapping="url"
        data-strict="1"
        data-reactions-enabled="1"
        data-emit-metadata="0"
        data-input-position="top"
        data-theme="https://giscus.app/themes/custom_example.css"
        data-lang="en"
        data-loading="lazy"
        crossorigin="anonymous"
        async>
</script>