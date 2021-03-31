package controllers

import (
	forms2 "gallerio/forms"
	models2 "gallerio/models"
	"gallerio/utils/rand"
	"gallerio/views"
	"log"
	"net/http"
)


func NewUsersController(us models2.UserService) *UsersController {
	return &UsersController{
		SignUpView: views.NewView("base", "accounts/signup"),
		SignInView: views.NewView("base", "accounts/signin"),
		us: us,
	}
}

type UsersController struct {
	SignUpView *views.View
	SignInView *views.View
	us         models2.UserService
}

// POST /signup
func (uc *UsersController) SignUp(w http.ResponseWriter, req *http.Request) {
	var data views.Data
	var form forms2.SignUpForm
	if err := forms2.ParseForm(req, &form); err != nil {
		log.Println(err)
		data.SetAlert(err)
		uc.SignUpView.Render(w, req, data)
		return
	}

	user := models2.User{
		Name: form.Name,
		Username: form.Username,
		Email: form.Email,
		Password: form.Password,
	}
	if err := uc.us.Create(&user); err != nil {
		log.Println(err)
		data.SetAlert(err)
		uc.SignUpView.Render(w, req, data)
		return
	}
	if err := uc.signInUser(w, &user); err != nil {
		http.Redirect(w, req, "/signin", http.StatusSeeOther)
		return
	}
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

// POST /signin
func (uc *UsersController) SignIn(w http.ResponseWriter, req *http.Request) {
	var data views.Data
	var form forms2.SignInForm
	if err := forms2.ParseForm(req, &form); err != nil {
		log.Println(err)
		data.SetAlert(err)
		uc.SignInView.Render(w, req, data)
		return
	}

	user, err := uc.us.Authenticate(form.Email, form.Password)
	if err != nil {
		log.Println(err)
		switch err {
		case models2.ErrNotFound:
			data.AlertError("Email Address is incorrect")
		default:
			data.SetAlert(err)
		}
		uc.SignInView.Render(w, req, data)
		return
	}

	if err := uc.signInUser(w, user); err != nil {
		log.Println(err)
		uc.SignInView.Render(w, req, data)
		return
	}
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func (uc *UsersController) signInUser(w http.ResponseWriter, user *models2.User) error {
	if user.RememberToken == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.RememberToken = token
		if err = uc.us.Update(user); err != nil {
			return err
		}
	}
	cookie := &http.Cookie{
		Name: "remember_token",
		Value: user.RememberToken,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	return nil
}