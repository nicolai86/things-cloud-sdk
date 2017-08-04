package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

// APIError represents any API error. Things cloud only returns "IsSyncronyErrorResponse" -
// I've extended this to include an error for easier debugging
type APIError struct {
	Error   string `json:"error"`
	IsError bool   `json:"IsSyncronyErrorResponse"`
}

// httpAccountHandler maps HTTP urls to AccountHandler methods
type httpAccountHandler struct {
	AccountHandler
}

func deleteAccount(h AccountHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		email := chi.URLParam(req, "email")
		err := h.Close(email)
		if err == nil {
			w.WriteHeader(http.StatusAccepted)
			return
		}

		handleError(w, err)
	}
}

func (h *httpAccountHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var err error

	email := chi.URLParam(req, "email")
	bs, readErr := ioutil.ReadAll(req.Body)
	if readErr != nil {
		// TODO error handling
		return
	}

	if isSignup(bs) {
		err = h.signup(bs, email)
		if err == nil {
			w.WriteHeader(http.StatusCreated)
			return
		}
	}

	if isConfirmation(bs) {
		err = h.confirm(bs, email)
		if err == nil {
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	if isPasswordChange(bs) {
		err = h.passwordChange(bs, email)
		if err == nil {
			w.WriteHeader(http.StatusOK)
			return
		}
	}
	handleError(w, err)
}

func handleError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	bs, _ := json.Marshal(&APIError{err.Error(), true})
	w.Write(bs)
}

// Handler routes requests in a thingscloud compatible manner
type Handler struct {
	m *chi.Mux
}

// HistoryResponse contains all information to interprete a sequence of items
type HistoryResponse struct {
	LatestSchemaVersion    int  `json:"latest-schema-version"`
	LatestTotalContentSize int  `json:"latest-total-content-size"`
	IsEmpty                bool `json:"is-empty"`
	LatestServerIndex      int  `json:"latest-server-index"`
}

type createHistoryResponse struct {
	Key string `json:"new-history-key"`
}

func deleteHistory(h HistoryHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		id := chi.URLParam(req, "id")
		err := h.Delete(id)
		if err == nil {
			w.WriteHeader(http.StatusAccepted)
			return
		}
		handleError(w, err)
	}
}

func createHistory(h HistoryHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		email := chi.URLParam(req, "email")
		history, err := h.Create(email)
		if err == nil {
			bs, _ := json.Marshal(&createHistoryResponse{Key: history.ID})
			w.Write(bs)
			return
		}
		handleError(w, err)
	}
}

func readHistory(h HistoryHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		id := chi.URLParam(req, "id")

		history, err := h.Get(id)
		if err == nil {
			bs, _ := json.Marshal(HistoryResponse{
				LatestSchemaVersion: history.LatestSchemaVersion,
				IsEmpty:             history.IsEmpty,
				LatestServerIndex:   history.LatestServerIndex,
			})
			w.Write(bs)
			return
		}
		handleError(w, err)
	}
}

func listHistories(h HistoryHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		email := chi.URLParam(req, "email")

		histories, err := h.List(email)
		if err == nil {
			bs, _ := json.Marshal(histories)
			w.Write(bs)
			return
		}
		handleError(w, err)
	}
}

type itemKind string
type itemAction int

// Item can be everything
type Item struct {
	P      json.RawMessage `json:"p"`
	Kind   itemKind        `json:"e"`
	Action itemAction      `json:"t"`
}

type itemsResponse struct {
	Items                  []map[string]Item `json:"items"`
	LatestTotalContentSize int               `json:"latest-total-content-size"`
	StartTotalContentSize  int               `json:"start-total-content-size"`
	EndTotalContentSize    int               `json:"end-total-content-size"`
	SchemaVersion          int               `json:"schema"`
	CurrentItemIndex       int               `json:"current-item-index"`
}

type writeRequest struct {
	AppID            string            `json:"app-id"`
	AppInstanceID    string            `json:"app-instance-id"`
	CurrentItemIndex int               `json:"current-item-index"`
	Items            []map[string]Item `json:"items"`
	PushPriority     int               `json:"push-priority"`
	Schema           int               `json:"schema"`
}

func listItems(h ItemsHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		id := chi.URLParam(req, "id")

		qs := req.URL.Query()
		startIndex := 0
		if vs := qs.Get("start-index"); vs != "" {
			v, err := strconv.Atoi(vs)
			if err == nil {
				startIndex = v
			}
		}

		history, items, err := h.List(id, startIndex)
		if err == nil {
			resp := itemsResponse{
				Items: items,
				StartTotalContentSize:  0, // TODO
				EndTotalContentSize:    0, // TODO
				LatestTotalContentSize: history.LatestTotalContentSize,
				SchemaVersion:          history.LatestSchemaVersion,
				CurrentItemIndex:       history.LatestServerIndex,
			}
			bs, _ := json.Marshal(resp)
			w.Write(bs)
			return
		}
		handleError(w, err)
	}
}

func writeItems(h ItemsHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		id := chi.URLParam(req, "id")

		bs, err := ioutil.ReadAll(req.Body)
		if err != nil {
			handleError(w, err)
			return
		}
		body := writeRequest{}
		err = json.Unmarshal(bs, &body)
		if err != nil {
			handleError(w, err)
			return
		}

		_, err = h.Write(id, body.Items)
		if err != nil {
			handleError(w, err)
			return
		}
	}
}

func fallback(w http.ResponseWriter, req *http.Request) {
	log.Printf("fallback: %q", req.URL)
}

// New creates a new thingscloud api server
func New(accountHandler AccountHandler, historyHandler HistoryHandler, itemsHandler ItemsHandler) *Handler {
	mux := chi.NewRouter()
	h := &Handler{mux}

	// TODO authentication middleware
	mux.Handle("/account/{email}", &httpAccountHandler{accountHandler})
	mux.Delete("/account/{email}", deleteAccount(accountHandler))

	mux.Get("/account/{email}/own-history-keys", listHistories(historyHandler))
	mux.Post("/account/{email}/own-history-keys", createHistory(historyHandler))
	mux.Delete("/account/{email}/own-history-keys/{id}", deleteHistory(historyHandler))
	mux.Get("/history/{id}", readHistory(historyHandler))

	mux.Get("/history/{id}/items", listItems(itemsHandler))
	mux.Post("/history/{id}/items", writeItems(itemsHandler))
	mux.Handle("/", http.HandlerFunc(fallback))
	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.m.ServeHTTP(w, req)
}
