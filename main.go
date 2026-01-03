package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/godbus/dbus/v5"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Printf(
			"Usage: powerprofile-dbus-hook <power-saver script> <balanced script> <performance script>\n",
		)
		os.Exit(1)
	}

	powersaverScript := strings.Fields(os.Args[1])
	balancedScript := strings.Fields(os.Args[2])
	performanceScript := strings.Fields(os.Args[3])
	// log.Println(powerSaverScript)
	// log.Println(balancedScript)
	// log.Println(performanceScript)

	conn, err := dbus.SystemBus()
	if err != nil {
		log.Fatal(err)
	}

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

		fmt.Println("ActiveProfile:", profile)

		// Run scripts per profile
		switch profile {
		case "power-saver":
			err := exec.Command(powersaverScript[0], powersaverScript[1:]...).Run()
			if err != nil {
				log.Fatal(err)
			}
		case "balanced":
			err := exec.Command(balancedScript[0], balancedScript[1:]...).Run()
			if err != nil {
				log.Fatal(err)
			}
		case "performance":
			err := exec.Command(performanceScript[0], performanceScript[1:]...).Run()
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
