package main

import (
	"encoding/json"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"golang.org/x/net/context"
	"io/ioutil"
	"net/http"
	"os"
	"scib-svr/auth"
	"scib-svr/configuration"
	"scib-svr/crm"
	"scib-svr/filestore"
	"scib-svr/inventory"
	"scib-svr/logging"
)

func main() {
	logger, router := configure()
	inventoryController := inventory.NewController(inventory.NewService(filestore.New(), logger), logger)
	customerController := crm.NewController(crm.NewService(logger), logger)
	authService := auth.NewService(logger)
	_ = configureUsers(authService, logger)
	authController := auth.NewController(authService, logger)

	// inventory routes
	router.GET("/", authController.Authenticate(defaultHandler))
	router.GET(inventory.RequestUri, inventoryController.Get)
	router.GET(inventory.RequestUri+"/:id", inventoryController.GetById)
	router.POST(inventory.RequestUri, authController.Authenticate(inventoryController.Create))
	router.PUT(inventory.RequestUri+"/:id", authController.Authenticate(inventoryController.Update))
	router.DELETE(inventory.RequestUri+"/:id", authController.Authenticate(inventoryController.Delete))
	router.POST(inventory.RequestUri+"/:id/images", authController.Authenticate(inventoryController.UploadImage))
	router.GET(inventory.RequestUri+"/:id/images/:imageId", inventoryController.GetImage)
	router.DELETE(inventory.RequestUri+"/:id/images/:imageId", authController.Authenticate(inventoryController.DeleteImage))

	// shop routes
	/*router.GET(makeUri(shopping.RequestUri, nil), shopping.Get)
	router.GET(makeUri(shopping.RequestUri, []string{"id"}), shopping.GetById)
	router.POST(makeUri(shopping.RequestUri, nil), shopping.Post)
	router.PUT(makeUri(shopping.RequestUri, []string{"id"}), shopping.Put)*/

	// crm routes
	router.GET(crm.RequestUri, authController.Authenticate(customerController.Get))
	router.POST(crm.RequestUri, authController.Authenticate(customerController.CreateOrUpdate))
	router.PUT(crm.RequestUri, authController.Authenticate(customerController.CreateOrUpdate))
	router.DELETE(crm.RequestUri, authController.Authenticate(customerController.Delete))

	// auth routes
	router.POST("/auth/token", authController.SignIn)

	// user routes
	router.GET(auth.RequestUri, authController.Authenticate(authController.Get))

	port := os.Getenv("SERVER_PORT")

	if port == "" {
		port = "8082"
		logger.Info(context.Background(), "Defaulting to port %s", port)
	}

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:     []string{"http://192.168.178.35:*", "http://localhost:*"},
		AllowedMethods:     []string{"OPTIONS", "HEAD", "GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:     []string{"Content-Type", "Accept", "Access-Control-Allow-Origin, Authorization"},
		AllowCredentials:   true,
		OptionsPassthrough: false,
		Debug:              false,
	}).Handler(router)

	logger.Info(context.Background(), "Listening on port %s", port)
	logger.Info(context.Background(), "Open http://localhost:%s in the browser", port)
	logger.Info(context.Background(), "db config %s", fmt.Sprintf("%s:%d/%s", configuration.MongoDbHost,
		configuration.MongoDbPort, configuration.MongoDbDatabase))
	logger.Critical(context.Background(), "%s", http.ListenAndServe(fmt.Sprintf(":%s", port), corsHandler))
}

func configureUsers(service *auth.Service, logger logging.Logger) error {
	usersFile, e := os.Open("users.json")
	if e != nil {
		return e
	}
	defer usersFile.Close()
	bytes, e := ioutil.ReadAll(usersFile)
	if e != nil {
		return e
	}
	var usersArray []auth.User
	e = json.Unmarshal(bytes, &usersArray)
	if e != nil {
		return e
	}
	for _, user := range usersArray {
		_, err := service.Save(context.Background(), &user)
		if err != nil {
			logger.Critical(context.Background(), "an error occurred when trying to save user: %v", err)
		}
	}
	return nil
}

func defaultHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	_, _ = fmt.Fprintf(w, "<H1>Welcome to SCIB</H1>")
}

func configure() (logger logging.Logger, router *httprouter.Router) {
	logger = logging.New()
	var cfg configuration.Config
	err := cleanenv.ReadConfig("config.yml", &cfg)
	if err != nil {
		panicConfig(logger, err)
	}
	err = os.Setenv("SERVER_PORT", cfg.Server.Port)
	err = os.Setenv("SECURITY_CLIENT_SECRET", cfg.Security.ClientSecret)
	if err != nil {
		panicConfig(logger, err)
	}
	router = httprouter.New()

	return
}

func panicConfig(logger logging.Logger, err error) {
	logger.Critical(context.Background(), "configuration problem %v", err)
	panic(err)
}

//
