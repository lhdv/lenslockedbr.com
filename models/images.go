package models

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type ImageService interface {
	Create(galleryID uint, r io.Reader, filename string) error
	ByGalleryID(galleryID uint) ([]string, error)
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

func (is *imageService) ByGalleryID(galleryID uint) ([]string, error) {

	path := is.imagePath(galleryID)
	strings, err := filepath.Glob(filepath.Join(path, "*"))
	if err != nil {
		return nil, err
	}

	// Adding a leading "/" to all image file paths to be easier
	// to build <img> tag
	for i := range strings {
		strings[i] = "/" + strings[i]
		strings[i] = filepath.ToSlash(strings[i])
	}

	return strings, nil
}

func (is *imageService) mkImagePath(galleryID uint) (string, error) {

	galleryPath := is.imagePath(galleryID)
	err := os.MkdirAll(galleryPath, 0755)
	if err != nil {
		return "", err
	}

	return galleryPath, nil
}

func (is *imageService) imagePath(galleryID uint) string {
	return filepath.Join("images", "galleries",
                             fmt.Sprintf("%v", galleryID))
}

