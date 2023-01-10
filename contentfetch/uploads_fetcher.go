package contentfetch

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	// "sort"
	// "time"

	"github.com/meemz/authentication"
	"github.com/meemz/database"
)

type Uploads struct {
	Id int
	FileName    string
	UploadTime string
	Pcomment   string
	Tags       string
	Credits    string
	FileId      string
	FileIndex   int
}

var db = database.Conn()

func check(s []Uploads, item Uploads) bool {
	for i := 0; i < len(s); i++ {
		if s[i].FileId == item.FileId && s[i].FileName != item.FileName {
			return true
		}
	}

	return false
}

func group_files(s []Uploads, item Uploads, index int) ([]Uploads, []Uploads) {
	var multiple, grouped_multiple []Uploads

	for i := 0; i < len(s); i++ {
		if s[i].FileId == item.FileId && s[i].FileName != item.FileName {
			grouped_multiple = append(grouped_multiple, s[i])
		} else {
			multiple = append(multiple, s[i])
		}
	}

	return multiple, grouped_multiple
}

func MyUploads(rw http.ResponseWriter, r *http.Request) {
	uid := authentication.ReadCookie(r)

	rows, err := db.Query("SELECT id,fileName,uploadTime,pComment,tags,credits,fileId,fileIndex FROM posts WHERE userId=?", uid)
	if err != nil {
		log.Println(err)
	}

	uploads := ExecQuery(rows)

	var sorted_uploads [][]Uploads

	for i := 0; i < len(uploads); i++ {
		isMultiple := check(uploads, uploads[i])
		if isMultiple {
			var grouped_multiple, multiple []Uploads
			grouped_multiple = append(grouped_multiple, uploads[i])
			uploads, multiple = group_files(uploads, uploads[i], i)
			grouped_multiple = append(grouped_multiple, multiple...)
			sorted_uploads = append(sorted_uploads, grouped_multiple)
		} else {
			sorted_uploads = append(sorted_uploads, []Uploads{uploads[i]})
		}
	}

	json.NewEncoder(rw).Encode(sorted_uploads)
}

func UsersUploads(rw http.ResponseWriter, r *http.Request) {
	username := authentication.FormReader(r).Username

	var uid string
	var profileImg string
	uid_row := db.QueryRow("SELECT userId, profileImg FROM users WHERE username=?", username)
	uid_row.Scan(&uid, &profileImg)

	rows, err := db.Query("SELECT fileName,possibleDuplicate,access,uploadTime,duplicateNum,pComment,tags FROM posts WHERE userId=? AND access=?", uid, "Public")
	if err != nil {
		log.Println(err)
	}

	uploads := ExecQuery(rows)

	json.NewEncoder(rw).Encode(uploads)
}

func ExecQuery(rows *sql.Rows) []Uploads {
	var uploads []Uploads
	for rows.Next() {
		var id int
		var fileName string
		var uploadTime string
		var pComment string
		var tags string
		var credits string
		var fileId string
		var fileIndex int

		err := rows.Scan(&id, &fileName, &uploadTime, &pComment, &tags, &credits, &fileId, &fileIndex)
		if err != nil {
			log.Println(err)
		}

		upload := Uploads{id, fileName, uploadTime, pComment, tags, authentication.TextParser(credits), fileId, fileIndex}

		uploads = append(uploads, upload)
	}

	defer rows.Close()

	return uploads
}
