package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/godbus/dbus/v5"
	"github.com/google/shlex"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Printf(
			"Usage: powerprofile-dbus-hook <power-saver script> <balanced script> <performance script>\n",
		)
		os.Exit(1)
	}

	powersaverScript, err := shlex.Split(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	balancedScript, err := shlex.Split(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	performanceScript, err := shlex.Split(os.Args[3])
	if err != nil {
		log.Fatal(err)
	}

	conn, err := dbus.SystemBus()
	if err != nil {
		log.Fatal(err)
	}

	// Grab the current power profile and apply it immediately
	current_property, err := conn.Object("net.hadess.PowerProfiles", "/net/hadess/PowerProfiles").
		GetProperty("net.hadess.PowerProfiles.ActiveProfile")
	if err != nil {
		log.Fatal(err)
	}

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
		log.Fatal(call.Err)
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
		fmt.Println("Running:", powersaverScript)
		err := exec.Command(powersaverScript[0], powersaverScript[1:]...).Run()
		if err != nil {
			log.Fatal(err)
		}
	case "balanced":
		fmt.Println("Running:", balancedScript)
		err := exec.Command(balancedScript[0], balancedScript[1:]...).Run()
		if err != nil {
			log.Fatal(err)
		}
	case "performance":
		fmt.Println("Running:", performanceScript)
		err := exec.Command(performanceScript[0], performanceScript[1:]...).Run()
		if err != nil {
			log.Fatal(err)
		}
	}
}
