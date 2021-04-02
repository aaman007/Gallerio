package models

import (
	"gallerio/utils/hash"
	"gallerio/utils/rand"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"strings"
	"time"
)

var (
	EmailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`)
)

type User struct {
	gorm.Model
	Name              string
	Username          string `gorm:"not null;unique_index"`
	Email             string `gorm:"not null;unique_index"`
	Password          string `gorm:"-"`
	PasswordHash      string `gorm:"not null"`
	RememberToken     string `gorm:"-"`
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
}

type UserService interface {
	Authenticate(email, password string) (*User, error)
	InitiateReset(email string) (string, error)
	CompleteReset(token, newPw string) (*User, error)
	UserDB
}

func NewUserService(db *gorm.DB, pepper, hmacKey string) UserService {
	ug := &userGorm{db}
	hmac := hash.NewHMAC(hmacKey)
	uv := newUserValidator(ug, hmac, pepper)
	
	return &userService{
		UserDB:          uv,
		passwordResetDB: newPasswordResetValidator(&passwordResetGorm{db}, hmac),
		pepper:          pepper,
	}
}

type userService struct {
	UserDB
	passwordResetDB passwordResetDB
	pepper          string
}

func (us *userService) Authenticate(email, password string) (*User, error) {
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}
	
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash),
		[]byte(password+us.pepper),
	)
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return nil, ErrPasswordIncorrect
		default:
			return nil, err
		}
	}
	
	return foundUser, nil
}

func (us *userService) InitiateReset(email string) (string, error) {
	user, err := us.ByEmail(email)
	if err != nil {
		return "", err
	}
	pwr := &passwordReset{UserID: user.ID}
	err = us.passwordResetDB.Create(pwr)
	if err != nil {
		return "", err
	}
	return pwr.Token, nil
}

func (us *userService) CompleteReset(token, newPw string) (*User, error) {
	pwr, err := us.passwordResetDB.ByToken(token)
	if err != nil {
		return nil, ErrTokenInvalid
	}
	if time.Now().Sub(pwr.CreatedAt) > (12 * time.Hour) {
		return nil, ErrTokenInvalid
	}
	if len(newPw) < 8 {
		return nil, ErrPasswordTooShort
	}
	user, err := us.ByID(pwr.ID)
	if err != nil {
		return nil, err
	}
	user.Password = newPw
	err = us.Update(user)
	if err != nil {
		return nil, err
	}
	us.passwordResetDB.Delete(pwr.ID)
	return user, nil
}

type userValFunc func(*User) error

func runUserValFuncs(user *User, fns ...userValFunc) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}
	}
	return nil
}

func newUserValidator(udb UserDB, hmac hash.HMAC, pepper string) *userValidator {
	return &userValidator{
		UserDB:     udb,
		hmac:       hmac,
		pepper:     pepper,
		emailRegex: EmailRegex,
	}
}

type userValidator struct {
	UserDB
	hmac       hash.HMAC
	emailRegex *regexp.Regexp
	pepper     string
}

func (uv *userValidator) ByEmail(email string) (*User, error) {
	user := &User{Email: email}
	err := runUserValFuncs(user,
		uv.emailNormalize,
		uv.emailFormat,
	)
	if err != nil {
		return nil, err
	}
	return uv.UserDB.ByEmail(user.Email)
}

func (uv *userValidator) ByRememberToken(token string) (*User, error) {
	user := &User{
		RememberToken: token,
	}
	if err := runUserValFuncs(user, uv.hashRememberToken); err != nil {
		return nil, err
	}
	return uv.UserDB.ByRememberToken(user.RememberTokenHash)
}

func (uv *userValidator) Create(user *User) error {
	err := runUserValFuncs(user,
		uv.passwordRequired,
		uv.passwordMinLength,
		uv.passwordBcrypt,
		uv.passwordHashRequired,
		uv.defaultRememberToken,
		uv.rememberTokenMinBytes,
		uv.hashRememberToken,
		uv.rememberTokenHashRequired,
		uv.emailNormalize,
		uv.emailRequired,
		uv.emailFormat,
		uv.emailAvailable,
	)
	if err != nil {
		return err
	}
	return uv.UserDB.Create(user)
}

func (uv *userValidator) Update(user *User) error {
	err := runUserValFuncs(user,
		uv.passwordMinLength,
		uv.passwordBcrypt,
		uv.passwordHashRequired,
		uv.hashRememberToken,
		uv.rememberTokenMinBytes,
		uv.hashRememberToken,
		uv.rememberTokenHashRequired,
		uv.emailNormalize,
		uv.emailRequired,
		uv.emailFormat,
		uv.emailAvailable,
	)
	if err != nil {
		return err
	}
	return uv.UserDB.Update(user)
}

func (uv *userValidator) Delete(id uint) error {
	var user User
	user.ID = id
	err := runUserValFuncs(&user, uv.idGreaterThan(0))
	if err != nil {
		return err
	}
	return uv.UserDB.Delete(id)
}

func (uv *userValidator) passwordBcrypt(user *User) error {
	if user.Password == "" {
		return nil
	}
	passwordBytes := []byte(user.Password + uv.pepper)
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

func (uv *userValidator) idGreaterThan(n uint) userValFunc {
	return func(user *User) error {
		if user.ID <= n {
			return ErrIDInvalid
		}
		return nil
	}
}

func (uv *userValidator) emailNormalize(user *User) error {
	user.Email = strings.ToLower(user.Email)
	user.Email = strings.TrimSpace(user.Email)
	return nil
}

func (uv *userValidator) emailRequired(user *User) error {
	if user.Email == "" {
		return ErrEmailRequired
	}
	return nil
}

func (uv *userValidator) emailFormat(user *User) error {
	if !uv.emailRegex.MatchString(user.Email) {
		return ErrEmailInvalid
	}
	return nil
}

func (uv *userValidator) emailAvailable(user *User) error {
	existing, err := uv.ByEmail(user.Email)
	if err == ErrNotFound {
		return nil
	}
	if err != nil {
		return err
	}
	if existing.ID != user.ID {
		return ErrEmailTaken
	}
	return nil
}

func (uv *userValidator) passwordRequired(user *User) error {
	if user.Password == "" {
		return ErrPasswordRequired
	}
	return nil
}

func (uv *userValidator) passwordHashRequired(user *User) error {
	if user.PasswordHash == "" {
		return ErrPasswordRequired
	}
	return nil
}

func (uv *userValidator) passwordMinLength(user *User) error {
	if user.Password == "" {
		return nil
	}
	if len(user.Password) < 8 {
		return ErrPasswordTooShort
	}
	return nil
}

func (uv *userValidator) rememberTokenMinBytes(user *User) error {
	if user.RememberToken == "" {
		return nil
	}
	b, err := rand.NBytes(user.RememberToken)
	if err != nil {
		return err
	}
	if b < 32 {
		return ErrRememberTokenTooShort
	}
	return nil
}

func (uv *userValidator) rememberTokenHashRequired(user *User) error {
	if user.RememberTokenHash == "" {
		return ErrRememberTokenRequired
	}
	return nil
}

type userGorm struct {
	db *gorm.DB
}

func (ug *userGorm) ByID(id uint) (*User, error) {
	var user User
	db := ug.db.Where("id = ?", id)
	err := First(db, &user)
	return &user, err
}

func (ug *userGorm) ByEmail(email string) (*User, error) {
	var user User
	db := ug.db.Where("email = ?", email)
	err := First(db, &user)
	return &user, err
}

func (ug *userGorm) ByRememberToken(hashedToken string) (*User, error) {
	var user User
	err := First(ug.db.Where("remember_token_hash = ?", hashedToken), &user)
	if err != nil {
		return nil, err
	}
	return &user, err
}

func (ug *userGorm) Create(user *User) error {
	return ug.db.Create(user).Error
}

func (ug *userGorm) Update(user *User) error {
	return ug.db.Save(user).Error
}

func (ug *userGorm) Delete(id uint) error {
	user := User{Model: gorm.Model{ID: id}}
	return ug.db.Delete(&user).Error
}
