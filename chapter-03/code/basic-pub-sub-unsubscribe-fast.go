package main

import (
	"log"
	"runtime"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect("nats://127.0.0.1:4222")
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	sub, err := nc.Subscribe("greeting", func(m *nats.Msg) {
		log.Printf("[Received] %s", string(m.Data))
	})
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	for i := 0; i < 5; i++ {
		nc.Publish("greeting", []byte("hello world!!!"))
	}
	nc.Flush()

	// Remove subscription
	sub.Unsubscribe()

	for i := 0; i < 5; i++ {
		nc.Publish("greeting", []byte("hello world!!!"))
	}

	runtime.Goexit()
}
