package datastore

import (
	"context"
	"fmt"
)

type (
	Datastore interface {
		Insert(c context.Context, collection string, id string, document interface{}) (string, error)
		Upsert(c context.Context, collection string, id string, document interface{}, replace bool) (int, error)
		Delete(c context.Context, collection string, id string) (map[string]interface{}, error)
		Read(c context.Context, collection string, id string, whereClause []WhereClause) ([]map[string]interface{}, error)
	}

	datastoreFactory func() Datastore
)

var (
	New datastoreFactory
)

type WhereClause struct {
	Path string
	Op string // The op argument must be one of "==", "<", "<=", ">", ">=", "in".
	Value interface{}
}

type ErrCode int

const (
	NotFound ErrCode = iota
)

type DsError struct {
	Code ErrCode
	Err error
}

func (r DsError) Error() string {
	return fmt.Sprintf("code %d: err %v", r.Code, r.Err)
}
