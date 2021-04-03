
package models

import (
	"github.com/jinzhu/gorm"
)

type GeneratedModel struct {
	gorm.Model
}

type GeneratedModelDB interface {
	Create(generatedModel *GeneratedModel) error
	Delete(id uint) error
}

func NewGeneratedModelService(db *gorm.DB) GeneratedModelService {
	return &generatedModelService{
		GeneratedModelDB: &generatedModelValidator{&generatedModelGorm{db}},
	}
}

type GeneratedModelService interface {
	GeneratedModelDB
}

type generatedModelService struct {
	GeneratedModelDB
}

type generatedModelValFunc func(generatedModel *GeneratedModel) error

func runGeneratedModelValFuncs(generatedModel *GeneratedModel, fns ...generatedModelValFunc) error {
	for _, fn := range fns {
		if err := fn(generatedModel); err != nil {
			return err
		}
	}
	return nil
}

type generatedModelValidator struct {
	GeneratedModelDB
}

func (mv *generatedModelValidator) Create(generatedModel *GeneratedModel) error {
	err := runGeneratedModelValFuncs(generatedModel)
	if err != nil {
		return err
	}
	return mv.GeneratedModelDB.Create(generatedModel)
}

func (mv *generatedModelValidator) Delete(id uint) error {
	if id <= 0 {
		return ErrIDInvalid
	}
	return mv.GeneratedModelDB.Delete(id)
}

type generatedModelGorm struct {
	db *gorm.DB
}

func (mg *generatedModelGorm) Create(generatedModel *GeneratedModel) error {
	return mg.db.Create(generatedModel).Error
}

func (mg *generatedModelGorm) Delete(id uint) error {
	generatedModel := GeneratedModel{Model: gorm.Model{ID: id}}
	return mg.db.Delete(&generatedModel).Error
}
