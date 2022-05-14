package authentication

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/schema"
)

type UserLength struct {
	Length int
}

type UserDetails struct {
	Username   string
	UserId     string
	Email      string
	Bio        string
	ProfileImg string
	IsVerified string
	Date       string
	Time       string
}

type UID struct {
	UID string
}

type ChatRoomMember struct {
	UID    string
	RoomId string
}

type MemberDetails struct {
	Username     string
	ProfileImg   string
	MemberStatus string
	JoinedOn     string
}

func UIDFetcher(r *http.Request) *UID {
	r.ParseForm()

	uid := new(UID)

	schema.NewDecoder().Decode(uid, r.PostForm)

	return uid
}

func MemberFetcher(r *http.Request) *ChatRoomMember {
	r.ParseForm()

	user := new(ChatRoomMember)

	schema.NewDecoder().Decode(user, r.PostForm)

	return user
}

func M2N(month string) string {
	months := []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}
	for i, m := range months {
		if m == month {
			return strconv.Itoa(i + 1)
		}
	}
	return ""
}

func CheckUser(rw http.ResponseWriter, r *http.Request) {
	user := FormReader(r)

	rows, err := db.Query("SELECT userId FROM users WHERE username=? OR email=?", user.Username, user.Email)
	if err != nil {
		log.Fatalln(err)
	}

	var resultLength int
	for rows.Next() {
		resultLength++
	}

	defer rows.Close()

	json.NewEncoder(rw).Encode(UserLength{resultLength})
}

func QuickCheck(username string, email string) int {
	rows, err := db.Query("SELECT userId FROM users WHERE username=? OR email=? AND username!='' AND email!=''", username, email)
	if err != nil {
		log.Fatalln(err)
	}

	var resultLength int
	for rows.Next() {
		resultLength++
	}

	defer rows.Close()

	return resultLength
}

func FetchUser(rw http.ResponseWriter, r *http.Request) {
	cookieVal := ReadCookie(r)
	UserFetcher(rw, cookieVal)
}

func FetchUserFromChat(rw http.ResponseWriter, r *http.Request) {
	cookieVal := DecodeUID_enc(UIDFetcher(r).UID)
	UserFetcher(rw, cookieVal)
}

func ZeroChecker(min int) string {
	if min < 10 {
		return "0" + strconv.Itoa(min) + ""
	}
	return strconv.Itoa(min)
}

func UserFetcher(rw http.ResponseWriter, uid string) {
	rows, err := db.Query("SELECT username,userId,email,bio,profileImg,isVerified FROM users WHERE userId=?", uid)
	if err != nil {
		log.Println(err)
	}

	var user *UserDetails
	t := time.Now()
	d := t.Day()
	m := M2N(t.Month().String())
	y := t.Year() - 2000
	hr := t.Hour()
	min := t.Minute()
	date := "" + strconv.Itoa(d) + "." + m + "." + strconv.Itoa(y) + ""
	time := "" + strconv.Itoa(hr) + ":" + ZeroChecker(min) + ""
	for rows.Next() {
		var username string
		var userId string
		var email string
		var bio string
		var profileImg string
		var isVerified string

		err := rows.Scan(&username, &userId, &email, &bio, &profileImg, &isVerified)
		if err != nil {
			log.Println("Error scanning rows")
		}
		user = &UserDetails{username, userId, email, bio, profileImg, isVerified, date, time}
	}

	defer rows.Close()

	json.NewEncoder(rw).Encode(user)
}

func FetchId(r *http.Request) string {
	cookieVal := ReadCookie(r)
	return cookieVal
}

func FetchMembershipDetails(rw http.ResponseWriter, r *http.Request) {
	user := MemberFetcher(r)
	uid := DecodeUID_enc(user.UID)

	row := db.QueryRow("SELECT users.username,users.profileImg,chatRoomMembers.memberStatus,chatRoomMembers.joinedOn FROM chatRoomMembers INNER JOIN users ON users.userId = chatRoomMembers.userId WHERE chatRoomMembers.userId=? AND chatRoomId=?", uid, user.RoomId)

	var username string
	var profileImg string
	var memberStatus string
	var joinedOn string

	if err := row.Scan(&username, &profileImg, &memberStatus, &joinedOn); err != nil {
		log.Println(err)
	}

	M_D := MemberDetails{username, profileImg, memberStatus, joinedOn}

	json.NewEncoder(rw).Encode(M_D)
}

func FetchAllMembers(rw http.ResponseWriter, r *http.Request) {
	user := MemberFetcher(r)
	uid := DecodeUID_enc(user.UID)

	rows, err := db.Query("SELECT users.username,users.profileImg,chatRoomMembers.memberStatus,chatRoomMembers.joinedOn FROM chatRoomMembers INNER JOIN users ON users.userId = chatRoomMembers.userId WHERE chatRoomMembers.userId != ? AND chatRoomId=?", uid, user.RoomId)
	if err != nil {
		log.Println(err)
	}

	var MemberDetailsArr []MemberDetails
	for rows.Next() {
		var username string
		var profileImg string
		var memberStatus string
		var joinedOn string

		if err := rows.Scan(&username, &profileImg, &memberStatus, &joinedOn); err != nil {
			log.Println(err)
		}

		M_D := MemberDetails{username, profileImg, memberStatus, joinedOn}
		MemberDetailsArr = append(MemberDetailsArr, M_D)
	}

	defer rows.Close()

	json.NewEncoder(rw).Encode(MemberDetailsArr)
}
