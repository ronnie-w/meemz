package profile

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/schema"
	"github.com/meemz/authentication"
	"github.com/meemz/database"
	"github.com/meemz/upload"
)

type Username struct {
	Username string
}

type Email struct {
	Email string
}

type Bio struct {
	Bio string
}

type ProfileImg struct {
	Dp string
}

func PostUsername(r *http.Request) *Username {
	r.ParseForm()

	username := new(Username)
	schema.NewDecoder().Decode(username, r.PostForm)

	return username
}

func PostEmail(r *http.Request) *Email {
	r.ParseForm()

	email := new(Email)
	schema.NewDecoder().Decode(email, r.PostForm)

	return email
}

func PostBio(r *http.Request) *Bio {
	r.ParseForm()

	bio := new(Bio)
	schema.NewDecoder().Decode(bio, r.PostForm)

	return bio
}

func PostDp(r *http.Request) *ProfileImg {
	r.ParseForm()

	img := new(ProfileImg)
	schema.NewDecoder().Decode(img, r.PostForm)

	return img
}

var db = database.Conn()

func PostUsernameToDb(rw http.ResponseWriter, r *http.Request) {
	username := PostUsername(r).Username

	uid := authentication.ReadCookie(r)

	update_row, _ := db.Query("UPDATE users SET username=? WHERE userId=?", username, uid)

	defer update_row.Close()
}

func PostEmailToDb(rw http.ResponseWriter, r *http.Request) {
	email := PostEmail(r).Email

	uid := authentication.ReadCookie(r)

	update_row, _ := db.Query("UPDATE users SET email=? WHERE userId=?", email, uid)

	defer update_row.Close()
}

func PostBioToDb(rw http.ResponseWriter, r *http.Request) {
	bio := PostBio(r).Bio

	uid := authentication.ReadCookie(r)

	update_row, err := db.Query("UPDATE users SET bio=? WHERE userId=?", bio, uid)
	if err != nil {
		log.Println(err)
	}

	defer update_row.Close()
}

func ProfileUpload(rw http.ResponseWriter, r *http.Request) {
	filename := upload.Uploader(r, "profile_img", "static/profile-pictures", "profile-image-*.jpg", 24)
	json.NewEncoder(rw).Encode(filename)
}

func ProfileUpdateImg(rw http.ResponseWriter, r *http.Request) {
	uid := authentication.ReadCookie(r)
	filename := PostDp(r).Dp
	rows, err := db.Query("UPDATE users SET profileImg=? WHERE userId=?", filename, uid)
	if err != nil {
		log.Println(err)
	}

	defer rows.Close()
}
