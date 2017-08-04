package main

import (
	"net/http/httptest"
	"testing"

	thingscloud "github.com/nicolai86/things-cloud-sdk"
)

func TestHTTPItemsHandler_List(t *testing.T) {
	email := "max@mustermann.de"
	password := "secret"
	history := History{
		ID:                  "A",
		LatestSchemaVersion: 300,
		LatestServerIndex:   20,
		IsEmpty:             false,
	}
	s := httptest.NewServer(New(nil, &fakeHistoryHandler{[]History{history}, nil, nil, nil}, &fakeItemsHandler{
		items: []map[string]Item{
			map[string]Item{
				"A": Item{},
			},
		},
	}))
	defer s.Close()
	c := thingscloud.New(s.URL, email, password)
	h := &thingscloud.History{
		Client: c,
		ID:     history.ID,
	}
	items, err := h.Items(thingscloud.ItemsOptions{StartIndex: 0})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err.Error())
	}
	if len(items) != 1 {
		t.Errorf("Expected %d items but got %d", 1, len(items))
	}
}

func TestHTTPItemsHandler_Write(t *testing.T) {
	email := "max@mustermann.de"
	password := "secret"
	history := History{
		ID:                  "A",
		LatestSchemaVersion: 300,
		LatestServerIndex:   0,
		IsEmpty:             false,
	}
	handler := &fakeItemsHandler{
		items: []map[string]Item{},
	}
	s := httptest.NewServer(New(nil, &fakeHistoryHandler{[]History{history}, nil, nil, nil}, handler))
	defer s.Close()
	c := thingscloud.New(s.URL, email, password)
	h := &thingscloud.History{
		Client: c,
		ID:     history.ID,
	}
	if err := h.Write(thingscloud.TaskActionItem{
		Item: thingscloud.Item{
			Kind:   thingscloud.ItemKindTask,
			Action: thingscloud.ItemActionDeleted,
			UUID:   "54152210-ABFA-4F9F-81AC-7F50FBDEDC2G",
		},
		P: thingscloud.TaskActionItemPayload{},
	}); err != nil {
		t.Fatalf("Unexpected error: %v", err.Error())
	}

	if len(handler.items) != 1 {
		t.Errorf("Expected write to add item, but didn't")
	}
}
