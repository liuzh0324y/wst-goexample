package wst

import(
	"log"
	"net/http"
	"html/template"
)

type indexController struct {

}

func (this *indexController)IndexAction(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("static/html/index.html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, nil)
	log.Println("indexController-->IndexController")
}