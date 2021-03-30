package accounts

import (
	errorUtil "go-web-dev-2/utils/error"
	formUtil "go-web-dev-2/utils/form"
	"go-web-dev-2/utils/rand"
	"go-web-dev-2/views"
	"net/http"
)


func NewUserController(us *UserService) *UserController {
	return &UserController{
		SignUpView: views.NewView("base", "accounts/signup"),
		SignInView: views.NewView("base", "accounts/signin"),
		us: us,
	}
}

type UserController struct {
	SignUpView *views.View
	SignInView *views.View
	us *UserService
}

func (uc *UserController) SignUp(w http.ResponseWriter, req *http.Request) {
	var form SignUpForm
	errorUtil.Must(formUtil.ParseForm(req, &form))

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
	if err := uc.signInUser(w, &user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func (uc *UserController) SignIn(w http.ResponseWriter, req *http.Request) {
	var form SignInForm
	errorUtil.Must(formUtil.ParseForm(req, &form))

	user, err := uc.us.Authenticate(form.Email, form.Password)
	if err != nil {
		switch err {
		case ErrNotFound:
			http.Error(w, "Email Address is incorrect", http.StatusBadRequest)
		case ErrInvalidPassword:
			http.Error(w, "Password is incorrect", http.StatusBadRequest)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if err := uc.signInUser(w, user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func (uc *UserController) signInUser(w http.ResponseWriter, user *User) error {
	if user.RememberToken == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		if err = uc.us.Update(user); err != nil {
			return err
		}
		user.RememberToken = token
	}
	cookie := &http.Cookie{
		Name: "remember_token",
		Value: user.RememberToken,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	return nil
}