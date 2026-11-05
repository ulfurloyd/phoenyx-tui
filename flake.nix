{
  description = "phoenyx-tui";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
  };
  outputs =
    { self, nixpkgs }:
    let
      systems = [
        "x86_64-linux"
        "aarch64-linux"
      ];
      forAllSystems =
        f:
        nixpkgs.lib.genAttrs systems (
          system: f (nixpkgs.legacyPackages.${system}.extend self.overlays.default)
        );
    in
    {
      packages = forAllSystems (pkgs: {
        default = pkgs.phoenyx-tui;
        phoenyx-tui = pkgs.phoenyx-tui;
      });

      overlays.default = final: _prev: {
        phoenyx-tui = final.callPackage ./packages/phoenyx-tui.nix {
          src = self;
        };
      };
      apps = forAllSystems (pkgs: {
        default = {
          type = "app";
          program = "${pkgs.phoenyx-tui}/bin/phoenyx";
          meta.description = "A terminal interface for managing systems and projects";
        };
      });
      checks = forAllSystems (pkgs: {
        default = pkgs.phoenyx-tui;
      });
      formatter = forAllSystems (pkgs: pkgs.nixfmt-tree);
    };
}
