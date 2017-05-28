package main

import (
	"fmt"
	"log"
	"os"

	thingscloud "github.com/nicolai86/things-cloud-sdk"
)

func printTag(tag *thingscloud.Tag, state *thingscloud.State, indent string) {
	fmt.Printf("%s-\t%s\n", indent, tag.Title)
	children := state.SubTags(tag)
	for _, child := range children {
		printTag(child, state, fmt.Sprintf("%s\t", indent))
	}
}

func printTask(task *thingscloud.Task, state *thingscloud.State, indent string) {
	fmt.Printf("%s-\t%s\n", indent, task.Title)
	checklist := state.CheckListItemsByTask(task)
	for _, item := range checklist {
		fmt.Printf("%s+%s\n", indent, item.Title)
	}
	children := state.Subtasks(task)
	for _, child := range children {
		printTask(child, state, fmt.Sprintf("%s\t", indent))
	}
}

func stringVal(s string) *string {
	return &s
}

func main() {
	c := thingscloud.New(thingscloud.APIEndpoint, os.Getenv("THINGS_USERNAME"), os.Getenv("THINGS_PASSWORD"))

	_, err := c.Verify()
	if err != nil {
		log.Fatalf("Login failed: %q\nPlease check your credentials.", err.Error())
	}
	fmt.Printf("User: %s\n", c.EMail)

	if hs, err := c.Histories(); err != nil {
		log.Fatalf("Failed to lookup histories: %q\n", err.Error())
	} else {
		fmt.Printf("Histories: %d\n", len(hs))

		if len(hs) > 0 {
			history := hs[0]
			history.Sync()

			state := thingscloud.NewState()

			pending := thingscloud.TaskStatusPending
			anytime := thingscloud.TaskScheduleAnytime
			yes := thingscloud.Boolean(true)
			if err := history.Write(thingscloud.TaskItem{
				Item: thingscloud.Item{
					Kind:   thingscloud.ItemKindTask,
					Action: thingscloud.ItemActionCreated,
					ID:     "54152210-ABFA-4F9F-81AC-7F50FBDEDC1G",
				},
				P: thingscloud.TaskItemPayload{
					Title:        stringVal("test 4"),
					Schedule:     &anytime,
					Status:       &pending,
					CreationDate: &thingscloud.Timestamp{},
					IsProject:    &yes,
				},
			}); err != nil {
				log.Fatalf("Write failed: %q\n", err.Error())
			}

			items, err := history.Items(thingscloud.ItemsOptions{StartIndex: 0})
			if err != nil {
				log.Fatalf("Failed to lookup items: %q\n", err.Error())
			}
			if err := state.Update(items...); err != nil {
				log.Fatalf("Failed to update state: %q\n", err.Error())
			}

			doneTasks := 0
			for _, task := range state.Tasks {
				if task.Status == thingscloud.TaskStatusCompleted {
					doneTasks = doneTasks + 1
				}
			}

			doneChecklistItems := 0
			for _, item := range state.CheckListItems {
				if item.Status == thingscloud.TaskStatusCompleted {
					doneChecklistItems = doneChecklistItems + 1
				}
			}
			fmt.Printf(`Summary:
Areas:          %d
Tasks:          %d (%d)
CheckListItems: %d (%d)
Tags:           %d
`, len(state.Areas),
				len(state.Tasks), doneTasks,
				len(state.CheckListItems), doneChecklistItems,
				len(state.Tags))

			fmt.Printf("Tags\n")
			for _, tag := range state.Tags {
				if len(tag.ParentTagIDs) != 0 {
					continue
				}
				printTag(tag, state, "")
			}
			fmt.Printf("\n\n")

			fmt.Printf("Areas\n")
			for _, area := range state.Areas {
				fmt.Printf("-\t%s\n", area.Title)

				for _, task := range state.TasksByArea(area) {
					printTask(task, state, "|")
				}
			}

			fmt.Printf("No Areas\n")
			for _, task := range state.TasksWithoutArea() {
				printTask(task, state, "|")
			}

			fmt.Printf("Today\n")
			for _, task := range state.Tasks {
				if task.Schedule != thingscloud.TaskScheduleToday {
					continue
				}
				if task.Status != thingscloud.TaskStatusPending {
					continue
				}
				printTask(task, state, "--")
			}
		}
	}

	// if h, err := c.CreateHistory(); err != nil {
	// 	log.Fatalf("Failed to create history: %q\n", err.Error())
	// } else {
	// 	fmt.Printf("Created new history…\n")

	// 	if err := h.Delete(); err != nil {
	// 		log.Fatalf("Failed to delete history: %q\n", err.Error())
	// 	} else {
	// 		fmt.Printf("Deleted new history…\n")
	// 	}
	// }
}
