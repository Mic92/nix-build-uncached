{ pkgs ? import <nixpkgs> {}, nixVersion ? null }:
let
   nix = if nixVersion == null then pkgs.nix else pkgs.nixVersions.${nixVersion};
in
pkgs.buildGoModule {
  pname = "nix-build-uncached";
  version = "1.1.2";
  src = ./.;

  vendorHash = null;

  nativeBuildInputs = [ pkgs.makeWrapper pkgs.delve ];

  shellHook = ''
    # needed for tests
    export PATH=${pkgs.lib.makeBinPath [ nix ]}:$PATH
  '';

  # requires nix, which we do not have in the sandbox
  doCheck = false;
}
