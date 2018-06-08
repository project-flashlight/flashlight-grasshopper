package main

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/vwdilab/flashlight-grasshopper/grasshopper"

	"github.com/julienschmidt/httprouter"
)

var appCtx context.Context
var router *httprouter.Router
var wg sync.WaitGroup

func init() {
	appCtx = context.Background()
	wg = sync.WaitGroup{}
	wg.Add(1)
	os.Setenv("CLOUDFOUNDRY_BASE_URL", "https://api.sys.emea.vwapps.io")
	os.Setenv("NEW_RELIC_BASE_URL", "https://api.newrelic.com")
	os.Setenv("NEW_RELIC_API_KEY", "67fdf7e0128161526d582f7545da8e932d576629c2229ad")
	os.Setenv("CLOUDFOUNDRY_TOKEN", "bearer eyJhbGciOiJSUzI1NiIsImtpZCI6ImtleS0xIiwidHlwIjoiSldUIn0.eyJqdGkiOiJhOTZkOGZiNTQ5NDU0MzVhOTU5MDEyN2JlZGM1NDEwMSIsInN1YiI6IjRlNmE3YzAxLWExMGItNGEyNC05MWRjLWIyZTA1NzAwNmY4NiIsInNjb3BlIjpbImNsb3VkX2NvbnRyb2xsZXIucmVhZCIsInBhc3N3b3JkLndyaXRlIiwiY2xvdWRfY29udHJvbGxlci53cml0ZSIsIm9wZW5pZCIsInVhYS51c2VyIl0sImNsaWVudF9pZCI6ImNmIiwiY2lkIjoiY2YiLCJhenAiOiJjZiIsImdyYW50X3R5cGUiOiJwYXNzd29yZCIsInVzZXJfaWQiOiI0ZTZhN2MwMS1hMTBiLTRhMjQtOTFkYy1iMmUwNTcwMDZmODYiLCJvcmlnaW4iOiJ1YWEiLCJ1c2VyX25hbWUiOiJwcm9qZWN0LWZsYXNobGlnaHRAZ29vZ2xlZ3JvdXBzLmNvbSIsImVtYWlsIjoicHJvamVjdC1mbGFzaGxpZ2h0QGdvb2dsZWdyb3Vwcy5jb20iLCJyZXZfc2lnIjoiM2VhNGUyZGIiLCJpYXQiOjE1MjgzNzg3MjMsImV4cCI6MTUyODM3OTMyMywiaXNzIjoiaHR0cHM6Ly91YWEuc3lzLmVtZWEudndhcHBzLmlvL29hdXRoL3Rva2VuIiwiemlkIjoidWFhIiwiYXVkIjpbImNsb3VkX2NvbnRyb2xsZXIiLCJwYXNzd29yZCIsImNmIiwidWFhIiwib3BlbmlkIl19.cF4XsTccc3DK-b6h9qQcUbUFru9fqWvOmSK8ndYSpqmbox5JpaGeQjFr-bMHmorySWMZvTJandH4zYb9-QswdQEENu8oDNi2bL3o4YLWBf1PjdExqQk1y6faJXpKQJ_rokc2LfxmMFNDrlr-XhNfyYCB6JiwLPb90fYBlwWYbYKwgsiMiqc-1cFwShaXnp1CK_4EGYIPT-_tqdVItCrcNJwfmwz5ZDwlFolKJ9f0cNpAKpanM0fOWFMuJpVcswUnn2jJ4QgNEXjX0ZirU5ORp9QdOrWc0LLG1RYP1jt2lWF3xgUr0VKPJtc0_81ovyFgNQo_1lh6UJwFDryMlM0cjie4gsvOEuqayh0oF4IvLii5yOi35H9YgDAhM4V8Y9AZDaozgaCa8zdPSVsVsBseofq0fV76mKQEzupBtlIQLIn8r4vjFOoRGOUok1CfMjhUpv6VEAGGB2ISDOGIphqn0PwfqXzKKZqP1SpZyzOpMxMjrqRb7cZWB90OFX21nX00KF3NFoBwLXcrXFTsqzdHkXcFcZJB0om5JdCdvQtBaVyM3ozR5If6e240wm-Bv8uAvcRp7e1DYjij0eUsOhd7usBk8T97RkH79MbZcjBVxLy1ZyU_AU_IDcloSn15oF_d2syJBtQ8nHUD3Y_mkm6xpMrDpmsnfVYGHu-z5uQVe8c")
}

func main() {
	cloudFoundryFetcher := grasshopper.NewCloudFoundryFetcher()
	apps, err := cloudFoundryFetcher.GetApps()
	if err != nil {
		panic(fmt.Sprintf("Error fetching Apps %s\n", err.Error()))
	}

	for _, entity := range apps.App {
		appName := entity.Entity.Name
		envVars := entity.Entity.Environment

		commitID := envVars["COMMIT_ID"]
		newRelicAppName := envVars["NEW_RELIC_APP_NAME"]

		if len(newRelicAppName) != 0 {
			fmt.Printf("App name: %s\n", appName)
			fmt.Printf("\tNew Relic App name: %s\n", newRelicAppName)
			fmt.Printf("\tCommitID: %s\n", commitID)
		}
	}

	newRelicFetcher := grasshopper.NewNewRelicFetcher()
	newRelicApps, err := newRelicFetcher.GetApps()
	if err != nil {
		panic(fmt.Sprintf("Error fetching Apps %s\n", err.Error()))
	}

	for _, entity := range newRelicApps.Applications {
		appId := entity.Id
		newrelicAppName := entity.Name
		fmt.Printf("New Relic App Id: %v\n", appId)
		fmt.Printf("\tNew Relic App name: %s\n", newrelicAppName)
	}
}

func loadPort() (string, bool) {
	p := os.Getenv("PORT")
	return p, p != ""
}
