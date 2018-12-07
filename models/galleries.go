package models

import (
	"github.com/jinzhu/gorm"
)

const (
	ErrUserIDRequired modelError = "models: user ID is required"
	ErrTitleRequired modelError = "models: title is required"
)

var _ GalleryDB = &galleryGorm{}

type Gallery struct {
	gorm.Model
	
	UserID uint `gorm:not_null;index`
	Title string `gorm:not_null`
}

type GalleryDB interface {
	Create(gallery *Gallery) error
}

type GalleryService interface {
	GalleryDB
}

func NewGalleryService(db *gorm.DB) GalleryService {
	return &galleryService {
		GalleryDB: &galleryValidator {
			GalleryDB: &galleryGorm {
				db: db,
			},
		},
	}
}

type galleryService struct {
	GalleryDB
}

type galleryValidator struct {
	GalleryDB
}

func (gv *galleryValidator) userIDRequired(g *Gallery) error {
	if g.UserID <= 0 {	
		return ErrUserIDRequired
	}

	return nil
}

func (gv *galleryValidator) titleRequired(g *Gallery) error {
	if g.Title == "" {
		return ErrTitleRequired
	}	

	return nil
}

func (gv *galleryValidator) Create(gallery *Gallery) error {

	err := runGalleryValFns(gallery,
		gv.userIDRequired,
		gv.titleRequired)

	if err != nil {
		return err
	}

	return gv.GalleryDB.Create(gallery)
}

type galleryGorm struct {
	db *gorm.DB
}

func (g *galleryGorm) Create(gallery *Gallery) error {
	return g.db.Create(gallery).Error
}

type galleryValFn func(*Gallery) error

func runGalleryValFns(gallery *Gallery, fns ...galleryValFn) error {
	for _, fn := range fns {
		if err := fn(gallery); err != nil {
			return err
		}
	}

	return nil
}
