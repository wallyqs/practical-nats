package main

import (
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	opts := nats.DefaultOptions

	// Arbitrarily small reconnecting buffer
	opts.ReconnectBufSize = 256
	nc, err := opts.Connect()
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	for range time.NewTicker(500 * time.Millisecond).C {
		// If disconnected for too long and buffer is full
		// then the client will receive a synchronous error.
		err := nc.Publish("numbers", []byte("4 8 15 16 23 42"))
		if err != nil {
			// nats: outbound buffer limit exceeded
			log.Fatalf("Error: %s", err)
		}
	}
}
