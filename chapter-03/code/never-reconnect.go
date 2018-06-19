package main


import (
	"log"
	"time"

	"github.com/nats-io/go-nats"
)

func main() {
	nc, err := nats.Connect("nats://127.0.0.1:4222", nats.NoReconnect())
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	// Since we disallow reconnection, this should report
	// an error quickly after stopping the NATS Server.
	for range time.NewTicker(1 * time.Second).C {
		err := nc.Publish("hello", []byte("world"))
		if err != nil {
			log.Fatalf("Error: %s", err)
		}
	}
}
