package main

// ItemsHandler is used to fetch history items as well as to write history items
type ItemsHandler interface {
	List(historyID string, startIndex int) (History, []map[string]Item, error)
	Write(historyID string, items []map[string]Item) (History, error)
}
