package main

import (
	"fmt"
	"log"
	"net/http"

	"lenslockedbr.com/controllers"
	"lenslockedbr.com/models"

	"github.com/gorilla/mux"
)

const (
	host     = "192.168.56.101"
	port     = 5432
	user     = "developer"
	password = "1234qwer"
	dbname   = "lenslockedbr_dev"
)

func main() {
	// Create a DB connection string and then use it to create
	// our model services.
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s " +  
                                "dbname=%s sslmode=disable",
		                 host, port, user, password, dbname)

	services, err := models.NewServices(psqlInfo)
	if err != nil {
		panic(err)
	}

	defer services.Close()
	services.AutoMigrate()

	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(services.User)
	galleriesC := controllers.NewGalleries(services.Gallery)

	r := mux.NewRouter()

	r.NotFoundHandler = http.HandlerFunc(staticC.PageNotFound.ServeHTTP)

	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.Handle("/faq", staticC.Faq).Methods("GET")

	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")

	// Gallery routes
	r.HandleFunc("/galleries/new", galleriesC.New).Methods("GET")
	r.HandleFunc("/galleries", galleriesC.Create).Methods("POST")

	r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")

	log.Println("Starting the server on :3000...")

	http.ListenAndServe(":3000", r)
}

