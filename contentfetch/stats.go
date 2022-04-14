package contentfetch

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/meemz/activities"
)

type Stats struct {
	Reaction1 int
	Reaction2 int
	Reaction3 int
	Reaction4 int
	Reaction5 int

	Comments int
	Views    int
}

func StatsHandler(rw http.ResponseWriter, r *http.Request) {
	ImgName := activities.ReactionsForm(r).ImgName

	rc1, _ := db.Query("SELECT * FROM reactions WHERE reactionType=? AND imgName=?", "fa-grin-tears", ImgName)
	rc2, _ := db.Query("SELECT * FROM reactions WHERE reactionType=? AND imgName=?", "fa-grin-tongue-squint", ImgName)
	rc3, _ := db.Query("SELECT * FROM reactions WHERE reactionType=? AND imgName=?", "fa-meh", ImgName)
	rc4, _ := db.Query("SELECT * FROM reactions WHERE reactionType=? AND imgName=?", "fa-sad-tear", ImgName)
	rc5, _ := db.Query("SELECT * FROM reactions WHERE reactionType=? AND imgName=?", "fa-angry", ImgName)

	comments, _ := db.Query("SELECT * FROM comments WHERE imgName=?", ImgName)
	views, _ := db.Query("SELECT * FROM regommend WHERE imgName=?", ImgName)

	var nrc1 int
	var nrc2 int
	var nrc3 int
	var nrc4 int
	var nrc5 int

	var comments_num int
	var views_num int

	for rc1.Next() {
		nrc1++
	}
	for rc2.Next() {
		nrc2++
	}
	for rc3.Next() {
		nrc3++
	}
	for rc4.Next() {
		nrc4++
	}
	for rc5.Next() {
		nrc5++
	}

	for comments.Next() {
		comments_num++
	}
	for views.Next() {
		views_num++
	}

	all_rows := []*sql.Rows{comments, views, rc1, rc2, rc3, rc4, rc5}
	for _, r := range all_rows {
		defer r.Close()
	}

	json.NewEncoder(rw).Encode(&Stats{Reaction1: nrc1, Reaction2: nrc2, Reaction3: nrc3, Reaction4: nrc4, Reaction5: nrc5, Comments: comments_num, Views: views_num})
}
