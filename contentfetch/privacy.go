package contentfetch

import (
	"log"
	"net/http"
	"os"

	"github.com/meemz/activities"
	"github.com/meemz/authentication"
)

func PublicizePost(rw http.ResponseWriter, r *http.Request) {
	uid := authentication.FetchId(r)
	imgName := activities.ReactionsForm(r).ImgName

	q, _ := db.Query("UPDATE posts SET access=? WHERE imgName=? AND userId=?", "Public", imgName, uid)
	defer q.Close()
}

func DeletePost(rw http.ResponseWriter, r *http.Request) {
	uid := authentication.FetchId(r)
	imgName := activities.ReactionsForm(r).ImgName

	q, _ := db.Query("DELETE FROM posts WHERE imgName=? AND userId=?", imgName, uid)
	defer q.Close()

	tables := []string{"reactions", "reports", "comments", "regommend"}
	for _, t := range tables {
		ExecDeleteQuery(t, imgName)
	}
	err := os.Remove("static/uploads/" + imgName + "")
	if err != nil {
		log.Println(err)
	}
}

func ExecDeleteQuery(row_name, imgName string) {
	q, _ := db.Query("DELETE FROM "+row_name+" WHERE imgName=?", imgName)
	defer q.Close()
}
