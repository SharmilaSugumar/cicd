package websocket

import (
	"encoding/json"
	"forgeflow/pkg/logger"
	"sync"

	"go.uber.org/zap"
)

type EventType string

const (
	EventJobStatusUpdate      EventType = "JOB_STATUS_UPDATE"
	EventPipelineStatusUpdate EventType = "PIPELINE_STATUS_UPDATE"
	EventLiveLog              EventType = "LIVE_LOG"
	EventWorkerStatusUpdate   EventType = "WORKER_STATUS_UPDATE"
	EventQueueUpdate          EventType = "QUEUE_UPDATE"
)

type Event struct {
	Type    EventType       `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			logger.Log.Info("WebSocket client registered")

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
			logger.Log.Info("WebSocket client unregistered")

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// BroadcastEvent is a helper for services to send typed events
func (h *Hub) BroadcastEvent(eventType EventType, payload interface{}) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		logger.Log.Error("Failed to marshal event payload", zap.Error(err))
		return
	}

	event := Event{
		Type:    eventType,
		Payload: payloadBytes,
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		logger.Log.Error("Failed to marshal event", zap.Error(err))
		return
	}

	h.broadcast <- eventBytes
}
