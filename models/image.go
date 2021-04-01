package models

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)


type Image struct {
	GalleryID uint
	Filename string
}

func (i *Image) Path() string {
	temp := url.URL{
		Path: "/" + i.RelativePath(),
	}
	return temp.String()
}

func (i *Image) RelativePath() string {
	return fmt.Sprintf("media/galleries/%v/%v", i.GalleryID, i.Filename)
}

func (i *Image) DeletePath() string {
	temp := url.URL{
		Path: fmt.Sprintf("/galleries/%v/images/%v/delete", i.GalleryID, i.Filename),
	}
	return temp.String()
}

type ImageService interface {
	// Mutations
	Create(galleryID uint, filename string, reader io.ReadCloser) error
	Delete(img *Image) error
	
	// Multiple queries
	ByGalleryID(galleryID uint) ([]Image, error)
}

func NewImageService() ImageService {
	return &imageService{}
}

type imageService struct {

}

func (is *imageService) Create(galleryID uint, filename string, reader io.ReadCloser) error {
	defer reader.Close()
	// Create Directory if does not exists
	path, err := is.mkImagePath(galleryID)
	if err != nil {
		return err
	}
	
	// Create a destination file
	dst, err := os.Create(path + filename)
	if err != nil {
		return err
	}
	defer dst.Close()
	
	// Copy reader data to destination file
	_, err = io.Copy(dst, reader)
	if err != nil {
		return err
	}
	return nil
}

func (is *imageService) Delete(img *Image) error {
	return os.Remove(img.RelativePath())
}

func (is *imageService) ByGalleryID(galleryID uint) ([]Image, error) {
	galleryImagePath := is.galleryImagePath(galleryID)
	files, err := filepath.Glob(galleryImagePath + "*")
	if err != nil {
		return nil, err
	}
	images := make([]Image, len(files))
	for idx := range files {
		files[idx] = strings.Replace(files[idx], galleryImagePath, "", 1)
		images[idx] = Image{
			GalleryID: galleryID,
			Filename: files[idx],
		}
	}
	return images, nil
}

func (is *imageService) mkImagePath(galleryID uint) (string, error) {
	galleryImagePath := is.galleryImagePath(galleryID)
	err := os.MkdirAll(galleryImagePath, 0755)
	if err != nil {
		return "", err
	}
	return galleryImagePath, nil
}

func (is *imageService) galleryImagePath(galleryID uint) string {
	return fmt.Sprintf("media/galleries/%v/", galleryID)
}