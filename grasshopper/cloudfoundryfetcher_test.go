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

func startCloudFoundryMocked(body string, expectedPath string, mockedRequests *MockedRequests) *httptest.Server {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, expectedPath) {
			panic("Endpoint url not mocked")
		}
		mockedRequests.counter++
		returnBody := body
		w.Write([]byte(returnBody))
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
	  }`, "/v2/apps", &MockedRequests{})
	defer server.Close()

	// When
	apps, err := subject.GetApps()

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, apps)
	assert.NotNil(t, apps.App)
	assert.Equal(t, 2, len(apps.App))
	assert.Equal(t, 2, apps.Results)
	assert.Equal(t, apps.App[0].Entity.RoutesURL, "/v2/apps/5d1a6bf3-82e0-4c38-abf6-f944c14e03b9/routes")
	assert.Equal(t, apps.App[1].Entity.RoutesURL, "/v2/apps/7d1a6bf3-82e0-4c38-abf6-f944c14e03b9/routes")

}

func Test_shouldReturnErroWhenMalformedJsonOnAppsFetching(t *testing.T) {
	// Given
	subject := grasshopper.NewCloudFoundryFetcher()

	server := startCloudFoundryMocked(`invalid json`, "/v2/apps", &MockedRequests{})
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
}`, path, &MockedRequests{})
	defer server.Close()

	inputApp := grasshopper.CloudFoundryEntities{
		Entity: grasshopper.CloudFoundryEntity{
			RoutesURL: path,
		},
	}

	// When
	app, err := subject.GetAppByRoute(inputApp)

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, app)
	assert.Equal(t, app.Resources[0].Entity.DomainURL, "/v2/shared_domains/e97d1675-894e-4808-8b58-82d1805b7368")
	assert.Equal(t, app.Resources[0].Entity.Host, "dashlight-acceptance")
}

func Test_shouldReturnErroWhenMalformedJsonOnRouteFetching(t *testing.T) {
	// Given
	subject := grasshopper.NewCloudFoundryFetcher()

	server := startCloudFoundryMocked(`invalid json`, "/v2", &MockedRequests{})
	defer server.Close()

	// When
	apps, err := subject.GetAppByRoute(grasshopper.CloudFoundryEntities{})

	// Then
	assert.Error(t, err)
	assert.Nil(t, apps)
}

func Test_shouldGetNameByDomain(t *testing.T) {
	// Given
	subject := grasshopper.NewCloudFoundryFetcher()

	mockedRequests := MockedRequests{}

	path := "/this/is/a/test/url"
	server := startCloudFoundryMocked(`{
		"entity": {
			"name": "apps.emea.vwapps.io"
		}
	}`, path, &mockedRequests)
	defer server.Close()

	// When
	domainName, err := subject.GetDomainName(path)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, "apps.emea.vwapps.io", domainName)
	assert.Equal(t, 1, mockedRequests.counter)
}

func Test_shouldReturnErrorWhenMalformedJsonOnDomainNameRouteFetching(t *testing.T) {
	// Given
	subject := grasshopper.NewCloudFoundryFetcher()

	mockedRequests := MockedRequests{}

	path := "/this/is/a/test/url"
	server := startCloudFoundryMocked(`invalid json`, path, &mockedRequests)
	defer server.Close()

	// When
	domainName, err := subject.GetDomainName(path)

	// Then
	assert.Error(t, err)
	assert.Equal(t, "", domainName)
}

type MockedRequests struct {
	counter int
}

func Test_shouldGetAppsWithEndpoint(t *testing.T) {
	// Given
	subject := grasshopper.NewCloudFoundryFetcher()

	mockedRequests := MockedRequests{}

	path := "/this/is/a/test/url"
	server := startCloudFoundryMocked(`invalid json`, path, &mockedRequests)
	defer server.Close()

	// When
	apps, err := subject.FetchCloudFoundryApps()

	// Then
	assert.NoError(t, err)
	assert.Len(t, 2, apps)
	appOne := apps[0]
	assert.Equal(t, "App1", appOne.AppName)
	assert.Equal(t, "https://someapp-acceptance.apps.emea.vwapps.io", appOne.AppUrl)
	assert.Equal(t, "acceptance", appOne.Stage)
}
