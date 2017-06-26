package thingscloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type accountRequestBody struct {
	Password           string `json:"password,omitempty"`
	SLAVersionAccepted string `json:"SLA-version-accepted"`
}

// SignUp creates a new thingscloud account and returns a configured client
func (c *Client) SignUp(email, password string) (*Client, error) {
	data, err := json.Marshal(accountRequestBody{
		Password:           password,
		SLAVersionAccepted: "https://thingscloud.appspot.com/sla/v5.html",
	})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PUT", fmt.Sprintf("/account/%s", email), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			return nil, ErrUnauthorized
		}
		return nil, fmt.Errorf("http response code: %s", resp.Status)
	}

	return New(c.Endpoint, email, password), nil
}

// SetAccountPassword allows you to change your account password.
// Because things does not work with sessions you need to create a new client instance after
// executing this method
func (c *Client) SetAccountPassword(newPassword string) (*Client, error) {
	data, err := json.Marshal(accountRequestBody{
		Password: newPassword,
	})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PUT", fmt.Sprintf("/account/%s", c.EMail), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			return nil, ErrUnauthorized
		}
		return nil, fmt.Errorf("http response code: %s", resp.Status)
	}

	return New(c.Endpoint, c.EMail, newPassword), nil
}
