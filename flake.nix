{
  description = "Configuration for Zynra";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-parts.url = "github:hercules-ci/flake-parts";
  };

  outputs = inputs @ {flake-parts, ...}:
    flake-parts.lib.mkFlake {inherit inputs;} {
      systems = ["x86_64-linux" "aarch64-linux" "aarch64-darwin" "x86_64-darwin"];

      perSystem = {
        config,
        self',
        inputs',
        pkgs,
        system,
        ...
      }: let
        name = "Zynra";
        vendorHash = "sha256-nZu9e6jKHZ2ndhWCjIbrGjy7qRBOsdrteM10btV6+qo="; # update whenever go.mod changes
      in {
        devShells = {
          default = pkgs.mkShell {
            inputsFrom = [self'.packages.default];
            buildInputs = with pkgs; [
              go_1_24
              just
            ];
          };
        };

        packages = {
          default = pkgs.buildGoModule {
            inherit name vendorHash;
            src = ./.;
            subPackages = ["cmd"];
            postBuild = ''
              mv $GOPATH/bin/cmd $GOPATH/bin/zyndra
            '';
          };
        };
      };
    };
}
