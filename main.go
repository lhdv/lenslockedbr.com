package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"lenslockedbr.com/controllers"
	"lenslockedbr.com/email"
	"lenslockedbr.com/middleware"
	"lenslockedbr.com/models"
	"lenslockedbr.com/rand"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

func main() {

	boolPtr := flag.Bool("prod", false, "Provide this flag in " +
				"production. This ensures that a " +
				".config file is provided before the " +
				"application starts.")
	flag.Parse()

	cfg := LoadConfig(*boolPtr)
	dbCfg := cfg.Database

	services, err := models.NewServices(
		models.WithGorm(dbCfg.Dialect(), dbCfg.ConnectionInfo()),
		models.WithLogMode(!cfg.IsProd()),
		models.WithUser(cfg.Pepper, cfg.HMACKey),
		models.WithGallery(),
		models.WithImage(),)
	if err != nil {
		panic(err)
	}

	defer services.Close()
	services.AutoMigrate()

	mgCfg := cfg.Mailgun
	emailer := email.NewClient(email.WithMailgun(mgCfg.Domain,
                                                     mgCfg.APIKey,
                                                     mgCfg.PublicAPIKey),
                                   email.WithSender("LensLockedBR Team",
                                                    "we@"+mgCfg.Domain))

	r := mux.NewRouter()

	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(services.User, emailer)
	galleriesC := controllers.NewGalleries(services.Gallery, 
                                               services.Image, r)

	//
	// Middleware setup
	//
	userMw := middleware.User {
		UserService: services.User,
	}
	requireUserMw := middleware.RequireUser{ }

	b, err := rand.Bytes(32)
	if err != nil {
		panic(err)
	}
	csrfMw := csrf.Protect(b, csrf.Secure(cfg.IsProd()))

	r.NotFoundHandler = http.HandlerFunc(staticC.PageNotFound.ServeHTTP)

	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.Handle("/faq", staticC.Faq).Methods("GET")

	//
	// User routes
	//

	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")
	r.Handle("/logout", 
		requireUserMw.ApplyFn(usersC.Logout)).Methods("POST")

	r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")
	r.Handle("/forgot", usersC.ForgotView).Methods("GET")
	r.HandleFunc("/forgot", usersC.InitiateReset).Methods("POST")
	r.Handle("/reset", usersC.ResettView).Methods("GET")
	r.HandleFunc("/reset", usersC.CompleteReset).Methods("POST")

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
	
	//
	// Image routes
	//
	imageHandler := http.FileServer(http.Dir("./images/"))
	imageHandler = http.StripPrefix("/images/", imageHandler)
	r.PathPrefix("/images/").Handler(imageHandler)

	r.HandleFunc("/galleries/{id:[0-9]+}/images/{filename}/delete",
		requireUserMw.ApplyFn(galleriesC.ImageDelete)).
                Methods("POST")

	//
	// Assets routes
	//
	assetHandler := http.FileServer(http.Dir("./assets/"))
	assetHandler = http.StripPrefix("/assets/", assetHandler)
	r.PathPrefix("/assets/").Handler(assetHandler)

	log.Printf("Starting the server on :%d...\n", cfg.Port)

	http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), 
                            csrfMw(userMw.Apply(r)))
}

