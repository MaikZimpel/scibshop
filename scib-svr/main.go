package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
	"net/http"
	"os"
	"scib-svr/configuration"
	"scib-svr/datastore"
	"scib-svr/inventory"
	"scib-svr/logging"
	"scib-svr/shopping"
)

func main() {
	router := httprouter.New()

	inventoryService := inventory.NewController(inventory.NewService(datastore.New()))

	// inventory routes
	router.GET("/", defaultHandler)
	router.GET(makeUri(inventory.RequestUri, nil), inventoryService.Get)
	router.GET(makeUri(inventory.RequestUri, []string{"id"}), inventoryService.GetById)
	router.POST(makeUri(inventory.RequestUri, nil), inventoryService.Create)
	router.PUT(makeUri(inventory.RequestUri, []string{"id"}), inventoryService.Update)
	router.POST(makeUri(inventory.RequestUri + "/images/upload", nil), inventoryService.SignedUrl)

	// shop routes
	/*router.GET(makeUri(shopping.RequestUri, nil), shopping.Get)
	router.GET(makeUri(shopping.RequestUri, []string{"id"}), shopping.GetById)
	router.POST(makeUri(shopping.RequestUri, nil), shopping.Post)*/
	router.PUT(makeUri(shopping.RequestUri, []string{"id"}), shopping.Put)

	port := os.Getenv("PORT")
	log := logging.New()
	if port == "" {
		port = "8082"
		log.Info(context.Background(),"Defaulting to port %s", port)
	}

	log.Info(context.Background(),"Listening on port %s", port)
	log.Info(context.Background(),"Open http://localhost:%s in the browser", port)
	log.Info(context.Background(), "credentials file in %s", os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	log.Info(context.Background(), "db config %s", fmt.Sprintf("%s:%d/%s",configuration.MongoDbHost,
		configuration.MongoDbPort, configuration.MongoDbDatabase))
	log.Critical(context.Background(), "%s", http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}

func defaultHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	_, _ = fmt.Fprintf(w, "<H1>Welcome to SCIB</H1>")
}

func makeUri(uri string, params []string) string {
	var paramStr string
	if params != nil {
		for _, param := range params {
			paramStr += "/:" + param
		}
	}
	return fmt.Sprintf("/%s%s",uri, paramStr)
}