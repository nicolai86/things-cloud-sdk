package thingscloud

import (
	"fmt"
	"testing"
)

func TestClient_Histories(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		server := fakeServer(fakeResponse{200, "histories-success.json"})
		defer server.Close()

		c := New(fmt.Sprintf("http://%s", server.Listener.Addr().String()), "martin@example.com", "")
		hs, err := c.Histories()
		if err != nil {
			t.Fatalf("Expected history request to succeed, but didn't: %q", err.Error())
		}
		if len(hs) != 1 {
			t.Errorf("Expected to receive %d histories, but got %d", 1, len(hs))
		}
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()
		server := fakeServer(fakeResponse{401, "error.json"})
		defer server.Close()

		c := New(fmt.Sprintf("http://%s", server.Listener.Addr().String()), "unknown@example.com", "")
		_, err := c.Histories()
		if err == nil {
			t.Error("Expected history request to fail, but didn't")
		}
	})
}

func TestClient_CreateHistory(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		server := fakeServer(fakeResponse{200, "create-history-success.json"})
		defer server.Close()

		c := New(fmt.Sprintf("http://%s", server.Listener.Addr().String()), "martin@example.com", "")
		h, err := c.CreateHistory()
		if err != nil {
			t.Fatalf("Expected request to succeed, but didn't: %q", err.Error())
		}
		if h.ID != "33333abb-bfe4-4b03-a5c9-106d42220c72" {
			t.Fatalf("Expected key %s but got %s", "33333abb-bfe4-4b03-a5c9-106d42220c72", h.ID)
		}
	})
}

func TestHistory_Delete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		server := fakeServer(fakeResponse{202, "create-history-success.json"})
		defer server.Close()

		c := New(fmt.Sprintf("http://%s", server.Listener.Addr().String()), "martin@example.com", "")
		h := History{Client: c, ID: "33333abb-bfe4-4b03-a5c9-106d42220c72"}
		err := h.Delete()
		if err != nil {
			t.Fatalf("Expected request to succeed, but didn't: %q", err.Error())
		}
	})
}

func TestHistory_Sync(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		server := fakeServer(fakeResponse{200, "history-success.json"})
		defer server.Close()

		c := New(fmt.Sprintf("http://%s", server.Listener.Addr().String()), "martin@example.com", "")
		h := History{Client: c, ID: "33333abb-bfe4-4b03-a5c9-106d42220c72"}
		err := h.Sync()
		if err != nil {
			t.Fatalf("Expected request to succeed, but didn't: %q", err.Error())
		}
		if h.LatestServerIndex != 27 {
			t.Errorf("Expected LatestServerIndex of %d, but got %d", 27, h.LatestServerIndex)
		}
	})
}
