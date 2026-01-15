package server

import "github.com/GabrielDCelery/chitchat/internal/protocol"

type Server struct {
	register   chan *Client
	unregister chan *Client
	encoder    protocol.Encoder
}
