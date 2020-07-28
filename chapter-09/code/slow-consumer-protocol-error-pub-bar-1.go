package main

import (
	"log"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect("nats://127.0.0.1:4222")
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	msg := make([]byte, 128)
	for i := 0; i < 1024; i++ {
		msg = append(msg, 'B')
	}

	for i := 0; i < 100000000; i++ {
		nc.Publish("bar", msg)
	}
	nc.Flush()
}
