package main

import (
	"fmt"
	"net/http"

	"lenslockedbr.com/views"
	"github.com/gorilla/mux"
)

var homeView *views.View
var contactView *views.View
var faqView *views.View

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(homeView.Render(w, nil))
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(contactView.Render(w, nil))
}

func faq(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(faqView.Render(w, nil))
}

func notfound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","text/html")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h2>WHOW! Page Not Found - 404!</h2>")
}

// A helper function that panics on any error
func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	homeView = views.NewView("bootstrap", "views/home.gohtml")
	contactView = views.NewView("bootstrap", "views/contact.gohtml")
	faqView = views.NewView("bootstrap_bggray", "views/faq.gohtml")

	r := mux.NewRouter()

	r.NotFoundHandler = http.HandlerFunc(notfound)

	r.HandleFunc("/", home)
	r.HandleFunc("/contact", contact)
	r.HandleFunc("/faq", faq)
	http.ListenAndServe(":3000", r)
}

