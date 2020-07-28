package main

import (
	"log"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("[Started]")

	inbox := nats.NewInbox()
	nc.Subscribe(inbox, func(m *nats.Msg) {
		log.Printf("Received message on inbox: %+v", m)
	})

	nc.Subscribe("collect", func(m *nats.Msg) {
		nc.Publish(m.Reply, []byte(inbox))
	})

	select {}
}
