// +build mongo

package datastore

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fatih/structs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"scib-svr/configuration"
	"scib-svr/logging"
	"time"
)

type mongostore struct {
	client *mongo.Client
	db     string
	log    logging.Logger
}

func (m mongostore) Insert(c context.Context, collection string, id string, document interface{}) (newId string, err error) {
	if id != "" {
		_, err = m.Upsert(c, collection, id, document, true)
		if err != nil {
			newId = id
		}
	} else {
		res, err := m.client.Database(m.db).Collection(collection).InsertOne(c, document)
		if err == nil {
			newId = fmt.Sprint(res.InsertedID)
		}
	}
	return
}

func (m mongostore) Upsert(c context.Context, collection string, id string, document interface{}, replace bool) (status int, err error) {
	if replace {
		m.log.Debug(c, "before findoneandreplace %#v", time.Now())
		sr := m.client.Database(m.db).Collection(collection).FindOneAndReplace(c, bson.D{{"_id", id}}, document)
		m.log.Debug(c, "after findoneandreplace %#v", time.Now())
		m.log.Debug(c, "result %#v", sr)
		if sr.Err() == mongo.ErrNoDocuments {
			status, err = m.update(c, collection, id, document)
		}
	} else {
		status, err = m.update(c, collection, id, document)
	}
	return
}

func (m mongostore) update(c context.Context, collection string, id string, document interface{}) (status int, err error) {
	updateOptions := options.Update().SetUpsert(true)
	filter := bson.D{{"_id", id}}
	uDoc, err := m.createUpdateDocument(document)
	ur, err := m.client.Database(m.db).Collection(collection).UpdateOne(c, filter, uDoc, updateOptions)
	if err == nil {
		if ur.UpsertedID != nil {
			status = http.StatusCreated
		} else {
			status = http.StatusNoContent
		}
	} else {
		status = http.StatusInternalServerError
	}
	return
}

func (m mongostore) createUpdateDocument(document interface{}) (b bson.D, err error) {
	dm := structs.Map(document)
	b = bson.D{
		{"$set", func() bson.M {
			m := make(bson.M, len(dm))
			for k, v := range dm {
				m[k] = v
			}
			return m
		}()},
	}
	return
}

func (m mongostore) Delete(c context.Context, collection string, id string) (res map[string]interface{}, err error) {
	sr := m.client.Database(m.db).Collection(collection).FindOneAndDelete(c, bson.D{{"_id", id}})
	if sr.Err() != nil && sr.Err() == mongo.ErrNoDocuments {
		err = &DsError{
			Code: NotFound,
			Err:  sr.Err(),
		}
	} else {
		var resultStruc interface{}
		err = sr.Decode(resultStruc)
		if err == nil {
			res = structs.Map(resultStruc)
		}
	}
	return
}

func (m mongostore) Read(c context.Context, collection string, id string, whereClauses []WhereClause) (resultVector []map[string]interface{}, err error) {
	if id != "" {
		singleResult := m.client.Database(m.db).Collection(collection).FindOne(c, bson.D{{"_id", id}})
		if singleResult.Err() == mongo.ErrNoDocuments {
			err = &DsError{
				Code: NotFound,
				Err:  singleResult.Err(),
			}
		} else {
			var resultStruc bson.D
			err = singleResult.Decode(&resultStruc)
			if err == nil {
				resultVector, err = appendDecoded(resultStruc, resultVector)
			}
		}
	} else {
		csr, err := m.client.Database(m.db).Collection(collection).Find(c, createBsonFilter(whereClauses))
		if err == nil {
			for csr.Next(c) {
				var b bson.D
				err = csr.Decode(&b)
				if err == nil {
					resultVector, err = appendDecoded(b, resultVector)
				}
			}
		}
	}
	if err != nil {
		m.log.Error(c, "%s", err.Error())
	}
	return
}

func appendDecoded(decoded bson.D, resultVector []map[string]interface{}) (rv []map[string]interface{}, err error) {
	decodedAsMap, err := bsonToMap(decoded)
	if err == nil {
		rv = append(resultVector, decodedAsMap)
	}
	return
}

func bsonToMap(b bson.D) (r map[string]interface{}, err error) {
	tempBytes, err := bson.MarshalExtJSON(b, true, true)
	if err == nil {
		err = json.Unmarshal(tempBytes, &r)
	}
	return
}

func createBsonFilter(clauses []WhereClause) bson.M {
	if len(clauses) == 0 {
		return bson.M{}
	}
	var b []bson.M
	for _, clause := range clauses {
		b = append(b, bson.M{clause.Path: bson.M{opToMongoOp(clause.Op): clause.Value}})
	}
	return bson.M{"$and": b}
}

func opToMongoOp(op string) (mop string) {
	switch op {
	case "==":
		mop = "$eq"
	case "<":
		mop = "$lt"
	case "<=":
		mop = "$lte"
	case ">":
		mop = "$gt"
	case ">=":
		mop = "$gte"
	case "in":
		mop = "$in"
	}
	return
}

func NewMongostore() Datastore {
	clientOptions := options.Client().ApplyURI(
		fmt.Sprintf("mongodb://%s:%d", configuration.MongoDbHost, configuration.MongoDbPort))
	client, err := mongo.Connect(context.Background(), clientOptions)
	log := logging.New()
	if err != nil {
		log.Critical(context.Background(), "unable to connect to mongodb %v", err)
	} else {
		log.Info(context.Background(), "connected to database at %s",
			fmt.Sprintf("mongodb://%s:%d", configuration.MongoDbHost, configuration.MongoDbPort))
	}
	return &mongostore{
		client: client,
		db:     configuration.MongoDbDatabase,
		log:    logging.New(),
	}
}

func init() {
	New = NewMongostore
}
