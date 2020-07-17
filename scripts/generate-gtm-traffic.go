// This script will evolve to generate traffic to a configurable GTM Property.
// Additional capabilities might include more throughput and the ability to
// generate traffic from different geographic locations.

package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func work() {
	ips, err := net.LookupIP("www.jdgri.me.")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get IPs: %v\n", err)
		os.Exit(1)
	}
	for _, ip := range ips {
		fmt.Printf("www.jdgri.me. IN A %s\n", ip.String())
	}
}

func main() {
	for {
		work()
		time.Sleep(30 * time.Second)
	}
}
