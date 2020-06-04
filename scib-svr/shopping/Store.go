package shopping

import (
	"scib-svr/inventory"
)

const COLLECTION = "store"

type StoreItem struct {
	Item      *inventory.Item `json:"item"`
	Sku         string   `json:"sku"`
	Cnt       int `json:"cnt"`
	Available bool `json:"available"`
}



func Stock(item *inventory.Item, cnt int, availability bool) (StoreItem, error) {

}

func Take(item *StoreItem, cnt int) (Pick, error) {

}
