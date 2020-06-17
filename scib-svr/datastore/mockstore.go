// +build mock

package datastore

import (
	"context"
	"github.com/google/uuid"
	"net/http"
)

type (
	mockDatastore struct {}
)

func (m mockDatastore) Insert(_ context.Context, _ string, id string, _ interface{}) (string, error) {
	if id == "" {
		return id, nil
	} else {
		return uuid.New().String(), nil
	}
}

func (m mockDatastore) Upsert(_ context.Context, _ string, _ string, _ interface{}, _ bool) (int, error) {
	return http.StatusNoContent, nil
}

func (m mockDatastore) Delete(c context.Context, collection string, id string) (map[string]interface{}, error) {
	mock := map[string]interface{}{
		"Id": id,
	}
	return mock, nil
}

func (m mockDatastore) Read(c context.Context, collection string, id string, whereClause []WhereClause) ([]map[string]interface{}, error) {
	mock := make([]map[string]interface{}, 0)
	mock = append(mock, map[string]interface{}{
		"Id": id,
	})
	return mock, nil
}

func NewMockDatastore () Datastore {
	return &mockDatastore{}
}

func init()  {
	New = NewMockDatastore
}




