package shopping

/*func TestPick(t *testing.T) {
	store := testStore()
	type exp struct {
		err error
		cart Cart
	}

	cartId := uuid.New().String()
	cart := Cart {
		Id: cartId,
		Picks: map[string]*Pick {},
	}
	tables := []struct {
		item *StoreItem
		cnt int
		expected struct{
			err error
			cart Cart
		}
	} {
		{store["1"], 3, exp {nil, Cart{cartId, getPicks(map[string]int{"1":3}, &store)}}},
		{store["2"], 4, exp {fmt.Errorf("not enough items of %s in insert. wanted: %d, available: %d", store["2"].Sku, 4, store["2"].Cnt), cart}},
		{store["3"], 1, exp {fmt.Errorf("item %s is unavaiable", store["3"].Sku), cart}},
	}

	for _, table := range tables {
		cart, err := pick(table.item, &cart, table.cnt)
		if err == nil && table.expected.err == nil {
			if str(cart) != str(&table.expected.cart) {
				t.Errorf("TestPick failed. Got: %+v wanted: %+v", str(cart), str(&table.expected.cart))
			}
		} else {
			if err != nil && table.expected.err != nil && err.Error() != table.expected.err.Error() {
				t.Errorf("TestPick failed. Got: %s, wanted: %s", err.Error(), table.expected.err.Error())
			}
		}

	}
}

func getPicks(ids map[string]int, store *map[string]*StoreItem) map[string]*Pick {
	result := map[string]*Pick {}
	for key, cnt := range ids {
		storeItem := (*store)[key]
		result[storeItem.Sku] = &Pick{
			Item: storeItem,
			Cnt:  cnt,
		}
	}
	return result
}

func str(cart *Cart) string  {
	bytes, _ := json.Marshal(&cart)
	return string(bytes)
}

func testStore() map[string]*StoreItem {
	store := map[string]*StoreItem{
		"1": {
			Item:      &inventory.Item{
				Id : "1",
				Name: "Badekappe",
				Brand: "Swimma",
			},
			Sku:       "1-1-Swimma",
			Cnt:       15,
			Available: true,
		},
		"2":{
			Item: &inventory.Item{
				Id: "2",
				Name: "Badeanzug",
				Brand: "Swimma",
			},
			Sku: "2-2-Swimma",
			Cnt: 3,
			Available: true,
		},
		"3":{
			Item: &inventory.Item{
				Id: "3",
				Name: "Sonnenbrille",
				Brand: "Adidas",
			},
			Sku: "3-3-Adidas",
			Cnt: 99,
			Available: false,
		},
	}
	return store
}*/