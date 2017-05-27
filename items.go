package thingscloud

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

//go:generate stringer -type ItemAction,TaskStatus

type Timestamp time.Time

func (t *Timestamp) UnmarshalJSON(bs []byte) error {
	var d float64
	if err := json.Unmarshal(bs, &d); err != nil {
		return err
	}
	*t = Timestamp(time.Unix(int64(d), 0))
	return nil
}

func (t *Timestamp) Format(layout string) string {
	return time.Time(*t).Format(layout)
}

type ItemAction int

const (
	// ActionCreated is used to indicate a new Item was created
	ActionCreated ItemAction = iota
	// ActionModified is used to indicate an existing Item was modified
	ActionModified ItemAction = 1
	// ActionDeleted is used as a tombstone for an Item
	ActionDeleted ItemAction = 2
)

type TaskStatus int

const (
	// TaskStatusPending indicates a new task
	TaskStatusPending TaskStatus = iota
	// TaskStatusCompleted indicates a completed task
	TaskStatusCompleted TaskStatus = 3
	// TaskStatusCanceled indicates a canceled task
	TaskStatusCanceled TaskStatus = 2
)

type ItemKind string

var (
	ItemKindChecklist ItemKind = "ChecklistItem"
	ItemKindTask      ItemKind = "Task3"
	ItemKindArea      ItemKind = "Area2"
	ItemKindSettings  ItemKind = "Settings3"
	ItemKindTag       ItemKind = "Tag3"
)

type Item struct {
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

type ItemsOptions struct {
	StartIndex int
}

func (h *History) Items(opts ItemsOptions) ([]map[string]Item, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("/history/%s/items", h.key), nil)

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
	return v.Items, nil
}
