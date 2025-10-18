{
  description = "Interactive TUI for viewing and searching Helix editor's health information";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages.default = pkgs.buildGo125Module (finalAttrs: {
          pname = "helix-health";
          version = "1.0.2";

          src = ./.;

          vendorHash = "sha256-BStOOu6GxqKToSA9cEyIzJdgK2T7PPhfReuacFnh2fU=";

          # Temporarily patch go.mod to work with available Go version
          postPatch = ''
            substituteInPlace go.mod \
              --replace "go 1.25.2" "go 1.25.1"
          '';

          nativeBuildInputs = [ pkgs.makeWrapper ];

          ldflags = [
            "-s"
            "-w"
          ];

          postInstall = ''
            wrapProgram $out/bin/helix-health \
              --prefix PATH : ${pkgs.lib.makeBinPath [ pkgs.helix ]}
          '';

          meta = {
            description = "Interactive TUI for viewing and searching Helix editor's health information";
            homepage = "https://github.com/gunererd/helix-health";
            license = pkgs.lib.licenses.mit;
            maintainers = with pkgs.lib.maintainers; [ ];
            mainProgram = "helix-health";
          };
        });

        packages.helix-health = self.packages.${system}.default;

        apps.default = {
          type = "app";
          program = "${self.packages.${system}.default}/bin/helix-health";
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            helix
            gopls
            gotools
            go-tools
          ];

          shellHook = ''
            echo "🚀 helix-health development environment"
            echo "Available commands:"
            echo "  go build          - Build the project"
            echo "  go run .          - Run directly"
            echo "  make build        - Use Makefile"
            echo "  helix --health    - Test helix integration"
          '';
        };
      });
}