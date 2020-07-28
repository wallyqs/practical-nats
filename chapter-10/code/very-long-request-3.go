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

	var i int
	var inbox string
	for ; i < 5; i++ {
		log.Println("[Inbox Request]")
		reply, err := nc.Request("very.long.request", []byte(""), 5*time.Second)
		if err != nil {
			log.Println("Retrying due to errors...")
			continue
		}
		if reply.Reply == "" {
			log.Println("Node replied with empty inbox, retry again later...")
			time.Sleep(1 * time.Second)
			continue
		}

		inbox = reply.Reply
		break
	}
	if i == 5 {
		log.Fatalf("No nodes available to reply!")
	}
	log.Println("[Detected node]", inbox)

	payload := []byte("hi...")
	response, err := nc.Request(inbox, payload, 30*time.Second)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("[Response]", string(response.Data))
}
