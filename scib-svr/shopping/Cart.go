package shopping

import "scib-svr/inventory"

type Pick struct {
	Item inventory.Item
	Cnt int
}

type Cart struct {
	Picks []Pick
}
