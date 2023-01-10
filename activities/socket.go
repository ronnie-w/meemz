package activities

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Message struct {
	Message string
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)
var upgrader = websocket.Upgrader{}

func ClientHandler(rw http.ResponseWriter, r *http.Request)  {
	go Broadcaster()
	ws, err := upgrader.Upgrade(rw, r, nil)
	if err != nil {
		log.Println(err)
	}

	defer ws.Close()

	for {
		var message Message
		err := ws.ReadJSON(&message)
		if err != nil {
			log.Println(err)
			delete(clients, ws)
			break
		}

		broadcast <- message
	}
}

func Broadcaster()  {
	for {
		message := <- broadcast
		for client := range clients {
			err := client.WriteJSON(message)
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
	}
}