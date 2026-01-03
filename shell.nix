{
  pkgs ? import <nixpkgs> { },
}:

with pkgs;

mkShell {
  nativeBuildInputs = [ ];

  buildInputs = [
    go
  ];
}
