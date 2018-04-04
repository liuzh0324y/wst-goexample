package wst

import (
	"html/template"
	"log"
	"net/http"
)

type liveController struct {
	path string
}

func (this *liveController) IndexAction(w http.ResponseWriter, h *http.Request) {
	t, err := template.ParseFiles(this.path + "/html/live/index.html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, nil)
}
