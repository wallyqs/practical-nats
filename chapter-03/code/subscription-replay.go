package main


import (
	"log"
	"time"

	"github.com/nats-io/go-nats"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL, nats.MaxReconnects(-1))
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	nc.Subscribe("hello", func(m *nats.Msg) {
		log.Printf("[Received] %s", string(m.Data))
	})

	nc.Subscribe("*", func(m *nats.Msg) {
		log.Printf("[Wildcard] %s", string(m.Data))
	})

	// If disconnected for too long and buffer is full
	// then the client will receive a synchronous error.
	for range time.NewTicker(1 * time.Second).C {
		err := nc.Publish("hello", []byte("hello world"))
		if err != nil {
			log.Printf("Error: %s", err)
		}
	}
}
