package crm

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"scib-svr/logging"
)

const (
	RequestUri = "/customer"
)

type Controller struct {
	s *Service
	l logging.Logger
}

func NewController(service *Service, logger logging.Logger) *Controller {
	log := logger
	return &Controller{service, log}
}

func (c *Controller) Get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	customers, err := c.s.queryCustomers(r.Context(), bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		out, err := json.Marshal(customers)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			_, _ = w.Write(out)
		}
	}
}

func (c *Controller) CreateOrUpdate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var cust = Customer{}
	err := json.NewDecoder(r.Body).Decode(&cust)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		resCode, id, err := c.s.save(r.Context(), &cust)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			w.WriteHeader(resCode)
			_, _ = fmt.Fprint(w, id)
		}
	}

}

func (c *Controller) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	_, err := c.s.delete(r.Context(), id)
	if err != nil {
		code := func() int {
			switch err {
			case mongo.ErrNoDocuments:
				{
					return http.StatusNotFound
				}
			default:
				return http.StatusInternalServerError
			}
		}()
		http.Error(w, err.Error(), code)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}



