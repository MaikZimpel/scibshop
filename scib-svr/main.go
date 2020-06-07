package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
	"scib-svr/inventory"
	"scib-svr/shopping"
)

func main() {
	router := httprouter.New()

	// inventory routes
	router.GET("/", defaultHandler)
	router.GET(makeUri(inventory.RequestUri, nil), inventory.Get)
	router.GET(makeUri(inventory.RequestUri, []string{"id"}), inventory.GetById)
	router.POST(makeUri(inventory.RequestUri, nil), inventory.Create)
	router.PUT(makeUri(inventory.RequestUri, []string{"id"}), inventory.Update)

	// shop routes
	router.GET(makeUri(shopping.RequestUri, nil), shopping.Get)
	router.GET(makeUri(shopping.RequestUri, []string{"id"}), shopping.GetById)
	router.POST(makeUri(shopping.RequestUri, nil), shopping.Post)
	router.PUT(makeUri(shopping.RequestUri, []string{"id"}), shopping.Put)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Printf("Open http://localhost:%s in the browser", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
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