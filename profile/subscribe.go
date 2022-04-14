package profile

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/meemz/authentication"
)

type Subs struct {
	Subscribers int
}

func Subscribe(rw http.ResponseWriter, r *http.Request) {
	Username := authentication.FormReader(r).Username

	my_uid := authentication.FetchId(r)
	time := time.Now().Format(time.RFC822)

	var uid string
	uid_row := db.QueryRow("SELECT userId FROM users WHERE username=?", Username)
	uid_row.Scan(&uid)

	sub, err := db.Query("INSERT INTO subs(userId,creatorId,subTime) values(?,?,?)", my_uid, uid, time)
	if err != nil {
		log.Println(err)
	}

	defer sub.Close()
}

func UnSubscribe(rw http.ResponseWriter, r *http.Request) {
	Username := authentication.FormReader(r).Username

	my_uid := authentication.FetchId(r)

	var uid string
	uid_row := db.QueryRow("SELECT userId FROM users WHERE username=?", Username)
	uid_row.Scan(&uid)

	sub, err := db.Query("DELETE FROM subs WHERE userId=? AND creatorId=?", my_uid, uid)
	if err != nil {
		log.Println(err)
	}

	defer sub.Close()
}

func FetchSubs(rw http.ResponseWriter, r *http.Request) {
	Username := authentication.FormReader(r).Username

	var creator_uid string
	uid_fetcher := db.QueryRow("SELECT userId FROM users WHERE username=?", Username)
	uid_fetcher.Scan(&creator_uid)

	var subs int
	subscribers, _ := db.Query("SELECT * FROM subs WHERE creatorId=?", creator_uid)

	for subscribers.Next() {
		subs++
	}

	defer subscribers.Close()

	json.NewEncoder(rw).Encode(&Subs{Subscribers: subs})
}
