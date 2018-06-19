package main


import (
	"log"
	"time"

	"sync"

	"github.com/nats-io/go-nats"
)

func main() {
	var busy bool
	var l sync.Mutex
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	myInbox := nats.NewInbox()
	nc.QueueSubscribe("very.long.request", "workers", func(m *nats.Msg) {
		l.Lock()
		shouldSkip := busy
		l.Unlock()

		if shouldSkip {
		        // Reply with empty inbox to signal that
			// was not available to process request.
			nc.PublishRequest(m.Reply, "", []byte(""))
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
