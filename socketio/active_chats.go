package socketio

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/schema"
)

type C_F struct {
	UserId     string
	ChatType   string
	Chat       string
	ChatDate   string
	ChatTime   string
	ChatRoomId string
}

type MsgPayload struct {
	Message []byte
}

type MsgRoomId struct {
	RoomId string
}

func MsgFetcherForm(r *http.Request) *MsgRoomId {
	r.ParseForm()

	user := new(MsgRoomId)

	schema.NewDecoder().Decode(user, r.PostForm)

	return user
}

func ChatForm(r *http.Request) *C_F {
	r.ParseForm()

	C_Fdetails := new(C_F)
	schema.NewDecoder().Decode(C_Fdetails, r.PostForm)

	return C_Fdetails
}

func PostChatToDb(rw http.ResponseWriter, r *http.Request) {
	CFD := ChatForm(r)
	stmt, err := db.Prepare("INSERT INTO activeChats(userId,chatType,chat,chatDate,chatTime,chatRoomId) values(?,?,?,?,?,?)")
	if err != nil {
		fmt.Println(err)
	}

	if _, err := stmt.Exec(CFD.UserId, CFD.ChatType, CFD.Chat, CFD.ChatDate, CFD.ChatTime, CFD.ChatRoomId); err != nil {
		fmt.Println(err)
	}

	defer stmt.Close()
}

func FetchMessages(rw http.ResponseWriter, r *http.Request) {
	room_id := MsgFetcherForm(r)

	rows, err := db.Query("SELECT chat FROM activeChats WHERE chatRoomId=? ORDER BY id DESC LIMIT 100", room_id.RoomId)
	if err != nil {
		log.Println(err)
	}

	var MSG_payload []MsgPayload
	for rows.Next() {
		var msg []byte

		if err := rows.Scan(&msg); err != nil {
			log.Println(err)
		}

		MSG_payload = append(MSG_payload, MsgPayload{msg})
	}

	defer rows.Close()

	json.NewEncoder(rw).Encode(MSG_payload)
}
