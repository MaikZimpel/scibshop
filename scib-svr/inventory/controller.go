package inventory

import (
	"encoding/json"
	"fmt"
	"github.com/google/martian/log"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"scib-svr/helpers"
	"strconv"
)

const (
	RequestUri = "inventory"
)

func Get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	queryValues := r.URL.Query()
	stockableOnly, _ := strconv.ParseBool(queryValues.Get("stockableOnly"))
	items, err := allItems(stockableOnly)
	if err != nil {
		log.Errorf("error retrieving inventory items: %+v", err)
		return
	}
	_, _ = fmt.Fprint(w, items)
}

func GetById(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
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

func Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	i, err := bodyToItem(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	createdItemId, err := insertItem(i)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_, _ = fmt.Fprintf(w, "http://"+r.Host+r.URL.String()+"/"+createdItemId)
}

func Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	i, err := bodyToItem(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id := ps.ByName("id")
	if id == "" || i.Id != id {
		http.Error(w, "id field can not be updated", http.StatusConflict)
		return
	}
	statusCode, refSelf, err := upsertItem(i, r.Host + "/" +RequestUri)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.WriteHeader(statusCode)
		if statusCode == http.StatusCreated {
			_, _ = fmt.Fprintf(w, refSelf)
		}
	}
}

func bodyToItem(r *http.Request) (Item, error) {
	var i Item
	err := json.NewDecoder(r.Body).Decode(&i)
	return i, err
}
