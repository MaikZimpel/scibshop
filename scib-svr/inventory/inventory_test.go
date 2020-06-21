package inventory

import (
	"net/http"
	"reflect"
	"scib-svr/datastore"
	"scib-svr/logging"
	"testing"
)

type InventoryServiceTest struct {
	service *Service
}

func NewInventoryServiceTest() *InventoryServiceTest {
	return &InventoryServiceTest{
		NewService(datastore.New(), logging.New()),
	}
}

func TestCRUD(t *testing.T) {
	test := NewInventoryServiceTest()
	item := NewItem()
	item.Color = "navy blue"
	item.Size = "XL"
	item.Brand = "Swimma"
	item.Categories = []string{"Swimwear", "Sport", "Clothing"}
	item.Description = "Navy Blue Swimsuit for women"
	item.Name = "Swimma One WC"
	item.Sku = "swi-one-wc-nav-blue-xl"
	item.Cnt = 15
	item.Available = true
	item.Stockable = true
	status, id, err := test.service.save(item)
	if status != http.StatusCreated {
		t.Errorf("save failed because of wrong status: expected: %v got %v", http.StatusCreated, status)
	}
	if err != nil {
		t.Errorf("SAVE failed because of an error %v", err)
	}
	item.Id = id
	readItem, err := test.service.itemById(id)
	if err != nil {
		t.Errorf("READ failed because of an error %v", err)
	}
	if !reflect.DeepEqual(item, readItem) {
		t.Errorf("READ failed expected: %v got %v", item, readItem)
	}
	item.Name = "Swimma Two AC"
	status, id, err = test.service.save(item)
	if status != http.StatusNoContent {
		t.Errorf("save failed because of wrong status: expected: %v got %v", http.StatusNoContent, status)
	}
	if err != nil {
		t.Errorf("SAVE failed because of an error %v", err)
	}
	if reflect.DeepEqual(item, readItem) {
		t.Errorf("READ failed expected: %v got %v", item, readItem)
	}
	readItem, err = test.service.itemById(id)
	if err != nil {
		t.Errorf("READ failed because of an error %v", err)
	}
	if !reflect.DeepEqual(item, readItem) {
		t.Errorf("READ failed expected: %v got %v", item, readItem)
	}
	deleted, err := test.service.delete(id)
	if err != nil {
		t.Errorf("READ failed because of an error %v", err)
	}
	if !reflect.DeepEqual(deleted, readItem) {
		t.Errorf("READ failed expected: %v got %v", item, readItem)
	}
	readItem, err = test.service.itemById(id)
	if err == nil {
		t.Errorf("DELETE failed expected: %v got %v", nil, readItem)
	}
}


