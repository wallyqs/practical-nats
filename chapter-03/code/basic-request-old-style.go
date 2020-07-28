package main

import (
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect("nats://127.0.0.1:4222",
		nats.UseOldRequestStyle(),
	)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	nc.Subscribe("help", func(m *nats.Msg) {
		log.Printf("[Received]: %s", string(m.Data))
		nc.Publish(m.Reply, []byte("I can help!!!"))
	})
	response, err := nc.Request("help", []byte("help!!"), 1*time.Second)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	log.Println("[Response]: " + string(response.Data))
}
