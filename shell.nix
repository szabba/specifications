{ pkgs ? import (builtins.fetchTarball {

  # Commit selected from the end of the nixos-unstable nixpkg branch on 2022-04-10.
  name = "nixos-unstable-42948b300670223ca8286aaf916bc381f66a5313";
  url =
    "https://github.com/NixOS/nixpkgs/archive/42948b300670223ca8286aaf916bc381f66a5313.tar.gz";

  # Use
  #
  #     nix-prefetch-url --unpack <url>
  #
  # to regenerate
  sha256 = "09nx6mmld7iag3ffcfz4ybk0w9j3sg2akjqmf41g1rajxgadc516";

}) { } }:

pkgs.mkShell {

  # Packages we need to run during development / build.
  nativeBuildInputs = [
    # Autoformatter for Nix
    pkgs.buildPackages.nixfmt
    # Go SDK
    pkgs.buildPackages.go_1_18
  ];

  shellHook = ''
    export CGO_ENABLED=0
  '';

}
