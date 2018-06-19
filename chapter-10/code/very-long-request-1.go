package main


import (
	"log"
	"time"

	"github.com/nats-io/go-nats"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}

	payload := []byte("hi...")
	log.Println("[Request]", string(payload))
	reply, err := nc.Request("very.long.request", payload, 30*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("[Response]", string(reply.Data))
}
