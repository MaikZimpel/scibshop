package shopping

type Cart struct {
	Id    string `json:"id"`
	Picks map[string]int `json:"picks"` // sku[qty]
}

/*func unstock(sku string) error {
	client, ctx := datastore.DataStoreClient()
	defer client.Close()
	_, err := client.Collection(datastore.StoreCollection).Doc(sku).Set(ctx, map[string]interface{}{"Available": false}, firestore.MergeAll)
	return err
}*/
