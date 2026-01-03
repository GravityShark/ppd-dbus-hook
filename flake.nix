{
  description = "Run commands when changing from \"balanced\", \"power-saver\" or \"performance\" from power-profiles-daemon using dbus.";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs =
    { self, nixpkgs }:
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs { inherit system; };
    in
    {
      packages.${system} =
        let
          ppd-dbus-hook = pkgs.callPackage ./package.nix { };
        in
        {
          ppd-dbus-hook = ppd-dbus-hook;
          default = ppd-dbus-hook;
        };
    };

}
