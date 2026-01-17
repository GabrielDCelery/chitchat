package server

import (
	"sync"

	"github.com/GabrielDCelery/chitchat/internal/protocol"
	"go.uber.org/zap"
)

type Room struct {
	id      string
	clients map[*Client]bool
	mu      sync.RWMutex
	logger  *zap.Logger
}

func NewRoom(id string, logger *zap.Logger) *Room {
	if logger == nil {
		logger = zap.NewNop()
	}
	logger = logger.With(zap.String("roomId", id))
	return &Room{
		id:      id,
		clients: make(map[*Client]bool),
		logger:  logger,
	}
}

func (r *Room) Add(client *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.clients[client] = true
	r.logger.Info("client added to room", zap.String("clientId", client.id))
}

func (r *Room) Remove(client *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.clients, client)
	r.logger.Info("client removed from room", zap.String("clientId", client.id))
}

func (r *Room) Broadcast(msg *protocol.Message) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	sent := 0
	dropped := 0
	for client := range r.clients {
		select {
		case client.send <- msg:
			sent++
		default:
			dropped++
			r.logger.Warn("client channel full, dropping message", zap.String("clientId", client.id))
		}
	}
	r.logger.Debug("broadcast complete", zap.String("messageType", msg.Type.String()), zap.Int("sent", sent), zap.Int("dropped", dropped), zap.Int("total", len(r.clients)))
}

func (r *Room) Size() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.clients)
}
