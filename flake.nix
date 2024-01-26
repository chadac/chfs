{
  description = "Versioned repositories.";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    systems.url = "github:nix-systems/default";
    flake-parts.url = "github:hercules-ci/flake-parts";
  };

  outputs = inputs@{ self, flake-parts, ... }: flake-parts.lib.mkFlake { inherit inputs; } {
    systems = import inputs.systems;
    perSystem = { pkgs, ... }: let
      versioned-repo = pkgs.callPackage ./. { };
    in {
      packages.default = versioned-repo;
      devShells.default = pkgs.mkShell {
        packages = with pkgs; [
          go
        ];
      };
    };
  };
}
