package server

type Server struct {
	register   chan *Client
	unregister chan *Client
}
