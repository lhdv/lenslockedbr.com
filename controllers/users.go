package controllers

import (
	"fmt"
	"net/http"
	"time"

	"lenslockedbr.com/context"
	"lenslockedbr.com/email"
	"lenslockedbr.com/models"
	"lenslockedbr.com/rand"
	"lenslockedbr.com/views"
)

type SignupForm struct {
	Name     string `schema:"name"`
	Age      int    `schema:"age"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

type ResetPwForm struct {
	Email    string `schema:"email"`
	Token    string `schema:"token"`
	Password string `schema:"password"`
}

type Users struct {
	NewView      *views.View
	LoginView    *views.View
	ForgotPwView *views.View
	ResetPwView  *views.View
	service      models.UserService
	emailer      *email.Client
}

func NewUsers(us models.UserService, emailer *email.Client) *Users {
	return &Users{
		NewView: views.NewView("bootstrap", false,
			"users/new"),

		LoginView: views.NewView("bootstrap", false,
			"users/login"),
		ForgotPwView: views.NewView("bootstrap", false,
			"users/forgot_pw"),
		ResetPwView: views.NewView("bootstrap", false,
			"users/reset_pw"),
		service: us,
		emailer: emailer,
	}
}

//
// New is used to render the form where a user can create a new
// user account.
//
// GET /signup
//
func (u *Users) New(w http.ResponseWriter, r *http.Request) {

	var form SignupForm

	parseURLParams(r, &form)

	u.NewView.Render(w, r, form)
}

//
// Create is used to process the signup form when a user
// tries to create a new user account.
//
// POST / signup
//
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form SignupForm

	vd.Yield = &form

	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		u.NewView.Render(w, r, vd)
		return
	}

	user := models.User{
		Name:     form.Name,
		Age:      form.Age,
		Email:    form.Email,
		Password: form.Password,
	}

	if err := u.service.Create(&user); err != nil {
		vd.SetAlert(err)
		u.NewView.Render(w, r, vd)
		return
	}

	err := u.signIn(w, &user)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	alert := views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "Welcome to LensLockedBR.com!",
	}

	views.RedirectAlert(w, r, "/galleries", http.StatusFound, alert)
}

// Login is used to process the login form when a user tries to log
// in as an existing user(via email & pwd).
//
// POST /login
//
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {

	var vd views.Data

	form := LoginForm{}
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		u.LoginView.Render(w, r, vd)
		return
	}

	user, err := u.service.Authenticate(form.Email,
		form.Password)
	if err != nil {

		switch err {
		case models.ErrNotFound:
			vd.AlertError("No user exists with that " +
				"email address")
		default:
			vd.SetAlert(err)
		}
		u.LoginView.Render(w, r, vd)
		return
	}

	err = u.signIn(w, user)
	if err != nil {
		vd.SetAlert(err)
		u.LoginView.Render(w, r, vd)
		return
	}

	alert := views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "Welcome back " + user.Name,
	}

	views.RedirectAlert(w, r, "/galleries", http.StatusFound, alert)
}

// Logout is used to delete a user's session cookie and invalidate
// theis current remember token, which will sign the current user out.
func (u *Users) Logout(w http.ResponseWriter, r *http.Request) {
	// First expire the user's cookie
	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)

	// Then we update the user with a new remember token
	user := context.User(r.Context())
	// We are ignoring errors for now because they are unlikely,
	// and even if they do occur we can't recover now that the
	// user doesn't have a valid cookie
	token, _ := rand.RememberToken()
	user.Remember = token
	u.service.Update(user)
	// Finally send the user to the home page
	alert := views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "Successfully logged out!",
	}

	views.RedirectAlert(w, r, "/galleries", http.StatusFound, alert)
}

// POST /forgot
func (u *Users) InitiateReset(w http.ResponseWriter, r *http.Request) {

	var vd views.Data
	var form ResetPwForm

	vd.Yield = &form
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		u.ForgotPwView.Render(w, r, vd)
		return
	}

	token, err := u.service.InitiateReset(form.Email)
	if err != nil {
		vd.SetAlert(err)
		u.ForgotPwView.Render(w, r, vd)
		return
	}

	err = u.emailer.ResetPw(form.Email, token)
	if err != nil {
		vd.SetAlert(err)
		u.ForgotPwView.Render(w, r, vd)
		return
	}

	v := views.Alert{
		Level: views.AlertLvlSuccess,
		Message: "Instructions for reseting your password have " +
			"been emailed to you.",
	}
	views.RedirectAlert(w, r, "/reset", http.StatusFound, v)
}

// ResetPw display the reset password form and has a method so that we
// can prefill the form data with a token provided via the URL query
// params.
//
// GET /reset
func (u *Users) ResetPw(w http.ResponseWriter, r *http.Request) {

	var vd views.Data
	var form ResetPwForm

	vd.Yield = &form
	if err := parseURLParams(r, &form); err != nil {
		vd.SetAlert(err)
	}

	u.ResetPwView.Render(w, r, vd)
}

// POST /reset
func (u *Users) CompleteReset(w http.ResponseWriter, r *http.Request) {

	var vd views.Data
	var form ResetPwForm

	vd.Yield = &form
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		u.ResetPwView.Render(w, r, vd)
		return
	}

	user, err := u.service.CompleteReset(form.Token, form.Password)
	if err != nil {
		vd.SetAlert(err)
		u.ResetPwView.Render(w, r, vd)
		return
	}

	u.signIn(w, user)

	v := views.Alert{
		Level: views.AlertLvlSuccess,
		Message: "Your password has been reset and you have " +
			"been logged in!",
	}
	views.RedirectAlert(w, r, "/galleries", http.StatusFound, v)
}

// CookieTest is used to display cookies set on the current user
func (u *Users) CookieTest(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("remember_cookie")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	user, err := u.service.ByRemember(cookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "User found is:", user)
}

/////////////////////////////////////////////////////////////////////
//
// HELPER METHODS
//
/////////////////////////////////////////////////////////////////////

// signIn is used to sign the given user in via cookies
func (u *Users) signIn(w http.ResponseWriter, user *models.User) error {

	// Set a remember token if none is found
	if user.Remember == "" {

		token, err := rand.RememberToken()
		if err != nil {
			return err
		}

		user.Remember = token

		err = u.service.Update(user)
		if err != nil {
			return err
		}

	}

	// Set a cookie with remember token from user
	cookie := http.Cookie{
		Name:     "remember_cookie",
		Value:    user.Remember,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	return nil
}
