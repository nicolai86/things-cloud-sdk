package thingscloud

import (
	"log"
	"os"
)

func ExampleHistory_Items() {
	client := New(APIEndpoint, os.Getenv("THINGS_USERNAME"), os.Getenv("THINGS_PASSWORD"))

	histories, err := client.Histories()
	if err != nil {
		log.Printf("Failed loading histories: %q", err.Error())
	}
	history := histories[0]

	items, _, err := history.Items(ItemsOptions{})
	if err != nil {
		log.Printf("Failed loading items: %q", err.Error())
	}

	log.Printf("got %d items.", len(items))
}
