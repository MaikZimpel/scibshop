package inventory

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"scib-svr/datastore"
	"scib-svr/logging"
	"strconv"
	"strings"
)

type Item struct {
	Id          string   `json:"id"`
	Upc         string   `json:"upc"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Categories  []string `json:"categories"`
	Brand       string   `json:"brand"`
	Size        string   `json:"size"`
	Color       string   `json:"color"`
	Price       int      `json:"price"`
	Images      []string `json:"images"`
	Supplier    string   `json:"supplier"`
	Sku         string   `json:"sku"`
	Cnt         int      `json:"cnt"`
	Stockable   bool     `json:"stockable"`
	Available   bool     `json:"available"`
}

const inventoryCollection = "inventory"

func (i *Item) String() string {
	b, _ := json.Marshal(i)
	return string(b)
}

func NewItem() *Item {
	return &Item{}
}

type Service struct {
	ds datastore.Datastore
	log logging.Logger
}

func NewService(datastore datastore.Datastore) *Service {
	log := logging.New()
	return &Service{datastore, log}
}

func (s *Service) save(item *Item) (int, string, error) {
	if item.Id == "" {
		item.Id = uuid.New().String()
	}
	sts, err := s.ds.Upsert(context.Background(), inventoryCollection, item.Id, &item, false)
	return sts, item.Id, err
}

func (s *Service) allItems(stockableOnly bool) (itemVector []*Item, err error) {
	s.log.Debug(context.Background(), "call allItems")
	whereClauses := func() []datastore.WhereClause {
		if stockableOnly {
			return []datastore.WhereClause{{Path: "Stockable", Op: "==", Value: stockableOnly}}
		} else {
			return make([]datastore.WhereClause, 0)
		}
	}()
	resultSlice, err := s.ds.Read(context.Background(), inventoryCollection, "", whereClauses)
	if err == nil && resultSlice != nil {
		for _, resultMap := range resultSlice {
			item, err := mapToItem(resultMap)
			if err == nil {
				itemVector = append(itemVector, item)
			}
		}
	}
	if err != nil {
		s.log.Error(context.Background(), "%s", err)
	}
	return
}

func (s *Service) itemById(id string) (item *Item, err error) {
	res, err := s.ds.Read(context.Background(), inventoryCollection, id, make([]datastore.WhereClause, 0))
	if err == nil && res != nil {
		item, err = mapToItem(res[0])
	}
	return
}

func (s *Service) delete(id string) (item *Item, err error) {
	res, err := s.ds.Delete(context.Background(), inventoryCollection, id)
	if err == nil && res != nil {
		item, err = mapToItem(res)
	}
	return
}

func mapToItem(in map[string]interface{}) (*Item, error) {
	var item Item
	for k, v := range in {
		err := set(k, v, &item)
		if err != nil {
			var errStrings []string
			for _, er := range err {
				errStrings = append(errStrings, er.Error())
			}
			return nil, fmt.Errorf(strings.Join(errStrings, "\n"))
		}
	}
	return &item, nil
}

func set(field string, value interface{}, item *Item) (err []error) {
	var e error
	log := logging.New()
	log.Debug(context.Background(), "%+v", value)
	switch field {
	case "Id":
		item.Id = fmt.Sprintf("%s", value)
	case "Upc":
		item.Upc = fmt.Sprintf("%s", value)
	case "Name":
		item.Name = fmt.Sprintf("%s", value)
	case "Description":
		item.Description = fmt.Sprintf("%s", value)
	case "Categories":
		if value != nil {
			item.Categories = toStringArray(value.([]interface{}))
		}
	case "Brand":
		item.Brand = fmt.Sprintf("%s", value)
	case "Size":
		item.Size = fmt.Sprintf("%s", value)
	case "Color":
		item.Color = fmt.Sprintf("%s", value)
	case "Images":
		if value != nil {
			item.Images = toStringArray(value.([]interface{}))
		}
	case "Supplier":
		item.Supplier = fmt.Sprintf("%s", value)
	case "Sku":
		item.Sku = fmt.Sprintf("%s", value)
	case "Price":
		if value != nil {
			item.Price, e = toInt(value.(map[string]interface{}))
			if e != nil {
				err = append(err, e)
			}
		}
	case "Cnt":
		if value != nil {
			item.Cnt, e = toInt(value.(map[string]interface{}))
			if e != nil {
				err = append(err, e)
			}
		}
	case "Stockable":
		if value != nil {
			item.Stockable = value.(bool)
		}
	case "Available":
		if value != nil {
			item.Available = value.(bool)
		}
	}
	return
}

func toStringArray(i []interface{}) []string {
	s := make([]string, len(i))
	for x, v := range i {
		s[x] = fmt.Sprintf("%s", v)
	}
	return s
}

func toInt(i map[string]interface{}) (int, error) {
	return strconv.Atoi(fmt.Sprintf("%s",i["$numberInt"]))
}
