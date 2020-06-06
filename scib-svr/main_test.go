package main

import (
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
