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
		Home: views.NewView("bootstrap", "static/home"),
                Contact: views.NewView("bootstrap","static/contact"),
	}
}
