package main

import (
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL,
		nats.DisconnectHandler(func(nc *nats.Conn) {
			log.Printf("Disconnected!\n")
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Printf("Reconnected to %v!\n", nc.ConnectedUrl())
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			log.Printf("Connection closed. Reason: %q\n", nc.LastError())
		}),
		nats.DiscoveredServersHandler(func(nc *nats.Conn) {
			log.Printf("Server discovered\n")
		}),
		nats.ErrorHandler(func(
			_ *nats.Conn,
			_ *nats.Subscription,
			err error,
		) {
			log.Printf("Async Error: %s", err)
		}),
	)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	nc.Subscribe("hello", func(m *nats.Msg) {
		log.Printf("[Received] %s", string(m.Data))
	})
	for range time.NewTicker(1 * time.Second).C {
		err := nc.Publish("hello", []byte("hello world"))
		if err != nil {
			log.Printf("Error: %s", err)
		}
	}
}
