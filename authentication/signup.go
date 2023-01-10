package authentication

import (
	"crypto/rand"
	"encoding/json"
	"log"
	"net/http"
	"net/mail"
	"time"

	"github.com/google/uuid"
	"github.com/meemz/database"
	"golang.org/x/crypto/bcrypt"
)

type Error struct {
	Err string
}

var db = database.Conn()

func Signup(rw http.ResponseWriter, r *http.Request) {
	user := FormReader(r)
	password, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

	random, _ := rand.Prime(rand.Reader, 20)
	uid := uuid.New()

	exists := QuickCheck(user.Username, user.Email)

	_, invalid_email := mail.ParseAddress(user.Email)

	joinDate := time.Now().Format(time.RFC822Z)

	if len(user.Username) < 3 {
		json.NewEncoder(rw).Encode(Error{"Username characters must be three or more"})
	} else if len(user.Password) < 4 {
		json.NewEncoder(rw).Encode(Error{"Password characters must be four or more"})
	} else if invalid_email != nil {
		json.NewEncoder(rw).Encode(Error{"The email you entered is invalid"})
	} else if exists == 0 {
		insert_row, _ := db.Query("INSERT INTO users(username,userId,email,password,serverCode,joinDate) values(?,?,?,?,?,?)", user.Username, uid.String(), user.Email, password, random.String(), joinDate)

		SendMail(user.Email, "Account Verification", "Your verification code is "+random.String()+"", `
		<html>
		  <head>
		  <meta http–equiv=“Content-Type” content=“text/html; charset=UTF-8” />
		  <meta http–equiv=“X-UA-Compatible” content=“IE=edge” />
		  <meta name=“viewport” content=“width=device-width, initial-scale=1.0 “ />
		  </head>
		  <body style="background-color: black; border-radius: 6px;">
		  <style>
			@font-face {
				font-family: 'Cabin Sketch';
				font-style: normal;
				font-weight: 700;
				src: url(https://fonts.gstatic.com/s/cabinsketch/v19/QGY2z_kZZAGCONcK2A4bGOj0I_1Y5tjz.woff2) format('woff2');
				unicode-range: U+0000-00FF, U+0131, U+0152-0153, U+02BB-02BC, U+02C6, U+02DA, U+02DC, U+2000-206F, U+2074, U+20AC, U+2122, U+2191, U+2193, U+2212, U+2215, U+FEFF, U+FFFD;
			}
		  </style>
		  <center>
			  <h1 style="color : white;font-family: 'Cabin Sketch', cursive;">Meemz</h1>
			  <p style="color : white;font-family: Courier, monospace;" id="willkommen">Welcome to Meemz 🥳. It's all about memes.</p>
			  <p style="color : white;font-family: Courier, monospace;">Copy the verification code below and paste it in the field.</p>
			  <p style="color : white;font-family: Courier, monospace;"><strong id="v_code">`+random.String()+`</strong></p>
		  </center>
		  </body>
		</html>`)
		CreateCookie(rw, r, uid.String(), "uid")

		defer insert_row.Close()
	} else {
		json.NewEncoder(rw).Encode(Error{"User already exists"})
	}

}

func DeleteUser(rw http.ResponseWriter, r *http.Request) {
	cookieVal := ReadCookie(r)
	_, err := db.Query("DELETE FROM users WHERE userId=?", cookieVal)
	if err != nil {
		log.Fatalln(err)
	}
}
