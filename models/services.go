package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Services struct {
	User UserService
	Gallery GalleryService
	Image ImageService
	db *gorm.DB
}

func NewServices(dialect, connectionInfo string) (*Services, error) {
	db, err := gorm.Open(dialect, connectionInfo)
	if err != nil {
		return nil, err
	}

	db.LogMode(true)

	// And next we need to construct services, but we can't construct
	// the UserService yet.
	return &Services {
		User: NewUserService(db),
		Gallery: NewGalleryService(db),
		Image: NewImageService(),
		db: db,
	}, nil
}

// Closes the database connection
func (s *Services) Close() error {
	return s.db.Close()
}

// Automigrate will attempt to automatically migrate all tables
func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(&User{}, &Gallery{}).Error
}

// DestructiveReset drops all tables and rebuilds them
func (s *Services) DestructiveReset() error {
	err := s.db.DropTableIfExists(&User{}, &Gallery{}).Error
	if err != nil {
		return err
	}

	return s.AutoMigrate()
}
