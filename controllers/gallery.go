package controllers

import (
	"fmt"
	"gallerio/forms"
	"gallerio/models"
	"gallerio/utils/context"
	"gallerio/views"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

var (
	ShowGalleryName = "show_gallery"
	EditGalleryName = "show_gallery"
	
	maxMemoryLimit int64 = 1 << 20 // 1MB
	galleryPath          = "media/galleries/"
)

func NewGalleriesController(gs models.GalleryService, is models.ImageService, router *mux.Router) *GalleriesController {
	return &GalleriesController{
		New:       views.NewView("base", "gallery/new"),
		IndexView: views.NewView("base", "gallery/index"),
		ShowView:  views.NewView("base", "gallery/show"),
		EditView:  views.NewView("base", "gallery/edit"),
		router:    router,
		gs:        gs,
		is:        is,
	}
}

type GalleriesController struct {
	New       *views.View
	IndexView *views.View
	ShowView  *views.View
	EditView  *views.View
	router    *mux.Router
	gs        models.GalleryService
	is        models.ImageService
}

// POST /galleries
func (gc *GalleriesController) Index(w http.ResponseWriter, req *http.Request) {
	user := context.User(req.Context())
	galleries, err := gc.gs.ByUserID(user.ID)
	if err != nil {
		http.Error(w, views.AlertMessageGeneric, http.StatusInternalServerError)
		return
	}
	data := views.Data{Content: galleries}
	gc.IndexView.Render(w, req, data)
}

// POST /galleries
func (gc *GalleriesController) Create(w http.ResponseWriter, req *http.Request) {
	var data views.Data
	var form forms.GalleryForm
	if err := forms.ParseForm(req, &form); err != nil {
		log.Println(err)
		data.SetAlert(err)
		gc.New.Render(w, req, data)
		return
	}
	
	user := context.User(req.Context())
	
	gallery := models.Gallery{
		Title:  form.Title,
		UserID: user.ID,
	}
	if err := gc.gs.Create(&gallery); err != nil {
		log.Println(err)
		data.SetAlert(err)
		gc.New.Render(w, req, data)
		return
	}
	url, err := gc.router.Get(EditGalleryName).URL("id", fmt.Sprintf("%v", gallery.ID))
	if err != nil {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	http.Redirect(w, req, url.Path, http.StatusSeeOther)
}

// GET /gallery/{id}
func (gc *GalleriesController) Show(w http.ResponseWriter, req *http.Request) {
	gallery, err := gc.galleryByID(w, req)
	if err != nil {
		return
	}
	data := views.Data{Content: gallery}
	gc.ShowView.Render(w, req, data)
}

// GET /galleries/{id}/edit
func (gc *GalleriesController) Edit(w http.ResponseWriter, req *http.Request) {
	gallery, err := gc.galleryByID(w, req)
	if err != nil {
		return
	}
	user := context.User(req.Context())
	if user.ID != gallery.UserID {
		http.Error(w, "Gallery Not Found", http.StatusNotFound)
		return
	}
	data := views.Data{Content: gallery}
	gc.EditView.Render(w, req, data)
}

// POST /galleries/{id}/update
func (gc *GalleriesController) Update(w http.ResponseWriter, req *http.Request) {
	gallery, err := gc.galleryByID(w, req)
	if err != nil {
		return
	}
	user := context.User(req.Context())
	if user.ID != gallery.UserID {
		http.Error(w, "Gallery Not Found", http.StatusNotFound)
		return
	}
	
	data := views.Data{Content: gallery}
	var form forms.GalleryForm
	if err := forms.ParseForm(req, &form); err != nil {
		log.Println(err)
		data.SetAlert(err)
		gc.EditView.Render(w, req, data)
		return
	}
	
	gallery.Title = form.Title
	err = gc.gs.Update(gallery)
	if err != nil {
		data.SetAlert(err)
		gc.EditView.Render(w, req, data)
		return
	}
	http.Redirect(w, req, "/gallery", http.StatusSeeOther)
}

// POST /galleries/{id}/images
func (gc *GalleriesController) UploadImages(w http.ResponseWriter, req *http.Request) {
	gallery, err := gc.galleryByID(w, req)
	if err != nil {
		return
	}
	user := context.User(req.Context())
	if user.ID != gallery.UserID {
		http.Error(w, "Gallery Not Found", http.StatusNotFound)
		return
	}
	
	data := views.Data{Content: gallery}
	err = req.ParseMultipartForm(maxMemoryLimit)
	if err != nil {
		data.SetAlert(err)
		gc.EditView.Render(w, req, data)
		return
	}
	
	files := req.MultipartForm.File["images"]
	galleryImagesPath := fmt.Sprintf("%v%v/", galleryPath, gallery.ID)
	for _, f := range files {
		file, err := f.Open()
		if err != nil {
			data.SetAlert(err)
			gc.EditView.Render(w, req, data)
			return
		}
		defer file.Close()
		
		err = gc.is.Create(galleryImagesPath, f.Filename, file)
		if err != nil {
			data.SetAlert(err)
			gc.EditView.Render(w, req, data)
			return
		}
	}
	fmt.Fprintln(w, "Uploaded Images on "+galleryImagesPath)
}

// POST /galleries/{id}/delete
func (gc *GalleriesController) Delete(w http.ResponseWriter, req *http.Request) {
	gallery, err := gc.galleryByID(w, req)
	if err != nil {
		return
	}
	user := context.User(req.Context())
	if user.ID != gallery.UserID {
		http.Error(w, "Gallery Not Found", http.StatusNotFound)
		return
	}
	
	data := views.Data{Content: gallery}
	err = gc.gs.Delete(gallery.ID)
	if err != nil {
		data.SetAlert(err)
		gc.EditView.Render(w, req, data)
		return
	}
	http.Redirect(w, req, "/gallery", http.StatusSeeOther)
}

func (gc *GalleriesController) galleryByID(w http.ResponseWriter, req *http.Request) (*models.Gallery, error) {
	params := mux.Vars(req)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid Gallery ID", http.StatusBadRequest)
		return nil, err
	}
	
	gallery, err := gc.gs.ByID(uint(id))
	if err != nil {
		switch err {
		case models.ErrNotFound:
			http.Error(w, "Gallery Not Found", http.StatusNotFound)
		default:
			http.Error(w, "Server Error", http.StatusInternalServerError)
		}
		return nil, err
	}
	return gallery, nil
}
