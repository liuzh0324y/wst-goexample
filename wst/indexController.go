package wst

import (
	"html/template"
	"log"
	"net/http"
)

type indexController struct {
	path string
}

func (this *indexController) IndexAction(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(this.path + "/html/index.html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, nil)
	log.Println("indexController-->IndexController")
}
