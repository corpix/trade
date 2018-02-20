with import <nixpkgs> {};
stdenv.mkDerivation {
  name = "nix-cage-shell";
  buildInputs = [
    go
    gocode
    go-bindata
    glide
    godef
  ];
  shellHook = ''
    export GOPATH=~/projects
  '';
}
