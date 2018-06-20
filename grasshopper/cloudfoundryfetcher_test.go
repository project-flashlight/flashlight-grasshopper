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

func startCloudFoundryMocked(body string, expectedPath string, expectedAuthToken string) *httptest.Server {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, expectedPath) {
			panic("Endpoint url not mocked")
		}

		if !strings.Contains(r.Header["Authorization"][0], expectedAuthToken) {
			panic("Authorization missing")
		}

		returnBody := body
		w.Write([]byte(returnBody))
	}))

	os.Setenv("SERVER_URL", srv.URL)

	return srv
}

func Test_shouldGetListOfApps(t *testing.T) {
	// Given
	authorization := "bearer myToken"
	os.Setenv("CLOUDFOUNDRY_TOKEN", authorization)
	server := startCloudFoundryMocked(`{ 
		"total_results": 2,
		"resources": [
			{
				"entity": {
					"routes_url": "/v2/apps/5d1a6bf3-82e0-4c38-abf6-f944c14e03b9/routes",
					"environment_json": {
						"ENV1": "BLA",
						"ENV2": "WULULU"
					}
				}
			},
			{
				"entity": {
					"routes_url": "/v2/apps/7d1a6bf3-82e0-4c38-abf6-f944c14e03b9/routes",
					"environment_json": {
						"ENV3": "BLABla"
					}
				}
			}
			]
			}`, "/v2/apps", authorization)
	defer server.Close()

	subject := grasshopper.NewCloudFoundryFetcher()

	// When
	apps, err := subject.GetApps()

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, apps)
	assert.NotNil(t, apps.App)
	assert.Equal(t, 2, len(apps.App))
	assert.Equal(t, 2, apps.Results)
	assert.Equal(t, apps.App[0].Entity.Environment["ENV1"], "BLA")
	assert.Equal(t, apps.App[0].Entity.Environment["ENV2"], "WULULU")
	assert.Equal(t, apps.App[1].Entity.Environment["ENV3"], "BLABla")
}

func Test_shouldReturnErroWhenMalformedJsonOnAppsFetching(t *testing.T) {
	// Given
	authorization := "bearer myToken"
	os.Setenv("CLOUDFOUNDRY_TOKEN", authorization)
	server := startCloudFoundryMocked(`invalid json`, "/v2/apps", authorization)
	defer server.Close()

	subject := grasshopper.NewCloudFoundryFetcher()

	// When
	apps, err := subject.GetApps()

	// Then
	assert.Error(t, err)
	assert.Nil(t, apps)
}
