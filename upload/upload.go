package upload

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/meemz/authentication"
	//"github.com/meemz/ocr"
	//"github.com/meemz/gcv"
)

type File struct {
	Name string
}

var wg sync.WaitGroup

var file_id, _ = rand.Prime(rand.Reader, 70)

func Uploader(r *http.Request, formFileName string, tempFileDir string, tempFileName string, fileNameSplice int) string {
	fileBytes := new(bytes.Buffer)

	r.ParseMultipartForm(100)

	file, handler, err := r.FormFile(formFileName)
	if err != nil {
		log.Println(err)
	}

	defer file.Close()

	file_ext := handler.Header.Values("Content-Type")[0][6:]

	fmt.Println(handler.Filename, handler.Size, handler.Header.Values("Content-Type")[0][6:])

	tempFile, err := ioutil.TempFile(tempFileDir, tempFileName)
	if err != nil {
		log.Println(err)
	}

	defer tempFile.Close()

	if (file_ext == "jpeg" || file_ext == "jpg") && handler.Size > 20*1024 {
		img, _, err := image.Decode(file)
		if err != nil {
			log.Println("Image decode error")
		}

		err = jpeg.Encode(fileBytes, img, &jpeg.Options{
			Quality: 20,
		})
		if err != nil {
			log.Println("JPEG encode err", err)
		}
	} else {
		f_io_bytes, _ := ioutil.ReadAll(file)
		fileBytes = bytes.NewBuffer(f_io_bytes)
	}

	if _, err := tempFile.Write(fileBytes.Bytes()); err != nil {
		log.Println(err)
	}

	fileName := tempFile.Name()[fileNameSplice:]

	return fileName
}

func VeemzUploader(r *http.Request, formFileName string, tempFileDir string, tempFileName string, fileNameSplice int) string {
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
	fileName := Uploader(r, "meemz_upload", "static/meemz_uploads", "meemz-*.jpg", 21)

	file_dir := "static/meemz_uploads/" + fileName + ""

	wg.Add(1)
	go func() {
		defer wg.Done()
		uploadTime := time.Now().Format(time.RFC3339)
		labelOcr := /*gcv.LabelOcr(file_dir)*/ "some functioning values"
		logoOcr := /*gcv.LogoOcr(file_dir)*/ "some functioning values"
		faceOcr := /*gcv.FaceOcr(file_dir)*/ "some functioning values"
		landmarkOcr := /*gcv.LandmarkOcr(file_dir)*/ "some functioning values"
		textOcr := /*gcv.TextOcr(file_dir)*/ "some functioning values"
		safeSearchOcr := "some functioning values"
		adultContent := "some functioning values"
		violentContent := /*gcv.SafeSearchOcr(file_dir)*/ "some functioning values"
		possibleDuplicate, duplicateNum := DuplicateCheck(r, textOcr)

		if adultContent == "VERY_LIKELY" || adultContent == "LIKELY" || violentContent == "VERY_LIKELY" || violentContent == "LIKELY" || possibleDuplicate == "Yes" {
			err := os.Remove(file_dir)
			if err != nil {
				log.Println(err)
			}
		} else {
			rows, err := db.Query("INSERT INTO posts(userId,fileName,labelOcr,logoOcr,faceOcr,landmarkOcr,textOcr,safeSearchOcr,possibleDuplicate,duplicateNum,uploadTime) values(?,?,?,?,?,?,?,?,?,?,?)", userId, fileName, labelOcr, logoOcr, faceOcr, landmarkOcr, textOcr, safeSearchOcr, possibleDuplicate, duplicateNum, uploadTime)
			if err != nil {
				log.Println(err)
			}

			defer rows.Close()
		}

		json.NewEncoder(rw).Encode(File{fileName})
	}()
	wg.Wait()
}

func UploadVeemz(rw http.ResponseWriter, r *http.Request) {
	userId := authentication.ReadCookie(r)
	fileName := VeemzUploader(r, "veemz_upload", "static/veemz_uploads", "veemz-*.mp4", 21)

	file_dir := "static/veemz_uploads/" + fileName + ""

	wg.Add(1)
	go func() {
		defer wg.Done()
		uploadTime := time.Now().Format(time.RFC3339)
		IsExplicit := /*gcv.ExplicitVideoContent(file_dir)*/ false

		if IsExplicit {
			if err := os.Remove(file_dir); err != nil {
				log.Println(err)
			}
		} else {
			rows, err := db.Query("INSERT INTO posts(userId,fileName,uploadTime) values(?,?,?)", userId, fileName, uploadTime)
			if err != nil {
				log.Println(err)
			}

			defer rows.Close()
		}

		json.NewEncoder(rw).Encode(File{fileName})
	}()
	wg.Wait()
}

func UpdateMeemzConfig(rw http.ResponseWriter, r *http.Request) {
	config := FormConfig(r)

	file_dir := "static/meemz_uploads/" + config.FileName + ""
	exists := CheckOriginalName(config.OriginalName)

	if !exists && config.UploadType == "multiple" {
		update_row, _ := db.Query("UPDATE posts SET tags=?, pComment=?, credits=?, originalName=?, fileId=?, fileIndex=? WHERE fileName=?", config.Tags, config.Pinned, config.Credits, config.OriginalName, file_id.String(), config.FileIndex, config.FileName)
		defer update_row.Close()
	} else if !exists && config.UploadType == "single" {
		new_img_id, _ := rand.Prime(rand.Reader, 70)
		update_row, _ := db.Query("UPDATE posts SET tags=?, pComment=?, credits=?, originalName=?, fileId=?, fileIndex=? WHERE fileName=?", config.Tags, config.Pinned, config.Credits, config.OriginalName, new_img_id.String(), config.FileIndex, config.FileName)
		defer update_row.Close()
	} else {
		err := os.Remove(file_dir)
		if err != nil {
			log.Println(err)
		}

		_ = db.QueryRow("DELETE FROM posts WHERE fileName=?", config.FileName)
	}

	json.NewEncoder(rw).Encode(File{"Upload complete"})
}

func UpdateVeemzConfig(rw http.ResponseWriter, r *http.Request) {
	config := FormConfig(r)

	file_dir := "static/veemz_uploads/" + config.FileName + ""
	exists := CheckOriginalName(config.OriginalName)

	if !exists && config.UploadType == "multiple" {
		update_row, _ := db.Query("UPDATE posts SET tags=?, pComment=?, credits=?, originalName=?, fileId=?, fileIndex=? WHERE fileName=?", config.Tags, config.Pinned, config.Credits, config.OriginalName, file_id.String(), config.FileIndex, config.FileName)
		defer update_row.Close()
	} else if !exists && config.UploadType == "single" {
		new_vid_id, _ := rand.Prime(rand.Reader, 70)

		update_row, _ := db.Query("UPDATE posts SET tags=?, pComment=?, credits=?, originalName=?, fileId=?, fileIndex=? WHERE fileName=?", config.Tags, config.Pinned, config.Credits, config.OriginalName, new_vid_id.String(), config.FileIndex, config.FileName)
		defer update_row.Close()
	} else {
		err := os.Remove(file_dir)
		if err != nil {
			log.Println(err)
		}

		_ = db.QueryRow("DELETE FROM posts WHERE fileName=?", config.FileName)
	}

	json.NewEncoder(rw).Encode(File{"Upload complete"})
}

func GenerateNewId(rw http.ResponseWriter, r *http.Request) {
	file_id, _ = rand.Prime(rand.Reader, 70)

	json.NewEncoder(rw).Encode(File{"GENERATED"})
}
