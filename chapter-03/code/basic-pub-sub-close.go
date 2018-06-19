package main


import (
	"log"
	"runtime"

	"github.com/nats-io/go-nats"
)

func main() {
	nc, err := nats.Connect("nats://127.0.0.1:4222")
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	nc.Subscribe("greeting", func(m *nats.Msg) {
		log.Printf("[Received] %s", string(m.Data))
	})

	for i := 0; i < 10; i++ {
		nc.Publish("greeting", []byte("hello world!!!"))
	}
	nc.Flush()

	// Terminate connection to NATS
	nc.Close()

	// Error: nats: connection closed
	err = nc.Publish("greeting", []byte("hello world!!!"))
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	runtime.Goexit()
}
