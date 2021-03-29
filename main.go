package main

import (
	"github.com/gorilla/mux"
	"go-web-dev-2/accounts"
	"go-web-dev-2/core"
	"net/http"
)

func main() {
	usersController := accounts.NewController()
	coreController := core.NewController()

	router := mux.NewRouter()
	router.Handle("/", coreController.HomeView).Methods("GET")
	router.Handle("/contact", coreController.ContactView).Methods("GET")
	router.HandleFunc("/signup", usersController.SignUp).Methods("GET")
	router.HandleFunc("/signup", usersController.Create).Methods("POST")
	http.ListenAndServe(":8000", router)
}
