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
      forAllSystems = nixpkgs.lib.genAttrs systems;
      pkgsFor = forAllSystems (
        system:
        import nixpkgs {
          inherit system;
        }
      );
    in
    {
      packages = forAllSystems (
        system:
        let
          pkgs = pkgsFor.${system};
        in
        {
          default =
            pkgs.buildGoModule.override
              {
                go = pkgs.go_1_27;
              }
              {
                pname = "phoenyx-tui";
                version = "0.1.0";
                src = self;
                subPackages = [ "cmd/phoenyx" ];
                vendorHash = "sha256-bxnj5fGcXAp1veo600mmq5f7SM2D4/WNekfWFYu8jOk=";
                doCheck = true;
                meta = {
                  description = "A terminal interface for managing systems and projects";
                  mainProgram = "phoenyx";
                  platforms = systems;
                };
              };
        }
      );
      apps = forAllSystems (system: {
        default = {
          type = "app";
          program = "${self.packages.${system}.default}/bin/phoenyx";
          meta = {
            description = "A terminal interface for managing systems and projects";
            mainProgram = "phoenyx";
            platforms = systems;
          };
        };
      });
      devShells = forAllSystems (
        system:
        let
          pkgs = pkgsFor.${system};
        in
        {
          default = pkgs.mkShell {
            packages = with pkgs; [
              go_1_27
              gopls
              gotools
              golangci-lint
            ];
          };
        }
      );
      checks = forAllSystems (system: {
        default = self.packages.${system}.default;
      });
    };
}
