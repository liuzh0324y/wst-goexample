package wst

import (
	"html/template"
	"log"
	"net/http"
)

type playbackController struct {
	path string
}

func (this *playbackController) IndexAction(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(this.path + "/html/playback/index.html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, nil)
}
