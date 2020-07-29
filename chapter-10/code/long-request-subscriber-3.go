package main

import (
	"log"
	"time"

	"sync"

	"github.com/nats-io/nats.go"
)

func main() {
	var busy bool
	var l sync.Mutex
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	myInbox := nats.NewInbox()
	nc.Subscribe("very.long.request", func(m *nats.Msg) {
		var shouldSkip bool

		l.Lock()
		shouldSkip = busy
		l.Unlock()

		// Only reply when not busy
		if shouldSkip {
			return
		}

		log.Println("[Processing] Announcing own inbox...")
		nc.PublishRequest(m.Reply, myInbox, []byte(""))
	})

	nc.Subscribe(myInbox, func(m *nats.Msg) {
		log.Println("[Processing] Message:", string(m.Data))
		l.Lock()
		busy = true
		l.Unlock()
		time.Sleep(20 * time.Second)

		l.Lock()
		busy = false
		l.Unlock()
		nc.Publish(m.Reply, []byte("done!"))
	})

	log.Println("[Started]")
	select {}
}
