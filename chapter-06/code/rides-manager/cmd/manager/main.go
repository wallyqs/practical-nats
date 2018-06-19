package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/nats-io/go-nats"
	"github.com/wallyqs/practical-nats/chapter-06/code/kit"
	"github.com/wallyqs/practical-nats/chapter-06/code/rides-manager"
)

func main() {
	var (
		showHelp    bool
		showVersion bool
		natsServers string
	)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: rides-manager [options...]\n\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
	}

	// Setup default flags
	flag.BoolVar(&showHelp, "help", false, "Show help")
	flag.BoolVar(&showVersion, "version", false, "Show version")
	flag.StringVar(&natsServers, "nats", nats.DefaultURL, "List of NATS Servers to connect")
	flag.Parse()

	switch {
	case showHelp:
		flag.Usage()
		os.Exit(0)
	case showVersion:
		fmt.Fprintf(os.Stderr, "NATS Rider Rides Manager Server v%s\n", ridesmanager.Version)
		os.Exit(0)
	}
	log.Printf("Starting NATS Rider Rides Manager version %s", ridesmanager.Version)

	comp := kit.NewComponent("rides-manager")
	err := comp.SetupConnectionToNATS(natsServers)
	if err != nil {
		log.Fatal(err)
	}

	s := ridesmanager.Server{
		Component: comp,
	}
	err = s.SetupSubscriptions()
	if err != nil {
		log.Fatal(err)
	}

	runtime.Goexit()
}
