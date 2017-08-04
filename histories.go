package thingscloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// History represents a synchronization stream. It's identified with a uuid v4
type History struct {
	ID                  string
	Client              *Client
	LatestServerIndex   int
	LatestSchemaVersion int
}

type historyResponse struct {
	LatestSchemaVersion    int  `json:"latest-schema-version"`
	LatestTotalContentSize int  `json:"latest-total-content-size"`
	IsEmpty                bool `json:"is-empty"`
	LatestServerIndex      int  `json:"latest-server-index"`
}

// Sync ensures the history object is able to write to things
func (h *History) Sync() error {
	req, err := http.NewRequest("GET", fmt.Sprintf("/history/%s", h.ID), nil)
	if err != nil {
		return err
	}
	resp, err := h.Client.do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http response code: %s", resp.Status)
	}

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var v historyResponse
	json.Unmarshal(bs, &v)
	h.LatestServerIndex = v.LatestServerIndex
	h.LatestSchemaVersion = v.LatestSchemaVersion
	return nil
}

// History requests a specific history
func (c *Client) History(id string) (*History, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("/history/%s", id), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			return nil, ErrUnauthorized
		}
		return nil, fmt.Errorf("http response code: %s", resp.Status)
	}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	h := historyResponse{}
	if err := json.Unmarshal(bs, &h); err != nil {
		return nil, err
	}

	return &History{
		Client:              c,
		ID:                  id,
		LatestServerIndex:   h.LatestServerIndex,
		LatestSchemaVersion: h.LatestSchemaVersion,
	}, nil
}

// Histories requests all known history keys
func (c *Client) Histories() ([]*History, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("/account/%s/own-history-keys", c.EMail), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			return nil, ErrUnauthorized
		}
		return nil, fmt.Errorf("http response code: %s", resp.Status)
	}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var keys []string
	json.Unmarshal(bs, &keys)

	var histories = make([]*History, len(keys))
	for i, key := range keys {
		histories[i] = &History{
			Client: c,
			ID:     key,
		}
	}
	return histories, nil
}

type createHistoryResponse struct {
	Key string `json:"new-history-key"`
}

// CreateHistory requests a new history key
func (c *Client) CreateHistory() (*History, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("/account/%s/own-history-keys", c.EMail), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			return nil, ErrUnauthorized
		}
		return nil, fmt.Errorf("http response code: %s", resp.Status)

	}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var v createHistoryResponse
	json.Unmarshal(bs, &v)
	return &History{
		Client: c,
		ID:     v.Key,
	}, nil
}

// Delete destroys a history
// Note that thingscloud will always return 202, even if the key is unknown
func (h *History) Delete() error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("/account/%s/own-history-keys/%s", h.Client.EMail, h.ID), nil)
	if err != nil {
		return err
	}
	resp, err := h.Client.do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("http response code: %s", resp.Status)
	}
	return nil
}

type writeRequest struct {
	AppID            string                   `json:"app-id"`
	AppInstanceID    string                   `json:"app-instance-id"`
	CurrentItemIndex int                      `json:"current-item-index"`
	Items            []map[string]interface{} `json:"items"`
	PushPriority     int                      `json:"push-priority"`
	Schema           int                      `json:"schema"`
}

type writeResponse struct {
	CurrentItemIndex int `json:"current-item-index"`
}

// Identifiable abstracts different thingscloud write requests. As we need to provide a map
// indexed by UUID, all we care about is the ID of the change, not the change itself
type Identifiable interface {
	UUID() string
}

func (h *History) Write(items ...Identifiable) error {
	var v = writeRequest{
		AppID:            "com.culturedcode.ThingsMac",
		AppInstanceID:    "-com.culturedcode.ThingsMac",
		CurrentItemIndex: h.LatestServerIndex,
		PushPriority:     10,
		Schema:           h.LatestSchemaVersion,
		Items:            []map[string]interface{}{},
	}
	for _, item := range items {
		m := map[string]interface{}{}
		m[item.UUID()] = item
		v.Items = append(v.Items, m)
	}
	bs, err := json.Marshal(v)
	fmt.Printf("\n\n%s\n\n", string(bs))
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("/history/%s/items", h.ID), bytes.NewReader(bs))
	if err != nil {
		return err
	}
	resp, err := h.Client.do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Write failed: %d", resp.StatusCode)
	}
	rs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var w writeResponse
	json.Unmarshal(rs, &w)
	h.LatestServerIndex = w.CurrentItemIndex
	return nil
}
