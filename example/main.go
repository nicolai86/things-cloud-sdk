package main

import (
	"log"
	"os"

	thingscloud "github.com/nicolai86/things-cloud-sdk"
)

func main() {
	c := thingscloud.New(thingscloud.APIEndpoint, os.Getenv("THINGS_USERNAME"), os.Getenv("THINGS_PASSWORD"))

	_, err := c.Verify()
	if err != nil {
		log.Fatalf("Login failed: %q\nPlease check your credentials.", err.Error())
	}
	log.Printf("User: %s\n", c.EMail)

	if hs, err := c.Histories(); err != nil {
		log.Fatalf("Failed to lookup histories: %q\n", err.Error())
	} else {
		log.Printf("Histories: %d\n", len(hs))

		if len(hs) > 0 {
			history := hs[0]
			state := thingscloud.NewState()

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
			log.Printf(`Summary:
Areas:          %d
Tasks:          %d (%d)
CheckListItems: %d (%d)
Tags:           %d
`, len(state.Areas), len(state.Tasks), doneTasks, len(state.CheckListItems), doneChecklistItems, len(state.Tags))
		}
	}

	if h, err := c.CreateHistory(); err != nil {
		log.Fatalf("Failed to create history: %q\n", err.Error())
	} else {
		log.Printf("Created new history…\n")

		if err := h.Delete(); err != nil {
			log.Fatalf("Failed to delete history: %q\n", err.Error())
		} else {
			log.Printf("Deleted new history…\n")
		}
	}
}
