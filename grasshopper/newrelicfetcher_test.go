package grasshopper_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/vwdilab/flashlight-grasshopper/grasshopper"
	"github.com/vwdilab/mango/assert"
)

func startNewRelicMocked(body string, expectedPath string, expectedApiKey string) *httptest.Server {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, expectedPath) {
			panic("Endpoint url not mocked")
		}

		if !strings.Contains(r.Header["X-Api-Key"][0], expectedApiKey) {
			panic("API Key missing")
		}

		returnBody := body
		w.Write([]byte(returnBody))
	}))

	os.Setenv("NEWRELIC_URL", srv.URL)

	return srv
}

func Test_shouldGetListOfNewRelicApps(t *testing.T) {
	// Given
	apiKey := "someKey"
	body := `{
		"applications": [
			{
				"id": 111245474,
				"name": "production-emea-chameleon"
			},
			{
				"id": 109719010,
				"name": "production-ema-blackhole"
			}
		]
	}`
	mockedServer := startNewRelicMocked(body, "v2/applications.json", apiKey)
	defer mockedServer.Close()

	baseURL := os.Getenv("NEWRELIC_URL")
	subject := grasshopper.NewNewRelicFetcher(baseURL, apiKey)

	// When
	apps, err := subject.GetApps()

	// Then
	assert.NoError(t, err)
	assert.Len(t, 2, apps.Applications)

	app1 := grasshopper.NewRelicApp{Id: 109719010, Name: "production-ema-blackhole"}
	assert.Contains(t, apps.Applications, app1)

	app2 := grasshopper.NewRelicApp{Id: 111245474, Name: "production-emea-chameleon"}
	assert.Contains(t, apps.Applications, app2)
}
