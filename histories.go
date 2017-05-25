package thingscloud

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// History represents a synchronization stream. It's identified with a uuid v4
type History struct {
	Client *Client

	key string
}

// Histories requests all known history keys
func (c *Client) Histories() ([]*History, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("/account/%s/own-history-keys", c.EMail), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Do(req)
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
	var keys []string
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
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
	resp, err := c.Do(req)
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
	var v createHistoryResponse
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
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
	resp, err := h.Client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("http response code: %s", resp.Status)
	}
	return nil
}
