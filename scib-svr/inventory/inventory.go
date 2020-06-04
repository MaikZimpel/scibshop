package inventory

import (
	"cloud.google.com/go/firestore"
	"encoding/json"
	"fmt"
	"google.golang.org/api/iterator"
	"scib-svr/helpers"
)

type Item struct {
	Id          string   `json:"id"`
	Upc         string   `json:"upc"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Categories  []string `json:"categories"`
	Brand       string   `json:"brand"`
	Size        string   `json:"size"`
	Price       int      `json:"price"`
	Images      []string `json:"images"`
	Stockable   bool     `json:"stockable"`
}

func (i *Item) String() string {
	b, _ := json.Marshal(i)
	return string(b)
}

func addItem(item Item) (string, error) {
	client, ctx := helpers.DataStoreClient()
	defer client.Close()
	_, _, err := client.Collection(helpers.INVENTORY_COLLECTION).Add(ctx, item)
	return item.Id, err
}

func allItems(stockableOnly bool) []*Item {
	client, ctx := helpers.DataStoreClient()
	defer client.Close()
	if stockableOnly {
		return transform(client.Collection(helpers.INVENTORY_COLLECTION).Where("Stockable", "==", stockableOnly).Documents(ctx))
	} else{
		return transform(client.Collection(helpers.INVENTORY_COLLECTION).Documents(ctx))
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
	items := transform(client.Collection(helpers.INVENTORY_COLLECTION).Where("Id", "==", id).Documents(ctx))
	if len(items) == 0 {
		return nil , helpers.Error(fmt.Sprintf("inventory item with %s not found", id), 404)
	}
	if len(items) > 1 {
		return nil , helpers.Error(fmt.Sprintf("id %s is not unique", id), 409)
	}
	return items[0], nil
}
