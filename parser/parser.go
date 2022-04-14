package parser

import (
	"html/template"
	"log"
	"net/http"
)

func Parser(rw http.ResponseWriter, temp_path string) {
	parseTemplate, _ := template.ParseFiles(temp_path)
	err := parseTemplate.Execute(rw, nil)
	if err != nil {
		log.Println(err)
	}
}

func BottomNav(rw http.ResponseWriter, r *http.Request) {
	Parser(rw, "templates/bottomNav.html")
}

func Layout(rw http.ResponseWriter, r *http.Request) {
	Parser(rw, "templates/layout.html")
}

func Login(rw http.ResponseWriter, r *http.Request) {
	Parser(rw, "templates/login.html")
}

func Signup(rw http.ResponseWriter, r *http.Request) {
	Parser(rw, "templates/signup.html")
}

func Verify(rw http.ResponseWriter, r *http.Request) {
	Parser(rw, "templates/verify.html")
}

func Terms(rw http.ResponseWriter, r *http.Request) {
	Parser(rw, "templates/terms.html")
}

func Home(rw http.ResponseWriter, r *http.Request) {
	Parser(rw, "templates/home.html")
}

func Profile(rw http.ResponseWriter, r *http.Request) {
	Parser(rw, "templates/profile.html")
}

func PrivatePosts(rw http.ResponseWriter, r *http.Request) {
	Parser(rw, "templates/privatePosts.html")
}

func ProfileEdit(rw http.ResponseWriter, r *http.Request) {
	Parser(rw, "templates/profileEdit.html")
}

func ImageStats(rw http.ResponseWriter, r *http.Request) {
	Parser(rw, "templates/imageStats.html")
}

func PublicStats(rw http.ResponseWriter, r *http.Request) {
	Parser(rw, "templates/publicStats.html")
}

func Convo(rw http.ResponseWriter, r *http.Request) {
	Parser(rw, "templates/convo.html")
}

func ConvoInit(rw http.ResponseWriter, r *http.Request) {
	Parser(rw, "templates/convoInit.html")
}

func ConvoCreate(rw http.ResponseWriter, r *http.Request) {
	Parser(rw, "templates/convoCreate.html")
}

func Create(rw http.ResponseWriter, r *http.Request) {
	Parser(rw, "templates/create.html")
}

func Creator(rw http.ResponseWriter, r *http.Request) {
	Parser(rw, "templates/creator.html")
}

func Account(rw http.ResponseWriter, r *http.Request) {
	Parser(rw, "templates/account.html")
}

func Search(rw http.ResponseWriter, r *http.Request) {
	Parser(rw, "templates/search.html")
}

func Notifications(rw http.ResponseWriter, r *http.Request) {
	Parser(rw, "templates/notifications.html")
}

func ForgotPassword(rw http.ResponseWriter, r *http.Request) {
	Parser(rw, "templates/forgotPassword.html")
}
