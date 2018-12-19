package controllers

import (
	"fmt"
	"net/http"
	"lenslockedbr.com/context"
	"lenslockedbr.com/models"
	"lenslockedbr.com/views"
)

type Galleries struct {
	NewView *views.View
        gs models.GalleryService
}

func NewGalleries(gs models.GalleryService) *Galleries {
	return &Galleries {
		NewView: views.NewView("bootstrap", false, 
                                       "galleries/new"),
		gs: gs,
	}
}

func (g *Galleries) New(w http.ResponseWriter, r *http.Request) {
	g.NewView.Render(w, nil)
}

func (g *Galleries) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form GalleryForm

	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		g.NewView.Render(w, vd)
		return
	}

	user := context.User(r.Context())

	gallery := models.Gallery {
		Title: form.Title,
		UserID: user.ID,
	}

	if err := g.gs.Create(&gallery); err != nil {
		vd.SetAlert(err)
		g.NewView.Render(w, vd)
		return
	}

	fmt.Fprintln(w, gallery)
}

type GalleryForm struct {
	Title string `schema:"title"`
}
