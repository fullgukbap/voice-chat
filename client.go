package main

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

func (c *Client) readPump(hub *Hub) {
	defer func() {
		hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		hub.broadcast <- message
	}
}

func (c *Client) handleSend(bytes []byte) {
	w, err := c.conn.NextWriter(websocket.BinaryMessage)
	if err != nil {
		log.Printf("%s가 nextWriter 하는데 실패했습니다 : %v\n", c.conn.RemoteAddr().String(), err)
		return
	}

	if _, err := w.Write(bytes); err != nil {
		log.Printf("%s가 write(전송) 하는데 실패했습니다. :%v\n", c.conn.RemoteAddr(), err)
		return
	}

	if err := w.Close(); err != nil {
		log.Printf("%s가 write를 close하는데 실패했습니다. : %v\n", c.conn.RemoteAddr(), err)
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()

	for {
		select {
		// RECEIVE SEND CHANNEL
		case bytes, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.handleSend(bytes)
		}
	}
}
