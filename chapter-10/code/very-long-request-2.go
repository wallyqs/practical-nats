package main

import (
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("[Inbox Request]")
	reply, err := nc.Request("very.long.request", []byte(""), 5*time.Second)
	if err != nil {
		log.Fatalf("No nodes available to reply: %s", err)
	}

	inbox := reply.Reply
	log.Println("[Detected node]", inbox)

	payload := []byte("hi...")
	response, err := nc.Request(inbox, payload, 30*time.Second)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("[Response]", string(response.Data))
}
