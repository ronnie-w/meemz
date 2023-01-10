package authentication

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/securecookie"
)

type OnCompletion struct {
	Done string
}

type SecureCookie struct {
	Name     string
	Value    string
	Expires  string
	SameSite string
}

type Authentication struct {
	AuthKey string
}

var cookieHandler *securecookie.SecureCookie

func init() {
	gen64 := securecookie.GenerateRandomKey(64)
	gen32 := securecookie.GenerateRandomKey(32)
	cookieHandler = securecookie.New(gen64, gen32)
}

func CreateCookie(rw http.ResponseWriter, r *http.Request, value string, name string) {
	base64enc, err := cookieHandler.Encode("R0nni3W3k35@M3em2", value)
	if err != nil {
		log.Fatalln("Error encoding cookie value")
	}

	expiry := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{
		Name:       name,
		Value:      base64enc,
		Expires:    expiry,
		SameSite:   http.SameSiteLaxMode,
	}

	http.SetCookie(rw, &cookie)
}

func ReadCookie(r *http.Request) string {
	cookie, cookie_err := r.Cookie("uid")

	var cookieValue string
	if cookie_err == nil {
		if err := cookieHandler.Decode("R0nni3W3k35@M3em2", cookie.Value, &cookieValue); err != nil {
			log.Println("Error decoding cookie")
		}
	}

	return cookieValue
}

func ReadCustomCookie(r *http.Request, val string) string {
	cookie, cookie_err := r.Cookie(val)

	var cookieValue string
	if cookie_err == nil {
		if err := cookieHandler.Decode("R0nni3W3k35@M3em2", cookie.Value, &cookieValue); err != nil {
			log.Println("Error decoding cookie")
		}
	}

	return cookieValue
}

func Logout(rw http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("uid")
	if err != nil {
		log.Println("Error reading cookie")
	}

	cookie.Expires = time.Now()
}
