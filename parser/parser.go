package parser

import (
	"html/template"
	"log"
	"net/http"
)

func Parser(temp_path string) func(http.ResponseWriter, *http.Request) {
	parser := func (rw http.ResponseWriter) *template.Template {		
		parseTemplate, _ := template.ParseFiles("templates/"+temp_path+".html")
		err := parseTemplate.Execute(rw, nil)
		if err != nil {
			log.Println(err)
		}

		return parseTemplate
	}

	template := func (rw http.ResponseWriter, r *http.Request)  {
		parser(rw)
		rw.Header().Set("Cache-Control", "max-age=604800")
	}

	return template
}

func Manifest(rw http.ResponseWriter, r *http.Request) {
	http.ServeFile(rw, r, "manifest.json")
}

func ServiceWorkerInit(rw http.ResponseWriter, r *http.Request) {
	http.ServeFile(rw, r, "sw_init.js")
}

func ServiceWorker(rw http.ResponseWriter, r *http.Request) {
	http.ServeFile(rw, r, "sw.js")
}

func AssetsLink(rw http.ResponseWriter, r *http.Request) {
	http.ServeFile(rw, r, ".well-known/assetslinks.json")
}