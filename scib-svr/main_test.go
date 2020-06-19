package main

import (
	"fmt"
	"scib-svr/inventory"
	"testing"
)

func TestMakeUri(t *testing.T) {

	tables := []struct{
		uri string
		params []string
		expected string
	} {
		{"hello",[]string{"world"},"/hello/:world"},
		{inventory.RequestUri, []string{"id"}, "/" + inventory.RequestUri + "/:id"},
		{inventory.RequestUri, []string{}, "/" + inventory.RequestUri},
		{inventory.RequestUri, nil, "/" + inventory.RequestUri},

	}

	for _, table := range tables {
		result := makeUri(table.uri, table.params)
		if result != table.expected {
			t.Errorf("TestMakeUri failed. Got: %s, wanted: %s.", result, table.expected)
		}
	}
}


func makeUri(uri string, params []string) string {
	var paramStr string
	if params != nil {
		for _, param := range params {
			paramStr += "/:" + param
		}
	}
	return fmt.Sprintf("/%s%s",uri, paramStr)
}
