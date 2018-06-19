package main


import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/go-nats"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	t := time.AfterFunc(1*time.Second, func() {
		cancel()
	})

	inbox := nats.NewInbox()
	replies := make([]string, 0)
	sub, err := nc.SubscribeSync(inbox)
	if err != nil {
		log.Fatal(err)
	}

	startTime := time.Now()
	nc.PublishRequest("_NATS_RIDER.discovery", inbox, []byte(""))
	for {
		msg, err := sub.NextMsgWithContext(ctx)
		if err != nil {
			break
		}
		id := string(msg.Data)
		log.Printf("Found component %q", id)
		replies = append(replies, id)

		// Extend deadline on each successful response.
		t.Reset(1*time.Second)
	}
	log.Printf("Received %d responses in %s", len(replies), time.Since(startTime))

	// Checking status of available components
	for _, componentID := range replies {
		statusSubject := fmt.Sprintf("_NATS_RIDER.%s.status", componentID)
		resp, err := nc.Request(statusSubject, []byte(""), 500*time.Millisecond)
		if err != nil {
			continue
		}
		log.Printf("Status of %q: %s", componentID, string(resp.Data[:50]))
	}
}
