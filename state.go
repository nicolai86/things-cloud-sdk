package thingscloud

import (
	"encoding/json"
	"fmt"
	"time"
)

type State struct {
	Areas          map[string]Area
	Tasks          map[string]Task
	Tags           map[string]Tag
	CheckListItems map[string]CheckListItem
}

func NewState() *State {
	return &State{
		Areas:          map[string]Area{},
		Tags:           map[string]Tag{},
		CheckListItems: map[string]CheckListItem{},
		Tasks:          map[string]Task{},
	}
}

// Area describes an Area inside things. An Area is a container for tasks
type Area struct {
	ID    string
	Title string
	Tags  []Tag
	Tasks []Task
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

type Task struct {
	ID               string
	CreationDate     time.Time
	ModificationDate *time.Time
	Status           TaskStatus
	Title            string
	Note             string
	DueDate          *time.Time
	CompletionDate   *time.Time

	SubTasks []Task
}

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

type CheckListItem struct {
	ID               string
	CreationDate     time.Time
	ModificationDate *time.Time
	Status           TaskStatus
	Title            string
	CompletionDate   *time.Time
}

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

// Update applies all items to update the aggregated state
func (s State) Update(items ...Item) error {
	for _, rawItem := range items {
		switch rawItem.Kind {
		case ItemKindTask:
			item := taskItem{Item: rawItem}
			if err := json.Unmarshal(rawItem.P, &item.P); err != nil {
				return err
			}

			switch item.Action {
			case ItemActionCreated:
				s.Tasks[item.ID] = Task{
					Title: *item.P.Title,
				}
			case ItemActionModified:
				if item.P.Title != nil {
					t := s.Tasks[item.ID]
					t.Title = *item.P.Title
					s.Tasks[item.ID] = t
				}
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
				s.CheckListItems[item.ID] = CheckListItem{
					Title: *item.P.Title,
				}
			case ItemActionModified:
				if item.P.Title != nil {
					cli := s.CheckListItems[item.ID]
					cli.Title = *item.P.Title
					s.CheckListItems[item.ID] = cli
				}
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
				s.Areas[item.ID] = Area{
					Title: *item.P.Title,
				}
			case ItemActionModified:
				if item.P.Title != nil {
					area := s.Areas[item.ID]
					area.Title = *item.P.Title
					s.Areas[item.ID] = area
				}
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
				s.Tags[item.ID] = Tag{
					Title: *item.P.Title,
				}
			case ItemActionModified:
				if item.P.Title != nil {
					tag := s.Tags[item.ID]
					tag.Title = *item.P.Title
					s.Tags[tag.ID] = tag
				}
			case ItemActionDeleted:
				delete(s.Tags, item.ID)
			default:
				fmt.Printf("Action %q on %q is not implemented yet", item.Action, rawItem.Kind)
			}

		default:
			fmt.Printf("%q is not implemented yet", rawItem.Kind)
		}
	}
	return nil
}

// merge what looks like an event sourcing approach
// func mergeItems(changes ...Item) (Item, bool) {
//   var base = changes[0]
//   var last = changes[len(changes)-1]
//   for _, change := range changes[1:] {
//     if change.P.Status != nil {
//       base.P.Status = change.P.Status
//     }
//     if change.P.DueDate != nil {
//       base.P.DueDate = change.P.DueDate
//     }
//     if change.P.CompletionDate != nil {
//       base.P.CompletionDate = change.P.CompletionDate
//     }
//     if change.P.ModificationDate != nil {
//       base.P.ModificationDate = change.P.ModificationDate
//     }
//     if change.P.Note != nil {
//       base.P.Note = change.P.Note
//     }
//     if change.P.Title != nil {
//       base.P.Title = change.P.Title
//     }
//     if change.P.TaskParent != nil {
//       base.P.TaskParent = change.P.TaskParent
//     }
//   }
//   return base, last.Action == ActionDeleted
// }
