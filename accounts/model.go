package accounts

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"go-web-dev-2/utils/hash"
	"go-web-dev-2/utils/rand"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrNotFound = errors.New("models: resource not found")
	ErrInvalidID = errors.New("models: ID provided was invalid")
	ErrInvalidPassword = errors.New("models: Incorrect password provided")
)

const applicationPasswordPepper = "asdhgs73ehgsahdahe36daghsdh3e"
const hmacSecretKey = "dshjrewedshjf38274gewrh"


type User struct {
	gorm.Model
	Name string
	Username string `gorm:"not null;unique_index"`
	Email string `gorm:"not null;unique_index"`
	Password string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	RememberToken string `gorm:"-"`
	RememberTokenHash string `gorm:"not null;unique_index"`
}

type UserDB interface {
	// Methods for single user queries
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByRememberToken(token string) (*User, error)

	// Methods for modifying user
	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error

	// Closing DB connection
	Close() error

	// Helper for migrations
	AutoMigrate() error
	DestructiveReset() error
}

type UserService interface {
	Authenticate(email, password string) (*User, error)
	UserDB
}

func NewUserService(connectionInfo string) (UserService, error) {
	ug, err := newUserGorm(connectionInfo)
	if err != nil {
		return nil, err
	}

	hmac := hash.NewHMAC(hmacSecretKey)
	return &userService{
		UserDB: &userValidator{
			UserDB: ug,
			hmac: hmac,
		},
	}, nil
}

type userService struct {
	UserDB
}

func (us *userService) Authenticate(email, password string) (*User, error) {
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash),
		[]byte(password+applicationPasswordPepper),
	)
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return nil, ErrInvalidPassword
		default:
			return nil, err
		}
	}

	return foundUser, nil
}

type userValidatorFunc func(*User) error

func runUserValidatorFuncs(user *User, fns ...userValidatorFunc) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}
	}
	return nil
}

type userValidator struct {
	UserDB
	hmac hash.HMAC
}

func (uv *userValidator) ByRememberToken(token string) (*User, error) {
	user := &User{
		RememberToken: token,
	}
	if err := runUserValidatorFuncs(user, uv.hashRememberToken); err != nil {
		return nil, err
	}
	return uv.UserDB.ByRememberToken(user.RememberToken)
}

func (uv *userValidator) Create(user *User) error {
	err := runUserValidatorFuncs(user,
		uv.bcryptPassword,
		uv.defaultRememberToken,
		uv.hashRememberToken,
	)
	if err != nil {
		return err
	}
	return uv.UserDB.Update(user)
}

func (uv *userValidator) Update(user *User) error {
	err := runUserValidatorFuncs(user, uv.bcryptPassword, uv.hashRememberToken)
	if err != nil {
		return err
	}
	return uv.UserDB.Update(user)
}

func (uv *userValidator) Delete(id uint) error {
	var user User
	user.ID = id
	err := runUserValidatorFuncs(&user, uv.idGreaterThan(0))
	if err != nil {
		return err
	}
	return uv.UserDB.Delete(id)
}

func (uv *userValidator) bcryptPassword(user *User) error {
	if user.Password == "" {
		return nil
	}
	passwordBytes := []byte(user.Password + applicationPasswordPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	return nil
}

func (uv *userValidator) hashRememberToken(user *User) error {
	if user.RememberToken == "" {
		return nil
	}
	user.RememberTokenHash = uv.hmac.Hash(user.RememberToken)
	return nil
}

func (uv *userValidator) defaultRememberToken(user *User) error {
	if user.RememberToken != "" {
		return nil
	}
	token, err := rand.RememberToken()
	if err != nil {
		return err
	}
	user.RememberToken = token
	return nil
}

func (uv *userValidator) idGreaterThan(n uint) userValidatorFunc {
	return func(user *User) error {
		if user.ID <= n {
			return ErrInvalidID
		}
		return nil
	}
}

func newUserGorm(connectionInfo string) (*userGorm, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)

	return &userGorm{
		DB: db,
	}, nil
}

type userGorm struct {
	DB *gorm.DB
}

func (ug *userGorm) ByID(id uint) (*User, error) {
	var user User
	db := ug.DB.Where("id = ?", id)
	err := first(db, &user)
	return &user, err
}

func (ug *userGorm) ByEmail(email string) (*User, error) {
	var user User
	db := ug.DB.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

func (ug *userGorm) ByRememberToken(hashedToken string) (*User, error) {
	var user User
	err := first(ug.DB.Where("remember_token_hash = ?", hashedToken), &user)
	if err != nil {
		return nil, err
	}
	return &user, err
}

func (ug *userGorm) Create(user *User) error {
	return ug.DB.Create(user).Error
}

func (ug *userGorm) Update(user *User) error {
	return ug.DB.Save(user).Error
}

func (ug *userGorm) Delete(id uint) error {
	user := User{Model: gorm.Model{ID: id}}
	return ug.DB.Delete(&user).Error
}

func (ug *userGorm) Close() error {
	return ug.DB.Close()
}

func (ug *userGorm) DestructiveReset() error {
	if err := ug.DB.DropTableIfExists(&User{}).Error; err != nil {
		return err
	}
	return ug.AutoMigrate()
}

func (ug *userGorm) AutoMigrate() error {
	if err := ug.DB.AutoMigrate(&User{}).Error; err != nil {
		return err
	}
	return nil
}

func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}
