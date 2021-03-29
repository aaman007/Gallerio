package accounts

import (
	"fmt"
	"go-web-dev-2/utils"
	"go-web-dev-2/views"
	"net/http"
)


func NewController() *Controller {
	return &Controller{
		SignUpView: views.NewView("base", "accounts/signup"),
	}
}

type Controller struct {
	SignUpView *views.View
}

func (uc *Controller) SignUp(w http.ResponseWriter, req *http.Request) {
	utils.Must(uc.SignUpView.Render(w, nil))
}

func (uc *Controller) Create(w http.ResponseWriter, req *http.Request) {
	var form SignUpForm
	utils.Must(utils.ParseForm(req, &form))
	fmt.Fprintln(w, form)
}