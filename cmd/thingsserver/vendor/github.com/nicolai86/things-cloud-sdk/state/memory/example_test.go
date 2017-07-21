package memory

import (
	"fmt"
	"log"
	"os"

	things "github.com/nicolai86/things-cloud-sdk"
)

func ExampleState_Update() {
	client := things.New(things.APIEndpoint, os.Getenv("THINGS_USERNAME"), os.Getenv("THINGS_PASSWORD"))

	histories, err := client.Histories()
	if err != nil {
		log.Printf("Failed loading histories: %q", err.Error())
	}
	history := histories[0]

	items, err := history.Items(things.ItemsOptions{})
	if err != nil {
		log.Printf("Failed loading items: %q", err.Error())
	}

	state := NewState()
	if err := state.Update(items...); err != nil {
		log.Printf("Failed aggregating state: %q", err.Error())
	}

	fmt.Printf(`Summary:
Areas:          %d
Tasks:          %d
CheckListItems: %d
Tags:           %d
`, len(state.Areas),
		len(state.Tasks),
		len(state.CheckListItems),
		len(state.Tags))
}
