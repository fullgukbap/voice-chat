package main

import (
	"log"
	"sync"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client

	mutext sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) handleBroadcast(bytes []byte) {
	h.mutext.Lock()
	defer h.mutext.Unlock()

	log.Println("[broadcast] data: ", bytes)
	for client := range h.clients {
		client.send <- bytes
	}

}

func (h *Hub) handleRegister(client *Client) {
	h.mutext.Lock()
	defer h.mutext.Unlock()

	h.clients[client] = true
}

func (h *Hub) handleUnregister(client *Client) {
	h.mutext.Lock()
	defer h.mutext.Unlock()

	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client.send)
	}
}

func (h *Hub) Run() {
	for {
		select {

		// REGISTER
		case client := <-h.register:
			h.handleRegister(client)

		// UNREGISTER
		case client := <-h.unregister:
			h.handleUnregister(client)

		// BROADCAST
		case bytes := <-h.broadcast:
			h.handleBroadcast(bytes)
		}
	}

}
