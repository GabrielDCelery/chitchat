package server

import (
	"context"
	"fmt"
	"time"

	"github.com/GabrielDCelery/chitchat/internal/protocol"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	sendChanSize   = 256                 // how many messages do I allow to to pile up on the client's end
	writeWait      = 10 * time.Second    // time allowed to send a message to the user
	pongWait       = 60 * time.Second    // time allowed to wait for the pong from the user
	pingPeriod     = (pongWait * 9) / 10 // Send ping at this interval
	maxMessageSize = 1024 * 1024         // 1MB
)

// This is the server's view of client who is connected
type Client struct {
	id       string
	username string
	conn     *websocket.Conn
	server   *Server
	send     chan *protocol.Message
	logger   *zap.Logger
}

func NewClient(id string, username string, conn *websocket.Conn, server *Server, logger *zap.Logger) *Client {
	if logger == nil {
		logger = zap.NewNop().Named(fmt.Sprintf("client-%s-%s", id, username))
	}
	logger = logger.With(zap.String("clientID", id), zap.String("username", username))
	return &Client{
		id:       id,
		username: username,
		conn:     conn,
		server:   server,
		send:     make(chan *protocol.Message, sendChanSize),
		logger:   logger,
	}
}

func (c *Client) readPump(ctx context.Context) {
	defer func() {
		c.server.unregister <- c // tell the server that the client is diconnecting
		c.conn.Close()           // close the connection
	}()
	// Loop forever, reading messages from the WebSocket
	// - Set read deadline (timeout)
	// - Read message from c.conn
	// - Decode using server.encoder
	// - Send to server.broadcast channel
	// - Handle errors (disconnect, invalid messages)
}

func (c *Client) writePump(ctx context.Context) {
	// start ticker that will periodically check for keeping the connection alive
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()  // stop the ticker
		c.conn.Close() // close the connection
	}()
	for {
		select {
		case <-ctx.Done():
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				c.logger.Error("failed to set write deadline", zap.Error(err))
				return
			}
			// context has been cancelled so we are just closing the connection
			closeMessage := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "server shutting down")
			if err := c.conn.WriteMessage(websocket.CloseMessage, closeMessage); err != nil {
				c.logger.Warn("failed to gracefully shut down client connection")
			}
			return
		case <-ticker.C:
			// time for pinging the client to see if the connection is still alive
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				c.logger.Error("failed to set write deadline", zap.Error(err))
				return
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.logger.Error("failed to ping the client", zap.Error(err))
				return
			}
		case msg, ok := <-c.send:
			// we got a message to send to the client
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				c.logger.Error("failed to set write deadline", zap.Error(err))
				return
			}
			if !ok {
				closeMessage := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "no more messages")
				if err := c.conn.WriteMessage(websocket.CloseMessage, closeMessage); err != nil {
					c.logger.Warn("failed to gracefully shut down client connection", zap.Error(err))
				}
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				c.logger.Error("failed to create writer", zap.Error(err))
				return
			}
			if err := c.server.encoder.Encode(w, msg); err != nil {
				w.Close()
				c.logger.Error("failed to encode message", zap.Error(err))
				return
			}
			if err := w.Close(); err != nil {
				c.logger.Error("failed to close writer", zap.Error(err))
				return
			}
		}
	}
}
