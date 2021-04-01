package main

import (
	"fmt"
	"log"
	"net/http"
	
	"gallerio/controllers"
	"gallerio/middlewares"
	"gallerio/models"
	
	"github.com/gorilla/mux"
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
	services, err := models.NewServices(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer services.Close()
	services.AutoMigrate()

	router := mux.NewRouter()
	usersController := controllers.NewUsersController(services.User)
	galleriesController := controllers.NewGalleriesController(services.Gallery, services.Image, router)
	coreController := controllers.NewStaticController()

	assignUserMw := middlewares.AssignUser{
		UserService: services.User,
	}
	loginRequiredMw := middlewares.LoginRequired{
		UserService: services.User,
	}

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
		loginRequiredMw.ApplyFunc(galleriesController.Index)).Methods("GET")
	router.HandleFunc("/galleries",
		loginRequiredMw.ApplyFunc(galleriesController.Create)).Methods("POST")
	router.HandleFunc("/galleries/{id:[0-9]+}",
		galleriesController.Show).Methods("GET").Name(controllers.ShowGalleryName)
	router.HandleFunc("/galleries/{id:[0-9]+}/edit",
		loginRequiredMw.ApplyFunc(galleriesController.Edit)).
		Methods("GET").Name(controllers.EditGalleryName)
	router.HandleFunc("/galleries/{id:[0-9]+}/update",
		loginRequiredMw.ApplyFunc(galleriesController.Update)).Methods("POST")
	router.HandleFunc("/galleries/{id:[0-9]+}/delete",
		loginRequiredMw.ApplyFunc(galleriesController.Delete)).Methods("POST")
	router.HandleFunc("/galleries/{id:[0-9]+}/images",
		loginRequiredMw.ApplyFunc(galleriesController.UploadImages)).Methods("POST")

	log.Fatal(http.ListenAndServe(":8005", assignUserMw.Apply(router)))
}
