with import <nixpkgs> {};
{
  attrSet = recurseIntoAttrs {
    test = pkgs.writeText "test2" "test";
  };
}
