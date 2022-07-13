package thingscloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type accountRequestBody struct {
	Password           string `json:"password,omitempty"`
	SLAVersionAccepted string `json:"SLA-version-accepted,omitempty"`
	ConfirmationCode   string `json:"confirmation-code,omitempty"`
}

// AccountService allows account specific interaction with thingscloud
type AccountService service

// Delete deletes your current thingscloud account. This cannot be reversed
func (s *AccountService) Delete() error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("/version/1/account/%s", s.client.EMail), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Password %s", s.client.password))
	resp, err := s.client.do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		if resp.StatusCode == http.StatusUnauthorized {
			return ErrUnauthorized
		}
		return fmt.Errorf("http response code: %s", resp.Status)
	}
	return nil
}

func (s *AccountService) AcceptSLA() error {
	data, err := json.Marshal(accountRequestBody{
		SLAVersionAccepted: "https://cloud.culturedcode.com/sla/v1.5-rich.html?language=en",
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", fmt.Sprintf("/version/1/account/%s", s.client.EMail), bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Password %s", s.client.password))
	resp, err := s.client.do(req)
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

// Confirm finishes the account creation by providing the email token send by thingscloud
func (s *AccountService) Confirm(code string) error {
	data, err := json.Marshal(accountRequestBody{
		ConfirmationCode: code,
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", fmt.Sprintf("/version/1/account/%s", s.client.EMail), bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Password %s", s.client.password))
	resp, err := s.client.do(req)
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

// SignUp creates a new thingscloud account and returns a configured client
func (s *AccountService) SignUp(email, password string) (*Client, error) {
	data, err := json.Marshal(accountRequestBody{
		Password:           password,
	})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PUT", fmt.Sprintf("/version/1/account/%s", email), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	resp, err := s.client.do(req)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("http response code: %s", resp.Status)
	}

	return New(s.client.Endpoint, email, password), nil
}

// ChangePassword allows you to change your account password.
// Because things does not work with sessions you need to create a new client instance after
// executing this method
func (s *AccountService) ChangePassword(newPassword string) (*Client, error) {
	data, err := json.Marshal(accountRequestBody{
		Password: newPassword,
	})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PUT", fmt.Sprintf("/version/1/account/%s", s.client.EMail), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	resp, err := s.client.do(req)
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

	return New(s.client.Endpoint, s.client.EMail, newPassword), nil
}
