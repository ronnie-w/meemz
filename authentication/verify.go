package authentication

import (
	"encoding/json"
	"log"
	"net/http"
)

func Verify(rw http.ResponseWriter, r *http.Request) {
	user := FormReader(r)
	cookieVal := ReadCookie(r)

	rows, err := db.Query("SELECT serverCode FROM users WHERE userId=?", cookieVal)
	if err != nil {
		log.Println(err)
	}

	var serverCode string
	for rows.Next() {
		err := rows.Scan(&serverCode)
		if err != nil {
			log.Println(err)
		}
	}

	if serverCode != user.VCode {
		json.NewEncoder(rw).Encode(Error{"The code you entered is incorrect"})
	} else {
		stmt, err := db.Prepare("UPDATE users SET isVerified=? WHERE userId=?")
		if err != nil {
			log.Println(err)
		}

		stmt.Exec("Yes", cookieVal)

		defer stmt.Close()
	}

	defer rows.Close()
}
