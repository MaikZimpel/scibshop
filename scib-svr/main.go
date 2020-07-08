package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"golang.org/x/net/context"
	"net/http"
	"os"
	"scib-svr/configuration"
	"scib-svr/filestore"
	"scib-svr/inventory"
	"scib-svr/logging"
)

func main() {
	router := httprouter.New()
	log := logging.New()
	inventoryController := inventory.NewController(
		inventory.NewService(filestore.New(), log),
		log)
	
	// inventory routes
	router.GET("/", defaultHandler)
	router.GET(inventory.RequestUri, inventoryController.Get)
	router.GET(inventory.RequestUri + "/:id", inventoryController.GetById)
	router.POST(inventory.RequestUri, inventoryController.Create)
	router.PUT(inventory.RequestUri + "/:id", inventoryController.Update)
	router.DELETE(inventory.RequestUri + "/:id", inventoryController.Delete)
	router.POST(inventory.RequestUri + "/:id/images", inventoryController.UploadImage)
	router.GET(inventory.RequestUri + "/:id/images/:imageId", inventoryController.GetImage)
	router.DELETE(inventory.RequestUri + "/:id/images/:imageId", inventoryController.DeleteImage)


	// shop routes
	/*router.GET(makeUri(shopping.RequestUri, nil), shopping.Get)
	router.GET(makeUri(shopping.RequestUri, []string{"id"}), shopping.GetById)
	router.POST(makeUri(shopping.RequestUri, nil), shopping.Post)
	router.PUT(makeUri(shopping.RequestUri, []string{"id"}), shopping.Put)*/

	port := os.Getenv("PORT")

	if port == "" {
		port = "8082"
		log.Info(context.Background(),"Defaulting to port %s", port)
	}

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:         []string{"http://localhost:*"},
		AllowedMethods:         []string{"OPTIONS","HEAD", "GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:         []string{"Content-Type", "Accept", "Access-Control-Allow-Origin, Authorization"},
		AllowCredentials:       true,
		OptionsPassthrough:     false,
		Debug:                  false,
	}).Handler(router)

	log.Info(context.Background(),"Listening on port %s", port)
	log.Info(context.Background(),"Open http://localhost:%s in the browser", port)
	log.Info(context.Background(), "credentials file in %s", os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	log.Info(context.Background(), "db config %s", fmt.Sprintf("%s:%d/%s",configuration.MongoDbHost,
		configuration.MongoDbPort, configuration.MongoDbDatabase))
	log.Critical(context.Background(), "%s", http.ListenAndServe(fmt.Sprintf(":%s", port), corsHandler))
}

func defaultHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	_, _ = fmt.Fprintf(w, "<H1>Welcome to SCIB</H1>")
}

//