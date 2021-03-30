package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"go-web-dev-2/accounts"
	"go-web-dev-2/core"
	"net/http"
)

const (
	dbHost = "localhost"
	dbPort = 5432
	dbUser = "robert"
	dbPassword = "password"
	dbName = "gallerio"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName,
	)

	us, err := accounts.NewService(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer us.Close()
	us.DestructiveReset()

	usersController := accounts.NewController(us)
	coreController := core.NewController()

	router := mux.NewRouter()
	router.Handle("/", coreController.HomeView).Methods("GET")
	router.Handle("/contact", coreController.ContactView).Methods("GET")
	router.HandleFunc("/signup", usersController.SignUp).Methods("GET")
	router.HandleFunc("/signup", usersController.Create).Methods("POST")
	http.ListenAndServe(":8000", router)
}
