package thingscloud

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// History represents a synchronization stream. It's identified with a uuid v4
type History struct {
	Client            *Client
	LatestServerIndex int

	key string
}

type historyResponse struct {
	LatestSchemaVersion    int  `json:"latest-schema-version"`
	LatestTotalContentSize int  `json:"latest-total-content-size"`
	IsEmpty                bool `json:"is-empty"`
	LatestServerIndex      int  `json:"latest-server-index"`
}

// Sync ensures the history object is able to write to things
func (h *History) Sync() error {
	req, err := http.NewRequest("GET", fmt.Sprintf("/history/%s", h.key), nil)
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
	return nil
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
		} else {
			return nil, fmt.Errorf("http response code: %s", resp.Status)
		}
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
			key:    key,
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
		} else {
			return nil, fmt.Errorf("http response code: %s", resp.Status)
		}
	}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var v createHistoryResponse
	json.Unmarshal(bs, &v)
	return &History{
		Client: c,
		key:    v.Key,
	}, nil
}

// Delete destroys a history
// Note that thingscloud will always return 202, even if the key is unknown
func (h *History) Delete() error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("/account/%s/own-history-keys/%s", h.Client.EMail, h.key), nil)
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
