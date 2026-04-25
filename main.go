package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/godbus/dbus/v5"
	"github.com/google/shlex"
)

var out bytes.Buffer

func main() {
	if len(os.Args) != 4 {
		fmt.Printf(
			"Usage: powerprofile-dbus-hook <power-saver script> <balanced script> <performance script>\n",
		)
		os.Exit(1)
	}

	powersaverScript, err := shlex.Split(os.Args[1])
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	balancedScript, err := shlex.Split(os.Args[2])
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	performanceScript, err := shlex.Split(os.Args[3])
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	conn, err := dbus.SystemBus()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	// Grab the current power profile and apply it immediately
	current_property, err := conn.Object("net.hadess.PowerProfiles", "/net/hadess/PowerProfiles").
		GetProperty("net.hadess.PowerProfiles.ActiveProfile")
	if err != nil {
		fmt.Print("Could not get propery net.hadess.PowerProfiles.ActiveProfile", err)
		os.Exit(1)
	}

	// Run it once
	scriptsPerProfile(
		current_property.Value().(string),
		powersaverScript,
		balancedScript,
		performanceScript,
	)

	// Watch for changes
	rule := "type='signal'," +
		"interface='org.freedesktop.DBus.Properties'," +
		"member='PropertiesChanged'," +
		"path='/net/hadess/PowerProfiles'," +
		"arg0='net.hadess.PowerProfiles'"
	call := conn.BusObject().Call(
		"org.freedesktop.DBus.AddMatch",
		0,
		rule,
	)
	if call.Err != nil {
		fmt.Print("org.freedesktop.DBus.AddMatch call failed", call.Err)
		os.Exit(1)
	}

	c := make(chan *dbus.Signal, 10)
	conn.Signal(c)

	fmt.Println("Watching ActiveProfile changes…")

	for sig := range c {
		if len(sig.Body) < 2 {
			continue
		}

		changed, ok := sig.Body[1].(map[string]dbus.Variant)
		if !ok {
			continue
		}

		v, ok := changed["ActiveProfile"]
		if !ok {
			continue
		}

		profile, ok := v.Value().(string)
		if !ok {
			continue
		}

		// Run scripts per profile
		scriptsPerProfile(profile, powersaverScript, balancedScript, performanceScript)
	}
}

func scriptsPerProfile(
	profile string,
	powersaverScript []string,
	balancedScript []string,
	performanceScript []string,
) {
	fmt.Println("ActiveProfile:", profile)
	switch profile {
	case "power-saver":
		fmt.Printf("Running: [%s]\n", powersaverScript)
		powersaverCommand := exec.Command(powersaverScript[0], powersaverScript[1:]...)
		powersaverCommand.Stdout = &out
		err := powersaverCommand.Run()
		if err != nil {
			fmt.Printf("PowerSaverScript Run Error: [%s] %s", err, out.String())
			os.Exit(1)
		}
	case "balanced":
		fmt.Printf("Running: [%s]\n", balancedScript)
		balancedCommand := exec.Command(balancedScript[0], balancedScript[1:]...)
		balancedCommand.Stdout = &out
		err := balancedCommand.Start()
		if err != nil {
			fmt.Printf("BalancedScript Run Error: [%s] %s", err, out.String())
			os.Exit(1)
		}
	case "performance":
		fmt.Printf("Running: [%s]\n", performanceScript)
		performanceCommand := exec.Command(performanceScript[0], performanceScript[1:]...)
		performanceCommand.Stdout = &out
		err := performanceCommand.Start()
		if err != nil {
			fmt.Printf("PerformanceScript Run Error: [%s] %s", err, out.String())
			os.Exit(1)
		}
	}
}
