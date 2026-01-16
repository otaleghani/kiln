{
  description = "otaleghani/kiln";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
  };

  outputs =
    { self, nixpkgs }:
    let
      system = "x86_64-linux";
      pkgs = nixpkgs.legacyPackages.${system};
    in
    {
      packages.${system}.default = pkgs.buildGoModule {
        pname = "kiln";
        version = "0.3.5";
        src = ./.;

        # Nix needs the hash of your go.sum dependencies to ensure reproducibility.
        # STEP 1: Set this to lib.fakeHash
        # STEP 2: Run 'nix build', let it fail, and copy the actual hash it shows you here.
        vendorHash = "sha256-3a1m+YB5TuIF+FgwUONljUZnY2/ROkKrqSKeV1SxmjY";

        # If you need to generate code (like templ) before compiling Go, do it here:
        nativeBuildInputs = [
          pkgs.tailwindcss_4
          pkgs.templ
        ];
        preBuild = ''
          tailwindcss -i ./assets/simple_style_input.css -o ./assets/simple_style.cs
          tailwindcss -i ./assets/default_style_input.css -o ./assets/default_style.css
          echo "sus"
        '';
      };

      apps.${system}.default = {
        type = "app";
        # This looks inside the package defined above for a binary named "kiln"
        program = "${self.packages.${system}.default}/bin/kiln";
      };

      devShells.${system}.default = pkgs.mkShell {
        packages = with pkgs; [
          go
          gopls
          gofumpt
          golines
          goimports-reviser
          templ
          air
          overmind

          nodejs_22
          vscode-langservers-extracted
          typescript-language-server
          typescript
          prettier
          djlint # Linter and formatter for templates
          tailwindcss_4
          tailwindcss-language-server
          watchman
        ];

        shellHook = ''
          export PATH="$PWD/node_modules/.bin:$PATH"
          if [ -z "$TMUX" ]; then
            tmux set-option -g default-command "nix develop --command zsh"
            tmux new-session -s kiln -d 'nvim' \; new-window
            tmux attach-session -t kiln
          fi
        '';
      };
    };
}
