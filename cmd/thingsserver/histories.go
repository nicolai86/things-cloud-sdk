package main

// HistoryHandler wraps various history handling tasks
type HistoryHandler interface {
	Create(email string) (History, error)
	List(email string) ([]string, error)
	Get(id string) (History, error)
	Delete(id string) error
}

// History contains meta data associated with every history
type History struct {
	ID                     string
	LatestSchemaVersion    int
	LatestTotalContentSize int
	IsEmpty                bool
	LatestServerIndex      int
}
