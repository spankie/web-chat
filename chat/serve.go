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
	send = make(chan Message)
}

// StartServer handles the message routing...
func StartServer() {
	for {
		select {

		case client := <-addClient:
			log.Println("Adding client...")
			clients[client] = true
			log.Println("Added client...")
		case m := <-send:
			// send the message to the recepient...
			log.Println("Server: Message received from client for sending to another client")
			var sent = false
			// first of all loop through the list of clients to find the recepient
			for c := range clients {
				if c.User.ID == m.recepient {
					// send the message to this user...
					log.Println("User found. sending the message...")
					c.send <- m.message
					sent = true
				}
			}
			if sent == false {
				// there is no user with this id...
				log.Println("No user with this id...")
				log.Println("Trying to send to invalid recepient", m)
			}
		}
	}
}
