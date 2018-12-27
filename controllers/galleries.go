package controllers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"lenslockedbr.com/context"
	"lenslockedbr.com/models"
	"lenslockedbr.com/views"
)

const (
	ShowGallery = "show_gallery"
)

type GalleryForm struct {
	Title string `schema:"title"`
}

type Galleries struct {
	NewView *views.View
	ShowView *views.View
	EditView *views.View
        gs models.GalleryService
	r *mux.Router
}

func NewGalleries(gs models.GalleryService, r *mux.Router) *Galleries {
	return &Galleries {
		NewView: views.NewView("bootstrap", false, 
                                       "galleries/new"),
		ShowView: views.NewView("bootstrap", false, 
                                       "galleries/show"),
                EditView: views.NewView("bootstrap", false,
                                       "galleries/edit"),
		gs: gs,
		r: r,
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

	url, err := g.r.Get(ShowGallery).URL("id",
		strconv.Itoa(int(gallery.ID)))
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	http.Redirect(w, r, url.Path, http.StatusFound)
}

func (g *Galleries) Show(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}

	var vd views.Data
	vd.Yield = gallery
	g.ShowView.Render(w, vd)
}

func (g *Galleries) Edit(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}

	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "You do not have permission to edit " +
                              "this gallery.",
                              http.StatusForbidden)
		return
	}

	var vd views.Data
	vd.Yield = gallery
	g.EditView.Render(w, vd)
}

func (g *Galleries) galleryByID(w http.ResponseWriter, r *http.Request) (*models.Gallery, error) {

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid gallery ID", http.StatusNotFound)
		return nil, err
	}

	gallery, err := g.gs.ByID(uint(id))
	if err != nil {
		switch err {
		case models.ErrNotFound:
			http.Error(w, "Gallery not found", 
                                      http.StatusNotFound)
		default:
			http.Error(w, "Whoops! Something went wrong.", 
                                      http.StatusInternalServerError)
		}

		return nil, err
	}

	return gallery, nil
}



