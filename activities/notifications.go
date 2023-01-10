package activities

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/meemz/authentication"
)

type Notifications struct {
	Username    string
	ProfileImg  string
	Notify      string
	ReceiveTime string
}

func Notify(r *http.Request, notification string, user_ids []string) {
	uid := authentication.FetchId(r)
	receive_time := time.Now().Format(time.RFC3339)

	for _, receiver_id := range user_ids{
		if receiver_id != uid {
			insert_receiver, _ := db.Query("INSERT INTO notifications(userId, receiverId, notify, receiveTime) values(?,?,?,?)", uid, receiver_id, notification, receive_time)
			defer insert_receiver.Close()			
		}
	}
}

func FetchNotifications(rw http.ResponseWriter, r *http.Request) {
	uid := authentication.FetchId(r)

	rows, err := db.Query("SELECT users.username, users.profileImg, notifications.notify, notifications.receiveTime FROM notifications INNER JOIN users ON users.userId = notifications.userId WHERE notifications.receiverId=?", uid)
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
		all_notifications = append(all_notifications, &Notifications{Username: username, ProfileImg: profileImg, Notify: notify, ReceiveTime: TimeStamp(receiveTime)})
	}

	json.NewEncoder(rw).Encode(all_notifications)

	defer rows.Close()
}

func DeleteNotifications(rw http.ResponseWriter, r *http.Request) {
	uid := authentication.FetchId(r)
	rows, _ := db.Query("DELETE FROM notifications WHERE receiverId=?", uid)

	defer rows.Close()
}
