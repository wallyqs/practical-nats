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

	var counter int
	var sub *nats.Subscription
	sub, err = nc.Subscribe("greeting", func(m *nats.Msg) {
		log.Printf("[Received] %s", string(m.Data))
		// Remove subscription after receiving a couple
		// of messages.
		counter++
		if counter == 2 {
			sub.Unsubscribe()
		}
	})
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	for i := 0; i < 5; i++ {
		nc.Publish("greeting", []byte("hello world!!!"))
	}
	nc.Flush()

	runtime.Goexit()
}
