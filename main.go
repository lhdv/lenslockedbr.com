package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

var homeTemplate *template.Template

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := homeTemplate.Execute(w, nil); err != nil {
		panic(err)
	}
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "To get in touch, please send an email " +
                      "to <a href=\"mailto:support@lenslockedbr.com\">"+
                      "support@lenslockedbr.com</a>.")
}

func faq(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<h1>F.A.Q.</h1>" +
                      "<p>You'll find here some answers which could be " +
		      "helpful for you.</p>" +
		      "<p><ol><li>Foobar</li>" +
                      "<li>XPTO</li>" +
                      "<li>Foo</li>" +
		      "<li>Bar</li></ol></p>" +
                      "<p><a href=\"/\">Home</a></p>")
}

func notfound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","text/html")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h2>WHOW! Page Not Found - 404!</h2>")
}

func main() {
	var err error

	homeTemplate, err = template.ParseFiles("views/home.gohtml")
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()

	r.NotFoundHandler = http.HandlerFunc(notfound)

	r.HandleFunc("/", home)
	r.HandleFunc("/contact", contact)
	r.HandleFunc("/faq", faq)
	http.ListenAndServe(":3000", r)
}
