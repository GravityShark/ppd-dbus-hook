Simple go program that subscribes to dbus to monitor for specifically only power-profiles-daemon
changes and then runs a command from that

## why?

cause when i try to use dbus-monitor or busctl, it goes back to eavesdropping
which according to shitgippity is wrong and recommended to just create a new
program, which it also gracefully provided (it just gave the dbus part i added
the argument part, if that even accounts for anything)

## installation

lowkey idk, its just a go program figure that out yourself lowkey

but i do have a flake

```nix

{
    inputs = {
        ppd-dbus-hook = {
          url = "github:GravityShark/ppd-dbus-hook";
          inputs.nixpkgs.follows = "nixpkgs";
        };
    }
}
```

now add it to your config in whichever way you like. personally i like inherting
inputs to nixos and home-manager and doing this

```nix
${inputs.ppd-dbus-hook.packages.${pkgs.stdenv.hostPlatform.system}.default}
```

## examples

you could like make it send notifications (i use [notify-desktop](https://github.com/nowrep/notify-desktop))

```shell
ppd-dbus-hook \
    "notify-desktop 'Power-saver mode enabled'" \
    "notify-desktop 'Balanced mode enabled'" \
    "notify-desktop 'Performance enabled'"
```

i personally want it to enable or disable [noctalia performance](https://github.com/noctalia-dev/noctalia-shell)

```shell
ppd-dbus-hook \
    "noctalia-shell ipc call powerProfile enableNoctaliaPerformance" \
    "noctalia-shell ipc call powerProfile disableNoctaliaPerformance" \
    "noctalia-shell ipc call powerProfile enableNoctaliaPerformance"
```

i also use it for my msi-ec shift mode

```shell
sudo ppd-dbus-hook \
    "sh -c 'echo eco > /sys/devices/platform/msi-ec/shift_mode'" \
    "sh -c 'echo comfort > /sys/devices/platform/msi-ec/shift_mode'" \
    "sh -c 'echo turbo > /sys/devices/platform/msi-ec/shift_mode'"
```

and i just run these scripts in a systemd service
