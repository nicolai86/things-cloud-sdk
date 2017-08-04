package main

import "fmt"

type signup struct {
	email    string
	password string
}

type confirmation struct {
	email string
	code  string
}

type fakeAccountHandler struct {
	*fakePasswordChangeHandler
	*fakeSignupHandler
	*fakeAccountCloseHandler
}

type fakePasswordChangeHandler struct {
	email       string
	newPassword string
	changeError error
}

func (f *fakePasswordChangeHandler) ChangePassword(email, password string) error {
	if f.changeError != nil {
		return f.changeError
	}

	f.email = email
	f.newPassword = password
	return nil
}

type fakeSignupHandler struct {
	signup        signup
	confirmation  confirmation
	signupError   error
	verifyError   error
	code          string
	deliveryError error
}

func (f *fakeSignupHandler) Signup(email, password string) (string, error) {
	if f.signupError != nil {
		return "", f.signupError
	}
	f.signup = signup{email, password}
	return f.code, nil
}

func (f *fakeSignupHandler) DeliverConfirmationCode(email, code string) error {
	if f.deliveryError != nil {
		return f.deliveryError
	}
	f.confirmation = confirmation{email, code}
	return nil
}

func (f *fakeSignupHandler) Confirm(email, code string) error {
	if f.verifyError != nil {
		return f.verifyError
	}
	f.confirmation = confirmation{email, code}
	return nil
}

type fakeAccountCloseHandler struct {
	email      string
	closeError error
}

func (f *fakeAccountCloseHandler) Close(email string) error {
	if f.closeError != nil {
		return f.closeError
	}
	f.email = email
	return nil
}

type fakeHistoryHandler struct {
	histories   []History
	listError   error
	getError    error
	deleteError error
}

func (f *fakeHistoryHandler) List(email string) ([]string, error) {
	if f.listError != nil {
		return nil, f.listError
	}
	ids := make([]string, len(f.histories))
	for i, h := range f.histories {
		ids[i] = h.ID
	}
	return ids, nil
}

func (f *fakeHistoryHandler) Get(id string) (History, error) {
	if f.getError != nil {
		return History{}, f.getError
	}
	for _, h := range f.histories {
		if h.ID == id {
			return h, nil
		}
	}
	return History{}, fmt.Errorf("Not found")
}

func (f *fakeHistoryHandler) Create(email string) (History, error) {
	f.histories = append(f.histories, History{
		ID: fmt.Sprintf("%d", len(f.histories)+1),
	})
	return f.histories[len(f.histories)-1], nil
}

func (f *fakeHistoryHandler) Delete(id string) error {
	if f.deleteError != nil {
		return f.deleteError
	}
	for i, h := range f.histories {
		if h.ID == id {
			f.histories, f.histories[len(f.histories)-1] = append(f.histories[:i], f.histories[i+1:]...), History{}
			return nil
		}
	}
	return fmt.Errorf("Not found")
}

type fakeItemsHandler struct {
	items      []map[string]Item
	history    History
	listError  error
	writeError error
}

func (f *fakeItemsHandler) List(historyID string, startIndex int) (History, []map[string]Item, error) {
	if f.listError != nil {
		return History{}, nil, f.listError
	}
	return f.history, f.items[startIndex:], nil
}

func (f *fakeItemsHandler) Write(historyID string, items []map[string]Item) (History, error) {
	if f.writeError != nil {
		return History{}, f.writeError
	}
	f.items = append(f.items, items...)
	return f.history, nil
}
