{ lib, stdenv, buildGoModule }:
let
  fs = lib.fileset;
in
buildGoModule {
  pname = "versioned-repo";
  version = "1.0.0";

  src = fs.toSource {
    root = ./.;
    fileset = fs.unions [
      ./src
    ];
  };
}
