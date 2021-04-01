package controllers

import (
	"gallerio/forms"
	"gallerio/models"
	"gallerio/utils/context"
	"gallerio/utils/email"
	"gallerio/utils/rand"
	"gallerio/views"
	"log"
	"net/http"
	"time"
)

func NewUsersController(us models.UserService, mg email.Client) *UsersController {
	return &UsersController{
		SignUpView: views.NewView("base", "user/signup"),
		SignInView: views.NewView("base", "user/signin"),
		us:         us,
		mg:         mg,
	}
}

type UsersController struct {
	SignUpView *views.View
	SignInView *views.View
	us         models.UserService
	mg         email.Client
}

// GET /signup
func (uc *UsersController) New(w http.ResponseWriter, req *http.Request) {
	var form forms.SignUpForm
	_ = forms.ParseURLParams(req, &form)
	uc.SignUpView.Render(w, req, form)
}

// POST /signup
func (uc *UsersController) SignUp(w http.ResponseWriter, req *http.Request) {
	var data views.Data
	var form forms.SignUpForm
	data.Content = &form
	if err := forms.ParseForm(req, &form); err != nil {
		log.Println(err)
		data.SetAlert(err)
		uc.SignUpView.Render(w, req, data)
		return
	}

	user := models.User{
		Name:     form.Name,
		Username: form.Username,
		Email:    form.Email,
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
	go uc.mg.Welcome(user.Name, user.Email)
	alert := views.Alert{
		Level:   views.AlertLevelSuccess,
		Message: "Welcome to Gallerio",
	}
	views.RedirectAlert(w, req, "/galleries", http.StatusSeeOther, alert)
}

// POST /signin
func (uc *UsersController) SignIn(w http.ResponseWriter, req *http.Request) {
	var data views.Data
	var form forms.SignInForm
	if err := forms.ParseForm(req, &form); err != nil {
		log.Println(err)
		data.SetAlert(err)
		uc.SignInView.Render(w, req, data)
		return
	}

	user, err := uc.us.Authenticate(form.Email, form.Password)
	if err != nil {
		log.Println(err)
		switch err {
		case models.ErrNotFound:
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
	http.Redirect(w, req, "/galleries", http.StatusSeeOther)
}

// POST /signout
func (uc *UsersController) SignOut(w http.ResponseWriter, req *http.Request) {
	cookie := &http.Cookie{
		Name:     "remember_token",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

	user := context.User(req.Context())
	token, _ := rand.RememberToken()
	user.RememberToken = token
	_ = uc.us.Update(user)

	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func (uc *UsersController) signInUser(w http.ResponseWriter, user *models.User) error {
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
		Name:     "remember_token",
		Value:    user.RememberToken,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	return nil
}
