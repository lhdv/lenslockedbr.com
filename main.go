package main

import (
	"fmt"
	"net/http"

	"lenslockedbr.com/controllers"

	"github.com/gorilla/mux"
)

func notfound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","text/html")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h2>WHOW! Page Not Found - 404!</h2>")
}

func main() {
	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers()
	galleryC := controllers.NewGalleries()

	r := mux.NewRouter()

	r.NotFoundHandler = http.HandlerFunc(notfound)

	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.Handle("/faq", staticC.Faq).Methods("GET")

	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")

	r.HandleFunc("/galleries/new", galleryC.New).Methods("GET")

	http.ListenAndServe(":3000", r)
}

