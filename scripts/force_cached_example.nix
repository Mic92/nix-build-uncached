let
  pkgs = import <nixpkgs> {};
in (pkgs.callPackage ./force_cached.nix {}) {
  hello = pkgs.writeScriptBin "hello" ''
    #!/bin/sh
    exec ${pkgs.hello}/bin/hello
  '';
}
