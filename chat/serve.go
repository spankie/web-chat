package chat

import "log"

// Message indicate the structure of a message...
type Message struct {
	Sender    string
	Datetime  string
	Message   string
	Recepient int
}

var (
	clients      map[*Client]int
	addClient    chan *Client
	removeClient chan *Client
	send         chan Message
)

func init() {
	clients = make(map[*Client]int, 10)
	addClient = make(chan *Client)
	removeClient = make(chan *Client)
	send = make(chan Message)
}

// StartServer handles the message routing...
func StartServer() {
	for {
		select {

		case client := <-addClient:
			for _, value := range clients {
				if value == client.User.ID {
					log.Println("User already in the websocket map.")
					break
				}
			}
			log.Println("Adding client...")
			clients[client] = client.User.ID
			log.Println("Added client...")
		case m := <-send:
			// send the message to the recepient...
			log.Println("Server: Message received from client for sending to another client")
			var sent = false
			// first of all loop through the list of clients to find the recepient
			for c := range clients {
				if c.User.ID == m.Recepient {
					// send the message to this user...
					log.Println("User found. sending the message...")
					c.send <- m
					sent = true
				}
			}
			if sent == false {
				// there is no user with this id...
				log.Println("No user with this id...")
				log.Println("Trying to send to invalid recepient", m)
			}
		case client := <-removeClient:
			if _, ok := clients[client]; ok {
				log.Println("ending this connection...Removing the client from the list...")
				delete(clients, client)
				close(client.send)
				log.Println("Client removed.")
			}
		}
	}
}
