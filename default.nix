{ pkgs ? import <nixpkgs> {} }:
with pkgs;
buildGoModule {
  pname = "nix-build-uncached";
  version = "1.0.0";
  src = ./.;

  modSha256 = "1fl0wb1xj4v4whqm6ivzqjpac1iwpq7m12g37gr4fpgqp8kzi6cn";
  vendorSha256 = null;

  nativeBuildInputs = [ makeWrapper delve ];

  shellHook = ''
    # needed for tests
    export PATH=$PATH:${lib.makeBinPath [ nix ]}
  '';

  # requires nix, which we do not have in the sandbox
  doCheck = false;

  goPackagePath = "github.com/Mic92/nix-build-uncached";
}
