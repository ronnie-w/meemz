package socketio

import (
	"crypto/rand"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/schema"
	"github.com/meemz/activities"
	"github.com/meemz/authentication"
	"github.com/meemz/database"
	"github.com/meemz/upload"
)

type ChatRoom struct {
	Topic      string
	Title      string
	MaxMembers string
	FileName   string
}

type Notification struct {
	Msg        string
	ChatRoomId string
}

var db = database.Conn()

func ConvoForm(r *http.Request) *ChatRoom {
	r.ParseForm()
	convo := new(ChatRoom)

	schema.NewDecoder().Decode(convo, r.PostForm)

	return convo
}

func CreateConvo(rw http.ResponseWriter, r *http.Request) {
	convo_details := ConvoForm(r)

	chatRoomId, _ := rand.Prime(rand.Reader, 100)

	uid := authentication.FetchId(r)

	time := time.Now().Format(time.RFC822)

	stmt1, err := db.Prepare("INSERT INTO chatRooms(userId,topic,title,topicProfile,authorizedNumber,chatRoomId,createdOn) values(?,?,?,?,?,?,?)")
	if err != nil {
		log.Println(err)
	}

	stmt2, err := db.Prepare("INSERT INTO chatRoomMembers(userId,chatRoomId,memberStatus,joinedOn) values(?,?,?,?)")
	if err != nil {
		log.Println(err)
	}

	if _, err := stmt1.Exec(uid, convo_details.Topic, convo_details.Title, convo_details.FileName, convo_details.MaxMembers, chatRoomId.String(), time); err != nil {
		log.Println(err)
	}

	if _, err := stmt2.Exec(uid, chatRoomId.String(), "host", time); err != nil {
		log.Println(err)
	}

	defer stmt1.Close()
	defer stmt2.Close()

	activities.Notify(r, "Just created a new chat room, "+convo_details.Title+". Go check it out")

	json.NewEncoder(rw).Encode(authentication.SecureCookie{Name: "roomid_enc", Value: chatRoomId.String(), Expires: "Thu, 10 Jan 2050 12:00:00 EAT", SameSite: "none"})
}

func ConvoBannerUpload(rw http.ResponseWriter, r *http.Request) {
	fileName := upload.Uploader(r, "convo-banner", "static/convo_banners", "banner-*.jpg", 21)

	json.NewEncoder(rw).Encode(Notification{Msg: fileName})
}
