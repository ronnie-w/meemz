package regommend

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/meemz/authentication"
	"github.com/meemz/database"
	"github.com/muesli/regommend"
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

type Veemz struct {
	Username   string
	ProfileImg string

	ImgName string

	Tags       string
	PComment   string
	UploadTime string

	Reaction1 string
	Reaction2 string
	Reaction3 string
	Reaction4 string
	Reaction5 string
}

type Meemz struct {
	ImgName interface{}
	PC      float64
}

var db = database.Conn()

var MainContent []*Posts

var VeemzMainContent []*Veemz

func TimeStamp(time_stamp string) string {
	t := time.Now().Format(time.RFC3339)
	t0, _ := time.Parse(time.RFC3339, t)
	t1, _ := time.Parse(time.RFC3339, time_stamp)

	return t0.Sub(t1).String()
}

func FetchRegommendedMeemz(r *http.Request, content_type string) []string {
	uid := authentication.FetchId(r)
	meemz := regommend.Table("meemz")

	viewed_by_me := make(map[interface{}]float64)
	viewed_by_others := make(map[interface{}]float64)

	vbm_rows, _ := db.Query("SELECT imgName, pc FROM regommend WHERE userId=? AND content_type=?", uid, content_type)
	vbo_rows, _ := db.Query("SELECT imgName, pc FROM regommend WHERE userId!=? AND content_type=?", uid, content_type)

	for vbm_rows.Next() {
		var imgName string
		var pc float64

		vbm_rows.Scan(&imgName, &pc)
		viewed_by_me[imgName] = pc
	}
	meemz.Add("viewed_by_me", viewed_by_me)

	for vbo_rows.Next() {
		var imgName string
		var pc float64

		vbo_rows.Scan(&imgName, &pc)
		viewed_by_others[imgName] = pc
	}
	meemz.Add("viewed_by_others", viewed_by_others)

	var regommended []string
	recs, _ := meemz.Recommend("viewed_by_me")

	for _, r := range recs {
		if r.Distance > 0 {
			regommended = append(regommended, fmt.Sprintf("%s", r.Key))
		}
	}

	defer vbm_rows.Close()
	defer vbo_rows.Close()

	return regommended
}

func Subscription(r *http.Request, table string, column string) []string {
	uid := authentication.FetchId(r)
	subs, _ := db.Query("SELECT creatorId FROM subs WHERE userId=?", uid)

	var creators []string
	for subs.Next() {
		var creator string

		subs.Scan(&creator)
		creators = append(creators, creator)
	}

	var subscription []string
	for _, c := range creators {
		subscriptions, _ := db.Query("SELECT "+column+" FROM "+table+" WHERE userId=? ORDER BY id DESC LIMIT 10", c)
		for subscriptions.Next() {
			var imgName string

			subscriptions.Scan(&imgName)
			subscription = append(subscription, imgName)
		}
		defer subscriptions.Close()
	}

	defer subs.Close()
	return subscription
}

func FetchViewedMeemz(r *http.Request, content_type string) []string {
	uid := authentication.FetchId(r)

	var viewed []string

	v_rows, err := db.Query("SELECT imgName FROM regommend WHERE userId=? AND content_type=?", uid, content_type)
	if err != nil {
		log.Println(err)
	}

	for v_rows.Next() {
		var imgName string

		if err := v_rows.Scan(&imgName); err != nil {
			log.Println(err)
		}

		viewed = append(viewed, imgName)
	}

	defer v_rows.Close()

	return viewed
}

func FetchAllMeemz(r *http.Request, table string, column string) []string {
	var selected_meemz []string

	rows, err := db.Query("SELECT "+column+" FROM "+table+" WHERE access=? ORDER BY "+table+".id DESC LIMIT 200", "Public")
	if err != nil {
		log.Println(err)
	}

	for rows.Next() {
		var imgName string

		rows.Scan(&imgName)

		selected_meemz = append(selected_meemz, imgName)
	}

	rand.Shuffle(len(selected_meemz), func(i, j int) {
		selected_meemz[i], selected_meemz[j] = selected_meemz[j], selected_meemz[i]
	})

	defer rows.Close()

	return selected_meemz
}

func Sieve(r *http.Request, imgName string, arr []string) bool {
	for _, i := range arr {
		if i == imgName {
			return true
		}
	}

	return false
}

func FetchMeemz(rw http.ResponseWriter, r *http.Request) {
	uid := authentication.FetchId(r)

	viewed := FetchViewedMeemz(r, "meem")
	subscription := Subscription(r, "posts", "imgName")
	regommended_meemz := FetchRegommendedMeemz(r, "meem")
	selected_meemz := FetchAllMeemz(r, "posts", "imgName")

	MainContent = []*Posts{}

	var sorted []string

	for _, i := range subscription {
		viewed := Sieve(r, i, viewed)
		if !viewed {
			sorted = append(sorted, i)
		}
	}

	for _, i := range regommended_meemz {
		var access string
		a := db.QueryRow("SELECT access FROM posts WHERE imgName=?", i)
		a.Scan(&access)

		isSorted := Sieve(r, i, sorted)

		if !isSorted && access == "Public" {
			sorted = append(sorted, i)
		}
	}

	for _, i := range selected_meemz {
		viewed := Sieve(r, i, viewed)
		isSorted := Sieve(r, i, sorted)
		if !viewed && !isSorted {
			sorted = append(sorted, i)
		}
	}

	for _, s := range sorted {
		rows, err := db.Query("SELECT users.username, users.profileImg, posts.imgName, posts.labelOcr, posts.logoOcr, posts.faceOcr, posts.landmarkOcr, posts.textOcr, posts.safeSearchOcr, posts.possibleDuplicate, posts.tags, posts.pComment, posts.uploadTime FROM posts INNER JOIN users ON users.userId = posts.userId WHERE imgName=? ORDER BY posts.id DESC LIMIT 5", s)
		if err != nil {
			log.Println(err)
		}

		ExecMeemzQuery(rows, uid)
	}

	fmt.Println(MainContent, len(MainContent))

	json.NewEncoder(rw).Encode(MainContent[0])
}

func FetchVeemz(rw http.ResponseWriter, r *http.Request) {
	uid := authentication.FetchId(r)

	viewed := FetchViewedMeemz(r, "veem")
	subscription := Subscription(r, "veemz", "vidName")
	regommended_veemz := FetchRegommendedMeemz(r, "veem")
	selected_veemz := FetchAllMeemz(r, "veemz", "vidName")

	fmt.Println(viewed)
	fmt.Println(subscription)
	fmt.Println(regommended_veemz)
	fmt.Println(selected_veemz)

	VeemzMainContent = []*Veemz{}

	var sorted []string

	for _, i := range subscription {
		viewed := Sieve(r, i, viewed)
		if !viewed {
			sorted = append(sorted, i)
		}
	}

	for _, i := range regommended_veemz {
		var access string
		a := db.QueryRow("SELECT access FROM veemz WHERE vidName=?", i)
		a.Scan(&access)

		isSorted := Sieve(r, i, sorted)

		if !isSorted && access == "Public" {
			sorted = append(sorted, i)
		}
	}

	for _, i := range selected_veemz {
		viewed := Sieve(r, i, viewed)
		isSorted := Sieve(r, i, sorted)
		if !viewed && !isSorted {
			sorted = append(sorted, i)
		}
	}

	for _, s := range sorted {
		rows, err := db.Query("SELECT users.username, users.profileImg, veemz.vidName, veemz.tags, veemz.pComment, veemz.uploadTime FROM veemz INNER JOIN users ON users.userId = veemz.userId WHERE vidName=? ORDER BY veemz.id DESC LIMIT 5", s)
		if err != nil {
			log.Println(err)
		}

		ExecVeemzQuery(rows, uid)
	}

	fmt.Println(VeemzMainContent, len(VeemzMainContent))

	json.NewEncoder(rw).Encode(VeemzMainContent[0])
}

func ExecMeemzQuery(rows *sql.Rows, uid string) {
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

		MainContent = append(MainContent, &Posts{Username: username, ProfileImg: profileImg, ImgName: imgName, LabelOcr: labelOcr, LogoOcr: logoOcr, FaceOcr: faceOcr, LandmarkOcr: landmarkOcr, TextOcr: textOcr, SafeSearchOcr: safeSearchOcr, PossibleDuplicate: possibleDuplicate, Tags: tags, PComment: pComment, UploadTime: TimeStamp(uploadTime), Reaction1: reaction1, Reaction2: reaction2, Reaction3: reaction3, Reaction4: reaction4, Reaction5: reaction5})
	}

	defer rows.Close()
}

func ExecVeemzQuery(rows *sql.Rows, uid string) {
	for rows.Next() {
		var username string
		var profileImg string

		var vidName string

		var tags string
		var pComment string
		var uploadTime string

		var reaction string
		var reaction1 string = "far fa-grin-tears"
		var reaction2 string = "far fa-grin-tongue-squint"
		var reaction3 string = "far fa-meh"
		var reaction4 string = "far fa-sad-tear"
		var reaction5 string = "far fa-angry"

		if err := rows.Scan(&username, &profileImg, &vidName, &tags, &pComment, &uploadTime); err != nil {
			log.Println(err)
		}

		row := db.QueryRow("SELECT reactionType FROM reactions WHERE userId=? AND imgName=?", uid, vidName)
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

		VeemzMainContent = append(VeemzMainContent, &Veemz{Username: username, ProfileImg: profileImg, ImgName: vidName, Tags: tags, PComment: pComment, UploadTime: TimeStamp(uploadTime), Reaction1: reaction1, Reaction2: reaction2, Reaction3: reaction3, Reaction4: reaction4, Reaction5: reaction5})
	}

	defer rows.Close()
}
