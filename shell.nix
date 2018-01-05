with import <nixpkgs> {};
stdenv.mkDerivation {
  name = "go-shell";
  buildInputs = [
    ncurses
    go
    gocode
    go-bindata
    glide
    godef
    bison
  ];
}
