package authentication

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/schema"
	"golang.org/x/crypto/bcrypt"
)

type UserPassChange struct {
	Username    string
	Email       string
	PassCode    string
	NewPassword string
	AuthKey     string
}

func UserPassForm(r *http.Request) *UserPassChange {
	r.ParseForm()
	user := new(UserPassChange)
	schema.NewDecoder().Decode(user, r.PostForm)

	return user
}

func SendVerificationCode(rw http.ResponseWriter, r *http.Request) {
	username := UserPassForm(r).Username

	user_row, _ := db.Query("SELECT email FROM users WHERE username=?", username)
	var exists int
	var email string
	for user_row.Next() {
		exists++
		user_row.Scan(&email)
	}

	random_code, _ := rand.Prime(rand.Reader, 20)
	if exists == 1 {
		update_row, _ := db.Query("UPDATE users SET passCode=? WHERE username=?", random_code.String(), username)
		defer update_row.Close()
		SendMail(email, "Password Reset", "Your verification code is "+random_code.String()+"", `
		<html>
		  <head>
		  <meta http–equiv=“Content-Type” content=“text/html; charset=UTF-8” />
		  <meta http–equiv=“X-UA-Compatible” content=“IE=edge” />
		  <meta name=“viewport” content=“width=device-width, initial-scale=1.0 “ />
		  </head>
		  <body style="background-color : black; border-radius: 6px;">
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
			  <p style="color : white;font-family: Courier, monospace;" id="willkommen">Looks like you forgot your password. Ignore this email if it wasn't you.</p>
			  <p style="color : white;font-family: Courier, monospace;">Copy the verification code below and paste it in the field.</p>
			  <p style="color : white;font-family: Courier, monospace;"><strong id="v_code">`+random_code.String()+`</strong></p>
		  </center>
		  </body>
		</html>`)

		CreateCookie(rw, r, username, "uid")

		json.NewEncoder(rw).Encode(UserPassChange{Email: email, AuthKey: username})
	} else {
		json.NewEncoder(rw).Encode(Error{Err: "User not found"})
	}
}

func ConfirmVerificationCode(rw http.ResponseWriter, r *http.Request) {
	username := ReadCookie(r)
	user_code := UserPassForm(r).PassCode

	code_row := db.QueryRow("SELECT passCode FROM users WHERE username=?", username)
	var pass_code string
	code_row.Scan(&pass_code)

	if pass_code != user_code {
		json.NewEncoder(rw).Encode(Error{Err: "The code you entered is incorrect"})
	} else {
		CreateCookie(rw, r, pass_code, "auth_key")
	}
}

func PasswordReset(rw http.ResponseWriter, r *http.Request) {
	username := ReadCookie(r)
	auth_key := ReadCustomCookie(r, "auth_key")

	fmt.Println(username, auth_key)

	code_row := db.QueryRow("SELECT passCode FROM users WHERE username=?", username)
	var pass_code string
	code_row.Scan(&pass_code)

	if pass_code == auth_key {
		new_pass, _ := bcrypt.GenerateFromPassword([]byte(UserPassForm(r).NewPassword), 10)

		update_pass, _ := db.Query("UPDATE users SET password=? WHERE username=?", new_pass, username)
		defer update_pass.Close()
	} else {
		json.NewEncoder(rw).Encode(Error{"Action unauthorized"})
	}
}
