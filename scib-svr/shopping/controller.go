package shopping

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

const (
	RequestUri = "store"
)



func Put(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {


}

func bodyToCart(r *http.Request) (Cart, error) {
	var c Cart
	err := json.NewDecoder(r.Body).Decode(&c)
	return c, err
}
