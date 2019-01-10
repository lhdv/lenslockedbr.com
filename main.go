package main

import (
	"fmt"
	"log"
	"net/http"

	"lenslockedbr.com/controllers"
	"lenslockedbr.com/middleware"
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

	r := mux.NewRouter()

	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(services.User)
	galleriesC := controllers.NewGalleries(services.Gallery, 
                                               services.Image, r)

	userMw := middleware.User {
		UserService: services.User,
	}
	requireUserMw := middleware.RequireUser{ }

	r.NotFoundHandler = http.HandlerFunc(staticC.PageNotFound.ServeHTTP)

	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.Handle("/faq", staticC.Faq).Methods("GET")

	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")

	//
	// Gallery routes
	//

	r.Handle("/galleries", 
		requireUserMw.ApplyFn(galleriesC.Index)).Methods("GET").
		Name(controllers.IndexGallery)

	r.Handle("/galleries/new", 
		requireUserMw.Apply(galleriesC.NewView)).Methods("GET")

	r.HandleFunc("/galleries/{id:[0-9]+}", 
		galleriesC.Show).Methods("GET").
                Name(controllers.ShowGallery)

	r.HandleFunc("/galleries", 
                 requireUserMw.ApplyFn(galleriesC.Create)).Methods("POST")

	r.HandleFunc("/galleries/{id:[0-9]+}/edit",
		requireUserMw.ApplyFn(galleriesC.Edit)).Methods("GET").
		Name(controllers.EditGallery)

	r.HandleFunc("/galleries/{id:[0-9]+}/update",
		requireUserMw.ApplyFn(galleriesC.Update)).Methods("POST")

	r.HandleFunc("/galleries/{id:[0-9]+}/delete",
		requireUserMw.ApplyFn(galleriesC.Delete)).Methods("POST")

	r.HandleFunc("/galleries/{id:[0-9]+}/images",
		requireUserMw.ApplyFn(galleriesC.ImageUpload)).
                Methods("POST")

	r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")

	log.Println("Starting the server on :3000...")

	http.ListenAndServe(":3000", userMw.Apply(r))
}

