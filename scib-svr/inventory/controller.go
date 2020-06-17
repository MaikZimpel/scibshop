package inventory

import (
	"cloud.google.com/go/storage"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"net/http"
	"os"
	"scib-svr/datastore"
	"scib-svr/httputil"
	"scib-svr/logging"
	"strconv"
	"time"
)

const (
	RequestUri = "inventory"
)

type Controller struct {
	s *Service
	l logging.Logger
}

func NewController(service *Service) *Controller {
	logger := logging.New()
	return &Controller{service, logger}
}

func (c *Controller) Get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	httputil.EnableCors(&w, "http://localhost:3001")
	c.l.Debug(context.Background(), "inventory requested by %s", r.RemoteAddr)
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

func (c *Controller) SignedUrl(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	httputil.EnableCors(&w, "http://localhost:3001")

	// Accepts only POST method.
	// Otherwise, this handler returns 405.
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Only POST is supported", http.StatusMethodNotAllowed)
		return
	}

	ct := r.FormValue("content_type")
	if ct == "" {
		http.Error(w, "content_type must be set", http.StatusBadRequest)
		return
	}

	// Generates an object key for use in new Cloud Storage Object.
	// It's not duplicate with any object keys because of UUID.
	key := uuid.New().String()
	if ext := r.FormValue("ext"); ext != "" {
		key += fmt.Sprintf(".%s", ext)
	}

	saKeyFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	saKey, err := ioutil.ReadFile(saKeyFile)
	if err != nil {
		c.l.Error(r.Context(), "%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	cfg, err := google.JWTConfigFromJSON(saKey)
	if err != nil {
		c.l.Error(r.Context(), "%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Generates a signed URL for use in the PUT request to GCS.
	// Generated URL should be expired after 15 mins.
	url, err := storage.SignedURL("scib-279217.appspot.com", key, &storage.SignedURLOptions{
		GoogleAccessID: cfg.Email,
		PrivateKey: cfg.PrivateKey,
		Method:         "PUT",
		Expires:        time.Now().Add(15 * time.Minute),
		ContentType:    ct,

	})
	if err != nil {
		c.l.Error(r.Context(), "sign: failed to sign, err = %v\n", err)
		http.Error(w, "failed to sign by internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, url)
}

func bodyToItem(r *http.Request) (item Item, err error) {
	err = json.NewDecoder(r.Body).Decode(&item)
	return
}
