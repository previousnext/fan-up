package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os/exec"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	cliInterface = kingpin.Flag("interface", "The interface we will attach fan to").Default("eth0").String()
	cliOverlay   = kingpin.Flag("overlay", "The overlay CIDR that will sit on top of the interface").Default("241.0.0.0").String()
)

// * Does our fan network interface exist.
// * If not create it
// * Check if it has come up and report
func main() {
	kingpin.Parse()

	var (
		bridge  = fmt.Sprintf("fan-%s", *cliInterface)
		limiter = time.Tick(time.Second * 15)
	)

	for {
		<-limiter

		// Check that the interface exists that we are going to start a fan network against.
		e, ip, err := find(*cliInterface)
		if err != nil {
			log.Println(err)
			continue
		}

		// Ensure the interface we are building against.
		if !e {
			log.Printf("Cannot find the interface: %s\n", *cliInterface)
			continue
		}

		// Run the command to bring up the interface.
		err = shellOut(fmt.Sprintf("fanctl up %s/8 %s/16 dhcp bridge %s", *cliOverlay, ip, bridge))
		if err != nil {
			log.Println(err)
			continue
		}

		// Query for the interface again.
		e, ip, err = find(bridge)
		if !e {
			log.Println("The interface has not come up")
			continue
		}
		log.Println("The interface has now come up")
	}
}

// Helper function to find our interface.
func find(n string) (bool, string, error) {
	ints, err := net.Interfaces()
	if err != nil {
		return false, "", err
	}

	for _, i := range ints {
		if i.Name == n {
			ip, err := ip(i)
			if err != nil {
				continue
			}
			return true, ip, nil
		}
	}

	return false, "", nil
}

func ip(i net.Interface) (string, error) {
	addrs, err := i.Addrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", errors.New("Cannot find IP address for interface")
}

// Helper function execute commands on the commandline.
func shellOut(cmd string) error {
	out, err := exec.Command("sh", "-c", cmd).CombinedOutput()
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to execute %v: %v, err: %v", cmd, string(out), err))
	}
	return nil
}
