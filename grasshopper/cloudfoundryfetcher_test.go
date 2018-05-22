package grasshopper_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vwdilab/flashlight-grasshopper/grasshopper"
)

func startCloudFoundryMocked(body string) *httptest.Server {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/v2/apps") {
			panic("Endpoint url not mocked")
		}
		returnBody := body
		// body := `{
		// 	"resources": [{"entity":{"routes_url":"/v2/apps/5d1a6bf3-82e0-4c38-abf6-f944c14e03b9/routes"}},{"entity":{"routes_url": "/v2/apps/7d1a6bf3-82e0-4c38-abf6-f944c14e03b9/routes"}}]}`
		w.Write([]byte(returnBody))
	}))

	os.Setenv("VARIABLE", srv.URL)
	return srv
}

func Test_shouldGetListOfApps(t *testing.T) {
	// Given
	subject := grasshopper.NewCloudFoundryFetcher()

	server := startCloudFoundryMocked(`{ 
		"total_results": 2,
		"resources": [
		  {
			"entity": {
			  "routes_url": "/v2/apps/5d1a6bf3-82e0-4c38-abf6-f944c14e03b9/routes"
			}
		  },
		  {
			"entity": {
			  "routes_url": "/v2/apps/7d1a6bf3-82e0-4c38-abf6-f944c14e03b9/routes"
			}
		  }
		]
	  }`)
	os.Setenv("SERVER_URL", server.URL)
	defer server.Close()

	// When
	apps, err := subject.GetApps()

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, apps)
	assert.NotNil(t, apps.App)
	assert.Equal(t, 2, len(apps.App))
	assert.Equal(t, 2, apps.Results)
	assert.Equal(t, apps.App[0].Entity.RoutesUrl, "/v2/apps/5d1a6bf3-82e0-4c38-abf6-f944c14e03b9/routes")
	assert.Equal(t, apps.App[1].Entity.RoutesUrl, "/v2/apps/7d1a6bf3-82e0-4c38-abf6-f944c14e03b9/routes")

}

func Test_shouldReturnErroWhenMalformedJson(t *testing.T) {
	// Given
	subject := grasshopper.NewCloudFoundryFetcher()

	server := startCloudFoundryMocked(`invalid json`)
	os.Setenv("SERVER_URL", server.URL)
	defer server.Close()

	// When
	apps, err := subject.GetApps()

	// Then
	assert.Error(t, err)
	assert.Nil(t, apps)

}
