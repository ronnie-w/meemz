package activities

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/schema"
	"github.com/meemz/authentication"
)

type Notifications struct {
	Username    string
	ProfileImg  string
	Notify      string
	ReceiveTime string
}

type NotificationStruct struct {
	UID      string
	RoomId   string
	Username string
}

type OnCompletion struct {
	Done string
}

func NotificationForm(r *http.Request) *NotificationStruct {
	r.ParseForm()
	details := new(NotificationStruct)
	schema.NewDecoder().Decode(details, r.PostForm)

	return details
}

func Notify(r *http.Request, notification string) {
	uid := authentication.FetchId(r)

	receiveTime := time.Now().Format(time.RFC822)

	subs_rows, _ := db.Query("SELECT userId FROM subs WHERE creatorId=?", uid)
	for subs_rows.Next() {
		var subscriber string

		subs_rows.Scan(&subscriber)

		var notify_num int
		exists_row, _ := db.Query("SELECT * FROM notifications WHERE userId=? AND notify=?", subscriber, notification)
		for exists_row.Next() {
			notify_num++
		}

		if notify_num == 0 {
			rows, _ := db.Query("INSERT INTO notifications(userId,senderId,notify,receiveTime) values(?,?,?,?)", subscriber, uid, notification, receiveTime)
			defer rows.Close()
		}
	}

	defer subs_rows.Close()
}

func NotifyInvite(rw http.ResponseWriter, r *http.Request) {
	uid_enc := NotificationForm(r).UID
	uid := authentication.DecodeUID_enc(uid_enc)
	room_id := NotificationForm(r).RoomId
	username := NotificationForm(r).Username

	receiveTime := time.Now().Format(time.RFC822)

	var user_uid string
	user_uid_row := db.QueryRow("SELECT userId FROM users WHERE username=?", username)
	user_uid_row.Scan(&user_uid)

	var room_name string
	room_name_row := db.QueryRow("SELECT title FROM chatRooms WHERE chatRoomId=?", room_id)
	room_name_row.Scan(&room_name)

	if user_uid != "" && room_name != "" {
		rows, _ := db.Query("INSERT INTO notifications(userId,senderId,notify,receiveTime) values(?,?,?,?)", user_uid, uid, "Come check out our chat room, you'll love it. Welcome to "+room_name+"", receiveTime)
		defer rows.Close()
		json.NewEncoder(rw).Encode(OnCompletion{Done: "Done"})
	} else {
		json.NewEncoder(rw).Encode(OnCompletion{Done: "Action unsuccessful!"})
	}
}

func FetchNotifications(rw http.ResponseWriter, r *http.Request) {
	uid := authentication.FetchId(r)

	rows, err := db.Query("SELECT users.username, users.profileImg, notifications.notify, notifications.receiveTime FROM notifications INNER JOIN users ON users.userId = notifications.senderId WHERE notifications.userId=?", uid)
	if err != nil {
		log.Println(err)
	}

	var all_notifications []*Notifications
	for rows.Next() {
		var username string
		var profileImg string
		var notify string
		var receiveTime string

		rows.Scan(&username, &profileImg, &notify, &receiveTime)
		all_notifications = append(all_notifications, &Notifications{Username: username, ProfileImg: profileImg, Notify: notify, ReceiveTime: receiveTime})
	}

	json.NewEncoder(rw).Encode(all_notifications)
	defer rows.Close()
}

func DeleteNotifications(rw http.ResponseWriter, r *http.Request) {
	uid := authentication.FetchId(r)
	rows, _ := db.Query("DELETE FROM notifications WHERE userId=?", uid)

	defer rows.Close()
}
