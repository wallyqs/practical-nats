package main


import (
	"context"
	"log"
	"time"

	"github.com/nats-io/go-nats"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	t := time.AfterFunc(2*time.Second, func() {
		cancel()
	})

	inbox := nats.NewInbox()
	replies := make([]interface{}, 0)
	sub, err := nc.SubscribeSync(inbox)
	if err != nil {
		log.Fatal(err)
	}

	startTime := time.Now()
	nc.PublishRequest("collect", inbox, []byte(""))
	for {
		msg, err := sub.NextMsgWithContext(ctx)
		if err != nil {
			break
		}
		replies = append(replies, msg)

		// Extend deadline on each successful response.
		t.Reset(2 * time.Second)
	}
	log.Printf("Received %d responses in %s", len(replies), time.Since(startTime))
}
