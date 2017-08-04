package main

import (
	"bytes"
	"encoding/json"
)

// AccountHandler wraps various account handling tasks
type AccountHandler interface {
	SignupHandler
	PasswordChangeHandler
	AccountCloseHandler
}

type changePasswordRequestBody struct {
	Password *string `json:"password"`
}

type signupRequestBody struct {
	Password           *string `json:"password"`
	SLAVersionAccepted *string `json:"SLA-version-accepted"`
}

type verificationRequestBody struct {
	ConfirmationCode *string `json:"confirmation-code"`
}

// SignupHandler handles new account creation and confirmation.
type SignupHandler interface {
	Signup(email, password string) (string, error)
	DeliverConfirmationCode(email, code string) error
	Confirm(email, code string) error
}

// PasswordChangeHandler handles existing account password modification
type PasswordChangeHandler interface {
	ChangePassword(email, password string) error
}

// AccountCloseHandler handles the deletion of existing accounts
type AccountCloseHandler interface {
	Close(email string) error
}

func isPasswordChange(body []byte) bool {
	return !isSignup(body) && bytes.Contains(body, []byte("password"))
}

func isSignup(body []byte) bool {
	return bytes.Contains(body, []byte("password")) && bytes.Contains(body, []byte("SLA-version-accepted"))
}

func isConfirmation(body []byte) bool {
	return bytes.Contains(body, []byte("confirmation-code"))
}

func (h *httpAccountHandler) signup(bs []byte, email string) error {
	body := signupRequestBody{}
	if err := json.Unmarshal(bs, &body); err != nil {
		return err
	}
	code, err := h.Signup(email, *body.Password)
	if err != nil {
		return err
	}
	return h.DeliverConfirmationCode(email, code)
}

func (h *httpAccountHandler) confirm(bs []byte, email string) error {
	body := verificationRequestBody{}
	if err := json.Unmarshal(bs, &body); err != nil {
		return err
	}
	err := h.Confirm(email, *body.ConfirmationCode)
	if err != nil {
		return err
	}
	return nil
}

func (h *httpAccountHandler) passwordChange(bs []byte, email string) error {
	body := changePasswordRequestBody{}
	if err := json.Unmarshal(bs, &body); err != nil {
		return err
	}
	err := h.ChangePassword(email, *body.Password)
	if err != nil {
		return err
	}
	return nil
}
