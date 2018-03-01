package wst 

import (
	"log"
	"net/http"
	"io"
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
	http.HandleFunc("/", indexHandler)
	log.Fatal(http.ListenAndServe(":8080",nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("wst index handler\n"))
	io.WriteString(w, "URL" + r.URL.String())
}