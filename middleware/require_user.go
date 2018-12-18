package middleware

import (
	"log"
	"net/http"
	"lenslockedbr.com/models"
)

type RequireUser struct {
	models.UserService
}

// ApplyFn will return an http.HandlerFunc that will check to see if a
// user is logged in and then either cll new(w, r)if they are, or
// redirect them to the login page if they are not.
func (mw *RequireUser) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	// We want to return a dynamically created
	// func(http.ResponseWriter, *http.Request) but we also need to
	// convert it into a http.HandlerFunc.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		cookie, err := r.Cookie("remember_cookie")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		user, err := mw.UserService.ByRemember(cookie.Value)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		log.Println("User found: ", user)
		next(w, r)
	})
}

func (mw *RequireUser) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}

