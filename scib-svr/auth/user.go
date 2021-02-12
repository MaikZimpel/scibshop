package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"scib-svr/configuration"
	"scib-svr/logging"
	"time"
)

type User struct {
	Username string `json:"username" bson:"username, omitempty"`
	Password string `json:"password" bson:"password, omitempty"`
	Roles []string `json:"roles" bson:"roles"`
	Token string `json:"token" bson:"token"`
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

func (s *Service) byUsername(c context.Context, username string) (*User, error) {
	var user = &User{}
	sr := s.collection(userCollection).FindOne(c, bson.M{"username": username})
	err := sr.Decode(user)
	return user, err
}

func (s *Service) SignIn(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username := ps.ByName("username")
	password := ps.ByName("password")
	user, err := s.byUsername(r.Context(), username)
	if err == nil && user.Password == password {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"roles": user.Roles,
		})
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, token)
	}
}

func (s *Service) Save(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user := User{}
	if err := json.NewDecoder(r.Body).Decode(&user); err == nil {
		var uOptions = &options.UpdateOptions{}
		uOptions.SetUpsert(true)
		if _, err := s.collection(userCollection).UpdateOne(r.Context(), bson.M{"username": user.Username},user, uOptions); err == nil {

		}
	}



}