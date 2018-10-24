package main

import (
	"fmt"
	"net/http"

	"lenslockedbr.com/controllers"
	"lenslockedbr.com/views"

	"github.com/gorilla/mux"
)

var faqView *views.View

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
	faqView = views.NewView("bootstrap_bggray", "views/faq.gohtml")

	usersC := controllers.NewUsers()
	staticC := controllers.NewStatic()

	r := mux.NewRouter()

	r.NotFoundHandler = http.HandlerFunc(notfound)

	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.HandleFunc("/faq", faq)
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	http.ListenAndServe(":3000", r)
}

