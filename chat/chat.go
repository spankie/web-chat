package chat

import (
	"log"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
	"github.com/spankie/web-chat/config"
	"github.com/spankie/web-chat/messages"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Chat handles incoming chat connection
func Chat(w http.ResponseWriter, r *http.Request) {
	log.Println("::CHAT::")
	ctx := r.Context()
	// accept the websocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Created connection...")
	wr, err := conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return
	}
	log.Println("New writer...")
	// check if there is a new claim
	claims, ok := ctx.Value("Claims").(jwt.MapClaims)
	if claims == nil || !ok {
		// log attempt to access unauthorized page...
		log.Println("No claims. sending error")
		wr.Write([]byte(messages.NotLoggedIn))
		// wr.WriteHeader(http.StatusOK)
		// messages.SendError(w, messages.NotLoggedIn)
		return
	}

	claimsUser := claims["User"].(map[string]interface{})
	claimID := int(claimsUser["ID"].(float64))
	log.Println("claimID:", claimID)

	db := config.Get().DB
	// create a client object for the user.
	client := &Client{User: db[claimID], Conn: conn, send: make(chan []byte)}
	log.Println("New client with id : ", claimID)
	// add the client to a map.
	addClient <- client
	log.Println("Sent client to addClient")
	// send a welcome message
	wr.Write([]byte("Welcome to DEE WEB-CHAT..."))
	if err = wr.Close(); err != nil {
		log.Println("error closing: ", err)
	}
	log.Println("Wrote: Welcome to DEE WEB-CHAT...")
	// launch a goroutine to handle reading and writing to the client
	go client.writePump()
	log.Println("launched writepump", claimID)
	log.Println("starting readpump", claimID)
	client.readPump()
	log.Println("started readpump", claimID)
	// log.Println("::::CHAT::::")
}
