package authentication

import (
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func Login(rw http.ResponseWriter, r *http.Request) {
	user := FormReader(r)
	exists := QuickCheck(user.Username, "")
	switch exists {
	case 1:
		rows, err := db.Query("SELECT password,userId FROM users WHERE username=?", user.Username)
		if err != nil {
			log.Println(err)
		}

		var password string
		var userId string
		for rows.Next() {
			err := rows.Scan(&password, &userId)
			if err != nil {
				log.Println(err)
			}
		}
		if err = bcrypt.CompareHashAndPassword([]byte(password), []byte(user.Password)); err != nil {
			json.NewEncoder(rw).Encode(Error{"The password you entered is incorrect"})
		} else {
			CreateCookie(rw, r, userId, "uid")
		}

		defer rows.Close()
	case 0:
		json.NewEncoder(rw).Encode(Error{"User does not exist"})
	}
}
