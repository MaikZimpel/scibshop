package shopping

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"scib-svr/helpers"
	"strconv"
)

const (
	RequestUri = "store"
)

func Get(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	helpers.EnableCors(&w, "http://localhost:3000")
	queryValues := r.URL.Query()
	availableOnly, _ := strconv.ParseBool(queryValues.Get("availableOnly"))
	_, _ = fmt.Fprint(w, allItems(availableOnly))
	println("called all items")
}

func GetById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	item, err := itemById(ps.ByName("id"))
	if err != nil {
		switch e := err.(type) {
		case helpers.QueryError:
			http.Error(w, e.Message, e.Code)
		default:
			http.Error(w, e.Error(), http.StatusInternalServerError)
		}
	} else {
		_, _ = fmt.Fprintf(w, item.String())
	}
}

func Post(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	storeItem, err := bodyToStoreItem(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	si, err, _ := insert(&storeItem)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_, _ = fmt.Fprintf(w, "http://"+r.Host+r.URL.String()+"/"+si.Sku)

}

func Put(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {


}

func bodyToCart(r *http.Request) (Cart, error) {
	var c Cart
	err := json.NewDecoder(r.Body).Decode(&c)
	return c, err
}

func bodyToStoreItem(r *http.Request) (StoreItem, error) {
	var si StoreItem
	err := json.NewDecoder(r.Body).Decode(&si)
	return si, err
}
