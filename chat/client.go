package chat

import (
	"time"

	"log"

	"bytes"

	"encoding/binary"

	"github.com/gorilla/websocket"
	"github.com/spankie/web-chat/models"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// Client type...
type Client struct {
	User models.User
	Conn *websocket.Conn
	send chan []byte
}

func (c *Client) readPump() {
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Println("Error: ", err)
			}
			break
		}
		// arrange the message by getting the recepient id
		message = bytes.TrimSpace(message)
		mm := bytes.SplitN(message, newline, 2)
		m := Message{
			recepient: int(binary.BigEndian.Uint64(mm[0])),
			message:   mm[1],
			Me:        c.User.ID,
		}
		send <- m
	}
}

func (c *Client) writePump() {
	for m := range c.send {
		log.Println("Message: ", string(m))
	}
	// for {
	// 	select {
	// 	case message := <-c.send:
	// 		// write the message to the user

	// 	}

	// }
}
