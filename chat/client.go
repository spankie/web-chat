package chat

import (
	"time"

	"log"

	"bytes"

	"encoding/json"

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
	send chan Message
}

func (c *Client) readPump() {
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Time{}) // pass zero value to prevent time out
	// c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Println("Error: ", err)
			}
			log.Println("readpump: normal Error: ", err)
			break
		}
		// arrange the message by getting the recepient id
		// log.Println("ReadPump: Got message from client...", string(message))
		message = bytes.TrimSpace(message)
		// mm := bytes.SplitN(message, newline, 2)
		// log.Println("mm: ", mm[0])
		// rpt, _ := strconv.Atoi(string(mm[0]))
		var m Message
		err = json.Unmarshal(message, &m)
		if err != nil {
			log.Println("Invalid message. Could not Unmarshal to struct message.")
			log.Println("readPUMP:: err: ", err)
			continue
			// send feed back of message not sent...
		}
		log.Println("ReadPump: Arranged the message...Sending to recepient server send channel: ", &m)
		send <- m
		log.Println("ReadPump: sent to the send channel...")
	}
}

func (c *Client) writePump() {
	// TODO:: Try and check if the server closed the channel
	for m := range c.send {
		// if !ok {
		// 	c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
		// 	return
		// }
		log.Println("New message received from server...")
		c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
		w, err := c.Conn.NextWriter(websocket.TextMessage)
		if err != nil {
			return
		}
		log.Println("Writing to self...")
		messageJSON, err := json.Marshal(m)
		if err != nil {
			// find a way to send feedback, the message could not be sent.
			log.Println("Could not encode the message to json.")
			continue
		}
		w.Write(messageJSON)
		log.Println("Wrote to self...")
		if err := w.Close(); err != nil {
			log.Println("Could not close the writer...")
			return
		}
		log.Println("Message: ", string(messageJSON))
	}
}
