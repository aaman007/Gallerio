package models

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/oauth2"
)

const (
	OAuthDropbox = "dropbox"
)

type OAuth struct {
	gorm.Model
	UserID uint `gorm:"not null;unique_index:user_id_provider"`
	Provider string `gorm:"not null;unique_index:user_id_provider"`
	oauth2.Token
}

type OAuthDB interface {
	Find(userID uint, provider string) (*OAuth, error)
	Create(oauth *OAuth) error
	Delete(id uint) error
}

func NewOAuthService(db *gorm.DB) OAuthService {
	return &oauthService{
		OAuthDB: &oauthValidator{&oauthGorm{db}},
	}
}

type OAuthService interface {
	OAuthDB
}

type oauthService struct {
	OAuthDB
}

type oauthValFunc func(oauth *OAuth) error

func runOAuthValFuncs(oauth *OAuth, fns ...oauthValFunc) error {
	for _, fn := range fns {
		if err := fn(oauth); err != nil {
			return err
		}
	}
	return nil
}

func (ov *oauthValidator) userIDRequired(oauth *OAuth) error {
	if oauth.UserID <= 0 {
		return ErrUserIDRequired
	}
	return nil
}

func (ov *oauthValidator) providerRequired(oauth *OAuth) error {
	if oauth.Provider == "" {
		return ErrProviderRequired
	}
	return nil
}

type oauthValidator struct {
	OAuthDB
}

func (ov *oauthValidator) Create(oauth *OAuth) error {
	err := runOAuthValFuncs(oauth,
		ov.userIDRequired,
		ov.providerRequired,
	)
	if err != nil {
		return err
	}
	return ov.OAuthDB.Create(oauth)
}

func (ov *oauthValidator) Delete(id uint) error {
	if id <= 0 {
		return ErrIDInvalid
	}
	return ov.OAuthDB.Delete(id)
}

type oauthGorm struct {
	db *gorm.DB
}

func (og *oauthGorm) Find(userID uint, provider string) (*OAuth, error) {
	var oauth OAuth
	db := og.db.Where("user_id = ?", userID).Where("provider = ?", provider)
	err := First(db, &oauth)
	if err != nil {
		return nil, err
	}
	return &oauth, nil
}

func (og *oauthGorm) Create(oauth *OAuth) error {
	return og.db.Create(oauth).Error
}

func (og *oauthGorm) Delete(id uint) error {
	oauth := OAuth{Model: gorm.Model{ID: id}}
	return og.db.Unscoped().Delete(&oauth).Error // Deletes permanently
}