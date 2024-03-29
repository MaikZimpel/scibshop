package inventory

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"scib-svr/configuration"
	"scib-svr/filestore"
	"scib-svr/logging"
	"strings"
	"time"
)

type Item struct {
	Id          string        `json:"id" bson:"_id, omitempty"`
	Upc         string        `json:"upc" bson:"upc, omitempty"`
	Name        string        `json:"name" bson:"name, omitempty"`
	Description string        `json:"description" bson:"description, omitempty"`
	Categories  []string      `json:"categories" bson:"categories, omitempty"`
	Brand       string        `json:"brand" bson:"brand, omitempty"`
	Price       float32       `json:"price" bson:"price, omitempty"`
	Images      []string      `json:"images" bson:"images, omitempty"`
	Supplier    string        `json:"supplier" bson:"supplier, omitempty"`
	Variants    []ItemVariant `json:"variants" bson:"variants, omitempty"`
}

type ItemVariant struct {
	Sku       string `json:"sku" bson:"sku, omitempty"`
	Color     string `json:"color" bson:"color, omitempty"`
	Image     string `json:"image" bson:"image, omitempty"`
	Size      string `json:"size" bson:"size, omitempty"`
	Cnt       int    `json:"cnt" bson:"cnt, omitempty"`
	Stockable bool   `json:"stockable" bson:"stockable, omitempty"`
	Available bool   `json:"available" bson:"available, omitempty"`
}

const inventoryCollection = "inventory"
const inventoryCollectionDeletedItems = "inventory-deleted"

func (i *Item) String() string {
	b, _ := json.Marshal(i)
	return string(b)
}

func (i *Item) addImage(id string) {
	i.Images = append(i.Images, id)
}

func (i *Item) removeImage(id string) bool {
	for x, v := range i.Images {
		if v == id {
			copy(i.Images[x:], i.Images[x+1:])
			i.Images[len(i.Images)-1] = ""
			i.Images = i.Images[:len(i.Images)-1]
			return true
		}
	}
	return false
}

func (i *Item) removeImageFromItemVariant(id string) bool {
	for x, v := range i.Variants {
		if v.Image == id {
			i.Variants[x].Image = ""
			return true
		}
	}
	return false
}

func (i *Item) addItemVariant(variant ItemVariant) {
	i.Variants = append(i.Variants, variant)
}

func (i *Item) removeItemVariant(variant ItemVariant) bool {
	for x, v := range i.Variants {
		if v.Sku == variant.Sku {
			copy(i.Variants[x:], i.Variants[x+1:])
			i.Variants = i.Variants[:len(i.Variants)-1]
			return true
		}
	}
	return false
}

func (i *Item) filterVariants(filter func(ItemVariant) bool) {
	var filteredVariants []ItemVariant
	for _, variant := range i.Variants {
		if filter(variant) {
			filteredVariants = append(filteredVariants, variant)
		}
	}
	i.Variants = filteredVariants
}

type Service struct {
	database *mongo.Database
	fs       filestore.Filestore
	log      logging.Logger
}

func (s *Service) collection(name string) *mongo.Collection {
	return s.database.Collection(name)
}

func NewService(filestore filestore.Filestore, logger logging.Logger) *Service {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	opts := options.Client()
	opts.ApplyURI(fmt.Sprintf("mongodb://%s:%d", configuration.MongoDbHost, configuration.MongoDbPort))
	opts.SetMaxPoolSize(5)
	client, err := mongo.Connect(ctx, opts)
	if client != nil {
		database := client.Database(configuration.MongoDbDatabase)
		if err == nil {
			return &Service{database, filestore, logger}
		} else {
			logger.Critical(context.Background(), "%v", err)
			return nil
		}
	}
	return nil
}

func (s *Service) save(c context.Context, item *Item) (resultCode int, itemId string, err error) {
	var uOptions = &options.UpdateOptions{}
	uOptions.SetUpsert(true)
	if item.Id == "" {
		item.Id = uuid.New().String()
	}
	for i, v := range item.Variants {
		if v.Sku == "" {
			sku := fmt.Sprintf("%s-%s-%s-%s-%d",
				strings.ToUpper(item.Brand[:4]),
				strings.ToUpper(item.Name[:4]),
				strings.ToUpper(v.Size),
				strings.ToUpper(v.Color),
				i)
			item.Variants[i].Sku = sku
		}
	}
	inventory := s.collection(inventoryCollection)
	result, err := inventory.UpdateOne(c, bson.M{"_id": bson.M{"$eq": fmt.Sprint(item.Id)}}, bson.M{"$set": item}, uOptions)
	if result != nil {
		if result.UpsertedCount > 0 {
			resultCode = http.StatusCreated
			itemId = fmt.Sprint(result.UpsertedID)
		} else {
			resultCode = http.StatusNoContent
		}
	}
	return
}

func (s *Service) byId(c context.Context, id string) (*Item, error) {
	var item = &Item{}
	sr := s.collection(inventoryCollection).FindOne(c, bson.M{"_id": id})
	err := sr.Decode(item)
	return item, err
}

func (s *Service) queryItems(c context.Context, filter bson.M) (itemVector []Item, err error) {
	inventory := s.collection(inventoryCollection)
	csr, err := inventory.Find(c, filter)
	if err == nil && csr != nil {
		if err = csr.All(c, &itemVector); err != nil {
			s.log.Error(c, "%v", err)
		}
	}
	return
}

func (s *Service) delete(c context.Context, id string) (item *Item, err error) {
	inventory := s.collection(inventoryCollection)
	delInventory := s.collection(inventoryCollectionDeletedItems)
	sr := inventory.FindOne(c, bson.M{"_id": id})
	if sr.Err() != nil && sr.Err() == mongo.ErrNoDocuments {
		err = sr.Err()
	} else {
		err = sr.Decode(&item)
		if err == nil {
			_, err := delInventory.InsertOne(c, item)
			if err == nil {
				_, err = inventory.DeleteOne(c, bson.M{"_id": id})
			}
		}
	}
	return
}

func (s *Service) uploadImage(c context.Context, item *Item, file []byte, ext string) (id string, err error) {
	id, err = s.fs.Upload(c, file, ext)
	if err == nil {
		item.addImage(id)
		_, _, err = s.save(c, item)
	}
	if err != nil {
		s.log.Error(c, "upload of image failed because of %v", err)
	}
	return
}

func (s *Service) uploadVImage(c context.Context, item *Item, index int, variant ItemVariant, file []byte, ext string) (id string, err error) {
	id, err = s.fs.Upload(c, file, ext)
	if err == nil {
		variant.Image = id
		if index >= len(item.Variants) {
			item.addItemVariant(variant)
		} else {
			item.Variants[index].Image = id
		}
		_, _, err = s.save(c, item)
	}
	if err != nil {
		s.log.Error(c, "upload of item variant image failed because of %v", err)
	}
	return
}

func (s *Service) downloadImage(c context.Context, id string) (size int64, file []byte, err error) {
	return s.fs.Download(c, id)
}

func (s *Service) deleteImage(c context.Context, id string, item *Item) (err error) {
	if item.removeImage(id) || item.removeImageFromItemVariant(id) {
		err = s.fs.Remove(c, id)
		if err != nil {
			s.log.Error(c, err.Error())
		}
		_, _, err = s.save(c, item)
	} else {
		err = fmt.Errorf("image [%s]not found", id)
	}
	if err != nil {
		s.log.Error(c, "removal of image [%s] failed because of %v", id, err)
	}
	return
}
