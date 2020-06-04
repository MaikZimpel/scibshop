package helpers

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"log"
)

const (
	INVENTORY_COLLECTION = "inventory"
	STORE_COLLECTION     = "store"
	projectID            = "scib-279217"
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

func DataStoreClient() (*firestore.Client, context.Context) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create data client: %v", err)
	}
	return client, ctx
}
