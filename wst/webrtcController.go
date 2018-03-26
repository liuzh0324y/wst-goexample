package wst 

import (
	"log"
	"net/http"
	"html/template"
)


type webrtcController struct {

}

func (this *webrtcController)IndexAction(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("static/html/webrtc/index.html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, nil)
	log.Println("webrtcController-->IndexAction")
}

func (this *webrtcController)GetUserMediaAction(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("static/html/webrtc/getusermedia/index.html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, nil)
	log.Println("webrtcController-->getUserMediaAction")
}

func (this *webrtcController)JsAction(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("static/js/webrtc/getUserMedia/getUserMedia.js")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, nil)
	log.Println("webrtcController-->JsAction")
}