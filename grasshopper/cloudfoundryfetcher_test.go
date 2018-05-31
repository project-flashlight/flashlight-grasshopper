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

func startCloudFoundryMocked(body string, expectedPath string) *httptest.Server {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, expectedPath) {
			panic("Endpoint url not mocked")
		}
		returnBody := body
		w.Write([]byte(returnBody))
	}))

	os.Setenv("SERVER_URL", srv.URL)

	return srv
}

func startCloudFoundryMockedBufferError(expectedPath string) *httptest.Server {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, expectedPath) {
			panic("Endpoint url not mocked")
		}
		w.Write([]byte(nil))
	}))

	os.Setenv("SERVER_URL", srv.URL)

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
	  }`, "/v2/apps")
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

func Test_shouldReturnErroWhenMalformedJsonOnAppsFetching(t *testing.T) {
	// Given
	subject := grasshopper.NewCloudFoundryFetcher()

	server := startCloudFoundryMocked(`invalid json`, "/v2/apps")
	defer server.Close()

	// When
	apps, err := subject.GetApps()

	// Then
	assert.Error(t, err)
	assert.Nil(t, apps)
}

func Test_shouldGetAppByRoute(t *testing.T) {
	// Given
	subject := grasshopper.NewCloudFoundryFetcher()
	path := "/v2/apps/e97d1675-894e-4808-8b58-82d1805b7368/routes"
	server := startCloudFoundryMocked(`{
    "resources": [
        {
            "entity": {
                "host": "dashlight-acceptance",
                "domain_url": "/v2/shared_domains/e97d1675-894e-4808-8b58-82d1805b7368"
            }
        }
    ]
}`, path)
	defer server.Close()

	inputApp := grasshopper.CloudFoundryApp{
		Entity: grasshopper.CloudFoundryEntity{
			RoutesUrl: path,
		},
	}

	// When
	app, err := subject.GetAppByRoute(inputApp)

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, app)
	assert.Equal(t, app.Resources[0].Entity.DomainUrl, "/v2/shared_domains/e97d1675-894e-4808-8b58-82d1805b7368")
	assert.Equal(t, app.Resources[0].Entity.Host, "dashlight-acceptance")
}

func Test_shouldReturnErroWhenMalformedJsonOnRouteFetching(t *testing.T) {
	// Given
	subject := grasshopper.NewCloudFoundryFetcher()

	server := startCloudFoundryMocked(`invalid json`, "/v2")
	defer server.Close()

	// When
	apps, err := subject.GetAppByRoute(grasshopper.CloudFoundryApp{})

	// Then
	assert.Error(t, err)
	assert.Nil(t, apps)
}

func Test_shouldReturnErroWhenBufferWronglyFormatted(t *testing.T) {
	// Given
	subject := grasshopper.NewCloudFoundryFetcher()

	server := startCloudFoundryMockedBufferError("/v2")
	defer server.Close()

	// When
	apps, err := subject.GetAppByRoute(grasshopper.CloudFoundryApp{})

	// Then
	assert.Error(t, err)
	assert.Nil(t, apps)
}
