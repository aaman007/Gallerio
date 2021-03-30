package accounts

import (
	"fmt"
	"go-web-dev-2/utils"
	"go-web-dev-2/views"
	"net/http"
)


func NewController(us *Service) *Controller {
	return &Controller{
		SignUpView: views.NewView("base", "accounts/signup"),
		us: us,
	}
}

type Controller struct {
	SignUpView *views.View
	us *Service
}

func (uc *Controller) SignUp(w http.ResponseWriter, req *http.Request) {
	utils.Must(uc.SignUpView.Render(w, nil))
}

func (uc *Controller) Create(w http.ResponseWriter, req *http.Request) {
	var form SignUpForm
	utils.Must(utils.ParseForm(req, &form))

	user := User{
		Name: form.Name,
		Username: form.Username,
		Email: form.Email,
		Password: form.Password,
	}
	if err := uc.us.Create(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, user)
}