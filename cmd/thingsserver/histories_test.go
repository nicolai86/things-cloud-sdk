package main

import (
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	thingscloud "github.com/nicolai86/things-cloud-sdk"
)

func TestHTTPHistoryHandler_Histories(t *testing.T) {
	email := "max@mustermann.de"
	password := "secret"
	historyIDs := []string{"A", "B"}
	histories := make([]History, len(historyIDs))
	for i, id := range historyIDs {
		histories[i] = History{ID: id}
	}
	s := httptest.NewServer(New(nil, &fakeHistoryHandler{histories, nil, nil, nil}, nil))
	defer s.Close()
	c := thingscloud.New(s.URL, email, password)
	hs, err := c.Histories()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err.Error())
	}
	hsIDs := make([]string, len(hs))
	for i, h := range hs {
		hsIDs[i] = h.ID
	}
	if !cmp.Equal(hsIDs, historyIDs) {
		t.Fatalf("Expected %v, but got %v", historyIDs, hsIDs)
	}
}

func TestHTTPHistoryHandler_History(t *testing.T) {
	email := "max@mustermann.de"
	password := "secret"
	history := History{
		ID:                  "A",
		LatestSchemaVersion: 300,
		LatestServerIndex:   20,
		IsEmpty:             false,
	}
	s := httptest.NewServer(New(nil, &fakeHistoryHandler{[]History{history}, nil, nil, nil}, nil))
	defer s.Close()
	c := thingscloud.New(s.URL, email, password)
	h, err := c.History(history.ID)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err.Error())
	}
	if h.ID != history.ID {
		t.Fatalf("Expected history %s but got %s", history.ID, h.ID)
	}
	if h.LatestSchemaVersion != history.LatestSchemaVersion {
		t.Errorf("Expected schema version %d but got %d", history.LatestSchemaVersion, h.LatestSchemaVersion)
	}
	if h.LatestServerIndex != history.LatestServerIndex {
		t.Errorf("Expected server index %d but got %d", history.LatestServerIndex, h.LatestServerIndex)
	}
}

func TestHTTPHistoryHandler_Sync(t *testing.T) {
	email := "max@mustermann.de"
	password := "secret"
	history := History{
		ID:                  "A",
		LatestSchemaVersion: 300,
		LatestServerIndex:   20,
		IsEmpty:             false,
	}
	s := httptest.NewServer(New(nil, &fakeHistoryHandler{[]History{history}, nil, nil, nil}, nil))
	defer s.Close()
	c := thingscloud.New(s.URL, email, password)
	h := &thingscloud.History{
		Client: c,
		ID:     history.ID,
	}
	if err := h.Sync(); err != nil {
		t.Fatalf("Unexpected error: %v", err.Error())
	}
	if h.LatestSchemaVersion != history.LatestSchemaVersion {
		t.Errorf("Expected schema version %d but got %d", history.LatestSchemaVersion, h.LatestSchemaVersion)
	}
	if h.LatestServerIndex != history.LatestServerIndex {
		t.Errorf("Expected server index %d but got %d", history.LatestServerIndex, h.LatestServerIndex)
	}
}

func TestHTTPHistoryHandler_Delete(t *testing.T) {
	email := "max@mustermann.de"
	password := "secret"
	history := History{
		ID:                  "A",
		LatestSchemaVersion: 300,
		LatestServerIndex:   20,
		IsEmpty:             false,
	}
	handler := &fakeHistoryHandler{[]History{history}, nil, nil, nil}
	s := httptest.NewServer(New(nil, handler, nil))
	defer s.Close()
	c := thingscloud.New(s.URL, email, password)
	h := &thingscloud.History{
		Client: c,
		ID:     history.ID,
	}

	if err := h.Delete(); err != nil {
		t.Fatalf("Unexpected error: %v", err.Error())
	}
	if len(handler.histories) != 0 {
		t.Errorf("Expected delete to remove history, but didn't")
	}
}

func TestHTTPHistoryHandler_Create(t *testing.T) {
	email := "max@mustermann.de"
	password := "secret"
	handler := &fakeHistoryHandler{[]History{}, nil, nil, nil}
	s := httptest.NewServer(New(nil, handler, nil))
	defer s.Close()
	c := thingscloud.New(s.URL, email, password)
	h, err := c.CreateHistory()
	_ = h
	if err != nil {
		t.Fatalf("Unexpected error: %v", err.Error())
	}
	if len(handler.histories) != 1 {
		t.Errorf("Expected create to add history, but didn't")
	}
	if h.ID != "1" {
		t.Errorf("Expected create to add history, but didn't")
	}
}
