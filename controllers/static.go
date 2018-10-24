package controllers

import (
	"lenslockedbr.com/views"
)

type Static struct {
	Home *views.View
	Contact *views.View
}

func NewStatic() *Static {
	return &Static {
		Home: views.NewView("bootstrap",
                                    "views/static/home.gohtml"),
		Contact: views.NewView("bootstrap",
                                       "views/static/contact.gohtml"),
	}
}
