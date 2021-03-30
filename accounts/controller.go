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
		SignInView: views.NewView("base", "accounts/signin"),
		us: us,
	}
}

type Controller struct {
	SignUpView *views.View
	SignInView *views.View
	us *Service
}

func (uc *Controller) SignUp(w http.ResponseWriter, req *http.Request) {
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

func (uc *Controller) SignIn(w http.ResponseWriter, req *http.Request) {
	var form SignInForm
	utils.Must(utils.ParseForm(req, &form))

	user, err := uc.us.Authenticate(form.Email, form.Password)
	switch err {
	case ErrNotFound:
		http.Error(w, "Email Address is incorrect", http.StatusBadRequest)
	case ErrInvalidPassword:
		http.Error(w, "Password is incorrect", http.StatusBadRequest)
	case nil:
		fmt.Fprintln(w, user)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}