package models

import (
	"gallerio/utils/hash"
	"gallerio/utils/rand"
	"github.com/jinzhu/gorm"
)

type passwordReset struct {
	gorm.Model
	UserID    uint   `gorm:"user_id"`
	Token     string `gorm:"-"`
	TokenHash string `gorm:"not null;unique_index"`
}

type passwordResetDB interface {
	ByToken(token string) (*passwordReset, error)
	
	Create(pwr *passwordReset) error
	Delete(id uint) error
}

type passwordResetValFunc func(*passwordReset) error

func runPasswordResetValFuncs(pwr *passwordReset, fns ...passwordResetValFunc) error {
	for _, fn := range fns {
		if err := fn(pwr); err != nil {
			return err
		}
	}
	return nil
}

func newPasswordResetValidator(db passwordResetDB, hmac hash.HMAC) *passwordResetValidator {
	return &passwordResetValidator{
		passwordResetDB: db,
		hmac:            hmac,
	}
}

type passwordResetValidator struct {
	passwordResetDB
	hmac hash.HMAC
}

func (pwrv *passwordResetValidator) ByToken(token string) (*passwordReset, error) {
	pwr := &passwordReset{Token: token}
	err := runPasswordResetValFuncs(pwr, pwrv.hashToken)
	if err != nil {
		return nil, err
	}
	return pwrv.passwordResetDB.ByToken(pwr.TokenHash)
}

func (pwrv *passwordResetValidator) Create(pwr *passwordReset) error {
	err := runPasswordResetValFuncs(pwr, pwrv.userIDRequired, pwrv.defaultToken, pwrv.hashToken)
	if err != nil {
		return err
	}
	return pwrv.passwordResetDB.Create(pwr)
}

func (pwrv *passwordResetValidator) Delete(id uint) error {
	if id <= 0 {
		return ErrIDInvalid
	}
	return pwrv.passwordResetDB.Delete(id)
}

func (pwrv *passwordResetValidator) userIDRequired(pwr *passwordReset) error {
	if pwr.UserID <= 0 {
		return ErrUserIDRequired
	}
	return nil
}

func (pwrv *passwordResetValidator) defaultToken(pwr *passwordReset) error {
	if pwr.Token != "" {
		return nil
	}
	token, err := rand.RememberToken()
	if err != nil {
		return err
	}
	pwr.Token = token
	return nil
}

func (pwrv *passwordResetValidator) hashToken(pwr *passwordReset) error {
	if pwr.Token == "" {
		return nil
	}
	pwr.TokenHash = pwrv.hmac.Hash(pwr.Token)
	return nil
}

type passwordResetGorm struct {
	db *gorm.DB
}

func (pwrg *passwordResetGorm) ByToken(tokenHash string) (*passwordReset, error) {
	var pwr passwordReset
	err := First(pwrg.db.Where("token_hash = ?", tokenHash), &pwr)
	if err != nil {
		return nil, err
	}
	return &pwr, err
}

func (pwrg *passwordResetGorm) Create(pwr *passwordReset) error {
	return pwrg.db.Create(pwr).Error
}

func (pwrg *passwordResetGorm) Delete(id uint) error {
	pwr := passwordReset{Model: gorm.Model{ID: id}}
	return pwrg.db.Delete(&pwr).Error
}
