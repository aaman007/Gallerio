package core

import (
	"gallerio/accounts"
	"gallerio/galleries"
	"github.com/jinzhu/gorm"
)


func NewServices(connectionInfo string) (*Services, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)

	return &Services{
		User: accounts.NewUserService(db),
		db: db,
	}, nil
}

type Services struct {
	User accounts.UserService
	Gallery galleries.Gallery
	db *gorm.DB
}

func (s *Services) Close() error {
	return s.db.Close()
}

func (s *Services) DestructiveReset() error {
	err := s.db.DropTableIfExists(&accounts.User{}, &galleries.Gallery{}).Error
	if err != nil {
		return err
	}
	return s.AutoMigrate()
}

func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(&accounts.User{}, &galleries.Gallery{}).Error
}