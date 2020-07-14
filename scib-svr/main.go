package main

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"golang.org/x/net/context"
	"net/http"
	"os"
	"scib-svr/auth"
	"scib-svr/configuration"
	"scib-svr/filestore"
	"scib-svr/inventory"
	"scib-svr/logging"
)

func main() {
	logger, router := configure()
	inventoryController := inventory.NewController(
		inventory.NewService(filestore.New(), logger),
		logger)


	
	// inventory routes
	router.GET("/", auth.WithAuthentication(defaultHandler))
	router.GET(inventory.RequestUri, inventoryController.Get)
	router.GET(inventory.RequestUri + "/:id", inventoryController.GetById)
	router.POST(inventory.RequestUri, auth.WithAuthentication(inventoryController.Create))
	router.PUT(inventory.RequestUri + "/:id", auth.WithAuthentication(inventoryController.Update))
	router.DELETE(inventory.RequestUri + "/:id", auth.WithAuthentication(inventoryController.Delete))
	router.POST(inventory.RequestUri + "/:id/images", auth.WithAuthentication(inventoryController.UploadImage))
	router.GET(inventory.RequestUri + "/:id/images/:imageId", inventoryController.GetImage)
	router.DELETE(inventory.RequestUri + "/:id/images/:imageId", auth.WithAuthentication(inventoryController.DeleteImage))


	// shop routes
	/*router.GET(makeUri(shopping.RequestUri, nil), shopping.Get)
	router.GET(makeUri(shopping.RequestUri, []string{"id"}), shopping.GetById)
	router.POST(makeUri(shopping.RequestUri, nil), shopping.Post)
	router.PUT(makeUri(shopping.RequestUri, []string{"id"}), shopping.Put)*/

	port := os.Getenv("SERVER_PORT")

	if port == "" {
		port = "8082"
		logger.Info(context.Background(),"Defaulting to port %s", port)
	}

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:         []string{"http://192.168.178.35:*","http://localhost:*"},
		AllowedMethods:         []string{"OPTIONS","HEAD", "GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:         []string{"Content-Type", "Accept", "Access-Control-Allow-Origin, Authorization"},
		AllowCredentials:       true,
		OptionsPassthrough:     false,
		Debug:                  false,
	}).Handler(router)

	logger.Info(context.Background(),"Listening on port %s", port)
	logger.Info(context.Background(),"Open http://localhost:%s in the browser", port)
	logger.Info(context.Background(), "db config %s", fmt.Sprintf("%s:%d/%s",configuration.MongoDbHost,
		configuration.MongoDbPort, configuration.MongoDbDatabase))
	logger.Critical(context.Background(), "%s", http.ListenAndServe(fmt.Sprintf(":%s", port), corsHandler))
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