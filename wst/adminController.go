package wst

import (
	"html/template"
	"log"
	"net/http"
)

type User struct {
	UserName string
}

type adminController struct {
	path string
}

func (this *adminController) IndexAction(w http.ResponseWriter, r *http.Request, user string) {
	t, err := template.ParseFiles(this.path + "/html/admin/index.html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, &User{user})
	log.Println("adminController-->IndexAction")
}
