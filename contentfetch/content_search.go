package contentfetch

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/meemz/authentication"
)

type Posts struct {
	Username   string
	ProfileImg string

	ImgName string

	LabelOcr      string
	LogoOcr       string
	FaceOcr       string
	LandmarkOcr   string
	TextOcr       string
	SafeSearchOcr string

	PossibleDuplicate string
	Tags              string
	PComment          string
	UploadTime        string

	Reaction1 string
	Reaction2 string
	Reaction3 string
	Reaction4 string
	Reaction5 string
}

var MeemzContent []*Posts
var TagsContent []*Posts

type SearchedUsers struct {
	Username     string
	ProfileImg   string
	Bio          string
	Subscription string
	BGC          string
	Color        string
}

func SearchMeemz(rw http.ResponseWriter, r *http.Request) {
	ImgOcr := authentication.FormReader(r).ImgOcr

	MeemzContent = []*Posts{}

	uid := authentication.FetchId(r)

	rows, err := db.Query("SELECT users.username, users.profileImg, posts.imgName, posts.labelOcr, posts.logoOcr, posts.faceOcr, posts.landmarkOcr, posts.textOcr, posts.safeSearchOcr, posts.possibleDuplicate, posts.tags, posts.pComment, posts.uploadTime FROM posts INNER JOIN users ON users.userId = posts.userId WHERE textOcr LIKE '%"+ImgOcr+"%' AND access=? ORDER BY posts.id DESC LIMIT 200", "Public")
	if err != nil {
		log.Println(err)
	}

	ExecQueryMeemz(rows, uid)

	json.NewEncoder(rw).Encode(MeemzContent)
}

func SearchUsers(rw http.ResponseWriter, r *http.Request) {
	Username := authentication.FormReader(r).Username

	uid := authentication.FetchId(r)

	var Users = []*SearchedUsers{}

	rows, err := db.Query("SELECT username, profileImg, userId, bio FROM users WHERE userId!=? AND username LIKE '%"+Username+"%'", uid)
	if err != nil {
		log.Println(err)
	}

	for rows.Next() {
		var username string
		var profileImg string
		var bio string
		var userId string
		var subscribed string
		var bgc string
		var color string

		if err := rows.Scan(&username, &profileImg, &userId, &bio); err != nil {
			log.Println(err)
		}

		var sub_status string
		sub := db.QueryRow("SELECT userId FROM subs WHERE userId=? AND creatorId=?", uid, userId)
		sub.Scan(&sub_status)

		if sub_status == "" {
			subscribed = "Subscribe"
			bgc = "white"
			color = "black"
		} else {
			subscribed = "Unsubscribe"
			bgc = "black"
			color = "wheat"
		}

		Users = append(Users, &SearchedUsers{Username: username, ProfileImg: profileImg, Bio: bio, Subscription: subscribed, BGC: bgc, Color: color})
	}

	defer rows.Close()

	json.NewEncoder(rw).Encode(Users)
}

func SearchTags(rw http.ResponseWriter, r *http.Request) {
	Tags := authentication.FormReader(r).Tags

	TagsContent = []*Posts{}

	uid := authentication.FetchId(r)

	rows, err := db.Query("SELECT users.username, users.profileImg, posts.imgName, posts.labelOcr, posts.logoOcr, posts.faceOcr, posts.landmarkOcr, posts.textOcr, posts.safeSearchOcr, posts.possibleDuplicate, posts.tags, posts.pComment, posts.uploadTime FROM posts INNER JOIN users ON users.userId = posts.userId WHERE tags LIKE '%"+Tags+"%' AND access=? ORDER BY posts.id DESC LIMIT 200", "Public")
	if err != nil {
		log.Println(err)
	}

	ExecQueryTags(rows, uid)

	json.NewEncoder(rw).Encode(TagsContent)
}

func ExecQueryMeemz(rows *sql.Rows, uid string) {
	for rows.Next() {
		var username string
		var profileImg string

		var imgName string

		var labelOcr string
		var logoOcr string
		var faceOcr string
		var landmarkOcr string
		var textOcr string
		var safeSearchOcr string

		var possibleDuplicate string
		var tags string
		var pComment string
		var uploadTime string

		var reaction string
		var reaction1 string = "far fa-grin-tears"
		var reaction2 string = "far fa-grin-tongue-squint"
		var reaction3 string = "far fa-meh"
		var reaction4 string = "far fa-sad-tear"
		var reaction5 string = "far fa-angry"

		if err := rows.Scan(&username, &profileImg, &imgName, &labelOcr, &logoOcr, &faceOcr, &landmarkOcr, &textOcr, &safeSearchOcr, &possibleDuplicate, &tags, &pComment, &uploadTime); err != nil {
			log.Println(err)
		}

		row := db.QueryRow("SELECT reactionType FROM reactions WHERE userId=? AND imgName=?", uid, imgName)
		row.Scan(&reaction)

		switch reaction {
		case "fa-grin-tears":
			reaction1 = "fas fa-grin-tears"
		case "fa-grin-tongue-squint":
			reaction2 = "fas fa-grin-tongue-squint"
		case "fa-meh":
			reaction3 = "fas fa-meh"
		case "fa-sad-tear":
			reaction4 = "fas fa-sad-tear"
		case "fa-angry":
			reaction5 = "fas fa-angry"
		}

		MeemzContent = append(MeemzContent, &Posts{Username: username, ProfileImg: profileImg, ImgName: imgName, LabelOcr: labelOcr, LogoOcr: logoOcr, FaceOcr: faceOcr, LandmarkOcr: landmarkOcr, TextOcr: textOcr, SafeSearchOcr: safeSearchOcr, PossibleDuplicate: possibleDuplicate, Tags: tags, PComment: pComment, UploadTime: uploadTime, Reaction1: reaction1, Reaction2: reaction2, Reaction3: reaction3, Reaction4: reaction4, Reaction5: reaction5})
	}

	defer rows.Close()
}

func ExecQueryTags(rows *sql.Rows, uid string) {
	for rows.Next() {
		var username string
		var profileImg string

		var imgName string

		var labelOcr string
		var logoOcr string
		var faceOcr string
		var landmarkOcr string
		var textOcr string
		var safeSearchOcr string

		var possibleDuplicate string
		var tags string
		var pComment string
		var uploadTime string

		var reaction string
		var reaction1 string = "far fa-grin-tears"
		var reaction2 string = "far fa-grin-tongue-squint"
		var reaction3 string = "far fa-meh"
		var reaction4 string = "far fa-sad-tear"
		var reaction5 string = "far fa-angry"

		if err := rows.Scan(&username, &profileImg, &imgName, &labelOcr, &logoOcr, &faceOcr, &landmarkOcr, &textOcr, &safeSearchOcr, &possibleDuplicate, &tags, &pComment, &uploadTime); err != nil {
			log.Println(err)
		}

		row := db.QueryRow("SELECT reactionType FROM reactions WHERE userId=? AND imgName=?", uid, imgName)
		row.Scan(&reaction)

		switch reaction {
		case "fa-grin-tears":
			reaction1 = "fas fa-grin-tears"
		case "fa-grin-tongue-squint":
			reaction2 = "fas fa-grin-tongue-squint"
		case "fa-meh":
			reaction3 = "fas fa-meh"
		case "fa-sad-tear":
			reaction4 = "fas fa-sad-tear"
		case "fa-angry":
			reaction5 = "fas fa-angry"
		}

		TagsContent = append(TagsContent, &Posts{Username: username, ProfileImg: profileImg, ImgName: imgName, LabelOcr: labelOcr, LogoOcr: logoOcr, FaceOcr: faceOcr, LandmarkOcr: landmarkOcr, TextOcr: textOcr, SafeSearchOcr: safeSearchOcr, PossibleDuplicate: possibleDuplicate, Tags: tags, PComment: pComment, UploadTime: uploadTime, Reaction1: reaction1, Reaction2: reaction2, Reaction3: reaction3, Reaction4: reaction4, Reaction5: reaction5})
	}

	defer rows.Close()
}
