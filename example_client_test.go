package thingscloud

import (
	"log"
	"os"
)

func ExampleClient_Verify() {
	client := New(APIEndpoint, os.Getenv("THINGS_USERNAME"), os.Getenv("THINGS_PASSWORD"))

	_, err := client.Verify()
	if err != nil {
		log.Printf("Invalid Credentials: %q", err.Error())
	}
}

func ExampleClient_Histories() {
	client := New(APIEndpoint, os.Getenv("THINGS_USERNAME"), os.Getenv("THINGS_PASSWORD"))

	histories, err := client.Histories()
	if err != nil {
		log.Printf("Failed loading histories: %q", err.Error())
	}

	for _, history := range histories {
		if err := history.Sync(); err != nil {
			log.Printf("Failed syncing history: %q", err.Error())
		}
	}
}
