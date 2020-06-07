package shopping

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/api/iterator"
	"scib-svr/external"
	"scib-svr/helpers"
	"scib-svr/inventory"
)

type StoreItem struct {
	Item      *inventory.Item    `json:"item"`
	Supplier  *external.Supplier `json:"supplier"`
	Sku       string             `json:"sku"`
	Cnt       int                `json:"cnt"`
	Available bool               `json:"available"`
}

func (i *StoreItem) String() string {
	b, _ := json.Marshal(i)
	return string(b)
}

type Pick struct {
	Item *StoreItem
	Cnt  int
}

type Cart struct {
	Id    string
	Picks map[string]*Pick
}

func allItems(availableOnly bool) []*StoreItem {
	client, ctx := helpers.DataStoreClient()
	defer client.Close()
	if availableOnly {
		return Transform(client.Collection(helpers.StoreCollection).Where("Available", "==", availableOnly).Documents(ctx))
	} else {
		return Transform(client.Collection(helpers.StoreCollection).Documents(ctx))
	}
}

func itemById(id string) (*StoreItem, error) {
	client, ctx := helpers.DataStoreClient()
	defer client.Close()
	docRef, err := client.Collection(helpers.StoreCollection).Doc(id).Get(ctx)
	if err != nil {
		msg, code := helpers.MapError(err)
		return nil, helpers.Error(msg, code)
	} else {
		var item *StoreItem
		err = docRef.DataTo(&item)
		return item, err
	}
}

func insert(item *StoreItem) (*StoreItem, error, bool) {
	client, ctx := helpers.DataStoreClient()
	defer client.Close()
	created := false
	if item.Sku == "" {
		item.Sku = generateSku(item.Supplier, item.Item, client, ctx)
		created = true
	}
	_, err := client.Collection(helpers.StoreCollection).Doc(item.Sku).Set(ctx, item)
	return item, err, created
}

func unstock(sku string) error {
	client, ctx := helpers.DataStoreClient()
	defer client.Close()
	_, err := client.Collection(helpers.StoreCollection).Doc(sku).Set(ctx, map[string]interface{}{"Available": false}, firestore.MergeAll)
	return err
}

func Transform(iter *firestore.DocumentIterator) []*StoreItem {
	var items []*StoreItem
	for {
		docRef, e := iter.Next()
		if e == iterator.Done {
			break
		}
		if e != nil {
			fmt.Println(e)
		}
		var item *StoreItem
		docRef.DataTo(&item)
		items = append(items, item)
	}
	return items
}

func generateSku(supplier *external.Supplier, item *inventory.Item, client *firestore.Client, ctx context.Context) string {
	items, _ := inventory.Transform(client.Collection(helpers.InventoryCollection).Documents(ctx))
	skuIndx := 0
	for _, item := range items {
		for _, sup := range item.Suppliers {
			if sup.Id == supplier.Id {
				skuIndx++
			}
		}
	}
	return fmt.Sprintf("%s-%04d-%s", supplier.SkuPrefix, skuIndx, item.SkuSuffix)
}
