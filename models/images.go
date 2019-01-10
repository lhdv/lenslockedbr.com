package models

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type ImageService interface {
	Create(galleryID uint, r io.Reader, filename string) error
}

func NewImageService() ImageService {
	return &imageService{}
}

type imageService struct {}

func (is *imageService) Create(galleryID uint, r io.Reader, filename string) error {

	path, err := is.mkImagePath(galleryID)
	if err != nil {
		return err
	}

	// Create a destination file
	dst, err := os.Create(filepath.Join(path, filename))
	if err != nil {
		return err
	}	
	defer dst.Close()

	// Copy uploaded file data to the destination file
	_, err = io.Copy(dst, r)
	if err != nil {
		return err
	}

	return nil
}

func (is *imageService) mkImagePath(galleryID uint) (string, error) {

	galleryPath := filepath.Join("images", "galleries",
                                     fmt.Sprintf("%v", galleryID))
	err := os.MkdirAll(galleryPath, 0755)
	if err != nil {
		return "", err
	}

	return galleryPath, nil
}

