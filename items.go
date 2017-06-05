package thingscloud

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Item is an event in thingscloud. Every action inside things generates an Item.
// Common items are the creation of a task, area or checklist, as well as modifying attributes
// or marking things as done.
type Item struct {
	UUID   string          `json:"-"`
	P      json.RawMessage `json:"p"`
	Kind   ItemKind        `json:"e"`
	Action ItemAction      `json:"t"`
}

type itemsResponse struct {
	Items                  []map[string]Item `json:"items"`
	LatestTotalContentSize int               `json:"latest-total-content-size"`
	StartTotalContentSize  int               `json:"start-total-content-size"`
	EndTotalContentSize    int               `json:"end-total-content-size"`
	SchemaVersion          int               `json:"schema"`
	CurrentItemIndex       int               `json:"current-item-index"`
}

// ItemsOptions allows a client to pickup changes from a specific index
type ItemsOptions struct {
	StartIndex int
}

// Items fetches changes from thingscloud. Every change contains multiple items which have been modified.
// The Items method unwraps these objects and returns a list instead.
//
// Note that if a item was changed multiple times it will be present multiple times in the result too.
func (h *History) Items(opts ItemsOptions) ([]Item, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("/history/%s/items", h.ID), nil)

	values := req.URL.Query()
	values.Set("start-index", strconv.Itoa(opts.StartIndex))
	req.URL.RawQuery = values.Encode()

	if err != nil {
		return nil, err
	}
	resp, err := h.Client.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http response code: %s", resp.Status)
	}

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var v itemsResponse
	if err := json.Unmarshal(bs, &v); err != nil {
		return nil, err
	}
	var items = []Item{}
	for _, m := range v.Items {
		for id, item := range m {
			item.UUID = id
			items = append(items, item)
		}
	}
	return items, nil
}
