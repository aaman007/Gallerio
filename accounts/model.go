package accounts

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	ErrNotFound = errors.New("models: resource not found")
	ErrInvalidID = errors.New("models: ID provided was invalid")
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
	db := us.DB.Where("id = ?", id)
	err := first(db, &user)
	return &user, err
}

func (us *Service) ByEmail(email string) (*User, error) {
	var user User
	db := us.DB.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}

func (us *Service) Create(user *User) error {
	return us.DB.Create(user).Error
}

func (us *Service) Update(user *User) error {
	return us.DB.Save(user).Error
}

func (us *Service) Delete(id uint) error {
	if id <= 0 {
		return ErrInvalidID
	}
	user := User{Model: gorm.Model{ID: id}}
	return us.DB.Delete(&user).Error
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
