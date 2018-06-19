package main


import (
	"log"

	"github.com/nats-io/go-nats"
)

func main() {
	nc, err := nats.Connect("nats://127.0.0.1:4222")
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	nc.Subscribe("foo", func(_ *nats.Msg) {
	        // Heavy processing
		for i := 0; i < 10000000000; i++ {
		}
	})
	nc.Subscribe("bar", func(_ *nats.Msg) {
	        // Not heavy processing
	})
	select {}
}
