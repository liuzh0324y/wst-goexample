package wst 

import (
	"log"
	"net/http"
	"html/template"
)

type loginController struct {
	
}

func (this *loginController)IndexAction(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("static/html/login/index.html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, nil)
	log.Println("loginController-->IndexAction")
}