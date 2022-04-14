package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/meemz/router"
)

func main() {
	methods := handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "HEAD"})
	origin := handlers.AllowedOrigins([]string{"https://meemzchat.cf"})

	logFile, err := os.OpenFile("server.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println(err)
	}

	http.ListenAndServe(":3000", handlers.LoggingHandler(logFile, handlers.CORS(methods, origin)(router.Routes())))
}
