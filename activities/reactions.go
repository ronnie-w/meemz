package activities

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/schema"
	"github.com/meemz/authentication"
	"github.com/meemz/database"
)

type Reactions struct {
	ImgName      string
	ReactionType string
}

type CommentsFetched struct {
	Username   string
	ProfileImg string
	Comment    string
	ImgName    string
}

var db = database.Conn()

func ReactionsForm(r *http.Request) *Reactions {
	r.ParseForm()
	r_f_details := new(Reactions)
	schema.NewDecoder().Decode(r_f_details, r.PostForm)

	return r_f_details
}

func update_regommendation(pc float64, imgName string, uid string) {
	row, err := db.Query("UPDATE regommend SET pc=? WHERE imgName=? AND userId=?", pc, imgName, uid)
	if err != nil {
		log.Println(err)
	}

	defer row.Close()
}

func PostReaction(rw http.ResponseWriter, r *http.Request) {
	uid := authentication.FetchId(r)
	reactions_form := ReactionsForm(r)

	insert_row, _ := db.Query("INSERT INTO reactions(userId,imgName,reactionType) values(?,?,?)", uid, reactions_form.ImgName, reactions_form.ReactionType)

	row := db.QueryRow("SELECT pc FROM regommend WHERE imgName=? AND userId=?", reactions_form.ImgName, uid)
	var pc float64
	row.Scan(&pc)

	switch reactions_form.ReactionType {
	case "fa-grin-tears":
		update_regommendation(pc+15.5, reactions_form.ImgName, uid)
	case "fa-grin-tongue-squint":
		update_regommendation(pc+10.5, reactions_form.ImgName, uid)
	case "fa-meh":
		update_regommendation(pc-5.5, reactions_form.ImgName, uid)
	case "fa-sad-tear":
		update_regommendation(pc-10.5, reactions_form.ImgName, uid)
	case "fa-angry":
		update_regommendation(pc-15.5, reactions_form.ImgName, uid)
	}

	defer insert_row.Close()
}

func DeleteReaction(rw http.ResponseWriter, r *http.Request) {
	uid := authentication.FetchId(r)
	reactions_form := ReactionsForm(r)

	row := db.QueryRow("SELECT reactionType FROM reactions WHERE userId=? AND imgName=?", uid, reactions_form.ImgName)
	var current_reaction string
	row.Scan(&current_reaction)

	pcrow := db.QueryRow("SELECT pc FROM regommend WHERE imgName=? AND userId=?", reactions_form.ImgName, uid)
	var pc float64
	pcrow.Scan(&pc)

	switch current_reaction {
	case "fa-grin-tears":
		update_regommendation(pc-15.5, reactions_form.ImgName, uid)
	case "fa-grin-tongue-squint":
		update_regommendation(pc-10.5, reactions_form.ImgName, uid)
	case "fa-meh":
		update_regommendation(pc+5.5, reactions_form.ImgName, uid)
	case "fa-sad-tear":
		update_regommendation(pc+10.5, reactions_form.ImgName, uid)
	case "fa-angry":
		update_regommendation(pc+15.5, reactions_form.ImgName, uid)
	}

	delete_row, _ := db.Query("DELETE FROM reactions WHERE userId=? AND imgName=?", uid, reactions_form.ImgName)

	defer delete_row.Close()
}

func PostReport(rw http.ResponseWriter, r *http.Request) {
	uid := authentication.FetchId(r)
	reactions_form := ReactionsForm(r)

	insert_row, _ := db.Query("INSERT INTO reports(userId,imgName,report) values(?,?,?)", uid, reactions_form.ImgName, reactions_form.ReactionType)

	row := db.QueryRow("SELECT pc FROM regommend WHERE imgName=? AND userId=?", reactions_form.ImgName, uid)
	var pc float64
	row.Scan(&pc)

	defer insert_row.Close()
	update_regommendation(0, reactions_form.ImgName, uid)
}

func DeleteReport(rw http.ResponseWriter, r *http.Request) {
	uid := authentication.FetchId(r)
	reactions_form := ReactionsForm(r)

	delete_row, _ := db.Query("DELETE FROM reports WHERE userId=? AND imgName=?", uid, reactions_form.ImgName)

	defer delete_row.Close()
}

func FetchMyComments(rw http.ResponseWriter, r *http.Request) {
	uid := authentication.FetchId(r)
	imgName := ReactionsForm(r).ImgName

	var FetchedComments []CommentsFetched
	rows, err := db.Query("SELECT users.Username , users.ProfileImg , comments.Comment , comments.imgName FROM comments INNER JOIN users ON users.userId = comments.userId WHERE comments.userId=? AND comments.imgName=?", uid, imgName)
	if err != nil {
		log.Println(err)
	}

	for rows.Next() {
		var username string
		var profileImg string
		var comment string
		var imgName string
		rows.Scan(&username, &profileImg, &comment, &imgName)
		FetchedComments = append(FetchedComments, CommentsFetched{Username: username, ProfileImg: profileImg, Comment: comment, ImgName: imgName})
	}

	defer rows.Close()

	json.NewEncoder(rw).Encode(FetchedComments)
}

func FetchOtherComments(rw http.ResponseWriter, r *http.Request) {
	uid := authentication.FetchId(r)
	imgName := ReactionsForm(r).ImgName

	var FetchedComments []CommentsFetched
	rows, err := db.Query("SELECT users.Username , users.ProfileImg , comments.Comment FROM comments INNER JOIN users ON users.userId = comments.userId WHERE comments.userId!=? AND comments.imgName=?", uid, imgName)
	if err != nil {
		log.Println(err)
	}

	for rows.Next() {
		var username string
		var profileImg string
		var comment string

		if err := rows.Scan(&username, &profileImg, &comment); err != nil {
			log.Println(err)
		}

		FetchedComments = append(FetchedComments, CommentsFetched{Username: username, ProfileImg: profileImg, Comment: comment})
	}

	defer rows.Close()

	json.NewEncoder(rw).Encode(FetchedComments)
}

func PostComment(rw http.ResponseWriter, r *http.Request) {
	uid := authentication.FetchId(r)
	imgName := ReactionsForm(r).ImgName
	comment := ReactionsForm(r).ReactionType

	insert_row, _ := db.Query("INSERT INTO comments(userId,imgName,comment) values(?,?,?)", uid, imgName, comment)

	row := db.QueryRow("SELECT pc FROM regommend WHERE imgName=? AND userId=?", imgName, uid)
	var pc float64
	row.Scan(&pc)

	defer insert_row.Close()
	update_regommendation(pc+25.5, imgName, uid)
}

func Viewed(rw http.ResponseWriter, r *http.Request) {
	uid := authentication.FetchId(r)
	img := ReactionsForm(r).ImgName

	row := db.QueryRow("SELECT imgName FROM regommend WHERE imgName=? AND userId=?", img, uid)
	var exists string
	row.Scan(&exists)

	if exists == "" {
		stmt, _ := db.Query("INSERT INTO regommend(userId,imgName,pc) values(?,?,?)", uid, img, 50)

		defer stmt.Close()
	}
}
