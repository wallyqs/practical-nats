package main


import (
	"log"
	"runtime"

	"github.com/nats-io/go-nats"
)

func main() {
	nc, err := nats.Connect("nats://127.0.0.1:4224")
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	log.Println("All Servers:", nc.Servers())
	log.Println("Discovered Servers:", nc.DiscoveredServers())

	nc.Subscribe("hi", func(m *nats.Msg) {
		log.Println("[Received] ", string(m.Data))
	})

	nc.Publish("hi", []byte("hello world"))

	runtime.Goexit()
}
