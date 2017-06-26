package thingscloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type setPasswordBody struct {
	Password string `json:"password"`
}

// SetAccountPassword allows you to change your account password.
// Because things does not work with sessions you need to create a new client instance after
// executing this method
func (c *Client) SetAccountPassword(newPassword string) error {
	data, err := json.Marshal(setPasswordBody{
		Password: newPassword,
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", fmt.Sprintf("/account/%s", c.EMail), bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	resp, err := c.do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			return ErrUnauthorized
		}
		return fmt.Errorf("http response code: %s", resp.Status)
	}

	return nil
}
