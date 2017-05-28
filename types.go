package thingscloud

import (
	"encoding/json"
	"time"
)

//go:generate stringer -type ItemAction,TaskStatus,TaskSchedule

// ItemAction describes possible actions on Items
type ItemAction int

const (
	// ItemActionCreated is used to indicate a new Item was created
	ItemActionCreated ItemAction = iota
	// ItemActionModified is used to indicate an existing Item was modified
	ItemActionModified ItemAction = 1
	// ItemActionDeleted is used as a tombstone for an Item
	ItemActionDeleted ItemAction = 2
)

// TaskSchedule describes when a task is scheduled
type TaskSchedule int

const (
	// TaskScheduleToday indicates tasks which should be completed today
	TaskScheduleToday TaskSchedule = 0
	// TaskScheduleAnytime indicates tasks which can be completed anyday
	TaskScheduleAnytime TaskSchedule = 1
	// TaskScheduleSomeday indicates tasks which might never be completed
	TaskScheduleSomeday TaskSchedule = 2
)

// TaskStatus describes if a thing is completed or not
type TaskStatus int

const (
	// TaskStatusPending indicates a new task
	TaskStatusPending TaskStatus = iota
	// TaskStatusCompleted indicates a completed task
	TaskStatusCompleted TaskStatus = 3
	// TaskStatusCanceled indicates a canceled task
	TaskStatusCanceled TaskStatus = 2
)

// ItemKind describes the different types things cloud supports
type ItemKind string

var (
	// ItemKindChecklistItem identifies a CheckList
	ItemKindChecklistItem ItemKind = "ChecklistItem"
	// ItemKindTask identifies a Task or Subtask
	ItemKindTask ItemKind = "Task3"
	// ItemKindArea identifies an Area
	ItemKindArea ItemKind = "Area2"
	// ItemKindSettings  identifies a setting
	ItemKindSettings ItemKind = "Settings3"
	// ItemKindTag identifies a Tag
	ItemKindTag ItemKind = "Tag3"
)

// Timestamp allows unix epochs represented as float or ints to be unmarshalled
// into time.Time objects
type Timestamp time.Time

// UnmarshalJSON takes a unix epoch from float/ int and creates a time.Time instance
func (t *Timestamp) UnmarshalJSON(bs []byte) error {
	var d float64
	if err := json.Unmarshal(bs, &d); err != nil {
		return err
	}
	*t = Timestamp(time.Unix(int64(d), 0))
	return nil
}

func (t *Timestamp) MarshalJSON() ([]byte, error) {
	var tt = time.Time(*t).Unix()
	return json.Marshal(tt)
}

// Format returns a textual representation of the time value formatted according to layout
func (t *Timestamp) Format(layout string) string {
	return time.Time(*t).Format(layout)
}

// Time returns the underlying time.Time instance
func (t *Timestamp) Time() *time.Time {
	tt := time.Time(*t)
	return &tt
}

// Boolean allows integers to be parsed into booleans, where 1 means true and 0 means false
type Boolean bool

// UnmarshalJSON takes an int and creates a boolean instance
func (b *Boolean) UnmarshalJSON(bs []byte) error {
	var d int
	if err := json.Unmarshal(bs, &d); err != nil {
		return err
	}
	*b = Boolean(d == 1)
	return nil
}

func (b *Boolean) MarshalJSON() ([]byte, error) {
	var d = 0
	if *b {
		d = 1
	}
	return json.Marshal(d)
}
