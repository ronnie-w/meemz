package authentication

import (
	"encoding/json"
	//"errors"
	"log"
	"net/http"
	"net/url"
	"strings"
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
	JoinDate   string
}

func TextParser(text string) string {
	str_arr := strings.Split(text, " ")

	for i := 0; i < len(str_arr); i++ {
		_, err := url.ParseRequestURI(str_arr[i])

		if err == nil {
			str_arr[i] = "<a href=" + str_arr[i] + " target=`_blank`>" + str_arr[i] + "</a>"
		}
	}

	join_arr := strings.Join(str_arr, " ")
	final_str := strings.ReplaceAll(join_arr, "\n", "<br>")

	return final_str
}

func QuickCheck(username string, email string) int {
	rows := db.QueryRow("SELECT COUNT(*) AS resultLength FROM users WHERE username=? OR email=?", username, email)

	var resultLength int
	rows.Scan(&resultLength)
	
	return resultLength
}

func FetchUser(rw http.ResponseWriter, r *http.Request) {
	cookieVal := ReadCookie(r)

	if cookieVal == "" {
		json.NewEncoder(rw).Encode(Error{"You need to login to proceed"})
	} else {
		UserFetcher(rw, cookieVal)
	}

}

func UserFetcher(rw http.ResponseWriter, uid string) {
	rows, err := db.Query("SELECT username,userId,email,bio,profileImg,isVerified,joinDate FROM users WHERE userId=?", uid)
	if err != nil {
		log.Println(err)
	}

	var user *UserDetails
	for rows.Next() {
		var username string
		var userId string
		var email string
		var bio string
		var profileImg string
		var isVerified string
		var joinDate string

		err := rows.Scan(&username, &userId, &email, &bio, &profileImg, &isVerified, &joinDate)
		if err != nil {
			log.Println("Error scanning rows")
		}

		bio = TextParser(bio)

		user = &UserDetails{username, userId, email, bio, profileImg, isVerified, joinDate}
	}

	defer rows.Close()

	json.NewEncoder(rw).Encode(user)
}

func FetchId(r *http.Request) string {
	cookieVal := ReadCookie(r)
	return cookieVal
}
