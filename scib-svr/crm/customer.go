package crm

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"scib-svr/configuration"
	"scib-svr/logging"
	"time"
)

type Customer struct {
	Id        string    `json:"id" bson:"_id, omitempty"`
	Email     string    `json:"email" bson:"email, omitempty"`
	Name      string    `json:"name" bson:"name, omitempty"`
	Addresses []Address `json:"addresses" bson:"addresses, omitempty"`
	Phone     string    `json:"phone" bson:"phone, omitempty"`
}

type Address struct {
	Id      string      `json:"id" bson:"id, omitempty"`
	Country string      `json:"country" bson:"country, omitempty"`
	City    string      `json:"city" bson:"city, omitempty"`
	Code    string      `json:"code" bson:"code, omitempty"`
	Street  string      `json:"street" bson:"street, omitempty"`
	Type    AddressType `json:"type" bson:"type, omitempty"`
}

type AddressType int

const (
	Shipping AddressType = iota
	Residence
	Other
)

func (at AddressType) String() string {
	return [...]string{"Shipping", "Residence", "Other"}[at]
}

const customerCollection = "customers"
const customerCollectionDeletedItems = "customersDeleted"

func (c *Customer) addAddress(a Address) {
	c.Addresses = append(c.Addresses, a)
}

func (c *Customer) removeAddress(id string) bool {
	for x, v := range c.Addresses {
		if v.Id == id {
			copy(c.Addresses[x:], c.Addresses[x+1:])
			c.Addresses[len(c.Addresses)-1] = Address{}
			c.Addresses = c.Addresses[:len(c.Addresses)-1]
			return true
		}
	}
	return false
}

func (c *Customer) addressByType(addressType AddressType) *Address {
	for _, v := range c.Addresses {
		if v.Type == addressType {
			return &v
		}
	}
	return nil
}

type Service struct {
	database *mongo.Database
	log      logging.Logger
}

func (s *Service) collection(name string) *mongo.Collection {
	return s.database.Collection(name)
}

func NewService(logger logging.Logger) *Service {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	opts := options.Client()
	opts.ApplyURI(fmt.Sprintf("mongodb://%s:%d", configuration.MongoDbHost, configuration.MongoDbPort))
	opts.SetMaxPoolSize(5)
	client, err := mongo.Connect(ctx, opts)
	if client != nil {
		database := client.Database(configuration.MongoDbDatabase)
		if err == nil {
			return &Service{database, logger}
		} else {
			logger.Critical(context.Background(), "%v", err)
			return nil
		}
	}
	return nil
}

func (s *Service) save(c context.Context, customer *Customer) (resultCode int, id string, err error) {
	collection := s.collection(customerCollection)
	var uOptions = &options.UpdateOptions{}
	uOptions.SetUpsert(true)
	if customer.Id == "" {
		customer.Id = uuid.New().String()
	}
	result, err := collection.UpdateOne(c, bson.M{"_id": bson.M{"$eq": fmt.Sprint(customer.Id)}}, bson.M{"$set": customer}, uOptions)
	if result != nil {
		if result.UpsertedCount > 0 {
			resultCode = http.StatusCreated
			id = fmt.Sprint(result.UpsertedID)
		} else {
			resultCode = http.StatusNoContent
		}
	}
	return
}

func (s *Service) byId(ctx context.Context, id string) (*Customer, error) {
	var c = &Customer{}
	sr := s.collection(customerCollection).FindOne(ctx, bson.M{"_id": id})
	err := sr.Decode(c)
	return c, err
}

func (s *Service) queryCustomers(ctx context.Context, filter bson.M) (customerVector []Customer, err error) {
	collection := s.collection(customerCollection)
	csr, err := collection.Find(ctx, filter)
	if err == nil && csr != nil {
		if err = csr.All(ctx, &customerVector); err != nil {
			s.log.Error(ctx, "%v", err)
		}
	}
	return
}

func (s *Service) delete(c context.Context, id string) (customer *Customer, err error) {
	collection := s.collection(customerCollection)
	delCollection := s.collection(customerCollectionDeletedItems)
	sr := collection.FindOne(c, bson.M{"_id": id})
	if sr.Err() != nil && sr.Err() == mongo.ErrNoDocuments {
		err = sr.Err()
	} else {
		err = sr.Decode(&customer)
		if err == nil {
			_, err := delCollection.InsertOne(c, customer)
			if err == nil {
				_, err = collection.DeleteOne(c, bson.M{"_id": id})
			}
		}
	}
	return
}
