package main

import (
	"flag"
	"fmt"
	"gallerio/configs"
	"gallerio/utils/email"
	"gallerio/utils/errors"
	"gallerio/utils/rand"
	"github.com/gorilla/csrf"
	"log"
	"net/http"
	
	"gallerio/controllers"
	"gallerio/middlewares"
	"gallerio/models"
	
	"github.com/gorilla/mux"
)


func main() {
	// To View list of flags
	// run: go build . && ./gallerio --help
	//
	// To run with prod flag
	// run: go build . && ./gallerio --prod
	boolPtr := flag.Bool("prod", false, "Provide this flag in production." +
		"This flag will ensure that a .config file is setup properly.")
	flag.Parse()
	
	cfg := configs.LoadConfig(*boolPtr)
	dbCfg := cfg.Database
	services, err := models.NewServices(
		models.WithGorm(dbCfg.Dialect(), dbCfg.ConnectionInfo()),
		models.WithLogMode(false),
		models.WithUser(cfg.Pepper, cfg.HMACKey),
		models.WithGallery(),
		models.WithImage(),
	)
	if err != nil {
		panic(err)
	}
	defer services.Close()
	services.AutoMigrate()

	mgCfg := cfg.Mailgun
	emailer := email.NewClient(
		email.WithSender("Gallerio Support",
		"support@sandboxfa300beae3034442af3cdd253f03c0c1.mailgun.org"),
		email.WithMailgun(mgCfg.Domain, mgCfg.PublicAPIKey),
	)

	router := mux.NewRouter()
	usersController := controllers.NewUsersController(services.User, emailer)
	galleriesController := controllers.NewGalleriesController(services.Gallery, services.Image, router)
	coreController := controllers.NewStaticController()
	
	b, err := rand.Bytes(32)
	errors.Must(err)
	csrfMw := csrf.Protect(b, csrf.Secure(cfg.IsProduction()))
	assignUserMw := middlewares.AssignUser{
		UserService: services.User,
	}
	loginRequiredMw := middlewares.LoginRequired{
		UserService: services.User,
	}
	alreadyLoggedInMw := middlewares.AlreadyLoggedIn{
		UserService: services.User,
	}

	// Static Routes
	router.Handle("/", coreController.HomeView).Methods("GET")
	router.Handle("/contact", coreController.ContactView).Methods("GET")

	// Accounts Routes
	router.Handle("/signin",
		alreadyLoggedInMw.Apply(usersController.SignInView)).Methods("GET")
	router.HandleFunc("/signin",
		alreadyLoggedInMw.ApplyFunc(usersController.SignIn)).Methods("POST")
	router.HandleFunc("/signup",
		alreadyLoggedInMw.ApplyFunc(usersController.New)).Methods("GET")
	router.HandleFunc("/signup",
		alreadyLoggedInMw.ApplyFunc(usersController.SignUp)).Methods("POST")
	router.HandleFunc("/signout",
		loginRequiredMw.ApplyFunc(usersController.SignOut)).Methods("POST")

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
		loginRequiredMw.ApplyFunc(galleriesController.UploadImage)).Methods("POST")
	router.HandleFunc("/galleries/{id:[0-9]+}/images/{filename}/delete",
		loginRequiredMw.ApplyFunc(galleriesController.DeleteImage)).Methods("POST")
	
	// Media Routes
	mediaHandler := http.FileServer(http.Dir("./media/"))
	router.PathPrefix("/media/").Handler(http.StripPrefix("/media/", mediaHandler))
	
	// Static Routes
	staticHandler := http.FileServer(http.Dir("./static/"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", staticHandler))
	
	fmt.Printf("Starting server on Port : %v\n", cfg.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", cfg.Port), csrfMw(assignUserMw.Apply(router))))
}
