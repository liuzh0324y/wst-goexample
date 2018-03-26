package wst 

import (
	"log"
	"net/http"
	"html/template"
	"strings"
	"reflect"
)

// Wst struct
type Wst struct {
	test int
}

// New Wst 
func New() *Wst {
	obj := &Wst{
		test : 0,
	}
	return obj
}

// Run attaches the router to a http.Server and starts listening and serving HTTP requests.
func (wst *Wst) Run() {

	http.Handle("/static/html/", http.FileServer(http.Dir("static/html")))
	http.Handle("/static/js/", http.FileServer(http.Dir("static/js")))
	http.Handle("/static/css/", http.FileServer(http.Dir("static/css")))
	http.Handle("/static/", http.FileServer(http.Dir("static")))
	
	http.HandleFunc("/index/", indexHandler)
	http.HandleFunc("/admin/", adminHandler)
	http.HandleFunc("/login/", loginHandler)
	http.HandleFunc("/ajax/", ajaxHandler)
	http.HandleFunc("/webrtc/", webrtcHandler)
	http.HandleFunc("/", notFoundHandler)

	log.Fatal(http.ListenAndServe(":8090",nil))
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/index/", http.StatusFound)
		// t, err := template.ParseFiles("static/html/index.html")
		// if err != nil {
		// 	log.Println(err)
		// }
		// t.Execute(w, nil)
	}

	t, err := template.ParseFiles("static/html/404.html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, nil)
	log.Println("not found handler")
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	var action = ""
	index := &indexController{}
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

func loginHandler(w http.ResponseWriter, r *http.Request) {
	pathInfo := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(pathInfo, "/")
	var action = ""
	if len(parts) > 1 {
		action = strings.Title(parts[1]) + "Action"
	}

	login := &loginController{}
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

func adminHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("admin_name")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login/index", http.StatusFound)
	}

	pathInfo := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(pathInfo, "/")
	var action = ""
	if len(parts) > 1 {
		action = strings.Title(parts[1])+"Action" 
	}

	admin := &adminController{}
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

func ajaxHandler(w http.ResponseWriter, r *http.Request) {
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

func webrtcHandler(w http.ResponseWriter, r *http.Request) {
	pathInfo := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(pathInfo, "/")
	var action = ""
	if len(parts) > 1 {
		action = strings.Title(parts[1]) + "Action"
	}
	log.Println("webrtc action: " + action)
	webrtc := &webrtcController{}
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