package server

import (
	"context"
	"sync"

	"github.com/GabrielDCelery/chitchat/internal/protocol"
	"go.uber.org/zap"
)

type Server struct {
	register   chan *Client
	unregister chan *Client
	rooms      map[string]*Room
	encoder    protocol.Encoder
	broadcast  chan *protocol.Message
	mu         sync.RWMutex
	logger     *zap.Logger
}

func NewServer(logger *zap.Logger) *Server {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &Server{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		rooms:      make(map[string]*Room),
		encoder:    protocol.NewJSONEncoder(),
		broadcast:  make(chan *protocol.Message),
		logger:     logger,
	}
}

func (s *Server) Run(ctx context.Context) {
	s.logger.Info("server started")
	defer s.logger.Info("server stopped")
	for {
		select {
		case <-ctx.Done():
			s.shutdown()
			s.logger.Info("server shutting down")
			return
		case client := <-s.register:
			s.handleRegister(client)
		case client := <-s.unregister:
			s.handleUnregister(client)
		case client := <-s.broadcast:
			s.handleBroadcast(client)
		}
	}
}

func (s *Server) shutdown() {}

func (s *Server) handleRegister(client *Client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	roomId := client.room
	room, ok := s.rooms[roomId]
	if !ok {
		room = NewRoom(roomId, s.logger)
		s.rooms[roomId] = room
		s.logger.Info("room created", zap.String("roomId", roomId))
	}
	room.Add(client)
	s.logger.Info("client added to room", zap.String("roomId", roomId), zap.String("clientId", client.id), zap.String("username", client.username))
}

func (s *Server) handleUnregister(client *Client) {}

func (s *Server) handleBroadcast(msg *protocol.Message) {}
