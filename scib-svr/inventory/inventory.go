package inventory

import (
	"cloud.google.com/go/firestore"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"scib-svr/external"
	"scib-svr/helpers"
)

type Item struct {
	Id          string               `json:"id"`
	Upc         string               `json:"upc"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Categories  []string             `json:"categories"`
	Brand       string               `json:"brand"`
	Size        string               `json:"size"`
	Price       int                  `json:"price"`
	Images      []string             `json:"images"`
	Stockable   bool                 `json:"stockable"`
	Suppliers   []*external.Supplier `json:"suppliers"`
}

func (i *Item) String() string {
	b, _ := json.Marshal(i)
	return string(b)
}

func insertItem(item Item) (string, error) {
	client, ctx := helpers.DataStoreClient()
	defer client.Close()
	item.Id = uuid.New().String()
	_, err := client.Collection(helpers.InventoryCollection).Doc(item.Id).Create(ctx, item)
	return item.Id, err
}

func upsertItem(item Item, path string) (int, string, error) {
	client, ctx := helpers.DataStoreClient()
	defer client.Close()
	var retCode = http.StatusNoContent
	var retValue = ""
	_, err := client.Collection(helpers.InventoryCollection).Doc(item.Id).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			retCode = http.StatusCreated
			retValue = fmt.Sprintf("http://%s/%s", path, item.Id)
		} else {
			return http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), err
		}
	}
	_, _ = client.Collection(helpers.InventoryCollection).Doc(item.Id).Set(ctx, item)
	return retCode, retValue, nil
}

func allItems(stockableOnly bool) []*Item {
	client, ctx := helpers.DataStoreClient()
	defer client.Close()
	if stockableOnly {
		return transform(client.Collection(helpers.InventoryCollection).Where("Stockable", "==", stockableOnly).Documents(ctx))
	} else {
		return transform(client.Collection(helpers.InventoryCollection).Documents(ctx))
	}
}

func transform(iter *firestore.DocumentIterator) []*Item {
	var items []*Item
	for {
		docRef, e := iter.Next()
		if e == iterator.Done {
			break
		}
		if e != nil {
			fmt.Println(e)
		}
		var item *Item
		docRef.DataTo(&item)
		items = append(items, item)
	}
	return items
}

func itemById(id string) (*Item, error) {
	client, ctx := helpers.DataStoreClient()
	defer client.Close()
	docRef, err := client.Collection(helpers.InventoryCollection).Doc(id).Get(ctx)
	if err != nil {
		msg, code := mapError(err)
		return nil, helpers.Error(msg, code)
	} else {
		var item *Item
		docRef.DataTo(&item)
		return item, nil
	}
}

func mapError(err error) (string, int) {
	switch status.Code(err) {
	case codes.NotFound:
		{
			return http.StatusText(http.StatusNotFound), http.StatusNotFound
		}
	default:
		return http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError

	}
}
