package main

import (
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect("nats://127.0.0.1:4222", nats.MaxReconnects(-1))
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	// If disconnected for too long and buffer is full
	// then the client will receive a synchronous error.
	for range time.NewTicker(1 * time.Second).C {
		err := nc.Publish("hello", []byte("world"))
		if err != nil {
			log.Fatalf("Error: %s", err)
		}
	}
}
