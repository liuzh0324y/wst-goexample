package wst 

import (
	"log"
	"net/http"
	"encoding/json"
)

type Result struct {
	Ret int
	Reason string
	Data interface{}
}

type ajaxController struct {}

func (this *ajaxController)LoginAction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	err := r.ParseForm()
	if err != nil {
		OutputJson(w, 0, "parameter error", nil)
		return
	}

	admin_name := r.FormValue("admin_name")
	admin_passwd := r.FormValue("admin_passwd")
	if admin_name == "" || admin_passwd == "" {
		OutputJson(w, 0, "user name or password is nuil!", nil)
		return
	}

	cookie := http.Cookie{Name: "admin_name", Value: admin_name, Path: "/"}
	http.SetCookie(w, &cookie)
	OutputJson(w, 1, "success", nil)
	log.Println("ajaxController-->LoginAction")
	log.Printf("%s,%s", admin_name, admin_passwd)
}

func OutputJson(w http.ResponseWriter, ret int, reason string, i interface{}) {
	out := &Result{ret, reason, i}
	b, err := json.Marshal(out)
	if err != nil {
		return
	}

	w.Write(b)
}