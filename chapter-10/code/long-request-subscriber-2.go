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
	myInbox := nats.NewInbox()
	nc.Subscribe("very.long.request", func(m *nats.Msg) {
		log.Println("[Processing] Announcing own inbox...")
		nc.PublishRequest(m.Reply, myInbox, []byte(""))
	})

	nc.Subscribe(myInbox, func(m *nats.Msg) {
		log.Println("[Processing] Message:", string(m.Data))
		time.Sleep(20 * time.Second)
		nc.Publish(m.Reply, []byte("done!"))
	})

	log.Println("[Started]")
	select {}
}
