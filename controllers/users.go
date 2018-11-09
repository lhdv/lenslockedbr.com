package controllers

import (
	"fmt"
	"net/http"

	"lenslockedbr.com/models"
	"lenslockedbr.com/views"
)

type SignupForm struct {
	Name     string `schema:"name"`
	Age      int    `schema:"age"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

type LoginForm struct {
	Email string `schema:"email"`
	Password string `schema:"password"`
}

type Users struct {
	NewView *views.View
	LoginView *views.View
	service *models.UserService
}

func NewUsers(us *models.UserService) *Users {
	return &Users{
		NewView: views.NewView("bootstrap", false,
			               "users/new"),

		LoginView: views.NewView("bootstrap", false,
			                 "users/login"),
		service: us,
	}
}

//
// New is used to render the form where a user can create a new
// user account.
//
// GET /signup
//
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	if err := u.NewView.Render(w, nil); err != nil {
		panic(err)
	}
}

//
// Create is used to process the signup form when a user
// tries to create a new user account.
//
// POST / signup
//
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var form SignupForm

	if err := parseForm(r, &form); err != nil {
		panic(err)
	}

	user := models.User{
		Name:     form.Name,
		Age:      form.Age,
		Email:    form.Email,
		Password: form.Password,
	}

	if err := u.service.Create(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "User is", user)
}

// Login is used to process the login form when a user tries to log
// in as an existing user(via email & pwd).
//
// POST /login
//
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {

	form := LoginForm{}
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}

	user, err := u.service.Authenticate(form.Email,
                                            form.Password)
	if err != nil {

		switch err {
		case models.ErrNotFound:
			fmt.Fprintln(w, "Invalid email address.")
		case models.ErrInvalidPassword:
			fmt.Fprintln(w, "Invalid password provided.")
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	cookie := http.Cookie {
		Name: "email",
		Value: user.Email,
	}

	http.SetCookie(w, &cookie)

	fmt.Fprintln(w, user)
}

// CookieTest is used to display cookies set on the current user
func (u *Users) CookieTest(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("email")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Email is:", cookie.Value)
}


