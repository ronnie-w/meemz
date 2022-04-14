package authentication

import (
	"net/http"

	"github.com/gorilla/schema"
)

type User struct {
	Username string
	Email    string
	Password string
	VCode    string
	ImgOcr   string
	Tags     string
}

func FormReader(r *http.Request) *User {
	r.ParseForm()
	user := new(User)
	schema.NewDecoder().Decode(user, r.PostForm)

	return user
}
