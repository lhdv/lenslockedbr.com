package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// ServicesConfig is really just a function, but I find using types like
// this are easier to read in my source code.
type ServicesConfig func(*Services) error

type Services struct {
	User    UserService
	Gallery GalleryService
	Image   ImageService
	db      *gorm.DB
}

func NewServices(cfgs ...ServicesConfig) (*Services, error) {

	var s Services

	// For each ServicesConfig function...
	for _, cfg := range cfgs {
		// Run the function passing in a pointer to our Services
		// object and catching any errors
		if err := cfg(&s); err != nil {
			return nil, err
		}
	}

	// Then finally return the result
	return &s, nil
}

// Closes the database connection
func (s *Services) Close() error {
	return s.db.Close()
}

// Automigrate will attempt to automatically migrate all tables
func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(&User{}, &Gallery{}, &pwReset{}).Error
}

// DestructiveReset drops all tables and rebuilds them
func (s *Services) DestructiveReset() error {
	err := s.db.DropTableIfExists(&User{}, &Gallery{}, &pwReset{}).
		Error
	if err != nil {
		return err
	}

	return s.AutoMigrate()
}

func WithGorm(dialect, connectionInfo string) ServicesConfig {
	return func(s *Services) error {
		db, err := gorm.Open(dialect, connectionInfo)
		if err != nil {
			return err
		}

		s.db = db

		return nil
	}
}

func WithLogMode(mode bool) ServicesConfig {
	return func(s *Services) error {
		s.db.LogMode(mode)
		return nil
	}
}

func WithUser(pepper, hmacKey string) ServicesConfig {
	return func(s *Services) error {
		s.User = NewUserService(s.db, pepper, hmacKey)
		return nil
	}
}

func WithGallery() ServicesConfig {
	return func(s *Services) error {
		s.Gallery = NewGalleryService(s.db)
		return nil
	}
}

func WithImage() ServicesConfig {
	return func(s *Services) error {
		s.Image = NewImageService()
		return nil
	}
}
