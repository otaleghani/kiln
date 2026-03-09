# Configuration File

Instead of passing flags on every command, you can place a `kiln.yaml` file in your project root. Kiln automatically discovers it and applies the values as defaults. Any CLI flag you pass explicitly will override the corresponding config value.

## Creating a config file

Run `kiln init` to scaffold a commented-out `kiln.yaml` in the current directory: