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
	nc.QueueSubscribe("greeting", "workers", func(m *nats.Msg) {
		log.Printf("[Received] %s", string(m.Data))
	})
	nc.Publish("greeting", []byte("hello world!!!"))
	runtime.Goexit()
}
