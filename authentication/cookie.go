package authentication

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/schema"
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
	base64enc, err := cookieHandler.Encode("R0nni3W3k35@", value)
	if err != nil {
		log.Fatalln("Error encoding cookie value")
	}

	expiry := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{
		Name:     name,
		Value:    base64enc,
		Expires:  expiry,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(rw, &cookie)
}

func CookieForm(r *http.Request) *Authentication {
	r.ParseForm()
	key := new(Authentication)
	schema.NewDecoder().Decode(key, r.PostForm)

	return key
}

func ReadCookie(r *http.Request) string {
	cookie, cookie_err := r.Cookie("uid")

	var cookieValue string
	if cookie_err == nil {
		if err := cookieHandler.Decode("R0nni3W3k35@", cookie.Value, &cookieValue); err != nil {
			log.Println("Error decoding cookie")
		}
	}

	return cookieValue
}

func DecodeUID_enc(cookie string) string {
	var cookieValue string
	if cookie != "" {
		if err := cookieHandler.Decode("R0nni3W3k35@", cookie, &cookieValue); err != nil {
			log.Println("Error decoding cookie")
		}
	}

	return cookieValue
}

func DeleteCookie(r *http.Request) {
	cookie, err := r.Cookie("uid")
	if err != nil {
		log.Println("Error reading cookie")
	}

	cookie.Expires = time.Now()
}
