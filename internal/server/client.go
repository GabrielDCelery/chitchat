package server

import (
	"context"

	"github.com/GabrielDCelery/chitchat/internal/protocol"
	"github.com/gorilla/websocket"
)

const (
	sendChanSize = 256
)

type Client struct {
	id       string
	username string
	conn     *websocket.Conn
	server   *Server
	send     chan *protocol.Message
}

func NewClient(id string, username string, conn *websocket.Conn, server *Server) *Client {
	return &Client{
		id:       id,
		username: username,
		conn:     conn,
		server:   server,
		send:     make(chan *protocol.Message, sendChanSize),
	}
}

func (c *Client) readPump(ctx context.Context) {
	defer func() {
		c.server.unregister <- c
		c.conn.Close()
	}()
	// Loop forever, reading messages from the WebSocket
	// - Set read deadline (timeout)
	// - Read message from c.conn
	// - Decode using server.encoder
	// - Send to server.broadcast channel
	// - Handle errors (disconnect, invalid messages)
}

func (c *Client) writePump() {
	// Setup ping ticker (keep connection alive)
	// Loop forever:
	//   - Wait for message on c.send channel
	//   - Set write deadline
	//   - Encode and write to c.conn
	//   - Send periodic pings
	//   - Handle errors
}
