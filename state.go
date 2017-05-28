package thingscloud

import (
	"encoding/json"
	"fmt"
	"time"
)

// State is created by applying all history items in order.
// Note that the hierarchy within the state (e.g. area > tasks > tasks > check list items)
// is modelled with pointers between the different maps, so concurrent modification
// is not safe.
type State struct {
	Areas          map[string]Area
	Tasks          map[string]Task
	Tags           map[string]Tag
	CheckListItems map[string]CheckListItem
}

// NewState creates a new, empty state
func NewState() *State {
	return &State{
		Areas:          map[string]Area{},
		Tags:           map[string]Tag{},
		CheckListItems: map[string]CheckListItem{},
		Tasks:          map[string]Task{},
	}
}

// Task describes a Task inside things.
type Task struct {
	ID               string
	CreationDate     time.Time
	ModificationDate *time.Time
	Status           TaskStatus
	Title            string
	Note             string
	DueDate          *time.Time
	CompletionDate   *time.Time

	SubTasks []*Task
}

// taskItem describes an event on a task
type taskItem struct {
	Item
	P struct {
		IX               *int        `json:"ix,omitempty"`
		CreationDate     *Timestamp  `json:"cd,omitempty"`
		ModificationDate *Timestamp  `json:"md,omitempty"`
		DueDate          *Timestamp  `json:"sr,omitempty"`
		CompletionDate   *Timestamp  `json:"sp,omitempty"`
		Status           *TaskStatus `json:"ss,omitempty"`
		TaskParent       *Boolean    `json:"tp,omitempty"`
		Title            *string     `json:"tt,omitempty"`
		Note             *string     `json:"nt,omitempty"`
		AreaIDs          []string    `json:"ar,omitempty"`
		ProjectIDs       []string    `json:"pr,omitempty"`
		TagIDs           []string    `json:"tg,omitempty"`
	} `json:"p"`
}

func (s *State) updateTask(item taskItem) Task {
	t, ok := s.Tasks[item.ID]
	if !ok {
		t = Task{ID: item.ID}
	}

	if item.P.Title != nil {
		t.Title = *item.P.Title
	}
	if item.P.Status != nil {
		t.Status = *item.P.Status
	}
	if item.P.DueDate != nil {
		t.DueDate = item.P.DueDate.Time()
	}
	if item.P.CompletionDate != nil {
		t.CompletionDate = item.P.CompletionDate.Time()
	}
	if item.P.CreationDate != nil {
		cd := item.P.CreationDate.Time()
		t.CreationDate = *cd
	}
	if item.P.ModificationDate != nil {
		t.ModificationDate = item.P.ModificationDate.Time()
	}
	if item.P.Note != nil {
		t.Note = *item.P.Note
	}
	if item.P.Title != nil {
		t.Title = *item.P.Title
	}

	return t
}

// CheckListItem describes a check list item
type CheckListItem struct {
	ID               string
	CreationDate     time.Time
	ModificationDate *time.Time
	Status           TaskStatus
	Title            string
	CompletionDate   *time.Time
}

// checkListItem describes an event on a check list item
type checkListItem struct {
	Item
	P struct {
		CreationDate     *Timestamp  `json:"cd,omitempty"`
		ModificationDate *Timestamp  `json:"md,omitempty"`
		IX               *int        `json:"ix"`
		Status           *TaskStatus `json:"ss,omitempty"`
		Title            *string     `json:"tt,omitempty"`
		CompletionDate   *Timestamp  `json:"sp,omitempty"`
		TaskIDs          []string    `json:"ts"`
	} `json:"p"`
}

func (s *State) updateCheckListItem(item checkListItem) CheckListItem {
	c, ok := s.CheckListItems[item.ID]
	if !ok {
		c = CheckListItem{ID: item.ID}
	}

	if item.P.CreationDate != nil {
		t := item.P.CreationDate.Time()
		c.CreationDate = *t
	}
	if item.P.ModificationDate != nil {
		c.ModificationDate = item.P.ModificationDate.Time()
	}
	if item.P.Title != nil {
		c.Title = *item.P.Title
	}
	if item.P.Status != nil {
		c.Status = *item.P.Status
	}

	return c
}

// Area describes an Area inside things. An Area is a container for tasks
type Area struct {
	ID    string
	Title string
	Tags  []*Tag
	Tasks []*Task
}

// areaItem describes an event on an area
type areaItem struct {
	Item
	P struct {
		IX     *int     `json:"ix"`
		Title  *string  `json:"tt"`
		TagIDs []string `json:"tg"`
	} `json:"p"`
}

func (s *State) updateArea(item areaItem) Area {
	a, ok := s.Areas[item.ID]
	if !ok {
		a = Area{ID: item.ID}
	}

	if item.P.Title != nil {
		a.Title = *item.P.Title
	}

	return a
}

// Tag describes the aggregated state of an Tag
type Tag struct {
	ID        string
	Title     string // tt
	ParentTag *Tag   // from `pm`
}

// tagItem describes an event on a tag
type tagItem struct {
	Item
	P struct {
		IX    *int     `json:"ix"`
		Title *string  `json:"tt"`
		SH    *string  `json:"sh"`
		PN    []string `json:"pn"`
	} `json:"p"`
}

func (s *State) updateTag(item tagItem) Tag {
	t, ok := s.Tags[item.ID]
	if !ok {
		t = Tag{ID: item.ID}
	}

	if item.P.Title != nil {
		t.Title = *item.P.Title
	}

	return t
}

// Update applies all items to update the aggregated state
func (s *State) Update(items ...Item) error {
	for _, rawItem := range items {
		switch rawItem.Kind {
		case ItemKindTask:
			item := taskItem{Item: rawItem}
			if err := json.Unmarshal(rawItem.P, &item.P); err != nil {
				return err
			}

			switch item.Action {
			case ItemActionCreated:
				fallthrough
			case ItemActionModified:
				s.Tasks[item.ID] = s.updateTask(item)
			case ItemActionDeleted:
				delete(s.Tasks, item.ID)
			default:
				fmt.Printf("Action %q on %q is not implemented yet", item.Action, rawItem.Kind)
			}

		case ItemKindChecklistItem:
			item := checkListItem{Item: rawItem}
			if err := json.Unmarshal(rawItem.P, &item.P); err != nil {
				return err
			}

			switch item.Action {
			case ItemActionCreated:
				fallthrough
			case ItemActionModified:
				s.CheckListItems[item.ID] = s.updateCheckListItem(item)
			case ItemActionDeleted:
				delete(s.CheckListItems, item.ID)
			default:
				fmt.Printf("Action %q on %q is not implemented yet", item.Action, rawItem.Kind)
			}

		case ItemKindArea:
			item := areaItem{Item: rawItem}
			if err := json.Unmarshal(rawItem.P, &item.P); err != nil {
				return err
			}

			switch item.Action {
			case ItemActionCreated:
				fallthrough
			case ItemActionModified:
				s.Areas[item.ID] = s.updateArea(item)

			case ItemActionDeleted:
				delete(s.Areas, item.ID)
			default:
				fmt.Printf("Action %q on %q is not implemented yet", item.Action, rawItem.Kind)
			}

		case ItemKindTag:
			item := tagItem{Item: rawItem}
			if err := json.Unmarshal(rawItem.P, &item.P); err != nil {
				return err
			}

			switch item.Action {
			case ItemActionCreated:
				fallthrough
			case ItemActionModified:
				s.Tags[item.ID] = s.updateTag(item)
			case ItemActionDeleted:
				delete(s.Tags, item.ID)
			default:
				fmt.Printf("Action %q on %q is not implemented yet", item.Action, rawItem.Kind)
			}

		default:
			fmt.Printf("%q is not implemented yet\n", rawItem.Kind)
		}
	}
	return nil
}
