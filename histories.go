package thingscloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"
)

// History represents a synchronization stream. It's identified with a uuid v4
type History struct {
	ID                     string
	Client                 *Client
	LatestServerIndex      int
	LoadedServerIndex      int
	LatestSchemaVersion    int
	EndTotalContentSize    int
	LatestTotalContentSize int
}

type historyResponse struct {
	LatestSchemaVersion    int  `json:"latest-schema-version"`
	LatestTotalContentSize int  `json:"latest-total-content-size"`
	IsEmpty                bool `json:"is-empty"`
	LatestServerIndex      int  `json:"latest-server-index"`
}

// Sync ensures the history object is able to write to things
func (h *History) Sync() error {
	req, err := http.NewRequest("GET", fmt.Sprintf("/version/1/history/%s/items", h.ID), nil)
	if err != nil {
		return err
	}
	query := req.URL.Query()
	query.Add("start-index", strconv.Itoa(h.LatestServerIndex))
	req.URL.RawQuery = query.Encode()
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
	var v itemsResponse
	json.Unmarshal(bs, &v)
	h.LatestServerIndex = v.CurrentItemIndex
	h.LatestSchemaVersion = v.SchemaVersion
	h.LatestTotalContentSize = v.LatestTotalContentSize
	return nil
}

// History requests a specific history
func (c *Client) History(id string) (*History, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("/version/1/history/%s", id), nil)
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
		Client:                 c,
		ID:                     id,
		LatestServerIndex:      h.LatestServerIndex,
		LatestSchemaVersion:    h.LatestSchemaVersion,
		LatestTotalContentSize: h.LatestTotalContentSize,
	}, nil
}

type v1historyResponse struct {
	Key                 string `json:"history-key"`
	LatestServerIndex   int    `json:"latest-server-index"`
	IsEmpty             bool   `json:"is-empty"`
	LatestSchemaVersion int    `json:"latest-schema-version"`
}

// OwnHistory returns the clients own history
func (c *Client) OwnHistory() (*History, error) {
	resp, err := c.Verify()
	if err != nil {
		return nil, err
	}

	return &History{
		Client: c,
		ID:     resp.HistoryKey,
	}, nil
}

// Histories requests all known history keys
func (c *Client) Histories() ([]*History, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("/version/1/account/%s/own-history-keys", c.EMail), nil)
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
	req, err := http.NewRequest("POST", fmt.Sprintf("/version/1/account/%s/own-history-keys", c.EMail), nil)
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
	req, err := http.NewRequest("DELETE", fmt.Sprintf("/version/1/account/%s/own-history-keys/%s", h.Client.EMail, h.ID), nil)
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

type commitResponse struct {
	ServerHeadIndex int `json:"server-head-index"`
}

// Identifiable abstracts different thingscloud write requests. As we need to provide a map
// indexed by UUID, all we care about is the ID of the change, not the change itself
type Identifiable interface {
	UUID() string
}

func (h *History) Write(items ...Identifiable) error {
	m := map[string]interface{}{}
	for _, item := range items {
		m[item.UUID()] = item
	}
	bs, err := json.Marshal(m)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("/version/1/history/%s/commit", h.ID), bytes.NewReader(bs))
	req.Header.Add("Schema", "301")
	req.Header.Add("Push-Priority", "5")
	req.Header.Add("App-Instance-Id", "-com.culturedcode.ThingsMac")
	req.Header.Add("App-Id", "com.culturedcode.ThingsMac")
	req.Header.Add("Content-Encoding", "UTF-8")
	req.Header.Add("Host", "cloud.culturedcode.com")
	req.Header.Add("Accept", "application/json")
	query := req.URL.Query()
	query.Add("ancestor-index", strconv.Itoa(h.LatestServerIndex))
	query.Add("_cnt", "1")
	req.URL.RawQuery = query.Encode()
	if err != nil {
		return err
	}
	resp, err := h.Client.do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		bs, _ := httputil.DumpResponse(resp, true)
		log.Println(string(bs))
		return fmt.Errorf("Write failed: %d", resp.StatusCode)
	}
	rs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var w commitResponse
	json.Unmarshal(rs, &w)
	h.LatestServerIndex = w.ServerHeadIndex
	return nil
}
