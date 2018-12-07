package models

import (
	"github.com/jinzhu/gorm"
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

type galleryGorm struct {
	db *gorm.DB
}

func (g *galleryGorm) Create(gallery *Gallery) error {
	return g.db.Create(gallery).Error
}

