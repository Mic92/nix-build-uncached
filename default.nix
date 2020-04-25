{ pkgs ? import <nixpkgs> {} }:
with pkgs;
buildGoModule {
  pname = "nix-build-uncached";
  version = "0.0.0";
  src = ./.;

  modSha256 = "1fl0wb1xj4v4whqm6ivzqjpac1iwpq7m12g37gr4fpgqp8kzi6cn";

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
