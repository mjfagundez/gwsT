package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

// Client structure to store ID and connection information
type Client struct {
	ID   int
	Conn *websocket.Conn
}

var clients []Client

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Assign an ID to the client
	clientID := len(clients) + 1

	// Add the new client to the list of clients
	client := Client{ID: clientID, Conn: conn}
	clients = append(clients, client)

	// Send a welcome message to the client with their ID
	welcomeMsg := fmt.Sprintf("Welcome, client #%d!", clientID)
	conn.WriteMessage(websocket.TextMessage, []byte(welcomeMsg))

	// Loop to read incoming messages
	for {
		// Read the message from the client
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		// Add the client ID to the message
		msgWithID := fmt.Sprintf("Client #%d: %s", clientID, string(msg))

		// Broadcast the message to all connected clients
		for _, c := range clients {
			err = c.Conn.WriteMessage(websocket.TextMessage, []byte(msgWithID))
			if err != nil {
				log.Println(err)
				return
			}
		}

		// Print the message to the server console
		fmt.Println(msgWithID)
	}
}
func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ws", wsHandler)
	fmt.Println("Server started on port 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
