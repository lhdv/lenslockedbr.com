package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"io"
	"encoding/json"
	"bytes"

	"golang.org/x/oauth2"

	llctx "lenslockedbr.com/context"
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

	cfg := LoadConfig(*boolPtr)
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
	userMw := middleware.User{
		UserService: services.User,
	}
	requireUserMw := middleware.RequireUser{}

	b, err := rand.Bytes(32)
	if err != nil {
		panic(err)
	}
	csrfMw := csrf.Protect(b, csrf.Secure(cfg.IsProd()))

	//
	// Dropbox OAuth Code - BEGIN
	//
        dbxOAuth := &oauth2.Config{
		ClientID: cfg.Dropbox.ID,
		ClientSecret: cfg.Dropbox.Secret,
		Endpoint: oauth2.Endpoint{
			AuthURL: cfg.Dropbox.AuthURL,
			TokenURL: cfg.Dropbox.TokenURL,
		},
		RedirectURL:"http://localhost:3001/oauth/dropbox/callback",
	}

	dbxRedirect := func(w http.ResponseWriter, r *http.Request) {
		state := csrf.Token(r)
		cookie := http.Cookie {
			Name: "oauth_state",
			Value: state,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
		url := dbxOAuth.AuthCodeURL(state)
		log.Println(state)
		http.Redirect(w, r, url, http.StatusFound)
	}

	r.HandleFunc("/oauth/dropbox/connect", 
                     requireUserMw.ApplyFn(dbxRedirect))

	dbxCallback := func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		fmt.Fprintln(w, "code: ", r.FormValue("code"),
                                " state: ", r.FormValue("state"))
		state := r.FormValue("state")
		cookie, err := r.Cookie("oauth_state")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		} else if cookie == nil || cookie.Value != state {
			http.Error(w, "Invalid state provided", 
                                   http.StatusBadRequest)
			return
		}
		cookie.Value = ""
		cookie.Expires = time.Now()
		http.SetCookie(w, cookie)

		code := r.FormValue("code")
		token, err := dbxOAuth.Exchange(context.TODO(), code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user := llctx.User(r.Context())
		existing, err := services.OAuth.Find(user.ID,
                                                     models.OAuthDropbox)
		if err == models.ErrNotFound {
			// noop
		} else if err != nil {
			http.Error(w, err.Error(), 
                                   http.StatusInternalServerError)
			return
		} else {
			services.OAuth.Delete(existing.ID)
		}

		userOAuth := models.OAuth {
			UserID: user.ID,
			Token: *token,
			Service: models.OAuthDropbox,
		}	
		err = services.OAuth.Create(&userOAuth)
		if err != nil {
			http.Error(w, err.Error(), 
                                   http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "%+v", token)
	}

	r.HandleFunc("/oauth/dropbox/callback", 
                     requireUserMw.ApplyFn(dbxCallback))


	dbxQuery := func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		path := r.FormValue("path")

		user := llctx.User(r.Context())
		userOAuth, err := services.OAuth.Find(user.ID,
                                                      models.OAuthDropbox)
		if err != nil {
			panic(err)
		} 

		token := userOAuth.Token
		client := dbxOAuth.Client(context.TODO(), &token)

		url := "https://api.dropboxapi.com/2/files/list_folder"

		data := struct {
			Path string `json:"path"`
		}{
			Path: path,
		}

		dataBytes, err := json.Marshal(data)
		if err != nil {
			panic(err)
		} 

		req, err := http.NewRequest(http.MethodPost,
                                            url,
                                            bytes.NewReader(dataBytes))
		if err != nil {
			panic(err)
		} 

		req.Header.Add("Content-Type", "application/json")
		
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		} 

		defer resp.Body.Close()

		io.Copy(w, resp.Body)
		
	}

	r.HandleFunc("/oauth/dropbox/test", 
                     requireUserMw.ApplyFn(dbxQuery))
	//
	// Dropbox OAuth Code - END
	//

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
