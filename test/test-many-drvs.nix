with import <nixpkgs> {};
lib.listToAttrs (map (i: let
  id = "test${toString i}";
in {
  name = id;
  value = pkgs.writeText id id;
}) (lib.range 1 999))
