package accounts

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	ErrNotFound = errors.New("model: resource not found")
)

func NewService(connectionInfo string) (*Service, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &Service{
		DB: db,
	}, nil
}

type Service struct {
	DB *gorm.DB
}

func (us *Service) ByID(id uint) (*User, error) {
	var user User
	err := us.DB.Where("id = ?", id).First(&user).Error
	switch err {
	case nil:
		return &user, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (us *Service) Create(user *User) error {
	return us.DB.Create(user).Error
}

func (us *Service) Close() error {
	return us.DB.Close()
}

func (us *Service) DestructiveReset() {
	us.DB.DropTableIfExists(&User{})
	us.DB.AutoMigrate(&User{})
}

type User struct {
	gorm.Model
	Name string
	Email string `gorm:"not null;unique_index"`
}
