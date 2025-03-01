{
  description = "Example of overridePythonAttrs";

  inputs = {
    flake-parts.url = "github:hercules-ci/flake-parts";
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";

    # newer python packages sources
    mdformat-src.url = "github:hukkin/mdformat";
    mdformat-src.flake = false;
  };

  outputs = inputs @ {flake-parts, ...}:
    flake-parts.lib.mkFlake {inherit inputs;} {
      imports = [];
      systems = ["x86_64-linux" "aarch64-linux" "aarch64-darwin" "x86_64-darwin"];
      perSystem = {
        config,
        self',
        inputs',
        pkgs,
        system,
        ...
      }: {
        # does not work/broken
        devShells.default = let
          myCustomPython = pkgs.python3.withPackages (ps:
            with ps; [
              (mdformat.overridePythonAttrs (old: {
                src = inputs.mdformat-src;
              }))
            ]);
        in
          pkgs.mkShell
          {
            nativeBuildInputs = with myCustomPython.pkgs; [
              mdformat
            ];
          };
        # create a custom python runtime with the mdformat package overridden
        devShells.solution1 = let
          myCustomPython = pkgs.python3.withPackages (ps:
            with ps; [
              (mdformat.overridePythonAttrs (old: {
                version = "0.7.22";
                src = inputs.mdformat-src;
              }))
            ]);
        in
          pkgs.mkShell
          {
            nativeBuildInputs = [
              myCustomPython
            ];
          };
        # override mdformat package
        devShells.solution2 = let
          # old is omitted per documentation as it is not accessed
          # https://ryantm.github.io/nixpkgs/using/overrides/#sec-pkg-overrideAttrs
          custom-mdformat = pkgs.python3Packages.mdformat.overridePythonAttrs {
            version = "0.7.22";
            src = inputs.mdformat-src;
          };
        in
          pkgs.mkShell
          {
            nativeBuildInputs = [
              custom-mdformat
            ];
          };
        # override packages in python3 derivation, use mdformat from new python3 definition (py3)
        devShells.solution3 = let
          py3 = let
            packageOverrides = final: prev: {
              mdformat = prev.mdformat.overridePythonAttrs (old: rec {
                version = "0.7.22";
                src = inputs.mdformat-src;
              });
            };
          in
            pkgs.python3.override {
              inherit packageOverrides;
              self = py3;
            };
        in
          pkgs.mkShell
          {
            nativeBuildInputs = [
              py3.pkgs.mdformat
            ];
          };
      };
    };
}
