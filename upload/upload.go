package upload

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/meemz/activities"
	"github.com/meemz/authentication"
	"github.com/meemz/gcv"
)

type File struct {
	Name string
}

var wg sync.WaitGroup

func Uploader(r *http.Request, formFileName string, tempFileDir string, tempFileName string, fileNameSplice int) string {
	r.ParseMultipartForm(100)

	file, handler, err := r.FormFile(formFileName)
	if err != nil {
		log.Println(err)
	}

	defer file.Close()
	fmt.Println(handler.Filename, handler.Size, handler.Header)

	tempFile, err := ioutil.TempFile(tempFileDir, tempFileName)
	if err != nil {
		log.Println(err)
	}

	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
	}

	if _, err := tempFile.Write(fileBytes); err != nil {
		log.Println(err)
	}

	fileName := tempFile.Name()[fileNameSplice:]

	return fileName
}

func UploadMeemz(rw http.ResponseWriter, r *http.Request) {
	userId := authentication.ReadCookie(r)
	fileName := Uploader(r, "meemz_upload", "static/uploads", "meemz-*.webp", 15)

	file_dir := "static/uploads/" + fileName + ""

	wg.Add(1)
	go func() {
		defer wg.Done()
		uploadTime := time.Now().Format(time.RFC822)
		labelOcr := gcv.LabelOcr(file_dir)
		logoOcr := gcv.LogoOcr(file_dir)
		faceOcr := gcv.FaceOcr(file_dir)
		landmarkOcr := gcv.LandmarkOcr(file_dir)
		textOcr := gcv.TextOcr(file_dir)
		safeSearchOcr, adultContent, violentContent := gcv.SafeSearchOcr(file_dir)
		possibleDuplicate, duplicateNum := DuplicateCheck(r, textOcr)

		if adultContent == "VERY_LIKELY" || adultContent == "LIKELY" || violentContent == "VERY_LIKELY" || violentContent == "LIKELY" || possibleDuplicate == "Yes" {
			err := os.Remove(file_dir)
			if err != nil {
				log.Println(err)
			}
		} else {
			rows, err := db.Query("INSERT INTO posts(userId,imgName,labelOcr,logoOcr,faceOcr,landmarkOcr,textOcr,safeSearchOcr,possibleDuplicate,duplicateNum,uploadTime) values(?,?,?,?,?,?,?,?,?,?,?)", userId, fileName, labelOcr, logoOcr, faceOcr, landmarkOcr, textOcr, safeSearchOcr, possibleDuplicate, duplicateNum, uploadTime)
			if err != nil {
				log.Println(err)
			}

			defer rows.Close()
		}

		activities.Notify(r, "Fresh Meemz just dropped. Catch up with the latest from your favourites")
		json.NewEncoder(rw).Encode(File{fileName})
	}()
	wg.Wait()
}

func UploadVeemz(rw http.ResponseWriter, r *http.Request)  {
	userId := authentication.ReadCookie(r)
	fileName := Uploader(r, "veemz_upload", "static/veemz_uploads", "veemz-*.mp4", 21)

	file_dir := "static/veemz_uploads/" + fileName + ""

	wg.Add(1)
	go func() {
		defer wg.Done()
		uploadTime := time.Now().Format(time.RFC822)
		IsExplicit := gcv.ExplicitVideoContent(file_dir)

		if IsExplicit {
			if err := os.Remove(file_dir); err != nil {
				log.Println(err)
			}
		} else {
			rows, err := db.Query("INSERT INTO veemz(userId,vidName,uploadTime) values(?,?,?)", userId, fileName, uploadTime)
			if err != nil {
				log.Println(err)
			}

			defer rows.Close()
		}

		activities.Notify(r, "Fresh Veemz just dropped. Catch up with the latest from your favourites")
		json.NewEncoder(rw).Encode(File{fileName})
	}()
	wg.Wait()
}

// func UploadVoiceNote(rw http.ResponseWriter , r *http.Request)  {
// 	fileName := Uploader(r , "meemz_voice_note" , "meemz/public/convo_uploads/voice_notes" , "meemz_voice_note-*.mp3" , 39)

// 	json.NewEncoder(rw).Encode(File{fileName})
// }

func UploadConvoImages(rw http.ResponseWriter, r *http.Request) {
	fileName := Uploader(r, "meemz_convo_images", "static/convo_uploads/image_uploads", "meemz_convo_images-*.png", 35)

	json.NewEncoder(rw).Encode(File{fileName})
}

func UpdateMeemzConfig(rw http.ResponseWriter, r *http.Request) {
	config := FormConfig(r)

	update_row, _ := db.Query("UPDATE posts SET tags=?, pComment=?, access=? WHERE imgName=?", config.Tags, config.Pinned, config.Access, config.FileName)

	defer update_row.Close()
}

func UpdateVeemzConfig(rw http.ResponseWriter, r *http.Request) {
	config := FormConfig(r)

	update_row, _ := db.Query("UPDATE veemz SET tags=?, pComment=?, access=? WHERE vidName=?", config.Tags, config.Pinned, config.Access, config.FileName)

	defer update_row.Close()
}