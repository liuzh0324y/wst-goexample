package wst

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"strings"
)

// Wst struct
type Wst struct {
	test     int
	confPath string
}

// New Wst
func New(path string) *Wst {
	obj := &Wst{
		test:     0,
		confPath: path,
	}
	return obj
}

// Run attaches the router to a http.Server and starts listening and serving HTTP requests.
func (wst *Wst) Run() {
	log.Println("wst server start.")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go wst.WstHttpServer()
	go wst.WstSignalServer()

	s := <-c
	log.Println("wst server quit.", s)
	// http.HandleFunc("/index/", wst.indexHandler)
	// http.HandleFunc("/admin/", wst.adminHandler)
	// http.HandleFunc("/login/", wst.loginHandler)
	// http.HandleFunc("/ajax/", wst.ajaxHandler)
	// http.HandleFunc("/webrtc/", wst.webrtcHandler)
	// http.HandleFunc("/live/", wst.liveHandler)
	// http.HandleFunc("/playback/", wst.playbackHandler)
	// http.HandleFunc("/", wst.notFoundHandler)

	// http.HandleFunc("/css/", wst.cssHandler)
	// http.HandleFunc("/js/", wst.jsHandler)
	// http.HandleFunc("/icon/", wst.iconHandler)

	// log.Fatal(http.ListenAndServeTLS(":8090", wst.confPath+"/key/cert.pem", wst.confPath+"/key/key.pem", nil))
}

func (wst *Wst) WstHttpServer() {
	httpserver := NewHttpServer()
	httpserver.Run()
	// h := http.FileServer(http.Dir(wst.confPath + "/html"))
	// http.Handle("/js/", http.FileServer(http.Dir(wst.confPath+"/js")))
	// http.Handle("/css/", http.FileServer(http.Dir(wst.confPath+"/css")))
	// err := http.ListenAndServeTLS(":8090", wst.confPath+"/key/cert.pem", wst.confPath+"/key/key.pem", h)
	// if err != nil {
	// 	log.Fatal("ListenAndServe: ", err)
	// }
}

func (wst *Wst) WstSignalServer() {
	wsserver := NewSignalServer()
	wsserver.Run()
}

func (wst *Wst) notFoundHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path == "/" {
		http.Redirect(w, r, "/index/", http.StatusFound)
	}

	t, err := template.ParseFiles(wst.confPath + "/html/404.html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, nil)
	log.Println("notFound handler")
}

func (wst *Wst) indexHandler(w http.ResponseWriter, r *http.Request) {
	var action = ""
	index := &indexController{path: wst.confPath}
	controller := reflect.ValueOf(index)
	method := controller.MethodByName(action)
	if !method.IsValid() {
		method = controller.MethodByName(strings.Title("index") + "Action")
	}
	requestValue := reflect.ValueOf(r)
	responseValue := reflect.ValueOf(w)
	method.Call([]reflect.Value{responseValue, requestValue})
	log.Println("index handler")
}

func (wst *Wst) loginHandler(w http.ResponseWriter, r *http.Request) {
	pathInfo := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(pathInfo, "/")
	var action = ""
	if len(parts) > 1 {
		action = strings.Title(parts[1]) + "Action"
	}

	login := &loginController{path: wst.confPath}
	controller := reflect.ValueOf(login)
	method := controller.MethodByName(action)
	if !method.IsValid() {
		method = controller.MethodByName(strings.Title("index") + "Action")
	}
	requestValue := reflect.ValueOf(r)
	responseValue := reflect.ValueOf(w)
	method.Call([]reflect.Value{responseValue, requestValue})
	log.Println("login handler")
}

func (wst *Wst) adminHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("admin_name")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login/index", http.StatusFound)
	}

	pathInfo := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(pathInfo, "/")
	var action = ""
	if len(parts) > 1 {
		action = strings.Title(parts[1]) + "Action"
	}

	admin := &adminController{path: wst.confPath}
	controller := reflect.ValueOf(admin)
	method := controller.MethodByName(action)
	if !method.IsValid() {
		method = controller.MethodByName(strings.Title("index") + "Action")
	}
	requestValue := reflect.ValueOf(r)
	responseValue := reflect.ValueOf(w)
	userValue := reflect.ValueOf(cookie.Value)
	method.Call([]reflect.Value{responseValue, requestValue, userValue})
	log.Println("admin handler")
}

func (wst *Wst) ajaxHandler(w http.ResponseWriter, r *http.Request) {
	pathInfo := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(pathInfo, "/")
	var action = ""
	if len(parts) > 1 {
		action = strings.Title(parts[1]) + "Action"
	}

	ajax := &ajaxController{}
	controller := reflect.ValueOf(ajax)
	method := controller.MethodByName(action)
	if !method.IsValid() {
		method = controller.MethodByName(strings.Title("index") + "Action")
	}
	requestValue := reflect.ValueOf(r)
	responseValue := reflect.ValueOf(w)
	method.Call([]reflect.Value{responseValue, requestValue})
	log.Println("ajax handler")
}

func (wst *Wst) webrtcHandler(w http.ResponseWriter, r *http.Request) {
	pathInfo := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(pathInfo, "/")
	var action = ""
	if len(parts) > 1 {
		action = strings.Title(parts[1]) + "Action"
	}
	log.Println("webrtc action: " + action)
	webrtc := &webrtcController{path: wst.confPath}
	controller := reflect.ValueOf(webrtc)
	method := controller.MethodByName(action)
	if !method.IsValid() {
		method = controller.MethodByName(strings.Title("index") + "Action")
	}
	requestValue := reflect.ValueOf(r)
	responseValue := reflect.ValueOf(w)
	method.Call([]reflect.Value{responseValue, requestValue})
	log.Println("webrtc handler")
}

func (wst *Wst) liveHandler(w http.ResponseWriter, r *http.Request) {
	pathInfo := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(pathInfo, "/")
	var action = ""
	if len(parts) > 1 {
		action = strings.Title(parts[1]) + "Action"
	}
	log.Println("live action: " + action)
	live := &liveController{path: wst.confPath}
	controller := reflect.ValueOf(live)
	method := controller.MethodByName(action)
	if !method.IsValid() {
		method = controller.MethodByName(strings.Title("index") + "Action")
	}
	requestValue := reflect.ValueOf(r)
	responseValue := reflect.ValueOf(w)
	method.Call([]reflect.Value{responseValue, requestValue})
	log.Println("live handler")
}

func (wst *Wst) playbackHandler(w http.ResponseWriter, r *http.Request) {
	pathInfo := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(pathInfo, "/")
	var action = ""
	if len(parts) > 1 {
		action = strings.Title(parts[1]) + "Action"
	}
	log.Println("playback action: " + action)
	playback := &playbackController{path: wst.confPath}
	controller := reflect.ValueOf(playback)
	method := controller.MethodByName(action)
	if !method.IsValid() {
		method = controller.MethodByName(strings.Title("index") + "Action")
	}
	requestValue := reflect.ValueOf(r)
	responseValue := reflect.ValueOf(w)
	method.Call([]reflect.Value{responseValue, requestValue})
	log.Println("playback handler")
}

func (wst *Wst) jsHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(path, "/")
	if len(parts) > 1 {
		t, err := template.ParseFiles(wst.confPath + "/js/" + parts[1])
		if err != nil {
			log.Println(err)
		}
		t.Execute(w, nil)
		log.Println("root js action: " + parts[1])
	}
}

func (wst *Wst) cssHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(path, "/")
	if len(parts) > 1 {
		t, err := template.ParseFiles(wst.confPath + "/css/" + parts[1])
		if err != nil {
			log.Println(err)
		}
		t.Execute(w, nil)
		log.Println("root css action: " + parts[1])
	}
}

func (wst *Wst) iconHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(path, "/")
	if len(parts) > 1 {
		t, err := template.ParseFiles(wst.confPath + "/icon/" + parts[1])
		if err != nil {
			log.Println(err)
		}
		t.Execute(w, nil)
		log.Println("root icon action: " + parts[1])
	}
}
