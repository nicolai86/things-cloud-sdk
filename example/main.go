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
		log.Printf("Failed to lookup histories: %q\n", err.Error())
	} else {
		log.Printf("Histories: %d\n", len(hs))

		if len(hs) > 0 {
			log.Printf("Syncing History 0… %s\n", hs[0].Sync())
		}
	}

	if h, err := c.CreateHistory(); err != nil {
		log.Printf("Failed to create history: %q\n", err.Error())
	} else {
		log.Printf("Created new history…\n")

		if err := h.Delete(); err != nil {
			log.Printf("Failed to delete history: %q\n", err.Error())
		} else {
			log.Printf("Deleted new history…\n")
		}
	}
}
