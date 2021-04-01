package models

import (
	"io"
	"os"
)


type ImageService interface {
	// Mutations
	Create(path, filename string, reader io.ReadCloser) error
	
	// Multiple queries
	ByPath(path string) ([]string, error)
}

func NewImageService() ImageService {
	return &imageService{}
}

type imageService struct {

}

func (is *imageService) Create(path, filename string, reader io.ReadCloser) error {
	defer reader.Close()
	// Create Directory if does not exists
	err := is.imagePath(path)
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

func (is *imageService) ByPath(path string) ([]string, error) {
	return nil, nil
}

func (is *imageService) imagePath(path string) error {
	return os.MkdirAll(path, 0700)
}