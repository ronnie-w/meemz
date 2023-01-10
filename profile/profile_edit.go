package profile

import (
	"encoding/json"
	"log"
	"net/http"
	"net/mail"
	"os"

	"github.com/gorilla/schema"
	"github.com/meemz/authentication"
	"github.com/meemz/database"
	"github.com/meemz/upload"
)

type profileData struct {
	Data string
}

type ProfileImg struct {
	Dp string
}

type FileIdStruct struct {
	FileId string
}

type Error struct {
	Err string
}

func postData(r *http.Request) *profileData {
	r.ParseForm()

	data := new(profileData)
	schema.NewDecoder().Decode(data, r.PostForm)

	return data
}

func PostDp(r *http.Request) *ProfileImg {
	r.ParseForm()

	img := new(ProfileImg)
	schema.NewDecoder().Decode(img, r.PostForm)

	return img
}

func GetFileId(r *http.Request) *FileIdStruct {
	r.ParseForm()

	img := new(FileIdStruct)
	schema.NewDecoder().Decode(img, r.PostForm)

	return img
}

var db = database.Conn()

func PostUsernameToDb(rw http.ResponseWriter, r *http.Request) {
	username := postData(r).Data

	uid := authentication.ReadCookie(r)

	exists := authentication.QuickCheck(username, "")

	if len(username) < 3 {
		json.NewEncoder(rw).Encode(Error{"Username characters must be three or more"})
	} else if exists > 0 {
		json.NewEncoder(rw).Encode(Error{"Username is taken"})
	} else {
		update_row, _ := db.Query("UPDATE users SET username=? WHERE userId=?", username, uid)

		defer update_row.Close()
	}
}

func PostEmailToDb(rw http.ResponseWriter, r *http.Request) {
	email := postData(r).Data

	_, invalid_email := mail.ParseAddress(email)

	uid := authentication.ReadCookie(r)

	exists := authentication.QuickCheck("", email)

	if invalid_email != nil {
		json.NewEncoder(rw).Encode(Error{"The email you entered is invalid"})
	} else if exists > 0 {
		json.NewEncoder(rw).Encode(Error{"The email you entered is already in use"})
	} else {
		update_row, _ := db.Query("UPDATE users SET email=? WHERE userId=?", email, uid)

		defer update_row.Close()
	}
}

func PostBioToDb(rw http.ResponseWriter, r *http.Request) {
	bio := postData(r).Data

	uid := authentication.ReadCookie(r)

	if len(bio) > 0 {
		update_row, err := db.Query("UPDATE users SET bio=? WHERE userId=?", bio, uid)
		if err != nil {
			log.Println(err)
		}
	
		defer update_row.Close()		
	} else {
		update_row, err := db.Query("UPDATE users SET bio=? WHERE userId=?", "No bio found", uid)
		if err != nil {
			log.Println(err)
		}
	
		defer update_row.Close()
	}
}

func ProfileUpload(rw http.ResponseWriter, r *http.Request) {
	filename := upload.Uploader(r, "profile_img", "static/profile-pictures", "profile-image-*.avif", 24)
	json.NewEncoder(rw).Encode(filename)
}

func ProfileUpdateImg(rw http.ResponseWriter, r *http.Request) {
	var old_dp string
	uid := authentication.ReadCookie(r)
	filename := PostDp(r).Dp

	dp_row := db.QueryRow("SELECT profileImg FROM users WHERE userId=?", uid)
	dp_row.Scan(&old_dp)

	if old_dp != "blank-profile-picture-973460_1280.png" {
		err := os.Remove("static/profile-pictures/" + old_dp + "")
		if err != nil {
			log.Println(err)
		}
	}
	rows, err := db.Query("UPDATE users SET profileImg=? WHERE userId=?", filename, uid)
	if err != nil {
		log.Println(err)
	}

	defer rows.Close()
}

func DeletePost(rw http.ResponseWriter, r *http.Request) {
	var file_name_arr []string

	uid := authentication.ReadCookie(r)
	file_id := GetFileId(r).FileId

	file_name_row, _ := db.Query("SELECT fileName FROM posts WHERE fileId=?", file_id)
	for file_name_row.Next() {
		var file_name string
		file_name_row.Scan(&file_name)

		file_name_arr = append(file_name_arr, file_name)
	}

	posts, _ := db.Query("DELETE FROM posts WHERE fileId=?", file_id)
	comments, _ := db.Query("DELETE FROM comments WHERE fileId=?", file_id)
	comment_likes, _ := db.Query("DELETE FROM commentReplyLikes WHERE fileId=?", file_id)
	replies, _ := db.Query("DELETE FROM replies WHERE fileId=?", file_id)
	//reports, _ = db.Query("DELETE FROM reports WHERE fileId=?", file_id)

	defer posts.Close()
	defer comments.Close()
	defer comment_likes.Close()
	defer replies.Close()

	for _, file := range file_name_arr {
		meemz_file_dir := "static/meemz_uploads/" + file + ""
		veemz_file_dir := "static/veemz_uploads/" + file + ""

		if uid != "" {
			err := os.Remove(meemz_file_dir)
			if err != nil {
				err := os.Remove(veemz_file_dir)
				if err != nil {
					json.NewEncoder(rw).Encode(Error{"Could not delete post. Try again"})
				}
			}
		} else {
			json.NewEncoder(rw).Encode(Error{"User authentication failed. You can't delete this post"})
		}
	}

	defer file_name_row.Close()
}
