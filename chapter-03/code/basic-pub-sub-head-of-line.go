package main

import (
	"log"
	"runtime"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect("nats://127.0.0.1:4222",
		nats.UserInfo("foo", "secret"),
	)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	nc.Subscribe("greeting", func(m *nats.Msg) {
		log.Printf("[Received] %s", string(m.Data))
	})
	nc.Subscribe(">", func(m *nats.Msg) {
		log.Printf("[Wildcard] %s", string(m.Data))
		time.Sleep(1 * time.Second)
	})

	for i := 0; i < 10; i++ {
		nc.Publish("greeting", []byte("hello world!"))
	}
	runtime.Goexit()
}
