package main

import (
	"net/http/httptest"
	"testing"

	thingscloud "github.com/nicolai86/things-cloud-sdk"
)

func TestHTTPAccountHandler_Signup(t *testing.T) {
	email := "max@mustermann.de"
	password := "nopass"
	handler := &fakeSignupHandler{
		code: "asd123",
	}
	s := httptest.NewServer(New(&fakeAccountHandler{nil, handler, nil}, nil, nil))
	defer s.Close()
	c := thingscloud.New(s.URL, "", "")
	_, err := c.Accounts.SignUp(email, password)
	if err != nil {
		t.Fatalf("ups: %#vs", err)
	}

	if handler.signup.email != email {
		t.Errorf("Expected signup email %q but got %q", email, handler.signup.email)
	}
	if handler.signup.password != password {
		t.Errorf("Expected signup password %q but got %q", password, handler.signup.password)
	}
	if handler.confirmation.code != handler.code {
		t.Errorf("Expected signup confirmation code %q but got %q", handler.confirmation.code, handler.code)
	}
}

func TestHTTPAccountHandler_Confirm(t *testing.T) {
	email := "max@mustermann.de"
	confirmationCode := "secret"
	handler := &fakeSignupHandler{}
	s := httptest.NewServer(New(&fakeAccountHandler{nil, handler, nil}, nil, nil))
	defer s.Close()
	c := thingscloud.New(s.URL, email, "")
	err := c.Accounts.Confirm(confirmationCode)
	if err != nil {
		t.Fatalf("ups: %#vs", err)
	}

	if handler.confirmation.email != email {
		t.Errorf("Expected confirmation email %q but got %q", email, handler.confirmation.email)
	}
	if handler.confirmation.code != confirmationCode {
		t.Errorf("Expected confirmation code %q but got %q", confirmationCode, handler.confirmation.code)
	}
}

func TestHTTPAccountHandler_PasswordChange(t *testing.T) {
	email := "max@mustermann.de"
	newPassword := "secret"
	handler := &fakePasswordChangeHandler{}
	s := httptest.NewServer(New(&fakeAccountHandler{handler, nil, nil}, nil, nil))
	defer s.Close()
	c := thingscloud.New(s.URL, email, "")
	_, err := c.Accounts.ChangePassword(newPassword)
	if err != nil {
		t.Fatalf("ups: %#vs", err)
	}

	if handler.email != email {
		t.Errorf("Expected email %q but got %q", email, handler.email)
	}
	if handler.newPassword != newPassword {
		t.Errorf("Expected password %q but got %q", newPassword, handler.newPassword)
	}
}

func TestHTTPAccountHandler_Delete(t *testing.T) {
	email := "max@mustermann.de"
	handler := &fakeAccountCloseHandler{}
	s := httptest.NewServer(New(&fakeAccountHandler{nil, nil, handler}, nil, nil))
	defer s.Close()
	c := thingscloud.New(s.URL, email, "")
	err := c.Accounts.Delete()
	if err != nil {
		t.Fatalf("ups: %#vs", err)
	}
	if handler.email != email {
		t.Errorf("Expected email %q but got %q", email, handler.email)
	}
}
