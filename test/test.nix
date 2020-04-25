{ pkgs ? import <nixpkgs> {} }:
{
  test = pkgs.writeText "test" "test";
  inherit (pkgs) hello;
}
