package inventory

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"net/http"
	"path/filepath"
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
	items, err := c.s.queryItems(r.Context(), bson.M{})
	if err != nil {
		_, _ = fmt.Fprintf(w, "an error occured: %+v", err.Error())
	} else {
		bytes, err := json.Marshal(items)
		if err == nil {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(bytes)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, "an error occured: %+v", err.Error())
		}
	}
}

func (c *Controller) GetById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	item, err := c.s.byId(r.Context(), ps.ByName("id"))
	if err == nil {
		_, _ = fmt.Fprintf(w, "%v", item)
	}
	if err != nil {
		_, _ = fmt.Fprintf(w, "an error occured: %+v", err.Error())
	}
}

func (c *Controller) getWithFilter(ctx context.Context, filter bson.M) (items []Item, err error)  {
	items, err = c.s.queryItems(ctx, filter)
	if err != nil {
		c.l.Error(context.Background(), "error retrieving inventory items: %+v", err)
	}
	return
}

func (c *Controller) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	i, err := c.bodyToItem(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sts, createdItemId, err := c.s.save(r.Context(), &i)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(sts)
	_, _ = fmt.Fprintf(w, "http://"+r.Host+r.URL.String()+"/"+createdItemId)
}

func (c *Controller) Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	i, err := c.bodyToItem(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id := ps.ByName("id")
	if id == "" || fmt.Sprint(i.Id) != id {
		http.Error(w, "id field can not be updated", http.StatusConflict)
		return
	}
	sts, itemId, err := c.s.save(r.Context(), &i)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.WriteHeader(sts)
		if sts == http.StatusCreated {
			_, _ = fmt.Fprintf(w, "http://"+r.Host+r.URL.String()+"/"+itemId)
		}
	}
}

func (c *Controller) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	_, err := c.s.delete(r.Context(), id)
	if err != nil {
		code := func() int {
			switch err {
			case mongo.ErrNoDocuments: {
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

func (c *Controller) UploadImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	itemId := ps.ByName("id")
	filter := bson.M {"_id": itemId}
	items, err := c.getWithFilter(r.Context(), filter)
	if err != nil || items == nil {
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
	defer func() {
		err = file.Close()
		if err != nil {
			c.l.Error(r.Context(), "%v", err)
		}
	}()
	var bytes = make([]byte, handler.Size)
	_, err = file.Read(bytes)
	if err != nil {
		c.l.Error(r.Context(), "%s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	imgId, err := func() (string, error) {
		if r.FormValue("isColorImage") == "true" {
			qty, _ := strconv.Atoi(r.FormValue("cnt"))
			stockable, _ := strconv.ParseBool(r.FormValue("stockable"))
			available, _ := strconv.ParseBool(r.FormValue("available"))
			index, _ := strconv.Atoi(r.FormValue("variantIndex"))
			variant := ItemVariant {
				Sku:       r.FormValue("sku"),
				Color:     r.FormValue("color"),
				Image:     "",
				Size:      r.FormValue("size"),
				Cnt:       qty,
				Stockable: stockable,
				Available: available,
			}
			return 	c.s.uploadVImage(r.Context(), &items[0], index, variant, bytes, filepath.Ext(handler.Filename))
		} else {
			return c.s.uploadImage(r.Context(), &items[0], bytes, filepath.Ext(handler.Filename))
		}
	}()

	if err != nil {
		c.l.Error(r.Context(), err.Error(), http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusCreated)
		_, _ = fmt.Fprintf(w, imgId)
	}
}

func (c *Controller) GetImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("imageId")
	_, image, err := c.s.downloadImage(r.Context(), id)
	if err != nil {
		c.l.Error(r.Context(), "file not found with id %s", id)
		http.Error(w, fmt.Errorf("image %s not found on server", id).Error(), http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(image)
	}
}

func (c *Controller) DeleteImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	items, err := c.getWithFilter(r.Context(), bson.M{"_id": ps.ByName("id")})
	if err != nil {
		c.l.Error(r.Context(), err.Error(), err)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	err = c.s.deleteImage(r.Context(), ps.ByName("imageId"), &items[0])
	if err != nil {
		c.l.Error(r.Context(), "file deletion failed %v", err)
		http.Error(w, "file deletion failed, ask your guy what went wrong", http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func (c *Controller) bodyToItem(r *http.Request) (item Item, err error) {
	err = json.NewDecoder(r.Body).Decode(&item)
	c.l.Debug(context.Background(), "item %#v", item)
	return
}
