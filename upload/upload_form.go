package upload

import (
	"log"
	"net/http"

	"github.com/gorilla/schema"
	"github.com/meemz/database"
)

var db = database.Conn()

type Config struct {
	Access   string
	Tags     string
	Pinned   string
	FileName string
}

func FormConfig(r *http.Request) *Config {
	r.ParseForm()
	config := new(Config)
	schema.NewDecoder().Decode(config, r.PostForm)

	return config
}

func DuplicateCheck(r *http.Request, textOcr string) (string, int) {
	rows, err := db.Query("SELECT textOcr FROM posts WHERE textOcr=?", textOcr)
	if err != nil {
		log.Println(err)
	}

	var length int
	for rows.Next() {
		length++
	}

	if length > 0 {
		return "Yes", length
	}

	defer rows.Close()

	return "No", length
}
