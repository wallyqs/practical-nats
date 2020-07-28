package driveragent

import (
	"encoding/json"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/wallyqs/practical-nats/chapter-06/code/kit"
)

const Version = "0.1.0"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Agent is the agent from the driver that provides rides.
type Agent struct {
	*kit.Component

	// mu is the lock from the agent.
	mu sync.Mutex

	// AgentType is the type of vehicle.
	AgentType string
}

// Type returns the type of vehicle from the driver.
func (s *Agent) Type() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.AgentType
}

// SetupSubscriptions prepares the NATS subscriptions.
func (s *Agent) SetupSubscriptions() error {
	nc := s.NATS()

	nc.Subscribe("drivers.rides", func(msg *nats.Msg) {
		if err := s.processDriveRequest(msg); err != nil {
			log.Printf("Error: %s\n", err)
			return
		}
	})

	return nil
}

func (s *Agent) processDriveRequest(msg *nats.Msg) error {
	var req *kit.DriverAgentRequest
	err := json.Unmarshal(msg.Data, &req)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return err
	}
	log.Printf("requestID:%s - Driver Ride Request: type:%s\n",
		req.RequestID, req.Type)

	if req.Type != s.Type() {
		// Skip request since this agent is of a different type.
		return nil
	}
	log.Printf("requestID:%s - Available to handle request", req.RequestID)

	// Randomly delay agent when receiving drive request
	// to simulate latency in replying.
	duration := time.Duration(rand.Int31n(1000)) * time.Millisecond
	log.Printf("requestID:%s - Backing off for %s before replying", req.RequestID, duration)
	time.Sleep(duration)

	// Replying back with own ID meaning that can help.
	return s.NATS().Publish(msg.Reply, []byte(s.ID()))
}
