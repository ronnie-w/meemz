package authentication

import (
	"crypto/rand"
	"encoding/json"
	"log"
	"net/http"

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
	uid, _ := rand.Prime(rand.Reader, 50)

	exists := QuickCheck(user.Username, user.Email)

	if exists == 0 && user.Username != "" {
		insert_row, _ := db.Query("INSERT INTO users(username,userId,email,password,serverCode) values(?,?,?,?,?)", user.Username, uid.String(), user.Email, password, random.String())

		SendMail(user.Email, "Account Verification", "Your verification code is "+random.String()+"", `
		<html>
		  <head>
		  <meta http–equiv=“Content-Type” content=“text/html; charset=UTF-8” />
		  <meta http–equiv=“X-UA-Compatible” content=“IE=edge” />
		  <meta name=“viewport” content=“width=device-width, initial-scale=1.0 “ />
		  </head>
		  <body style="background-color : black;">
		  <center>
			  <h1 style="color : white;font-family: ‘Palatino Linotype’, ‘Book Antiqua’, Palatino, serif;">Meemz</h1>
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
