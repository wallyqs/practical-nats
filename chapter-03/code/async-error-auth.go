package main

import (
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL,
		nats.ErrorHandler(func(
			_ *nats.Conn,
			_ *nats.Subscription,
			err error,
		) {
			log.Printf("Async Error: %s", err)
		}))
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	nc.Publish("_SYS.hi", []byte("hi"))
	nc.Flush()
	time.Sleep(1 * time.Second)
}
