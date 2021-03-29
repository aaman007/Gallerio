package main

import (
	"github.com/gorilla/mux"
	"go-web-dev-2/utils"
	"go-web-dev-2/views"
	"net/http"
)

var (
	homeView *views.View
	contactView *views.View
	signupView *views.View
)

func HomeHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	utils.Must(homeView.Render(w, nil))
}

func ContactHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	utils.Must(contactView.Render(w, nil))
}

func SignUpHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	utils.Must(signupView.Render(w, nil))
}

func main() {
	homeView = views.NewView("base", "views/core/home.gohtml")
	contactView = views.NewView("base", "views/core/contact.gohtml")
	signupView = views.NewView("base", "views/accounts/signup.gohtml")

	router := mux.NewRouter()
	router.HandleFunc("/", HomeHandler)
	router.HandleFunc("/contact", ContactHandler)
	router.HandleFunc("/signup", SignUpHandler)
	http.ListenAndServe(":8000", router)
}
