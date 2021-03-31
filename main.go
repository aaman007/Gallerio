package main

import (
	"fmt"
	"gallerio/accounts"
	"gallerio/core"
	"github.com/gorilla/mux"
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

	services, err := core.NewServices(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer services.Close()
	services.AutoMigrate()

	usersController := accounts.NewUserController(services.User)
	coreController := core.NewStaticController()

	router := mux.NewRouter()
	router.Handle("/", coreController.HomeView).Methods("GET")
	router.Handle("/contact", coreController.ContactView).Methods("GET")
	router.Handle("/signup", usersController.SignUpView).Methods("GET")
	router.HandleFunc("/signup", usersController.SignUp).Methods("POST")
	router.Handle("/signin", usersController.SignInView).Methods("GET")
	router.HandleFunc("/signin", usersController.SignIn).Methods("POST")
	http.ListenAndServe(":8000", router)
}
