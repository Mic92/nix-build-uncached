{ pkgs ? import <nixpkgs> {}, nixVersion ? null }:
let
   nix = if nixVersion == null then pkgs.nix else pkgs.nixVersions.${nixVersion};
in
pkgs.buildGoModule {
  pname = "nix-build-uncached";
  version = "1.0.0";
  src = ./.;

  modSha256 = "1fl0wb1xj4v4whqm6ivzqjpac1iwpq7m12g37gr4fpgqp8kzi6cn";
  vendorSha256 = null;

  nativeBuildInputs = [ pkgs.makeWrapper pkgs.delve ];

  shellHook = ''
    # needed for tests
    export PATH=${pkgs.lib.makeBinPath [ nix ]}:$PATH
  '';

  # requires nix, which we do not have in the sandbox
  doCheck = false;
}
