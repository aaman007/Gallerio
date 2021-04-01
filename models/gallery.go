package models

import (
	"github.com/jinzhu/gorm"
)

type Gallery struct {
	gorm.Model
	UserID uint     `gorm:"not_null;index"`
	Title  string   `gorm:"not_null"`
	Images []Image `gorm:"-"`
}

func (g *Gallery) ImageSplitN(n int) [][]Image {
	ret := make([][]Image, n)
	for i:=0; i<n; i++ {
		ret[i] = make([]Image, 0)
	}
	
	for i, img := range g.Images {
		bucket := i % n
		ret[bucket] = append(ret[bucket], img)
	}
	return ret
}

type GalleryService interface {
	GalleryDB
}

type GalleryDB interface {
	// Methods for multiple gallery queries
	ByUserID(id uint) ([]Gallery, error)
	
	// Methods for single gallery queries
	ByID(id uint) (*Gallery, error)
	
	// Methods for modifying Gallery
	Create(gallery *Gallery) error
	Update(gallery *Gallery) error
	Delete(id uint) error
}

func NewGalleryService(db *gorm.DB) GalleryService {
	return &galleryService{
		GalleryDB: &galleryValidator{&galleryGorm{db}},
	}
}

type galleryService struct {
	GalleryDB
}

type galleryValFunc func(gallery *Gallery) error

func runGalleryValFuncs(user *Gallery, fns ...galleryValFunc) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}
	}
	return nil
}

type galleryValidator struct {
	GalleryDB
}

func (gv *galleryValidator) Create(gallery *Gallery) error {
	err := runGalleryValFuncs(gallery,
		gv.userIDRequired,
		gv.titleRequired,
	)
	if err != nil {
		return err
	}
	return gv.GalleryDB.Create(gallery)
}

func (gv *galleryValidator) Update(gallery *Gallery) error {
	err := runGalleryValFuncs(gallery,
		gv.userIDRequired,
		gv.titleRequired,
	)
	if err != nil {
		return err
	}
	return gv.GalleryDB.Update(gallery)
}

func (gv *galleryValidator) Delete(id uint) error {
	if id <= 0 {
		return ErrIDInvalid
	}
	return gv.GalleryDB.Delete(id)
}

func (gv *galleryValidator) userIDRequired(gallery *Gallery) error {
	if gallery.UserID <= 0 {
		return ErrUserIDRequired
	}
	return nil
}

func (gv *galleryValidator) titleRequired(gallery *Gallery) error {
	if gallery.Title == "" {
		return ErrTitleRequired
	}
	return nil
}

var _ GalleryDB = &galleryGorm{}

type galleryGorm struct {
	db *gorm.DB
}

func (gg *galleryGorm) ByUserID(userId uint) ([]Gallery, error) {
	var galleries []Gallery
	gg.db.Where("user_id = ?", userId).Find(&galleries)
	return galleries, nil
}

func (gg *galleryGorm) ByID(id uint) (*Gallery, error) {
	var gallery Gallery
	db := gg.db.Where("id = ?", id)
	err := First(db, &gallery)
	return &gallery, err
}

func (gg *galleryGorm) Create(gallery *Gallery) error {
	return gg.db.Create(gallery).Error
}

func (gg *galleryGorm) Update(gallery *Gallery) error {
	return gg.db.Save(gallery).Error
}

func (gg *galleryGorm) Delete(id uint) error {
	gallery := Gallery{Model: gorm.Model{ID: id}}
	return gg.db.Delete(&gallery).Error
}
