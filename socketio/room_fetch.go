package socketio

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/schema"
	"github.com/meemz/activities"
	"github.com/meemz/authentication"
)

type ChatRoomInfo struct {
	UID string
	Id  string
}

type Details struct {
	HostName         string
	HostDp           string
	Topic            string
	Title            string
	TopicDp          string
	AuthorizedNumber int
	NumberOfMembers  int
	ChatRoomId       string
	CreatedOn        string
}

type Msg struct {
	Msg string
}

type SearchedRoom struct {
	ChatRoom string
	Username string
}

type SearchedRoomDetails struct {
	Topic        string
	Title        string
	TopicProfile string
	ChatRoomId   string
	UID          string
}

func RoomInfo(r *http.Request) *ChatRoomInfo {
	r.ParseForm()

	info := new(ChatRoomInfo)

	schema.NewDecoder().Decode(info, r.PostForm)

	return info
}

func SearchRoomForm(r *http.Request) *SearchedRoom {
	r.ParseForm()

	info := new(SearchedRoom)

	schema.NewDecoder().Decode(info, r.PostForm)

	return info
}

func FetchRoom(rw http.ResponseWriter, r *http.Request) {
	room_id := RoomInfo(r).Id

	row, err := db.Query("SELECT userId FROM chatRooms WHERE chatRoomId=?", room_id)
	if err != nil {
		log.Println(err)
	}

	var roomLength int
	for row.Next() {
		roomLength++
	}

	switch roomLength {
	case 1:
		CheckMember(rw, r, room_id)
	case 0:
		json.NewEncoder(rw).Encode(Msg{Msg: "Room does not exist"})
	}

	defer row.Close()
}

func CheckMember(rw http.ResponseWriter, r *http.Request, id string) {
	uid_enc := RoomInfo(r).UID

	uid := authentication.DecodeUID_enc(uid_enc)

	row, err := db.Query("SELECT userId FROM chatRoomMembers WHERE chatRoomId=? AND userId=?", id, uid)
	if err != nil {
		log.Println("Check member err", err)
	}

	var memberLength int
	for row.Next() {
		memberLength++
	}

	switch memberLength {
	case 1:
		RoomDetails(rw, id)
	case 0:
		AutoJoin(rw, r, id)
	}

	defer row.Close()
}

func AutoJoin(rw http.ResponseWriter, r *http.Request, id string) {
	isFull := MaxChecker(rw, id)

	switch isFull {
	case true:
		json.NewEncoder(rw).Encode(Msg{Msg: "Room is full"})
	case false:
		uid_enc := RoomInfo(r).UID

		uid := authentication.DecodeUID_enc(uid_enc)

		time := time.Now().Format(time.RFC822)

		stmt, err := db.Prepare("INSERT INTO chatRoomMembers(userId,chatRoomId,memberStatus,joinedOn) values(?,?,?,?)")
		if err != nil {
			log.Println("auto join err", err)
		}

		if _, err := stmt.Exec(uid, id, "member", time); err != nil {
			log.Println(err)
		}

		defer stmt.Close()

		RoomDetails(rw, id)
	}
}

func RoomDetails(rw http.ResponseWriter, id string) {
	row := db.QueryRow("SELECT users.username , users.profileImg , chatRooms.topic , chatRooms.title , chatRooms.topicProfile , chatRooms.authorizedNumber , chatRooms.chatRoomId , chatRooms.createdOn FROM chatRooms INNER JOIN users ON users.userId = chatRooms.userId WHERE chatRoomId=?", id)
	num, _ := db.Query("SELECT * FROM chatRoomMembers WHERE chatRoomId=?", id)

	var numOfMembers int
	for num.Next() {
		numOfMembers++
	}

	var hostName string
	var hostDp string
	var topic string
	var title string
	var topicDp string
	var authorizedNumber int
	var chatRoomId string
	var createdOn string
	err := row.Scan(&hostName, &hostDp, &topic, &title, &topicDp, &authorizedNumber, &chatRoomId, &createdOn)
	if err != nil {
		log.Println("Room details", err)
	}

	Room_Details := Details{hostName, hostDp, topic, title, topicDp, authorizedNumber, numOfMembers, chatRoomId, createdOn}

	json.NewEncoder(rw).Encode(Room_Details)

	defer num.Close()
}

func MaxChecker(rw http.ResponseWriter, id string) bool {
	authQuery := db.QueryRow("SELECT authorizedNumber FROM chatRooms WHERE chatRoomId=?", id)

	var authorizedNumber int
	err := authQuery.Scan(&authorizedNumber)
	if err != nil {
		log.Println(err)
	}

	row, err := db.Query("SELECT * FROM chatRoomMembers WHERE chatRoomId=?", id)
	if err != nil {
		log.Println(err)
	}

	var numOfMembers int
	for row.Next() {
		numOfMembers++
	}

	switch {
	case authorizedNumber == numOfMembers:
		return true
	case authorizedNumber > numOfMembers:
		return false
	}

	defer row.Close()

	return true
}

func SearchRoomTopics(rw http.ResponseWriter, r *http.Request) {
	chat_room := SearchRoomForm(r).ChatRoom
	rows, err := db.Query("SELECT topic , title , topicProfile , chatRoomId FROM chatRooms WHERE topic LIKE '%" + chat_room + "%' ORDER BY id DESC")
	if err != nil {
		log.Println(err)
	}

	SearchedRoomsArr := exec_query(rw, r, rows)

	defer rows.Close()
	json.NewEncoder(rw).Encode(SearchedRoomsArr)
}

func SearchedRoomTitles(rw http.ResponseWriter, r *http.Request) {
	chat_room := SearchRoomForm(r).ChatRoom
	rows, err := db.Query("SELECT topic , title , topicProfile , chatRoomId FROM chatRooms WHERE title LIKE '%" + chat_room + "%' ORDER BY id DESC")
	if err != nil {
		log.Println(err)
	}

	SearchedRoomsArr := exec_query(rw, r, rows)

	defer rows.Close()
	json.NewEncoder(rw).Encode(SearchedRoomsArr)
}

func exec_query(rw http.ResponseWriter, r *http.Request, rows *sql.Rows) []SearchedRoomDetails {
	uid := authentication.FetchId(r)
	var SearchedRoomsArr []SearchedRoomDetails
	for rows.Next() {
		var topic string
		var title string
		var topicProfile string
		var chatRoomId string

		if err := rows.Scan(&topic, &title, &topicProfile, &chatRoomId); err != nil {
			log.Println(err)
		}

		isFull := MaxChecker(rw, chatRoomId)

		if !isFull {
			SearchedRoomsArr = append(SearchedRoomsArr, SearchedRoomDetails{Topic: topic, Title: title, TopicProfile: topicProfile, ChatRoomId: chatRoomId, UID: uid})
		}
	}

	return SearchedRoomsArr
}

func ProfileRooms(rw http.ResponseWriter, r *http.Request) {
	user_id := authentication.FetchId(r)
	username := SearchRoomForm(r).Username

	var uid string
	uid_row := db.QueryRow("SELECT userId FROM users WHERE username=?", username)
	uid_row.Scan(&uid)

	rows, err := db.Query("SELECT chatRooms.topic , chatRooms.title , chatRooms.topicProfile , chatRooms.chatRoomId FROM chatRooms INNER JOIN chatRoomMembers ON chatRooms.chatRoomId = chatRoomMembers.chatRoomId WHERE chatRoomMembers.userId=? ORDER BY chatRoomMembers.id DESC", uid)
	if err != nil {
		log.Println(err)
	}

	var SearchedRoomsArr []SearchedRoomDetails
	for rows.Next() {
		var topic string
		var title string
		var topicProfile string
		var chatRoomId string

		if err := rows.Scan(&topic, &title, &topicProfile, &chatRoomId); err != nil {
			log.Println(err)
		}

		if user_id != uid {
			isFull := MaxChecker(rw, chatRoomId)
			if !isFull {
				SearchedRoomsArr = append(SearchedRoomsArr, SearchedRoomDetails{Topic: topic, Title: title, TopicProfile: topicProfile, ChatRoomId: chatRoomId, UID: uid})
			}
		} else {
			SearchedRoomsArr = append(SearchedRoomsArr, SearchedRoomDetails{Topic: topic, Title: title, TopicProfile: topicProfile, ChatRoomId: chatRoomId, UID: uid})
		}

	}

	defer rows.Close()

	json.NewEncoder(rw).Encode(SearchedRoomsArr)
}

func LeaveRoom(rw http.ResponseWriter, r *http.Request) {
	uid_enc := activities.NotificationForm(r).UID
	uid := authentication.DecodeUID_enc(uid_enc)
	room_id := activities.NotificationForm(r).RoomId

	row, _ := db.Query("DELETE FROM chatRoomMembers WHERE userId=? AND chatRoomId=?", uid, room_id)
	defer row.Close()
}
