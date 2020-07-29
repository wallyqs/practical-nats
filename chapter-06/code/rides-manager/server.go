package ridesmanager

import (
	"encoding/json"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/wallyqs/practical-nats/chapter-06/code/kit"
)

const (
	Version = "0.1.0"
)

type Server struct {
	*kit.Component
}

// SetupSubscriptions registers interest to the subjects that the
// Rides Manager will be handling.
func (s *Server) SetupSubscriptions() error {
	nc := s.NATS()

	// Helps finding an available driver to accept a drive request.
	nc.QueueSubscribe("drivers.find", "manager", func(msg *nats.Msg) {
		var req *kit.DriverAgentRequest
		err := json.Unmarshal(msg.Data, &req)
		if err != nil {
			log.Printf("Error: %v\n", err)
			return
		}
		log.Printf("requestID:%s - Driver Find Request\n", req.RequestID)
		response := &kit.DriverAgentResponse{}

		// Find an available driver that can handle the user request.
		m, err := nc.Request("drivers.rides", msg.Data, 2*time.Second)
		if err != nil {
			response.Error = "No drivers available found, sorry!"
			resp, err := json.Marshal(response)
			if err != nil {
				log.Printf("requestID:%s - Error preparing response: %s",
					req.RequestID, err)
				return
			}

			// Reply with error response
			nc.Publish(msg.Reply, resp)
			return
		}
		response.ID = string(m.Data)

		resp, err := json.Marshal(response)
		if err != nil {
			response.Error = "No drivers available found, sorry!"
			resp, err := json.Marshal(response)
			if err != nil {
				log.Printf("requestID:%s - Error preparing response: %s",
					req.RequestID, err)
				return
			}

			// Reply with error response
			nc.Publish(msg.Reply, resp)
			return
		}
		log.Printf("requestID:%s - Driver Find Response: %+v\n",
			req.RequestID, string(m.Data))
		nc.Publish(msg.Reply, resp)
	})

	return nil
}
