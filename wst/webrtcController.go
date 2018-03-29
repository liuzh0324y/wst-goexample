package wst

import (
	"html/template"
	"log"
	"net/http"
	"strings"
)

type webrtcController struct {
	path string
}

func (this *webrtcController) IndexAction(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(this.path + "/html/webrtc/index.html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, nil)
	log.Println("webrtcController-->IndexAction")
}

func (this *webrtcController) GetUserMediaAction(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(this.path + "/html/webrtc/getusermedia.html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, nil)
	log.Println("webrtcController-->getUserMediaAction")
}

func (this *webrtcController) GetUserMediaCanvasAction(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(this.path + "/html/webrtc/getusermediacanvas.html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, nil)
	log.Println("webrtcController-->getUserMediaCanvasAction")
}

func (this *webrtcController) VideoToVideoAction(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(this.path + "/html/webrtc/videotovideo.html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, nil)
	log.Println("webrtcController-->VideoToVideoAction")
}

func (this *webrtcController) JsAction(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(path, "/")
	if len(parts) > 2 {
		t, err := template.ParseFiles(this.path + "/js/webrtc/" + parts[2])
		if err != nil {
			log.Println(err)
		}
		t.Execute(w, nil)
	}

	log.Println("webrtcController-->JsAction")
}
