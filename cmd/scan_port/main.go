package main

import (
	"fmt"

	"github.com/zerodot618/go-huang/cmd/scan_port/port"
)

func main() {
	fmt.Println("Port Scanning")
	results := port.InitialScan("localhost")
	fmt.Println(results)

	widescanresults := port.WideScan("localhost")
	fmt.Println(widescanresults)
}
