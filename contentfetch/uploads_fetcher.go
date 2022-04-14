package contentfetch

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/meemz/authentication"
	"github.com/meemz/database"
)

type Uploads struct {
	ImgName           string
	PossibleDuplicate string
	Access            string
	UploadTime        string
	DuplicateNum      string
	Pcomment          string
	Tags              string
}

var db = database.Conn()

func MyUploads(rw http.ResponseWriter, r *http.Request) {
	uid := authentication.ReadCookie(r)

	rows, err := db.Query("SELECT imgName,possibleDuplicate,access,uploadTime,duplicateNum,pComment,tags FROM posts WHERE userId=?", uid)
	if err != nil {
		log.Println(err)
	}

	uploads := ExecQuery(rows)

	json.NewEncoder(rw).Encode(uploads)
}

func FetchMyPrivateUploads(rw http.ResponseWriter, r *http.Request) {
	uid := authentication.ReadCookie(r)

	rows, err := db.Query("SELECT imgName,possibleDuplicate,access,uploadTime,duplicateNum,pComment,tags FROM posts WHERE userId=? AND access=?", uid, "Private")
	if err != nil {
		log.Println(err)
	}

	uploads := ExecQuery(rows)

	json.NewEncoder(rw).Encode(uploads)
}

func UsersUploads(rw http.ResponseWriter, r *http.Request) {
	username := authentication.FormReader(r).Username

	var uid string
	var profileImg string
	uid_row := db.QueryRow("SELECT userId, profileImg FROM users WHERE username=?", username)
	uid_row.Scan(&uid, &profileImg)

	rows, err := db.Query("SELECT imgName,possibleDuplicate,access,uploadTime,duplicateNum,pComment,tags FROM posts WHERE userId=? AND access=?", uid, "Public")
	if err != nil {
		log.Println(err)
	}

	uploads := ExecQuery(rows)

	json.NewEncoder(rw).Encode(uploads)
}

func ExecQuery(rows *sql.Rows) []Uploads {
	var uploads []Uploads
	for rows.Next() {
		var imgName string
		var possibleDuplicate string
		var access string
		var uploadTime string
		var duplicateNum string
		var pComment string
		var tags string

		err := rows.Scan(&imgName, &possibleDuplicate, &access, &uploadTime, &duplicateNum, &pComment, &tags)
		if err != nil {
			log.Println(err)
		}

		upload := Uploads{imgName, possibleDuplicate, access, uploadTime, duplicateNum, pComment, tags}

		uploads = append(uploads, upload)
	}

	defer rows.Close()

	return uploads
}
