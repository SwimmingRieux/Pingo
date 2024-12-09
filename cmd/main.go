package main

import (
	"fmt"
	"pingo/internal/config_collector"
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
			sub := config_collector.SubscriptionLoader{}
			res, err := sub.GetSub("")
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(res)
			}
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
