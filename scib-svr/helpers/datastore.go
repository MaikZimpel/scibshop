package helpers

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net/http"
)

const (
	InventoryCollection = "inventory"
	StoreCollection     = "store"
	SupplierCollection  = "suppliers"
	CartCollection      = "carts"
	projectID           = "scib-279217"
)

type QueryError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (q QueryError) Error() string {
	b, _ := json.Marshal(q)
	return string(b)
}

func Error(message string, code int) error {
	return QueryError{message, code}
}

func MapError(err error) (string, int) {
	switch status.Code(err) {
	case codes.NotFound:
		{
			return http.StatusText(http.StatusNotFound), http.StatusNotFound
		}
	default:
		return http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError

	}
}

func DataStoreClient() (*firestore.Client, context.Context) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create data client: %v", err)
	}
	return client, ctx
}
