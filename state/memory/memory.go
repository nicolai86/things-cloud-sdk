package memory

import (
	"encoding/json"
	"fmt"
	"sort"

	things "github.com/nicolai86/things-cloud-sdk"
)

// State is created by applying all history items in order.
// Note that the hierarchy within the state (e.g. area > tasks > tasks > check list items)
// is modelled with pointers between the different maps, so concurrent modification
// is not safe.
type State struct {
	Areas          map[string]*things.Area
	Tasks          map[string]*things.Task
	Tags           map[string]*things.Tag
	CheckListItems map[string]*things.CheckList
}

// NewState creates a new, empty state
func NewState() *State {
	return &State{
		Areas:          map[string]*things.Area{},
		Tags:           map[string]*things.Tag{},
		CheckListItems: map[string]*things.CheckList{},
		Tasks:          map[string]*things.Task{},
	}
}

func (s *State) updateTask(item things.TaskActionItem) *things.Task {
	t, ok := s.Tasks[item.UUID()]
	if !ok {
		t = &things.Task{
			Schedule: things.TaskScheduleAnytime,
		}
	}
	t.UUID = item.UUID()

	if item.P.Title != nil {
		t.Title = *item.P.Title
	}
	if item.P.IsProject != nil {
		t.IsProject = bool(*item.P.IsProject)
	}
	if item.P.Status != nil {
		t.Status = *item.P.Status
	}
	if item.P.Index != nil {
		t.Index = *item.P.Index
	}
	if item.P.InTrash != nil {
		t.InTrash = *item.P.InTrash
	}
	if item.P.Schedule != nil {
		t.Schedule = *item.P.Schedule
	}
	if item.P.ScheduledDate != nil {
		t.ScheduledDate = item.P.ScheduledDate.Time()
	}
	if item.P.CompletionDate != nil {
		t.CompletionDate = item.P.CompletionDate.Time()
	}
	if item.P.DeadlineDate != nil {
		t.DeadlineDate = item.P.DeadlineDate.Time()
	}
	if item.P.CreationDate != nil {
		cd := item.P.CreationDate.Time()
		t.CreationDate = *cd
	}
	if item.P.ModificationDate != nil {
		t.ModificationDate = item.P.ModificationDate.Time()
	}
	if item.P.AreaIDs != nil {
		ids := *item.P.AreaIDs
		t.AreaIDs = ids
	}
	if item.P.ParentTaskIDs != nil {
		ids := *item.P.ParentTaskIDs
		t.ParentTaskIDs = ids
	}
	if item.P.Note != nil {
		t.Note = *item.P.Note
	}
	if item.P.Title != nil {
		t.Title = *item.P.Title
	}

	return t
}

func (s *State) updateCheckListItem(item things.CheckListActionItem) *things.CheckList {
	c, ok := s.CheckListItems[item.UUID()]
	if !ok {
		c = &things.CheckList{}
	}
	c.UUID = item.UUID()

	if item.P.CreationDate != nil {
		t := item.P.CreationDate.Time()
		c.CreationDate = *t
	}
	if item.P.ModificationDate != nil {
		c.ModificationDate = item.P.ModificationDate.Time()
	}
	if item.P.Index != nil {
		c.Index = *item.P.Index
	}
	if item.P.Title != nil {
		c.Title = *item.P.Title
	}
	if item.P.Status != nil {
		c.Status = *item.P.Status
	}
	if item.P.TaskIDs != nil {
		ids := *item.P.TaskIDs
		c.TaskIDs = ids
	}

	return c
}

func (s *State) updateArea(item things.AreaActionItem) *things.Area {
	a, ok := s.Areas[item.UUID()]
	if !ok {
		a = &things.Area{}
	}
	a.UUID = item.UUID()

	if item.P.Title != nil {
		a.Title = *item.P.Title
	}

	return a
}

func (s *State) updateTag(item things.TagActionItem) *things.Tag {
	t, ok := s.Tags[item.UUID()]
	if !ok {
		t = &things.Tag{}
	}
	t.UUID = item.UUID()

	if item.P.Title != nil {
		t.Title = *item.P.Title
	}
	if item.P.ShortHand != nil {
		t.ShortHand = *item.P.ShortHand
	}
	if item.P.ParentTagIDs != nil {
		var ids = *item.P.ParentTagIDs
		t.ParentTagIDs = ids
	}

	return t
}

// Update applies all items to update the aggregated state
func (s *State) Update(items ...things.Item) error {
	for _, rawItem := range items {
		switch rawItem.Kind {
		case things.ItemKindTask:
			item := things.TaskActionItem{Item: rawItem}
			if err := json.Unmarshal(rawItem.P, &item.P); err != nil {
				return err
			}

			switch item.Action {
			case things.ItemActionCreated:
				fallthrough
			case things.ItemActionModified:
				s.Tasks[item.UUID()] = s.updateTask(item)
			case things.ItemActionDeleted:
				delete(s.Tasks, item.UUID())
			default:
				fmt.Printf("Action %q on %q is not implemented yet", item.Action, rawItem.Kind)
			}

		case things.ItemKindChecklistItem:
			item := things.CheckListActionItem{Item: rawItem}
			if err := json.Unmarshal(rawItem.P, &item.P); err != nil {
				return err
			}

			switch item.Action {
			case things.ItemActionCreated:
				fallthrough
			case things.ItemActionModified:
				s.CheckListItems[item.UUID()] = s.updateCheckListItem(item)
			case things.ItemActionDeleted:
				delete(s.CheckListItems, item.UUID())
			default:
				fmt.Printf("Action %q on %q is not implemented yet", item.Action, rawItem.Kind)
			}

		case things.ItemKindArea:
			item := things.AreaActionItem{Item: rawItem}
			if err := json.Unmarshal(rawItem.P, &item.P); err != nil {
				return err
			}

			switch item.Action {
			case things.ItemActionCreated:
				fallthrough
			case things.ItemActionModified:
				s.Areas[item.UUID()] = s.updateArea(item)

			case things.ItemActionDeleted:
				delete(s.Areas, item.UUID())
			default:
				fmt.Printf("Action %q on %q is not implemented yet", item.Action, rawItem.Kind)
			}

		case things.ItemKindTag:
			item := things.TagActionItem{Item: rawItem}
			if err := json.Unmarshal(rawItem.P, &item.P); err != nil {
				return err
			}

			switch item.Action {
			case things.ItemActionCreated:
				fallthrough
			case things.ItemActionModified:
				s.Tags[item.UUID()] = s.updateTag(item)
			case things.ItemActionDeleted:
				delete(s.Tags, item.UUID())
			default:
				fmt.Printf("Action %q on %q is not implemented yet", item.Action, rawItem.Kind)
			}

		default:
			fmt.Printf("%q is not implemented yet\n", rawItem.Kind)
		}
	}
	return nil
}

// Subtasks returns tasks grouped together with under a root task
func (s *State) Subtasks(root *things.Task) []*things.Task {
	tasks := []*things.Task{}
	for _, task := range s.Tasks {
		if task.Status == things.TaskStatusCompleted {
			continue
		}
		if task == root {
			continue
		}
		if task.InTrash {
			continue
		}
		isChild := false
		for _, taskID := range task.ParentTaskIDs {
			isChild = isChild || taskID == root.UUID
		}
		if isChild {
			tasks = append(tasks, task)
		}
	}
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Index < tasks[j].Index
	})
	return tasks
}

func hasArea(task *things.Task, state *State) bool {
	if len(task.AreaIDs) != 0 {
		return true
	}
	if len(task.ParentTaskIDs) == 0 {
		return false
	}
	for _, taskID := range task.ParentTaskIDs {
		if hasArea(state.Tasks[taskID], state) {
			return true
		}
	}
	return false
}

// TasksWithoutArea looks up top level tasks not assigned to any area, e.g. just created and placed in today
func (s *State) TasksWithoutArea() []*things.Task {
	tasks := []*things.Task{}
	for _, task := range s.Tasks {
		if task.Status == things.TaskStatusCompleted {
			continue
		}
		if len(task.ParentTaskIDs) != 0 {
			continue
		}
		if task.InTrash {
			continue
		}
		if !hasArea(task, s) {
			tasks = append(tasks, task)
		}
	}
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Index < tasks[j].Index
	})
	return tasks
}

// TasksByArea returns tasks associated with a given area
func (s *State) TasksByArea(area *things.Area) []*things.Task {
	tasks := []*things.Task{}
	for _, task := range s.Tasks {
		if task.Status == things.TaskStatusCompleted {
			continue
		}
		if task.InTrash {
			continue
		}
		isChild := false
		for _, areaID := range task.AreaIDs {
			isChild = isChild || areaID == area.UUID
		}
		if isChild {
			tasks = append(tasks, task)
		}
	}
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Index < tasks[j].Index
	})
	return tasks
}

// CheckListItemsByTask returns check lists associated with a particular item
func (s *State) CheckListItemsByTask(task *things.Task) []*things.CheckList {
	items := []*things.CheckList{}
	for _, item := range s.CheckListItems {
		if item.Status == things.TaskStatusCompleted {
			continue
		}
		isChild := false
		for _, taskID := range item.TaskIDs {
			isChild = isChild || task.UUID == taskID
		}
		if isChild {
			items = append(items, item)
		}
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].Index < items[j].Index
	})
	return items
}

// SubTags returns all child tags for a given root, ensuring sort order is kept intact
func (s *State) SubTags(root *things.Tag) []*things.Tag {
	children := []*things.Tag{}
	for _, tag := range s.Tags {
		if tag == root {
			continue
		}

		isChild := false
		for _, parentID := range tag.ParentTagIDs {
			isChild = isChild || parentID == root.UUID
		}
		if isChild {
			children = append(children, tag)
		}
	}
	sort.Slice(children, func(i, j int) bool {
		return children[i].ShortHand < children[j].ShortHand
	})
	return children
}
