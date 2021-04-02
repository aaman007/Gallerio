package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type ServicesConfig func(*Services) error

func WithGorm(dialect, connectionInfo string) ServicesConfig {
	return func(services *Services) error {
		db, err := gorm.Open(dialect, connectionInfo)
		if err != nil {
			return err
		}
		services.db = db
		return nil
	}
}

func WithLogMode(logMode bool) ServicesConfig {
	return func(services *Services) error {
		services.db.LogMode(logMode)
		return nil
	}
}

func WithUser(pepper, hmacKey string) ServicesConfig {
	return func(services *Services) error {
		services.User = NewUserService(services.db, pepper, hmacKey)
		return nil
	}
}

func WithGallery() ServicesConfig {
	return func(services *Services) error {
		services.Gallery = NewGalleryService(services.db)
		return nil
	}
}

func WithImage() ServicesConfig {
	return func(services *Services) error {
		services.Image = NewImageService()
		return nil
	}
}

func NewServices(cfgs ...ServicesConfig) (*Services, error) {
	var services Services
	for _, cfg := range cfgs {
		if err := cfg(&services); err != nil {
			return nil, err
		}
	}
	return &services, nil
}

type Services struct {
	User    UserService
	Gallery GalleryService
	Image   ImageService
	db      *gorm.DB
}

func (s *Services) Close() error {
	return s.db.Close()
}

func (s *Services) DestructiveReset() error {
	err := s.db.DropTableIfExists(&User{}, &Gallery{}, &passwordReset{}).Error
	if err != nil {
		return err
	}
	return s.AutoMigrate()
}

func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(&User{}, &Gallery{}, &passwordReset{}).Error
}
