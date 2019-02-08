package models

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"

	"lenslockedbr.com/hash"

	"log"
)

// Image is used to represent images stored in a Gallery.
// Image is NOT stored in the database, and instead references data
// stored on disk.
type Image struct {
	GalleryID uint
	Filename  string
	Hash string
}

// Path is used to build the absolute path used to reference this image
// via web request.
func (i *Image) Path() string {
	temp := url.URL{
		Path: "/" + i.RelativePath(),
	}
	return temp.String()
}

// RelativePath is used to build the path to this image on our local
// disk, relative to where our Go application is run from.
func (i *Image) RelativePath() string {
	// Convert the gallery ID to a string
	galleryID := fmt.Sprintf("%v", i.GalleryID)

	return filepath.ToSlash(filepath.Join("images",
		"galleries",
		galleryID,
		i.Filename))
}

type ImageService interface {
	Create(galleryID uint, r io.Reader, filename string) error
	ByGalleryID(galleryID uint) ([]Image, error)
	Delete(i *Image) error
}

func NewImageService() ImageService {
	return &imageService{}
}

type imageService struct{}

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

func (is *imageService) Delete(i *Image) error {
	return os.Remove(i.RelativePath())
}

func (is *imageService) ByGalleryID(galleryID uint) ([]Image, error) {

	path := is.imagePath(galleryID)
	strings, err := filepath.Glob(filepath.Join(path, "*"))
	if err != nil {
		return nil, err
	}

	// Setup the Image slice we are returning
	ret := make([]Image, len(strings))
	for i, imgStr := range strings {
		ret[i] = Image{
			GalleryID: galleryID,
			Filename:  filepath.Base(imgStr),
			Hash: is.hashFile(imgStr),
		}
	}

	return ret, nil
}

/////////////////////////////////////////////////////////////////////
//
// Helper Methods
//
/////////////////////////////////////////////////////////////////////

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

func (is *imageService) hashFile(path string) string {

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	return hash.ReaderSHA256(file)
}
