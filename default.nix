{ pkgs ? import <nixpkgs> {} }:
with pkgs;
buildGoPackage {
  pname = "nix-build-uncached";
  version = "0.0.0";
  src = ./.;

  nativeBuildInputs = [ makeWrapper ];

  postInstall = ''
    wrapProgram $bin/bin/nix-build-uncached \
      --prefix PATH ":" ${lib.makeBinPath [ nix ]}
  '';

  goPackagePath = "github.com/Mic92/nix-build-uncached";

  shellHook = ''
    unset GOROOT GOPATH
  '';
}
