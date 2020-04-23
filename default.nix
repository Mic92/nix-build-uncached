{ pkgs ? import <nixpkgs> {} }:
with pkgs;
buildGoPackage {
  pname = "ci-nix-build";
  version = "0.0.0";
  src = ./.;

  nativeBuildInputs = [ makeWrapper ];

  postInstall = ''
    wrapProgram $bin/bin/ci-nix-build \
      --prefix PATH ":" ${lib.makeBinPath [ nix ]}
  '';

  goPackagePath = "github.com/Mic92/ci-nix-build";

  shellHook = ''
    unset GOROOT GOPATH
  '';
}
