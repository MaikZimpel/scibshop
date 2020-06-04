package inventory

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"scib-svr/helpers"
	"strconv"
)

const(
	REQUEST_URI = "inventory"
)

func Get(w http.ResponseWriter, r *http.Request, _ httprouter.Params)  {
	queryValues := r.URL.Query()
	stockableOnly, _ := strconv.ParseBool(queryValues.Get("stockableOnly"))
	fmt.Fprint(w, allItems(stockableOnly))
}

func GetById(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	item, err := itemById(ps.ByName("id"))
	if err != nil {
		qErr := err.(helpers.QueryError)
		http.Error(w, qErr.Message, qErr.Code)
	} else {
		_, _ = fmt.Fprintf(w, item.String())
	}
}

func Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params)  {
	var i Item
	err := json.NewDecoder(r.Body).Decode(&i)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		i.Id = uuid.New().String()
		createdItemId, addErr := addItem(i)
		if addErr != nil {
			http.Error(w, addErr.Error(), http.StatusInternalServerError)
		} else {
			_, _ = fmt.Fprintf(w, "http://" + r.Host + r.URL.String() + "/" + createdItemId)
		}
	}
}