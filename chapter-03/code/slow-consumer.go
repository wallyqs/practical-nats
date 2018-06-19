package main


import (
	"log"

	"github.com/nats-io/go-nats"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	msg := []byte("Hello World!")
	nc.Subscribe("greetings", func(_ *nats.Msg) {})
	for i := 0; i < 100000000; i++ {
		nc.Publish("greetings", msg)
	}
}
