{ pkgs ? import <nixpkgs> {} }:
with pkgs;
buildGoModule {
  pname = "nix-build-uncached";
  version = "0.1.0";
  src = ./.;

  modSha256 = "1fl0wb1xj4v4whqm6ivzqjpac1iwpq7m12g37gr4fpgqp8kzi6cn";
  vendorSha256 = null;

  nativeBuildInputs = [ makeWrapper delve ];

  postInstall = ''
    wrapProgram $out/bin/nix-build-uncached \
      --prefix PATH ":" ${lib.makeBinPath [ nix ]}
  '';

  # requires nix, which we do not have in the sandbox
  doCheck = false;

  goPackagePath = "github.com/Mic92/nix-build-uncached";
}
