package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/oauth2"

	"lenslockedbr.com/controllers"
	"lenslockedbr.com/email"
	"lenslockedbr.com/middleware"
	"lenslockedbr.com/models"
	"lenslockedbr.com/rand"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

func main() {

	boolPtr := flag.Bool("prod", false, "Provide this flag in "+
		"production. This ensures that a "+
		".config file is provided before the "+
		"application starts.")
	flag.Parse()

	//
	// Load config file
	//
	cfg := LoadConfig(*boolPtr)

	//
	// Database configuration
	//
	dbCfg := cfg.Database

	services, err := models.NewServices(
		models.WithGorm(dbCfg.Dialect(), dbCfg.ConnectionInfo()),
		models.WithLogMode(!cfg.IsProd()),
		models.WithUser(cfg.Pepper, cfg.HMACKey),
		models.WithGallery(),
		models.WithImage(),
		models.WithOAuth())
	if err != nil {
		panic(err)
	}

	defer services.Close()
	services.AutoMigrate()

	//
	// Mailing configuration
	//
	mgCfg := cfg.Mailgun
	emailer := email.NewClient(email.WithMailgun(mgCfg.Domain,
		mgCfg.APIKey,
		mgCfg.PublicAPIKey),
		email.WithSender("LensLockedBR Team",
			"we@"+mgCfg.Domain))

	//
	// OAuth configuration
	//
	oauthCfgs := make(map[string]*oauth2.Config)
        oauthCfgs[models.OAuthDropbox] = &oauth2.Config{
		ClientID: cfg.Dropbox.ID,
		ClientSecret: cfg.Dropbox.Secret,
		Endpoint: oauth2.Endpoint{
			AuthURL: cfg.Dropbox.AuthURL,
			TokenURL: cfg.Dropbox.TokenURL,
		},
		RedirectURL:"http://localhost:3001/oauth/dropbox/callback",
	}

	r := mux.NewRouter()

	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(services.User, emailer)
	galleriesC := controllers.NewGalleries(services.Gallery,
		services.Image, r)
	oauthsC := controllers.NewOAuths(services.OAuth, oauthCfgs)

	//
	// Middleware setup
	//
	userMw := middleware.User{
		UserService: services.User,
	}
	requireUserMw := middleware.RequireUser{}

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

	r.Handle("/forgot", usersC.ForgotPwView).Methods("GET")
	r.HandleFunc("/forgot", usersC.InitiateReset).Methods("POST")
	r.HandleFunc("/reset", usersC.ResetPw).Methods("GET")
	r.HandleFunc("/reset", usersC.CompleteReset).Methods("POST")

	r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")

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

	r.HandleFunc("/galleries/{id:[0-9]+}/images/{filename}/delete",
		requireUserMw.ApplyFn(galleriesC.ImageDelete)).
		Methods("POST")

	r.HandleFunc("/galleries/{id:[0-9]+}/images/link",
		requireUserMw.ApplyFn(galleriesC.ImageViaLink)).
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
	// DropboxAPI routes
	//
	r.HandleFunc("/oauth/{service:[a-z]+}/connect", 
                     requireUserMw.ApplyFn(oauthsC.Connect))
	r.HandleFunc("/oauth/{service:[a-z]+}/callback", 
                     requireUserMw.ApplyFn(oauthsC.Callback))
	r.HandleFunc("/oauth/{service:[a-z]+}/test", 
                     requireUserMw.ApplyFn(oauthsC.DropboxTest))

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
