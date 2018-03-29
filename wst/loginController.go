package wst

import (
	"html/template"
	"log"
	"net/http"
)

type loginController struct {
	path string
}

func (this *loginController) IndexAction(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(this.path + "/html/login/index.html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, nil)
	log.Println("loginController-->IndexAction")
}
