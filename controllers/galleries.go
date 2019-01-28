package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/mux"

	"lenslockedbr.com/context"
	"lenslockedbr.com/models"
	"lenslockedbr.com/views"
)

const (
	IndexGallery = "index_galleries"
	ShowGallery  = "show_gallery"
	EditGallery  = "edit_gallery"

	maxMultipartMem = 1 << 20 // 1 megabyte
)

type GalleryForm struct {
	Title string `schema:"title"`
}

type Galleries struct {
	NewView   *views.View
	ShowView  *views.View
	EditView  *views.View
	IndexView *views.View
	gs        models.GalleryService
	is        models.ImageService
	r         *mux.Router
}

func NewGalleries(gs models.GalleryService, is models.ImageService,
	r *mux.Router) *Galleries {
	return &Galleries{
		NewView: views.NewView("bootstrap", false,
			"galleries/new"),
		ShowView: views.NewView("bootstrap", false,
			"galleries/show"),
		EditView: views.NewView("bootstrap", false,
			"galleries/edit"),
		IndexView: views.NewView("bootstrap", false,
			"galleries/index"),
		gs: gs,
		is: is,
		r:  r,
	}
}

func (g *Galleries) New(w http.ResponseWriter, r *http.Request) {
	g.NewView.Render(w, r, nil)
}

func (g *Galleries) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form GalleryForm

	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		g.NewView.Render(w, r, vd)
		return
	}

	user := context.User(r.Context())

	gallery := models.Gallery{
		Title:  form.Title,
		UserID: user.ID,
	}

	if err := g.gs.Create(&gallery); err != nil {
		vd.SetAlert(err)
		g.NewView.Render(w, r, vd)
		return
	}

	url, err := g.r.Get(EditGallery).URL("id",
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
	g.ShowView.Render(w, r, vd)
}

func (g *Galleries) Edit(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}

	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "You do not have permission to edit "+
			"this gallery.",
			http.StatusForbidden)
		return
	}

	var vd views.Data
	vd.Yield = gallery
	g.EditView.Render(w, r, vd)
}

func (g *Galleries) Update(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}

	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "Gallery not found.",
			http.StatusForbidden)
		return
	}

	var vd views.Data
	vd.Yield = gallery

	var form GalleryForm

	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		g.EditView.Render(w, r, vd)
		return
	}

	gallery.Title = form.Title

	err = g.gs.Update(gallery)
	if err != nil {
		vd.SetAlert(err)
	} else {
		vd.Alert = &views.Alert{
			Level:   views.AlertLvlSuccess,
			Message: "Gallery updated successfully!",
		}
	}

	g.EditView.Render(w, r, vd)
}

func (g *Galleries) Delete(w http.ResponseWriter, r *http.Request) {

	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}

	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "You do not have permission to delete "+
			"this gallery.",
			http.StatusForbidden)
		return
	}

	var vd views.Data

	err = g.gs.Delete(gallery.ID)
	if err != nil {
		vd.SetAlert(err)
		vd.Yield = gallery
		g.EditView.Render(w, r, vd)
	}

	url, err := g.r.Get(IndexGallery).URL()
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	http.Redirect(w, r, url.Path, http.StatusFound)
}

func (g *Galleries) Index(w http.ResponseWriter, r *http.Request) {

	user := context.User(r.Context())

	galleries, err := g.gs.ByUserID(user.ID)
	if err != nil {
		http.Error(w, "Something went wrong.",
			http.StatusInternalServerError)
		return
	}

	var vd views.Data
	vd.Yield = galleries
	g.IndexView.Render(w, r, vd)
}

func (g *Galleries) ImageUpload(w http.ResponseWriter, r *http.Request) {

	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}

	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "Gallery not found.",
			http.StatusForbidden)
		return
	}

	var vd views.Data
	vd.Yield = gallery

	err = r.ParseMultipartForm(maxMultipartMem)
	if err != nil {
		vd.SetAlert(err)
		g.EditView.Render(w, r, vd)
		return
	}

	files := r.MultipartForm.File["images"]
	for _, f := range files {

		// Open the uploaded file
		file, err := f.Open()
		if err != nil {
			vd.SetAlert(err)
			g.EditView.Render(w, r, vd)
			return
		}
		defer file.Close()

		// Create image
		err = g.is.Create(gallery.ID, file, f.Filename)
		if err != nil {
			vd.SetAlert(err)
			g.EditView.Render(w, r, vd)
			return
		}
	}

	url, err := g.r.Get(EditGallery).
		URL("id", fmt.Sprintf("%v", gallery.ID))
	if err != nil {
		http.Redirect(w, r, "/galleries", http.StatusFound)
		return
	}

	http.Redirect(w, r, url.Path, http.StatusFound)
}

func (g *Galleries) ImageDelete(w http.ResponseWriter, r *http.Request) {

	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}

	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "You do not have permission to edit "+
			"this gallery or image", http.StatusForbidden)
		return
	}

	// Get the filname from the path
	filename := mux.Vars(r)["filename"]
	// Build the Image model
	i := models.Image{
		GalleryID: gallery.ID,
		Filename:  filename,
	}

	// Try to delete the image
	err = g.is.Delete(&i)
	if err != nil {
		// Render the edit page with any error
		var vd views.Data
		vd.Yield = gallery
		vd.SetAlert(err)
		g.EditView.Render(w, r, vd)
		return
	}

	// If all goes well, redirect to the edit gallery page
	url, err := g.r.Get(EditGallery).
		URL("id", fmt.Sprintf("%v", gallery.ID))
	if err != nil {
		http.Redirect(w, r, "/galleries", http.StatusFound)
	}

	http.Redirect(w, r, url.Path, http.StatusFound)
}

func (g *Galleries) ImageViaLink(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}

	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "Gallery not found.",
			http.StatusForbidden)
		return
	}

	var vd views.Data
	vd.Yield = gallery

	err = r.ParseForm()
	if err != nil {
		vd.SetAlert(err)
		g.EditView.Render(w, r, vd)
		return
	}

	files := r.PostForm["files"]

	var wg sync.WaitGroup
	wg.Add(len(files))

	for _, fileURL := range files {

		go func(url string) {
			defer wg.Done()
			resp, err := http.Get(url)
			if err != nil {
				log.Println("Failed to download the image from:", url)
				return
			}

			defer resp.Body.Close()

			pieces := strings.Split(url, "/")
			filename := pieces[len(pieces)-1]
			err = g.is.Create(gallery.ID, resp.Body, filename)
			if err != nil {
				log.Println("Failed to create the image from:", url)
			}
		} (fileURL)
	}

	wg.Wait()

	url, err := g.r.Get(EditGallery).
		URL("id", fmt.Sprintf("%v", gallery.ID))
	if err != nil {
		http.Redirect(w, r, "/galleries", http.StatusFound)
		return
	}

	http.Redirect(w, r, url.Path, http.StatusFound)
}

/////////////////////////////////////////////////////////////////////
//
// Helper methods
//
/////////////////////////////////////////////////////////////////////

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

	images, _ := g.is.ByGalleryID(gallery.ID)
	gallery.Images = images

	return gallery, nil
}
