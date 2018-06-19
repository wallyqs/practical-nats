package main


import (
	"log"
	"runtime"

	"github.com/nats-io/go-nats"
)

func main() {
	nc, err := nats.Connect("nats://127.0.0.1:4222",
		nats.Name("practical-nats-client"),
	)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	nc.Subscribe("greeting", func(m *nats.Msg) {
		log.Printf("[Received] %s", string(m.Data))
	})
	nc.Publish("greeting", []byte("hello world"))
	runtime.Goexit()
}
