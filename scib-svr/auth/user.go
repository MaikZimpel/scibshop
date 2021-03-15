package auth

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"scib-svr/configuration"
	"scib-svr/logging"
	"time"
)

type User struct {
	Id       string   `json:"id" bson:"_id omitempty"`
	Username string   `json:"username" bson:"username, omitempty"`
	Password string   `json:"password" bson:"password, omitempty"`
	Roles    []string `json:"roles" bson:"roles"`
	Token    string   `json:"token" bson:"token"`
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
	if err != nil {
		logger.Critical(context.Background(), "%v", err)
		return nil
	}
	database := client.Database(configuration.MongoDbDatabase)
	return &Service{database, logger}
}

func (s *Service) byUsername(c context.Context, username string) (*User, error) {
	var user = &User{}
	sr := s.collection(userCollection).FindOne(c, bson.M{"username": username})
	err := sr.Decode(user)
	return user, err
}

func (s *Service) SignIn(c context.Context, username string, password string) (*User, error) {
	user, err := s.byUsername(c, username)
	if err == nil && user.Password == password {
		return user, nil
	} else {
		return nil, fmt.Errorf("authentication failed")
	}
}

func (s *Service) Save(c context.Context, user *User) (int64, error) {
	var uOptions = &options.UpdateOptions{}
	uOptions.SetUpsert(true)
	if user.Id == "" {
		user.Id = uuid.New().String()
	}
	updateResult, err := s.collection(userCollection).UpdateOne(c, bson.M{"_id": bson.M{"$eq": fmt.Sprint(user.Id)}}, bson.M{"$set": user}, uOptions)
	s.log.Info(c, "update result: %v", updateResult)
	return 1, err
}

func (s *Service) queryUsers(c context.Context, filter bson.M) (userVector []User, err error) {
	users := s.collection(userCollection)
	csr, err := users.Find(c, filter)
	if err == nil && csr != nil {
		if err = csr.All(c, &userVector); err != nil {
			s.log.Error(c, "%v", err)
		}
	}
	return
}
