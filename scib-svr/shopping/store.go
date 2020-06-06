package shopping

import (
	"fmt"
	"github.com/google/uuid"
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

func stock(item *inventory.Item, sku string, cnt int, available bool) (*StoreItem, error) {
	si := StoreItem {
		Item:      item,
		Sku:       sku,
		Cnt:       cnt,
		Available: available,
	}

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
	if cart.Id == "" {
		cart.Id = uuid.New().String()
	}
	_, _ = client.Collection(helpers.StoreCollection).Doc(cart.Id).Set(ctx, cart)
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
	_, _ = client.Collection(helpers.StoreCollection).Doc(cart.Id).Set(ctx, cart)
	return cart, nil
}

func discard(cart *Cart) error {
	client, ctx := helpers.DataStoreClient()
	_, err := client.Collection(helpers.StoreCollection).Doc(cart.Id).Delete(ctx)
	return err
}
