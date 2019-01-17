package controllers

import (
	"lenslockedbr.com/views"
)

type Static struct {
	Home         *views.View
	Contact      *views.View
	Faq          *views.View
	PageNotFound *views.View
}

func NewStatic() *Static {
	return &Static{
		Home: views.NewView("bootstrap", false,
			"static/home"),
		Contact: views.NewView("bootstrap", false,
			"static/contact"),
		Faq: views.NewView("bootstrap_bggray", false,
			"static/faq"),
		PageNotFound: views.NewView("bootstrap", true,
			"static/page_not_found"),
	}
}
