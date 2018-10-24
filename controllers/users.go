package controllers

import (
	"fmt"
	"net/http"

	"lenslockedbr.com/views"

	"github.com/gorilla/schema"
)

type SignupForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

type Users struct {
	NewView *views.View
}

func NewUsers() *Users {
	return &Users{
		NewView: views.NewView("bootstrap", 
                                       "views/users/new.gohtml"),
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
	if err := r.ParseForm() ; err != nil {
		panic(err)
	}

	dec := schema.NewDecoder()
	form := SignupForm{}
	if err := dec.Decode(&form, r.PostForm); err != nil {
		panic(err)
	}
	fmt.Fprintln(w, form)
}

