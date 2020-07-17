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
	// ToDO: make the GTM property name configurable
	ips, err := net.LookupIP("mirror-failover.edgedns.zone.akadns.net.")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get IPs: %v\n", err)
		os.Exit(1)
	}
	for _, ip := range ips {
		fmt.Printf("mirror-failover.edgedns.zone.akadns.net. IN A %s\n", ip.String())
	}
}

func main() {
	for {
		work()
		time.Sleep(30 * time.Second)
	}
}
