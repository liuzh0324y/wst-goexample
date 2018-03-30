package wst

import (
	"html/template"
	"log"
	"net/http"
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

	// http.Handle("/static/html/", http.FileServer(http.Dir("static/html")))
	// http.Handle("/static/js/", http.FileServer(http.Dir("static/js")))
	// http.Handle("/static/css/", http.FileServer(http.Dir("static/css")))
	// http.Handle("/static/", http.FileServer(http.Dir("static")))

	http.HandleFunc("/index/", wst.indexHandler)
	http.HandleFunc("/admin/", wst.adminHandler)
	http.HandleFunc("/login/", wst.loginHandler)
	http.HandleFunc("/ajax/", wst.ajaxHandler)
	http.HandleFunc("/webrtc/", wst.webrtcHandler)
	http.HandleFunc("/", wst.notFoundHandler)

	// log.Fatal(http.ListenAndServe(":8090", nil))
	log.Fatal(http.ListenAndServeTLS(":8090", wst.confPath+"/key/cert.pem", wst.confPath+"/key/key.pem", nil))
}

func (wst *Wst) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
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
