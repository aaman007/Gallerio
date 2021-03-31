package galleries

import (
	"fmt"
	"gallerio/utils/context"
	"gallerio/utils/forms"
	"gallerio/utils/models"
	"gallerio/views"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

var (
	ShowGalleryName = "show_gallery"
)

func NewGalleryController(gs GalleryService, router *mux.Router) *GalleryController {
	return &GalleryController{
		New: views.NewView("base", "galleries/new"),
		ShowView: views.NewView("base", "galleries/show"),
		EditView: views.NewView("base", "galleries/edit"),
		router: router,
		gs: gs,
	}
}

type GalleryController struct {
	New *views.View
	ShowView *views.View
	EditView *views.View
	router *mux.Router
	gs GalleryService
}

// POST /galleries
func (gc *GalleryController) Create(w http.ResponseWriter, req *http.Request) {
	var data views.Data
	var form GalleryForm
	if err := forms.ParseForm(req, &form); err != nil {
		log.Println(err)
		data.SetAlert(err)
		gc.New.Render(w, data)
		return
	}

	user := context.User(req.Context())

	gallery := Gallery{
		Title: form.Title,
		UserID: user.ID,
	}
	if err := gc.gs.Create(&gallery); err != nil {
		log.Println(err)
		data.SetAlert(err)
		gc.New.Render(w, data)
		return
	}
	url, err := gc.router.Get(ShowGalleryName).URL("id", fmt.Sprintf("%v", gallery.ID))
	if err != nil {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	http.Redirect(w, req, url.Path, http.StatusSeeOther)
}

// GET /galleries/{id}
func (gc *GalleryController) Show(w http.ResponseWriter, req *http.Request) {
	gallery, err := gc.galleryByID(w, req)
	if err != nil {
		return
	}
	data := views.Data{Content: gallery}
	gc.ShowView.Render(w, data)
}

// GET /galleries/{id}/edit
func (gc *GalleryController) Edit(w http.ResponseWriter, req *http.Request) {
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
	gc.EditView.Render(w, data)
}

// POST /galleries/{id}/update
func (gc *GalleryController) Update(w http.ResponseWriter, req *http.Request) {
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
	var form GalleryForm
	if err := forms.ParseForm(req, &form); err != nil {
		log.Println(err)
		data.SetAlert(err)
		gc.EditView.Render(w, data)
		return
	}

	gallery.Title = form.Title
	err = gc.gs.Update(gallery)
	if err != nil {
		data.SetAlert(err)
		gc.EditView.Render(w, data)
		return
	}
	data.Alert = &views.Alert{
		Level: views.AlertLevelSuccess,
		Message: "Successfully updated gallery",
	}
	gc.EditView.Render(w, data)
}

// POST /galleries/{id}/delete
func (gc *GalleryController) Delete(w http.ResponseWriter, req *http.Request) {
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
		gc.EditView.Render(w, data)
		return
	}
	fmt.Fprintln(w, "deleted")
}

func (gc *GalleryController) galleryByID(w http.ResponseWriter, req *http.Request) (*Gallery, error) {
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
