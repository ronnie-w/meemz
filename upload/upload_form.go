package upload

import (
	"log"
	"net/http"

	"github.com/gorilla/schema"
	"github.com/meemz/database"
)

var db = database.Conn()

type Config struct {
	Credits      string
	Tags         string
	Pinned       string
	FileName     string
	OriginalName string
	UploadType   string
	FileIndex   int
}

func FormConfig(r *http.Request) *Config {
	r.ParseForm()
	config := new(Config)
	schema.NewDecoder().Decode(config, r.PostForm)

	return config
}

func CheckOriginalName(originalName string) bool {
	rows, err := db.Query("SELECT originalName FROM posts WHERE originalName=?", originalName)
	if err != nil {
		log.Println(err)
	}

	var length int
	for rows.Next() {
		length++
	}

	if length > 0 {
		return true
	}

	defer rows.Close()

	return false
}

func DuplicateCheck(r *http.Request, textOcr string) (string, int) {
	// rows, err := db.Query("SELECT textOcr FROM posts WHERE textOcr=?", textOcr)
	// if err != nil {
	// 	log.Println(err)
	// }

	// var length int
	// for rows.Next() {
	// 	length++
	// }

	// if length > 0 {
	// 	return "Yes", length
	// }

	// defer rows.Close()

	return "No" /*length*/, 0
}
