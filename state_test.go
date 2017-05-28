package thingscloud

import (
	"encoding/json"
	"testing"
)

var newAreaPayload = `{
  "tt": "New Area"
}`

var newTagPayload = `{
  "ix": 4,
  "tt": "‼",
  "sh": "",
  "pn": []
}`

var newCheckListItemPayload = `{
  "md": 1495650756.177454,
  "ix": -543,
  "ss": 0,
  "tt": "A CheckListItem",
  "sp": null,
  "ts": [
    "1452D14B-3099-4251-B7C8-8B4C2B9BA334"
  ],
  "cd": 1495650723.3711741
}`

var newTaskPayload = `{
  "ix": -485,
  "cd": 1495650251.230479,
  "icsd": null,
  "ar": [],
  "tir": null,
  "rt": [],
  "rr": null,
  "icc": 0,
  "tt": "A \"Task\"",
  "tr": false,
  "tp": 0,
  "acrd": null,
  "ti": 0,
  "tg": [],
  "icp": false,
  "nt": null,
  "do": 0,
  "dl": [],
  "lai": null,
  "dd": null,
  "pr": [
    "BB7323E0-36E5-4DE7-8203-6A27B67C1CD4"
  ],
  "md": 1495650259.683249,
  "ss": 0,
  "sr": null,
  "sp": null,
  "st": 1,
  "dds": null,
  "ato": null,
  "sb": 0,
  "agr": []
}`

func TestState_Update(t *testing.T) {
	t.Run("Empty State", func(t *testing.T) {
		t.Parallel()
		t.Run("Create Area", func(t *testing.T) {
			t.Parallel()
			s := NewState()
			if err := s.Update(Item{
				Action: ItemActionCreated,
				Kind:   ItemKindArea,
				P:      json.RawMessage(newAreaPayload),
			}); err != nil {
				t.Fatal(err.Error())
			}

			if len(s.Areas) != 1 {
				t.Fatal("Expected to have a new area")
			}
		})

		t.Run("Create Tag", func(t *testing.T) {
			t.Parallel()
			s := NewState()
			if err := s.Update(Item{
				Action: ItemActionCreated,
				Kind:   ItemKindTag,
				P:      json.RawMessage(newTagPayload),
			}); err != nil {
				t.Fatal(err.Error())
			}

			if len(s.Tags) != 1 {
				t.Fatal("Expected to have a new tag")
			}
		})

		t.Run("Create CheckListItem", func(t *testing.T) {
			t.Parallel()
			s := NewState()
			if err := s.Update(Item{
				Action: ItemActionCreated,
				Kind:   ItemKindChecklistItem,
				P:      json.RawMessage(newCheckListItemPayload),
			}); err != nil {
				t.Fatal(err.Error())
			}

			if len(s.CheckListItems) != 1 {
				t.Fatal("Expected to have a new check list item")
			}
		})

		t.Run("Create Task", func(t *testing.T) {
			t.Parallel()
			s := NewState()
			if err := s.Update(Item{
				Action: ItemActionCreated,
				Kind:   ItemKindTask,
				P:      json.RawMessage(newTaskPayload),
			}); err != nil {
				t.Fatal(err.Error())
			}

			if len(s.Tasks) != 1 {
				t.Fatal("Expected to have a new task")
			}
		})
	})

	t.Run("Existing State", func(t *testing.T) {
		t.Parallel()
		newState := func() *State {
			s := NewState()
			s.Update(Item{
				Action: ItemActionCreated,
				Kind:   ItemKindArea,
				P:      json.RawMessage(newAreaPayload),
			}, Item{
				Action: ItemActionCreated,
				Kind:   ItemKindTag,
				P:      json.RawMessage(newTagPayload),
			}, Item{
				Action: ItemActionCreated,
				Kind:   ItemKindChecklistItem,
				P:      json.RawMessage(newCheckListItemPayload),
			}, Item{
				Action: ItemActionCreated,
				Kind:   ItemKindTask,
				P:      json.RawMessage(newTaskPayload),
			})
			return s
		}

		t.Run("Update Task", func(t *testing.T) {
			t.Parallel()
			s := newState()
			if err := s.Update(Item{
				Action: ItemActionModified,
				Kind:   ItemKindTask,
				P:      json.RawMessage(`{"tt": "Modified Title"}`),
			}); err != nil {
				t.Fatal(err.Error())
			}
			if s.Tasks[""].Title != "Modified Title" {
				t.Fatal("Expected title to be updated")
			}
		})

		t.Run("Delete Task", func(t *testing.T) {
			t.Parallel()
			s := newState()
			if err := s.Update(Item{
				Action: ItemActionDeleted,
				Kind:   ItemKindTask,
				P:      json.RawMessage(`{}`),
			}); err != nil {
				t.Fatal(err.Error())
			}
			if len(s.Tasks) != 0 {
				t.Fatal("Expected to have no more tasks")
			}
		})

		t.Run("Update CheckListItem", func(t *testing.T) {
			t.Parallel()
			s := newState()
			if err := s.Update(Item{
				Action: ItemActionModified,
				Kind:   ItemKindChecklistItem,
				P:      json.RawMessage(`{"tt": "Modified Title"}`),
			}); err != nil {
				t.Fatal(err.Error())
			}
			if s.CheckListItems[""].Title != "Modified Title" {
				t.Fatal("Expected title to be updated")
			}
		})

		t.Run("Delete ChecklistItem", func(t *testing.T) {
			t.Parallel()
			s := newState()
			if err := s.Update(Item{
				Action: ItemActionDeleted,
				Kind:   ItemKindChecklistItem,
				P:      json.RawMessage(`{}`),
			}); err != nil {
				t.Fatal(err.Error())
			}
			if len(s.CheckListItems) != 0 {
				t.Fatal("Expected to have no more check list items")
			}
		})

		t.Run("Update Tag", func(t *testing.T) {
			t.Parallel()
			s := newState()
			if err := s.Update(Item{
				Action: ItemActionModified,
				Kind:   ItemKindTag,
				P:      json.RawMessage(`{"tt": "Modified Tag"}`),
			}); err != nil {
				t.Fatal(err.Error())
			}
			if s.Tags[""].Title != "Modified Tag" {
				t.Fatal("Expected title to be updated")
			}
		})

		t.Run("Delete Tag", func(t *testing.T) {
			t.Parallel()
			s := newState()
			if err := s.Update(Item{
				Action: ItemActionDeleted,
				Kind:   ItemKindTag,
				P:      json.RawMessage(`{}`),
			}); err != nil {
				t.Fatal(err.Error())
			}
			if len(s.Tags) != 0 {
				t.Fatal("Expected to have no more tags")
			}
		})

		t.Run("Update Area", func(t *testing.T) {
			t.Parallel()
			s := newState()

			if err := s.Update(Item{
				Action: ItemActionModified,
				Kind:   ItemKindArea,
				P:      json.RawMessage(`{"tt": "Modified Area"}`),
			}); err != nil {
				t.Fatal(err.Error())
			}

			if s.Areas[""].Title != "Modified Area" {
				t.Fatal("Expected title to be updated")
			}
		})

		t.Run("Delete Area", func(t *testing.T) {
			t.Parallel()
			s := newState()

			if err := s.Update(Item{
				Action: ItemActionDeleted,
				Kind:   ItemKindArea,
				P:      json.RawMessage(`{}`),
			}); err != nil {
				t.Fatal(err.Error())
			}

			if len(s.Areas) != 0 {
				t.Fatal("Expected to have no more areas")
			}
		})
	})
}

func TestState_updateTag(t *testing.T) {
	s := NewState()
	a := &Tag{
		ID:    "CC-Things-Tag-High",
		Title: "High",
	}
	b := &Tag{
		ID:    "CC-Things-Tag-Priority",
		Title: "!!",
	}
	c := &Tag{
		ID:    "CC-Things-Tag-Errand",
		Title: "Errand",
	}
	s.Tags[a.ID] = a
	s.Tags[b.ID] = b
	s.Tags[c.ID] = c

	t.Run("sets Title", func(t *testing.T) {
		tag := s.updateTag(TagItem{
			Item: Item{ID: a.ID},
			P: TagItemPayload{
				Title: stringVal("a title"),
			},
		})
		if tag.Title != "a title" {
			t.Fatalf("Expected Title %q but got %q", "a title", tag.Title)
		}
	})

	t.Run("hierarchy", func(t *testing.T) {
		t.Run("create", func(t *testing.T) {
			a2 := s.updateTag(TagItem{
				Item: Item{ID: a.ID},
				P: TagItemPayload{
					ParentTagIDs: &[]string{b.ID},
				},
			})
			if a2.ParentTagIDs[0] != b.ID {
				t.Fatalf("Expected parent of %q, but got %q", b.ID, a2.ParentTagIDs)
			}
		})
		t.Run("change", func(t *testing.T) {
			a2 := s.updateTag(TagItem{
				Item: Item{ID: a.ID},
				P: TagItemPayload{
					ParentTagIDs: &[]string{c.ID},
				},
			})
			if a2.ParentTagIDs[0] != c.ID {
				t.Fatalf("Expected parent of %q, but got %q", c.ID, a2.ParentTagIDs)
			}
		})
		t.Run("no change", func(t *testing.T) {
			a2 := s.updateTag(TagItem{
				Item: Item{ID: a.ID},
				P:    TagItemPayload{},
			})
			if a2.ParentTagIDs[0] != c.ID {
				t.Fatalf("Expected parent of %q, but got %q", c.ID, a2.ParentTagIDs)
			}
		})
		t.Run("delete", func(t *testing.T) {
			a2 := s.updateTag(TagItem{
				Item: Item{ID: a.ID},
				P: TagItemPayload{
					ParentTagIDs: &[]string{},
				},
			})
			if len(a2.ParentTagIDs) != 0 {
				t.Fatalf("Expected parent of nil, but got %q", a2.ParentTagIDs)
			}
		})
	})
}
