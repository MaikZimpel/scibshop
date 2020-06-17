// +build firestore

package datastore

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/fatih/structs"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strconv"
)

type (
	firestoreDatastore struct{}
)

func NewFirestoreDatastore() Datastore {
	return &firestoreDatastore{}
}

func init()  {
	New = NewFirestoreDatastore
}

func (f *firestoreDatastore) Read(c context.Context, collection string, id string, whereClauses []WhereClause) (resultVector []map[string]interface{}, err error) {
	client, err := dataStoreClient(c)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	colRef := client.Collection(collection)
	if id != "" {
		docSnap, err := colRef.Doc(id).Get(c)
		if err != nil {
			return nil, &DsError {
				Code: NotFound,
				Err:  err,
			}
		} else {
			return append(resultVector, docSnap.Data()), nil
		}
	} else {
		if len(whereClauses) > 0 {
			addWhereClause(colRef.Query, whereClauses)
		}
		docIter := colRef.Documents(c)
		for {
			doc, err := docIter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return nil, err
			}
			resultVector = append(resultVector, doc.Data())
		}
	}
	return resultVector, nil
}

func addWhereClause(q firestore.Query, whereClauses []WhereClause) firestore.Query {
	if len(whereClauses) > 1 {
		return addWhereClause(q, whereClauses[:len(whereClauses)-1])
	} else {
		wc := whereClauses[0]
		return q.Where(wc.Path, wc.Op, wc.Value)
	}
}

// Insert a new document into the collection. The newly created document will be created with the given id in the firestore.
// If the given id is empty a new UUID is created and used as id for the document.
// If the id is not empty and it already exists insert returns an error.
func (f *firestoreDatastore) Insert(c context.Context, collection string, id string, document interface{}) (ID string, error error) {
	if id == "" {
		ID = uuid.New().String()
	} else {
		ID = id
	}
	client, err := dataStoreClient(c)
	if err != nil {
		return "", err
	}
	defer client.Close()
	_, err = client.Collection(collection).Doc(ID).Create(c, document)
	if err != nil {
		return "", err
	} else {
		return ID, nil
	}
}

// Updates or creates the given document. If the id does not exist yet, a new document is inserted similar to above insert operation.
// If the documents already exist it will be replaced if replace is set to true. Otherwise the supplied document will
// be merged into the existing one i.e. all fields that are present in the supplied document will be updated in the existing
// one, new fields will be added while fields that only exist in the present document remain untouched.
func (f *firestoreDatastore) Upsert(c context.Context, collection string, id string, document interface{}, replace bool) (int, error) {
	if id == "" {
		return http.StatusConflict, fmt.Errorf("%s", "id can not be empty")
	}
	client, err := dataStoreClient(c)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer client.Close()
	_, err = client.Collection(collection).Doc(id).Get(c)
	returnCode := http.StatusNoContent
	if err != nil {
		if status.Code(err) == codes.NotFound {
			returnCode = http.StatusCreated
		} else {
			return http.StatusInternalServerError, err
		}
	}
	if replace {
		_, err = client.Collection(collection).Doc(id).Set(c, document)
	} else {
		m := structs.Map(document)
		_, err = client.Collection(collection).Doc(id).Set(c, m, firestore.MergeAll)
		if err != nil {
			return http.StatusInternalServerError, err
		}
	}
	return returnCode, nil

}

func (f *firestoreDatastore) Delete(c context.Context, collection string, id string) (returnData map[string]interface{}, err error) {
	if id == "" {
		return nil, fmt.Errorf("%s", "id can not be empty")
	}
	client, err := dataStoreClient(c)
	if err != nil {
		return
	}
	defer client.Close()
	docRef, err := client.Collection(collection).Doc(id).Get(c)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, nil
		} else {
			return
		}
	} else {
		returnData = docRef.Data()
		_, err = client.Collection(collection).Doc(id).Delete(c)
		if err != nil {
			return nil, err
		} else {
			return
		}
	}

}

func dataStoreClient(c context.Context) (client *firestore.Client, error error) {
	size, err := strconv.Atoi("5")
	if err != nil {
		// TODO print configuration warning
		size = 2
	}
	client, error = firestore.NewClient(c, "scib-279217", option.WithGRPCConnectionPool(size))
	return client, error
}
