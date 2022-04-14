package socketio

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	gubrak "github.com/novalagung/gubrak/v2"
)

type WebSocketConnection struct {
	Conn *websocket.Conn
	Uid  string
}

type SocketPayload struct {
	Username   string `json:"Username"`
	RoomId     string `json:"RoomId"`
	ProfileImg string `json:"ProfileImg"`
	Message    string `json:"Message"`
}

var connections = make([]*WebSocketConnection, 0)
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func SocketIo(rw http.ResponseWriter, r *http.Request) {
	currentGorillaConn, err := upgrader.Upgrade(rw, r, nil)
	if err != nil {
		log.Println(err)
	}

	uid := r.URL.Query().Get("uid")
	currentConn := WebSocketConnection{Conn: currentGorillaConn, Uid: uid}
	connections = append(connections, &currentConn)

	go HandleIo(&currentConn, connections)
}

func HandleIo(currentConn *WebSocketConnection, connections []*WebSocketConnection) {
	for {
		payload := SocketPayload{}
		err := currentConn.Conn.ReadJSON(&payload)
		if err != nil {
			log.Println(err)
			EjectConnection(currentConn)
			break
		}

		Broadcaster(currentConn, &payload)
	}
}

func Broadcaster(currentConn *WebSocketConnection, payload *SocketPayload) {
	for _, eachConn := range connections {
		if eachConn.Conn == currentConn.Conn {
			continue
		}

		eachConn.Conn.WriteJSON(&payload)
	}
}

func EjectConnection(currentConn *WebSocketConnection) {
	filtered := gubrak.From(connections).Reject(func(each *WebSocketConnection) bool {
		return each == currentConn
	}).Result()

	connections = filtered.([]*WebSocketConnection)
}
