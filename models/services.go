package models

import (
	"github.com/jinzhu/gorm"
)

func NewServices(connectionInfo string) (*Services, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	
	return &Services{
		User:    NewUserService(db),
		Gallery: NewGalleryService(db),
		Image:   NewImageService(),
		DB:      db,
	}, nil
}

type Services struct {
	User    UserService
	Gallery GalleryService
	Image   ImageService
	DB      *gorm.DB
}

func (s *Services) Close() error {
	return s.DB.Close()
}

func (s *Services) DestructiveReset() error {
	err := s.DB.DropTableIfExists(&User{}, &Gallery{}).Error
	if err != nil {
		return err
	}
	return s.AutoMigrate()
}

func (s *Services) AutoMigrate() error {
	return s.DB.AutoMigrate(&User{}, &Gallery{}).Error
}
