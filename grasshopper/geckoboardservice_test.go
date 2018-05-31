package grasshopper_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/vwdilab/mango/assert"

	"github.com/vwdilab/flashlight-grasshopper/grasshopper"
)

func startGeckoboardServerMocked(callback func(bs string)) *httptest.Server {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Endpoint hit")
		if !strings.Contains(r.URL.Path, "data") {
			panic("Endpoint url not mocked")
		}

		bodyBytes, err := ioutil.ReadAll(r.Body)

		if err != nil {
			panic("There was an error parsing the request: ")

		}

		bodyString := string(bodyBytes)
		callback(bodyString)
		fmt.Printf("Body: " + bodyString)
	}))

	os.Setenv("VARIABLE", srv.URL)
	return srv
}

func Test_shouldSendDataSetToGeckoBoard(t *testing.T) {
	// given
	actualBody := ""
	geckoboardMockServer := startGeckoboardServerMocked(func(bodyString string) {
		actualBody = bodyString
	})

	defer geckoboardMockServer.Close()

	geckoboardService := grasshopper.NewGeckoboardService()
	appStatus := grasshopper.AppStatus{
		AppName:  "dummy",
		CommitID: "asdasda",
		Date:     "2018-05-03T12:00:00Z",
		Stage:    "Production",
		Status:   "up",
	}

	// make sure to use this string to catch json marshalling issues
	expectedRequestBody := `{"data":[{"app_name":"dummy","commit_id":"asdasda","date":"2018-05-03T12:00:00Z","stage":"Production","status":"up"}]}`

	// when
	resp, err := geckoboardService.PublishStatus(appStatus)

	// then
	assert.Equal(t, resp.StatusCode, 200)
	assert.NoError(t, err)
	assert.Equal(t, expectedRequestBody, actualBody)
}
