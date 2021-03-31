package controllers

import (
	"gallerio/models"
	"github.com/jinzhu/gorm"
)


func NewServices(connectionInfo string) (*Services, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	// db.LogMode(true)

	return &Services{
		User:    models.NewUserService(db),
		Gallery: models.NewGalleryService(db),
		db:      db,
	}, nil
}

type Services struct {
	User    models.UserService
	Gallery models.GalleryService
	db      *gorm.DB
}

func (s *Services) Close() error {
	return s.db.Close()
}

func (s *Services) DestructiveReset() error {
	err := s.db.DropTableIfExists(&models.User{}, &models.Gallery{}).Error
	if err != nil {
		return err
	}
	return s.AutoMigrate()
}

func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(&models.User{}, &models.Gallery{}).Error
}