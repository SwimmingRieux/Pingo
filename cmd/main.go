package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello!")
	// Run the pinging service thread
	// Run the domain sorter service thread
	var command string
	fmt.Scan(&command)
	for command != "q" {
		switch command {
		case "update":
			// run the update subscription links service thread
		case "connect":
			// connect to vpn
		case "reconnect":
			// push the current config to the back and connect to the next one
		case "disconnect":
			// disconnect
		case "ping":
			// ping the current connection
		}
		fmt.Scan(&command)
	}
}
