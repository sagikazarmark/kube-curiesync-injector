{
  description = "Kubernetes mutating webhook for injecting Curiefense config sync component into Pods";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};

        buildDeps = with pkgs; [ git go_1_17 gnumake ];
        devDeps = with pkgs; buildDeps ++ [
          golangci-lint
          gotestsum
        ];
      in
      { devShell = pkgs.mkShell { buildInputs = devDeps; }; });
}
