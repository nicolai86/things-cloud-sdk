package thingscloud

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// AccountStatus describes possible thingscloud account statuses
type AccountStatus string

const (
	// AccountStatusActive is for active accounts
	AccountStatusActive AccountStatus = "SYAccountStatusActive"
)

// VerifyResponse contains details about your account
type VerifyResponse struct {
	SLAVersionAccepted string          `json:"SLA-version-accepted"`
	Issues             json.RawMessage `json:"issues"`
	Status             AccountStatus   `json:"status"`
}

// Verify checks that the provided API credentials are valid.
func (c *Client) Verify() (*VerifyResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("/account/%s", c.EMail), nil)
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
	var v VerifyResponse
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(bs, &v)
	return &v, nil
}
