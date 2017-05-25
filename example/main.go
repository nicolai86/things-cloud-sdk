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
		log.Printf("Login failed: %q\nPlease check your credentials.", err.Error())
	}
}
