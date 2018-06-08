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
	os.Setenv("CLOUDFOUNDRY_TOKEN", "bearer eyJhbGciOiJSUzI1NiIsImtpZCI6ImtleS0xIiwidHlwIjoiSldUIn0.eyJqdGkiOiI0YTBkMDY1ZTlhNWY0OTQ3YTMxYzYzOTVjNmRmNzAzYyIsInN1YiI6IjRlNmE3YzAxLWExMGItNGEyNC05MWRjLWIyZTA1NzAwNmY4NiIsInNjb3BlIjpbImNsb3VkX2NvbnRyb2xsZXIucmVhZCIsInBhc3N3b3JkLndyaXRlIiwiY2xvdWRfY29udHJvbGxlci53cml0ZSIsIm9wZW5pZCIsInVhYS51c2VyIl0sImNsaWVudF9pZCI6ImNmIiwiY2lkIjoiY2YiLCJhenAiOiJjZiIsImdyYW50X3R5cGUiOiJwYXNzd29yZCIsInVzZXJfaWQiOiI0ZTZhN2MwMS1hMTBiLTRhMjQtOTFkYy1iMmUwNTcwMDZmODYiLCJvcmlnaW4iOiJ1YWEiLCJ1c2VyX25hbWUiOiJwcm9qZWN0LWZsYXNobGlnaHRAZ29vZ2xlZ3JvdXBzLmNvbSIsImVtYWlsIjoicHJvamVjdC1mbGFzaGxpZ2h0QGdvb2dsZWdyb3Vwcy5jb20iLCJyZXZfc2lnIjoiM2VhNGUyZGIiLCJpYXQiOjE1Mjk0OTUyNTcsImV4cCI6MTUyOTQ5NTg1NywiaXNzIjoiaHR0cHM6Ly91YWEuc3lzLmVtZWEudndhcHBzLmlvL29hdXRoL3Rva2VuIiwiemlkIjoidWFhIiwiYXVkIjpbImNsb3VkX2NvbnRyb2xsZXIiLCJwYXNzd29yZCIsImNmIiwidWFhIiwib3BlbmlkIl19.OWfEu8r65M8N1ljEqD7KikwPzpcToEB5_wVlf2nr9FTidtHRy39yQ5NnAI-gEoAf6cDTIVIEbdrjQwIHXxHDcSIZVB1IDHf_h4ql2DwYFSn7U1n6Kd_SB_ZDD3XqBuKaO4EvaX-CgaVvQ8R6LSdX1yPdUXU9-HfJAJpEwWfISwVFLxHDUQgeEiFU_AOrvf_yYV7i8gA7RsfQfoTQ1mQeEyCOx2obln5uC8hQQF3OHPIGW4IdmWLwzqEWsv9tix07mHKsMhBArBB1dUaZJhcs6vpKmyo3rqXu-5vW14Ft_y2xm0Sdtbrj0PD5HnoQ1_nwDmyM3WxZkcfj6-H49tpJaPb-nVcuZCSMPNYdW27Qr7rO2ZAdC94b2AV3oDCO8vSlzo_gvThqiQbdtjj_sHID2MU716cQAnc3FDtY8AArMz0659FesEa3UwPsrUPrIz9BTlFYo-pM1QQRiathXAVAxCltN29opYWx6601B5GdWaKGZvBX3fVg_dW7XygAP76uZ4npKFKnXE85nbe94_JiOfbctvjbm2eeWgthjVMFP7ZqN-2N1Pf8wVR48lV-EysLoUd0cVFvG1MVr4K3EYIyauuhtkVvTjtU7bIeZ__TgtdF0ZgB-QdBRrbCv-GGsPyOY8WDp87GK0Y7g0H8nzOy_Uq3QZMHX1zhqMHkJlNK6XE")
}

func main() {

	cloudFoundryFetcher := grasshopper.NewCloudFoundryFetcher()
	apps, err := cloudFoundryFetcher.GetApps()
	if err != nil {
		panic(fmt.Sprintf("Error fetching Apps %s\n", err.Error()))
	}

	fmt.Println("apps: %v", apps)

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
	newRelicFetcher := grasshopper.NewNewRelicFetcher(baseURL)
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
