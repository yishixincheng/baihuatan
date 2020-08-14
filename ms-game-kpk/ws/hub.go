package ws

import (
	"log"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)


// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// registered clients
	clients map[*Client]bool

	// Inbound message from the clients.
	broadcasts chan []byte

	// register requests from the clients.
	register chan *Client

	// unregister reqeusts from the clients
	unregister chan *Client
}

// newHub -
func newHub() *Hub {
	return &Hub{
		broadcasts: make(chan []byte),
		register: make(chan *Client),
		unregister: make(chan *Client),
		clients: make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcasts:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}		
		}
	}
}