package thingscloud

//go:generate stringer -type ItemAction,TaskStatus,TaskSchedule

import (
	"encoding/json"
	"time"
)

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
	*t = Timestamp(time.Unix(int64(d), 0).UTC())
	return nil
}

// MarshalJSON convers a timestamp into unix nano representation
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

// MarshalJSON takes a boolean and serializes it as an integer
func (b *Boolean) MarshalJSON() ([]byte, error) {
	var d = 0
	if *b {
		d = 1
	}
	return json.Marshal(d)
}

// Task describes a Task inside things.
// 0|uuid|TEXT|0||1
// 1|userModificationDate|REAL|0||0
// 2|creationDate|REAL|0||0
// 3|trashed|INTEGER|0||0
// 4|type|INTEGER|0||0
// 5|title|TEXT|0||0
// 6|notes|TEXT|0||0
// 7|dueDate|REAL|0||0
// 8|dueDateOffset|INTEGER|0||0
// 9|status|INTEGER|0||0
// 10|stopDate|REAL|0||0
// 11|start|INTEGER|0||0
// 12|startDate|REAL|0||0
// 13|index|INTEGER|0||0
// 14|todayIndex|INTEGER|0||0
// 15|area|TEXT|0||0
// 16|project|TEXT|0||0
// 17|repeatingTemplate|TEXT|0||0
// 18|delegate|TEXT|0||0
// 19|recurrenceRule|BLOB|0||0
// 20|instanceCreationStartDate|REAL|0||0
// 21|instanceCreationPaused|INTEGER|0||0
// 22|instanceCreationCount|INTEGER|0||0
// 23|afterCompletionReferenceDate|REAL|0||0
// 24|actionGroup|TEXT|0||0
// 25|untrashedLeafActionsCount|INTEGER|0||0
// 26|openUntrashedLeafActionsCount|INTEGER|0||0
// 27|checklistItemsCount|INTEGER|0||0
// 28|openChecklistItemsCount|INTEGER|0||0
// 29|startBucket|INTEGER|0||0
// 30|alarmTimeOffset|REAL|0||0
// 31|lastAlarmInteractionDate|REAL|0||0
// 32|todayIndexReferenceDate|REAL|0||0
// 33|nextInstanceStartDate|REAL|0||0
// 34|dueDateSuppressionDate|REAL|0||0
type Task struct {
	UUID             string
	CreationDate     time.Time
	ModificationDate *time.Time
	Status           TaskStatus
	Title            string
	Note             string
	ScheduledDate    *time.Time
	CompletionDate   *time.Time
	DeadlineDate     *time.Time
	Index            int
	AreaIDs          []string
	ParentTaskIDs    []string
	InTrash          bool
	Schedule         TaskSchedule
	IsProject        bool
}

// TaskActionItemPayload describes the payload for modifying Tasks, and also Projects,
// as projects are special kind of Tasks
type TaskActionItemPayload struct {
	Index             *int          `json:"ix,omitempty"`
	CreationDate      *Timestamp    `json:"cd,omitempty"`
	ModificationDate  *Timestamp    `json:"md,omitempty"` // ok
	ScheduledDate     *Timestamp    `json:"sr,omitempty"`
	CompletionDate    *Timestamp    `json:"sp,omitempty"`
	DeadlineDate      *Timestamp    `json:"dd,omitempty"` //
	Status            *TaskStatus   `json:"ss,omitempty"`
	IsProject         *Boolean      `json:"tp,omitempty"`
	Title             *string       `json:"tt,omitempty"`
	Note              *string       `json:"nt,omitempty"`
	AreaIDs           *[]string     `json:"ar,omitempty"`
	ParentTaskIDs     *[]string     `json:"pr,omitempty"`
	TagIDs            []string      `json:"tg,omitempty"`
	InTrash           *bool         `json:"tr,omitempty"`
	RecurrenceTaskIDs *[]string     `json:"rt,omitempty"`
	Schedule          *TaskSchedule `json:"st,omitempty"`
	//  {
	//      "acrd": null,
	//      "agr": [],
	//      "ar": [],
	//      "ato": null,
	//      "cd": 1495662927.014228,
	//      "dd": null,
	//      "dds": null,
	//      "dl": [],
	//      "do": 0,
	//      "icc": 0,
	//      "icp": false,
	//      "icsd": null, instance creation start date
	//      "ix": 0,
	//      "lai": null,
	//      "md": 1495662933.606909,
	//      "nt": "<note xml:space=\"preserve\">test body pm</note>",
	//      "pr": [],
	//      "rr": null,
	//      "rt": [],
	//      "sb": 0,
	//      "sp": null,
	//      "sr": 1495584000,
	//      "ss": 0,
	//      "st": 1,
	//      "tg": [],
	//      "ti": 0,
	//      "tir": 1495584000,
	//      "tp": 0,
	//      "tr": false,
	//      "tt": "test"
	//  },
	// "rr": {
	//   "fu": 16,
	//   "sr": 1495670400,
	//   "of": [
	//     {
	//       "dy": 0
	//     }
	//   ],
	//   "ts": -3,
	//   "tp": 0,
	//   "fa": 1,
	//   "rc": 2,
	//   "ia": 1495929600,
	//   "rrv": 4
	// },
}

// TaskActionItem describes an event on a Task
type TaskActionItem struct {
	Item
	P TaskActionItemPayload `json:"p"`
}

// UUID returns the UUID of the modified Task
func (t TaskActionItem) UUID() string {
	return t.Item.UUID
}

// Tag describes the aggregated state of an Tag
// 0|uuid|TEXT|0||1
// 1|title|TEXT|0||0
// 2|shortcut|TEXT|0||0
// 3|usedDate|REAL|0||0
// 4|parent|TEXT|0||0
// 5|index|INTEGER|0||0
type Tag struct {
	UUID         string
	Title        string
	ParentTagIDs []string
	ShortHand    string
}

// TagActionItemPayload describes the payload for modifying Areas
type TagActionItemPayload struct {
	IX           *int      `json:"ix"`
	Title        *string   `json:"tt"`
	ShortHand    *string   `json:"sh"`
	ParentTagIDs *[]string `json:"pn"`
}

// TagActionItem describes an event on a tag
type TagActionItem struct {
	Item
	P TagActionItemPayload `json:"p"`
}

// UUID returns the UUID of the modified Tag
func (t TagActionItem) UUID() string {
	return t.Item.UUID
}

// Setting describes things settings
// 0|uuid|TEXT|0||1
// 1|logInterval|INTEGER|0||0
// 2|manualLogDate|REAL|0||0
// 3|groupTodayByParent|INTEGER|0||0
type Setting struct{}

// Area describes an Area inside things. An Area is a container for tasks
// 0|uuid|TEXT|0||1
// 1|title|TEXT|0||0
// 2|visible|INTEGER|0||0
// 3|index|INTEGER|0||0
type Area struct {
	UUID  string
	Title string
	Tags  []*Tag
	Tasks []*Task
}

// AreaActionItemPayload describes the payload for modifying Areas
type AreaActionItemPayload struct {
	IX     *int     `json:"ix,omitempty"`
	Title  *string  `json:"tt,omitempty"`
	TagIDs []string `json:"tg,omitempty"`
}

// AreaActionItem describes an event on an Area
type AreaActionItem struct {
	Item
	P AreaActionItemPayload `json:"p"`
}

// UUID returns the UUID of the modified Area
func (item AreaActionItem) UUID() string {
	return item.Item.UUID
}

// CheckListItem describes a check list item
//0|uuid|TEXT|0||1
//1|userModificationDate|REAL|0||0
//2|creationDate|REAL|0||0
//3|title|TEXT|0||0
//4|status|INTEGER|0||0
//5|stopDate|REAL|0||0
//6|index|INTEGER|0||0
//7|task|TEXT|0||0
type CheckListItem struct {
	UUID             string
	CreationDate     time.Time
	ModificationDate *time.Time
	Status           TaskStatus
	Title            string
	Index            int
	CompletionDate   *time.Time
	TaskIDs          []string
}

// CheckListActionItemPayload describes the payload for modifying CheckListItems
type CheckListActionItemPayload struct {
	CreationDate     *Timestamp  `json:"cd,omitempty"`
	ModificationDate *Timestamp  `json:"md,omitempty"`
	Index            *int        `json:"ix"`
	Status           *TaskStatus `json:"ss,omitempty"`
	Title            *string     `json:"tt,omitempty"`
	CompletionDate   *Timestamp  `json:"sp,omitempty"`
	TaskIDs          *[]string   `json:"ts,omitempty"`
}

// CheckListActionItem describes an event on a check list item
type CheckListActionItem struct {
	Item
	P CheckListActionItemPayload `json:"p"`
}

// UUID returns the UUID of the modified CheckListItem
func (item CheckListActionItem) UUID() string {
	return item.Item.UUID
}
