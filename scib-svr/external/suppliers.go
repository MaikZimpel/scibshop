package external

import (
	"cloud.google.com/go/firestore"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
	"google.golang.org/api/people/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"scib-svr/helpers"
)

type Supplier struct {
	Id string
	Name string
	Address people.Address
	Email people.EmailAddress
	SkuPrefix string
}

func suppliers() []*Supplier {
	client, ctx := helpers.DataStoreClient()
	defer client.Close()
	return transform(client.Collection(helpers.SupplierCollection).Documents(ctx))
}

func transform(iter *firestore.DocumentIterator) []*Supplier {
	var suppliers []*Supplier
	for {
		docRef, e := iter.Next()
		if e == iterator.Done {
			break
		}
		if e != nil {
			fmt.Println(e)
		}
		var supplier *Supplier
		_ = docRef.DataTo(&supplier)
		suppliers = append(suppliers, supplier)
	}
	return suppliers
}

func insert(supplier Supplier) (string, error) {
	client, ctx := helpers.DataStoreClient()
	defer client.Close()
	supplier.Id = uuid.New().String()
	_, err := client.Collection(helpers.InventoryCollection).Doc(supplier.Id).Create(ctx, supplier)
	return supplier.Id, err
}

func upsert(supplier Supplier, path string) (int, string, error) {
	client, ctx := helpers.DataStoreClient()
	defer client.Close()
	var retCode = http.StatusNoContent
	var retValue = ""
	_, err := client.Collection(helpers.SupplierCollection).Doc(supplier.Id).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			retCode = http.StatusCreated
			retValue = fmt.Sprintf("http://%s/%s", path, supplier.Id)
		} else {
			return http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), err
		}
	}
	_, _ = client.Collection(helpers.SupplierCollection).Doc(supplier.Id).Set(ctx, supplier)
	return retCode, retValue, nil
}
