package middleware

import (
	"net/http"
	"strings"

	"lenslockedbr.com/context"
	"lenslockedbr.com/models"
)

type RequireUser struct{}

func (mw *RequireUser) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter,
		r *http.Request) {

		user := context.User(r.Context())
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		next(w, r)
	})
}

func (mw *RequireUser) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}

// User middleware will lookup the current user via their remember_token
// cookie using the UserService. If the user is found, they will be set
// on the request context.
// Regardless, the next handler is always called.
type User struct {
	models.UserService
}

func (mw *User) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter,
		r *http.Request) {
		path := r.URL.Path
		if strings.HasPrefix(path, "/assets/") ||
			strings.HasPrefix(path, "/images/") {
			next(w, r)
			return
		}

		cookie, err := r.Cookie("remember_cookie")
		if err != nil {
			next(w, r)
			return
		}

		user, err := mw.UserService.ByRemember(cookie.Value)
		if err != nil {
			next(w, r)
			return
		}

		ctx := r.Context()
		ctx = context.WithUser(ctx, user)
		r = r.WithContext(ctx)
		next(w, r)
	})
}

func (mw *User) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}
