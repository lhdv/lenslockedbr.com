package models

import (
	"github.com/jinzhu/gorm"
)

type Services struct {
	Gallery GalleryService
	User UserService
}

func NewServices(connectionInfo string) (*Services, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}

	db.LogMode(true)

	// And next we need to construct services, but we can't construct
	// the UserService yet.
	return &Services {
		User: NewUserService(db),
		Gallery: &galleryGorm{},
	}, nil
}
