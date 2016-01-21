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

		i, err := exists(bridge)
		if err != nil {
			log.Println("Cannot find interface, will try again soon.")
			continue
		}
		log.Println("The interface already exists, no action needs to be taken.")

		// Get the IP address of the interface.
		address, err := ip(i)
		if err != nil {
			log.Println(err)
			continue
		}

		// Run the command to bring up the interface.
		err = shellOut(fmt.Sprintf("fanctl up %s/8 %s/16 dhcp bridge %s", *cliOverlay, address, bridge))
		if err != nil {
			log.Println(err)
			continue
		}

		// Query for the interface again.
		i, err = exists(bridge)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("The interface has now come up")
	}
}

// Helper function to find our interface.
func exists(n string) (net.Interface, error) {
	ints, err := net.Interfaces()
	if err != nil {
		return net.Interface{}, err
	}

	for _, i := range ints {
		if i.Name == n {
			return i, nil
		}
	}

	return net.Interface{}, errors.New("Cannot find interface.")
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
