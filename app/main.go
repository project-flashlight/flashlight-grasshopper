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
	os.Setenv("SERVER_URL", "https://api.sys.emea.vwapps.io")
	os.Setenv("NEW_RELIC_BASE_URL", "https://api.newrelic.com")
	os.Setenv("NEW_RELIC_API_KEY", "67fdf7e0128161526d582f7545da8e932d576629c2229ad")
	os.Setenv("CLOUDFOUNDRY_TOKEN", "bearer eyJhbGciOiJSUzI1NiIsImtpZCI6ImtleS0xIiwidHlwIjoiSldUIn0.eyJqdGkiOiI3NWY4MzVjM2U0ZTE0YjAwYjQ3YjYxNGZiYmYwYzQ4NCIsInN1YiI6IjRlNmE3YzAxLWExMGItNGEyNC05MWRjLWIyZTA1NzAwNmY4NiIsInNjb3BlIjpbImNsb3VkX2NvbnRyb2xsZXIucmVhZCIsInBhc3N3b3JkLndyaXRlIiwiY2xvdWRfY29udHJvbGxlci53cml0ZSIsIm9wZW5pZCIsInVhYS51c2VyIl0sImNsaWVudF9pZCI6ImNmIiwiY2lkIjoiY2YiLCJhenAiOiJjZiIsImdyYW50X3R5cGUiOiJwYXNzd29yZCIsInVzZXJfaWQiOiI0ZTZhN2MwMS1hMTBiLTRhMjQtOTFkYy1iMmUwNTcwMDZmODYiLCJvcmlnaW4iOiJ1YWEiLCJ1c2VyX25hbWUiOiJwcm9qZWN0LWZsYXNobGlnaHRAZ29vZ2xlZ3JvdXBzLmNvbSIsImVtYWlsIjoicHJvamVjdC1mbGFzaGxpZ2h0QGdvb2dsZWdyb3Vwcy5jb20iLCJyZXZfc2lnIjoiM2VhNGUyZGIiLCJpYXQiOjE1Mjk1MDU2NzksImV4cCI6MTUyOTUwNjI3OSwiaXNzIjoiaHR0cHM6Ly91YWEuc3lzLmVtZWEudndhcHBzLmlvL29hdXRoL3Rva2VuIiwiemlkIjoidWFhIiwiYXVkIjpbImNsb3VkX2NvbnRyb2xsZXIiLCJwYXNzd29yZCIsImNmIiwidWFhIiwib3BlbmlkIl19.VZgImHoWA_dI67HE02rcDYq9J8vRzeX8cTmX8JqgwyJMncF3h8l7nvvPcfT5N5tJgU94nZ4Zbx05CdN5tJFCcdz1e6GwJL2pc4D1uIGrMB0KjRCyueZMJobj5sMmu0N-86St7B-rYniDpsIiwn9q893G5RGru0Mm5sfq-SzoA7MlNHQRDDuhh7E-jKcefVpxd-Ad7g8OwSTLfzFtUdQRStG3rzpGycuZaZRclSYvwUsa3EmOP8_SwRdIfUyhcGX9fDrdczkFmrww2DqjTMAPBJ0SgLVz-LTF91YU0gedrHSZSjpjblrW9vEVHNxNrMozI5fZtFssYbC69STIzNgWH5AR20J079o7LRZG-h2BTvGTX1iKn3_HwIUaYe8CWkr8QDV0YNQK2X8J5dF4gUrpUVIA0ssHkTPEBRSMVb86IhWF0dGhQw3-InsRy1TpZq6-bhIjpzbUghQfYzrwH41LEKbDjlrXYo1vl8u9JA6ad3O8033aE0UMSEa3sjTG3ihobkqMJzZ-M4mMH0kCwzEkUBbSQCMdkA2XRXAM7yGGn0jR_sqZujUBt7BjvtM6vHyjk32ZiBDzCVauGSm3vfxUoHtmyL63E9vCqWwcDV5E2CdSjL6JUI2kg-G9C3nNaJ-rxHQCA4ZFZJceLUQZgg-KyQWuAcP13wM2A_TI1KgQb30")
}

func main() {

	cloudFoundryFetcher := grasshopper.NewCloudFoundryFetcher()
	apps, err := cloudFoundryFetcher.GetApps()
	if err != nil {
		panic(fmt.Sprintf("Error fetching Apps %s\n", err.Error()))
	} else if apps.Results == 0 {
		panic(fmt.Sprintf("Error fetching Apps, probably authorization token is not valid\n"))

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

	baseURL := os.Getenv("NEW_RELIC_BASE_URL")
	apiKey := os.Getenv("NEW_RELIC_API_KEY")
	newRelicFetcher := grasshopper.NewNewRelicFetcher(baseURL, apiKey)
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
