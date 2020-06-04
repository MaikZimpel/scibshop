package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
	"scib-svr/inventory"
)

func main() {
	router := httprouter.New()

	router.GET("/", defaultHandler)
	router.GET("/" + inventory.REQUEST_URI + "/", inventory.Get)
	router.GET("/" + inventory.REQUEST_URI + "/:id", inventory.GetById)
	router.POST("/" + inventory.REQUEST_URI, inventory.Create)

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