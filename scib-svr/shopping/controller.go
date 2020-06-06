package shopping

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

const (
	RequestUri = "shop"
)

func AddToCart(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	
}