package main

import (
	"flag"
	"log"
	"os"
	"pingo/configs"
)

func main() {
	configuration, err := configs.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	device := flag.String("device", "", "")
	flag.Parse()

	if *device == "" {
		log.Fatal(err)
	}

	port := flag.String("port", "", "")
	flag.Parse()

	if *port == "" {
		log.Fatal(err)
	}

	err = os.Setenv("PINGO_DEFAULT_RECORDING_DEVICE", *device)
	if err != nil {
		log.Fatal(err)
	}
	err = os.Setenv("PINGO_DEFAULT_PORT", *port)
	if err != nil {
		log.Fatal(err)
	}

}
