package grasshopper_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/vwdilab/flashlight-grasshopper/grasshopper"
)

var interactionsWithGrasshopper = 0

type PublishResponse struct {
	Status string `json:"status"`
}

func startGrasshopperServerMocked() *httptest.Server {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Endpoint hit\n")
		fmt.Printf("Endpoint hit2\n")
		if !strings.Contains(r.URL.Path, "data") {
			panic("Endpoint url not mocked: " + r.URL.Path)
		}
		res := PublishResponse{
			Status: "ok",
		}
		bytes, err := json.Marshal(res)
		if err != nil {
			panic("Could not marshal response message")
		}
		interactionsWithGrasshopper++
		w.Write(bytes)
	}))

	os.Setenv("VARIABLE", srv.URL)
	return srv
}

// func init() {
// 	server := initServer()
// 	fmt.Println(server)
// }

func testStartserver(t *testing.T) func() {
	// init mocks
	geckoboardService := grasshopper.NewGeckoboardService()

	// create subject
	s := grasshopper.NewServer(*geckoboardService)
	geckoboardMockServer := startGrasshopperServerMocked()

	s.Start(":2345")
	testWaitForServer()

	return func() {
		s.Stop()
		geckoboardMockServer.Close()
		testWaitForServer()
	}
}

func testWaitForServer() {
	if os.Getenv("IS_PIPELINE") != "" {
		time.Sleep(time.Second * 5)
	}
	time.Sleep(time.Millisecond * 200)
}

func Test_Something(t *testing.T) {
	// given
	interactionsWithGrasshopper = 0
	defer testStartserver(t)()

	// when
	req, _ := http.NewRequest("POST", "http://0.0.0.0:2345/publish", nil)
	resp, err := http.DefaultClient.Do(req)

	// then
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, interactionsWithGrasshopper, 1)

}
