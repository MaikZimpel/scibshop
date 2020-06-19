package inventory

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
	"net/http"
	"path/filepath"
	"scib-svr/datastore"
	"scib-svr/httputil"
	"scib-svr/logging"
	"strconv"
)

const (
	RequestUri = "/inventory"
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
	httputil.EnableCors(&w, "http://localhost:3001")
	queryValues := r.URL.Query()
	stockableOnly, _ := strconv.ParseBool(queryValues.Get("stockableOnly"))
	items, err := c.s.allItems(stockableOnly)
	if err != nil {
		c.l.Error(context.Background(),"error retrieving inventory items: %+v", err)
		return
	}
	err = json.NewEncoder(w).Encode(items)
	if err != nil {
		_, _ = fmt.Fprintf(w, "an error occured: %+v", err.Error())
	}
}

func (c *Controller) GetById(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	httputil.EnableCors(&w, "http://localhost:3001")
	item, err := c.s.itemById(ps.ByName("id"))
	if err != nil {
		code := func() int {
			switch e:= err.(type) {
			case *datastore.DsError:
				if e.Code == datastore.NotFound {
					return http.StatusNotFound
				} else {
					return http.StatusInternalServerError
				}
			default:
				return http.StatusInternalServerError
			}
		}()
		http.Error(w, err.Error(),code)
	} else {
		_, _ = fmt.Fprintf(w, item.String())
	}
}

func (c *Controller) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	httputil.EnableCors(&w, "http://localhost:3001")
	i, err := bodyToItem(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sts, createdItemId, err := c.s.save(&i)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(sts)
	_, _ = fmt.Fprintf(w, "http://"+r.Host+r.URL.String()+"/"+createdItemId)
}

func (c *Controller) Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	httputil.EnableCors(&w, "http://localhost:3001")
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
	sts, itemId, err := c.s.save(&i)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.WriteHeader(sts)
		if sts == http.StatusCreated {
			_, _ = fmt.Fprintf(w, "http://"+r.Host+r.URL.String()+"/"+itemId)
		}
	}
}

func (c *Controller) UploadImages(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	httputil.EnableCors(&w, "http://localhost:3001")
	itemId := ps.ByName("id")
	item, err := c.s.itemById(itemId)
	if err != nil || item == nil {
		if err == nil {
			err = fmt.Errorf("item not found with id %s", itemId)
		}
		c.l.Error(r.Context(), "item not found with id %s", itemId)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	err = r.ParseMultipartForm(50000000) // max 500KB per pic
	if err != nil {
		c.l.Error(r.Context(), "%s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	file, handler, err := r.FormFile("originalFile")
	if err != nil {
		c.l.Error(r.Context(), "%s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()
	var bytes = make([]byte, handler.Size)
	_, err = file.Read(bytes)
	if err != nil {
		c.l.Error(r.Context(), "%s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	imgName, err := c.s.uploadImage(item, bytes, filepath.Ext(handler.Filename))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusCreated)
		_, _ = fmt.Fprintf(w, imgName)
	}
}

func (c *Controller) GetImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	itemId := ps.ByName("id")
	fileName := ps.ByName("fileName")
	image, err := c.s.downloadImage(itemId, fileName)
	if err != nil {
		c.l.Error(r.Context(), "file not found with filename %s", fileName)
		http.Error(w, err.Error(), http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(image)
	}
}

func bodyToItem(r *http.Request) (item Item, err error) {
	err = json.NewDecoder(r.Body).Decode(&item)
	return
}
