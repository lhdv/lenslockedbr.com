package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func home(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<h1>Welcome to my awesome site!</h1>")
}

func contact(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "To get in touch, please send an email " +
                      "to <a href=\"mailto:support@lenslockedbr.com\">"+
                      "support@lenslockedbr.com</a>.")
}

func faq(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
	r := httprouter.New()

	r.NotFound = http.HandlerFunc(notfound)

	r.GET("/", home)
	r.GET("/contact", contact)
	r.GET("/faq", faq)
	http.ListenAndServe(":3000", r)
}
