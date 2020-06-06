package shopping

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/google/uuid"
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

type Pick struct {
	Item *StoreItem
	Cnt  int
}

type Cart struct {
	Id    string
	Picks map[string]*Pick
}

func stock(item *inventory.Item, supplier *external.Supplier, cnt int, available bool) (*StoreItem, error) {
	client, ctx := helpers.DataStoreClient()
	defer client.Close()
	sisRef := client.Collection(helpers.StoreCollection).
		Where("item.Id", "==", item.Id).
		Where("supplier.Id", "==", supplier.Id).
		Documents(ctx)
	si := StoreItem{
		Item:      item,
		Sku:       sku(Transform(sisRef), supplier, item, client, ctx),
		Supplier:  supplier,
		Cnt:       cnt,
		Available: available,
	}
	_, err := client.Collection(helpers.StoreCollection).Doc(si.Sku).Set(ctx, si)
	return &si, err
}

func unstock(sku string) error {
	client, ctx := helpers.DataStoreClient()
	defer client.Close()
	_, err := client.Collection(helpers.StoreCollection).Doc(sku).Delete(ctx)
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

func sku(si []*StoreItem, supplier *external.Supplier, item *inventory.Item, client *firestore.Client, ctx context.Context) string {
	if len(si) == 0 {
		return generateSku(supplier, item, client, ctx)
	} else {
		return si[0].Sku
	}
}

func generateSku(supplier *external.Supplier, item *inventory.Item, client *firestore.Client, ctx context.Context) string {
	items := inventory.Transform(client.Collection(helpers.InventoryCollection).Documents(ctx))
	sukIndx := 0
	for _, item := range items {
		for _, sup := range item.Suppliers {
			if sup.Id == supplier.Id {
				sukIndx++
			}
		}
	}
	return fmt.Sprintf("%s-%04d-%s", supplier.SkuPrefix, sukIndx, item.SkuSuffix)
}

func pick(item *StoreItem, cart *Cart, cnt int) (*Cart, error) {
	if item.Available == false {
		return cart, fmt.Errorf("item %s is unavaiable", item.Sku)
	}
	if item.Cnt < 1 || item.Cnt < cnt {
		return cart, fmt.Errorf("not enough items of %s in stock. wanted: %d, available: %d", item.Sku, cnt, item.Cnt)
	}
	pick := cart.Picks[item.Sku]
	if pick == nil {
		pick := Pick{
			Item: item,
			Cnt:  cnt,
		}
		cart.Picks[item.Sku] = &pick
	} else {
		pick.Cnt += cnt
	}
	client, ctx := helpers.DataStoreClient()
	defer client.Close()
	if cart.Id == "" {
		cart.Id = uuid.New().String()
	}
	_, _ = client.Collection(helpers.CartCollection).Doc(cart.Id).Set(ctx, cart)
	return cart, nil
}

func unpick(sku string, cart *Cart, cnt int) (*Cart, error) {
	pick := cart.Picks[sku]
	if pick == nil {
		return cart, fmt.Errorf("no such item in shopping cart: %s", sku)
	}
	if pick.Cnt <= cnt {
		delete(cart.Picks, sku)
	} else {
		pick.Cnt -= cnt
	}
	client, ctx := helpers.DataStoreClient()
	defer client.Close()
	_, _ = client.Collection(helpers.CartCollection).Doc(cart.Id).Set(ctx, cart)
	return cart, nil
}

func discard(cart *Cart) error {
	client, ctx := helpers.DataStoreClient()
	defer client.Close()
	_, err := client.Collection(helpers.CartCollection).Doc(cart.Id).Delete(ctx)
	return err
}
