package thingscloud

import (
	"encoding/json"
	"time"
)

//go:generate stringer -type ItemAction,TaskStatus

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
	// ItemKindChecklist identifies a CheckList
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

type Timestamp time.Time

func (t *Timestamp) UnmarshalJSON(bs []byte) error {
	var d float64
	if err := json.Unmarshal(bs, &d); err != nil {
		return err
	}
	*t = Timestamp(time.Unix(int64(d), 0))
	return nil
}

func (t *Timestamp) Format(layout string) string {
	return time.Time(*t).Format(layout)
}

func (t *Timestamp) Time() *time.Time {
	tt := time.Time(*t)
	return &tt
}

type Boolean bool

func (b *Boolean) UnmarshalJSON(bs []byte) error {
	var d int
	if err := json.Unmarshal(bs, &d); err != nil {
		return err
	}
	*b = Boolean(d == 1)
	return nil
}
