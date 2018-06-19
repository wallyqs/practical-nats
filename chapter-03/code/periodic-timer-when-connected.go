package main


import (
	"log"
	"time"

	"github.com/nats-io/go-nats"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	for range time.NewTicker(500 * time.Millisecond).C {
		if nc.IsClosed() {
			log.Fatalf("Disconnected forever! Exiting...")
		}
		if nc.IsReconnecting() {
			log.Println("Disconnected temporarily, skipping for now...")
			continue
		}
		err := nc.Publish("numbers", []byte("4 8 15 16 23 42"))
		if err != nil {
			// nats: outbound buffer limit exceeded
			log.Fatalf("Error: %s", err)
		}
	}
}
