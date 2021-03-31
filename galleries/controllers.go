package galleries

import (
	"fmt"
	"gallerio/utils/context"
	"gallerio/utils/forms"
	"gallerio/views"
	"log"
	"net/http"
)


func NewGalleryController(gs GalleryService) *GalleryController {
	return &GalleryController{
		New: views.NewView("base", "galleries/new"),
		gs: gs,
	}
}

type GalleryController struct {
	New *views.View
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
	fmt.Fprintln(w, gallery)
	// http.Redirect(w, req, "/", http.StatusSeeOther)
}
