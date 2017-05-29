package chat

import "log"

// Message indicate the structure of a message...
type Message struct {
	recepient int
	message   []byte
	Me        int
}

var (
	clients   map[*Client]bool
	addClient chan *Client
	send      chan Message
)

func init() {
	clients = make(map[*Client]bool, 10)
	addClient = make(chan *Client)
}

// StartServer handles the message routing...
func StartServer() {
	for {
		select {

		case client := <-addClient:
			clients[client] = true

		case m := <-send:
			// send the message to the recepient...
			var sent bool
			// first of all loop through the list of clients to find the recepient
			for c := range clients {
				if c.User.ID == m.recepient {
					// send the message to this user...
					c.send <- m.message
					sent = true
				}
			}
			if sent == false {
				// there is no user with this id...
				log.Println("Trying to send to invalid recepient", m)
			}
		}
	}
}
