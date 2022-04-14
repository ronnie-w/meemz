package regommend

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

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

type Meemz struct {
	ImgName interface{}
	PC      float64
}

var db = database.Conn()

var MainContent []*Posts

func FetchRegommendedMeemz(r *http.Request) []string {
	uid := authentication.FetchId(r)
	meemz := regommend.Table("meemz")

	viewed_by_me := make(map[interface{}]float64)
	viewed_by_others := make(map[interface{}]float64)

	vbm_rows, _ := db.Query("SELECT imgName, pc FROM regommend WHERE userId=?", uid)
	vbo_rows, _ := db.Query("SELECT imgName, pc FROM regommend WHERE userId!=?", uid)

	for vbm_rows.Next() {
		var imgName string
		var pc int

		vbm_rows.Scan(&imgName, &pc)
		viewed_by_me[imgName] = float64(pc)
	}
	meemz.Add("viewed_by_me", viewed_by_me)

	for vbo_rows.Next() {
		var imgName string
		var pc int

		vbo_rows.Scan(&imgName, &pc)
		viewed_by_others[imgName] = float64(pc)
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

func Subscription(r *http.Request) []string {
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
		subscriptions, _ := db.Query("SELECT imgName FROM posts WHERE userId=? ORDER BY id DESC LIMIT 25", c)
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

func FetchViewedMeemz(r *http.Request) []string {
	uid := authentication.FetchId(r)

	var viewed []string

	v_rows, err := db.Query("SELECT imgName FROM regommend WHERE userId=?", uid)
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

func FetchAllMeemz(r *http.Request) []string {
	var selected_meemz []string

	rows, err := db.Query("SELECT imgName FROM posts WHERE access=? ORDER BY posts.id DESC LIMIT 250", "Public")
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

	viewed := FetchViewedMeemz(r)
	subscription := Subscription(r)
	regommended_meemz := FetchRegommendedMeemz(r)
	selected_meemz := FetchAllMeemz(r)

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
		rows, err := db.Query("SELECT users.username, users.profileImg, posts.imgName, posts.labelOcr, posts.logoOcr, posts.faceOcr, posts.landmarkOcr, posts.textOcr, posts.safeSearchOcr, posts.possibleDuplicate, posts.tags, posts.pComment, posts.uploadTime FROM posts INNER JOIN users ON users.userId = posts.userId WHERE imgName=? ORDER BY posts.id DESC", s)
		if err != nil {
			log.Println(err)
		}

		ExecQuery(rows, uid)
	}

	fmt.Println(MainContent, len(MainContent))

	json.NewEncoder(rw).Encode(MainContent)
}

func ExecQuery(rows *sql.Rows, uid string) {
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
		var reaction1 string = "fal fa-grin-tears"
		var reaction2 string = "fal fa-grin-tongue-squint"
		var reaction3 string = "fal fa-meh"
		var reaction4 string = "fal fa-sad-tear"
		var reaction5 string = "fal fa-angry"

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

		MainContent = append(MainContent, &Posts{Username: username, ProfileImg: profileImg, ImgName: imgName, LabelOcr: labelOcr, LogoOcr: logoOcr, FaceOcr: faceOcr, LandmarkOcr: landmarkOcr, TextOcr: textOcr, SafeSearchOcr: safeSearchOcr, PossibleDuplicate: possibleDuplicate, Tags: tags, PComment: pComment, UploadTime: uploadTime, Reaction1: reaction1, Reaction2: reaction2, Reaction3: reaction3, Reaction4: reaction4, Reaction5: reaction5})
	}

	defer rows.Close()
}
