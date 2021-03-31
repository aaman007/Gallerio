package main

import (
	"fmt"
	"gallerio/accounts"
	"gallerio/core"
	"gallerio/galleries"
	"gallerio/middlewares"
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
	galleriesController := galleries.NewGalleryController(services.Gallery)
	coreController := core.NewStaticController()

	loginRequiredMw := middlewares.LoginRequired{
		UserService: services.User,
	}

	router := mux.NewRouter()

	// Static Routes
	router.Handle("/", coreController.HomeView).Methods("GET")
	router.Handle("/contact", coreController.ContactView).Methods("GET")

	// Accounts Routes
	router.Handle("/signin", usersController.SignInView).Methods("GET")
	router.HandleFunc("/signin", usersController.SignIn).Methods("POST")
	router.Handle("/signup", usersController.SignUpView).Methods("GET")
	router.HandleFunc("/signup", usersController.SignUp).Methods("POST")

	// Galleries Routes
	router.Handle("/galleries/new",
		loginRequiredMw.Apply(galleriesController.New)).Methods("GET")
	router.HandleFunc("/galleries",
		loginRequiredMw.ApplyFunc(galleriesController.Create)).Methods("POST")

	http.ListenAndServe(":8000", router)
}
