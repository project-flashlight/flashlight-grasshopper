package grasshopper_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/vwdilab/flashlight-grasshopper/grasshopper"
	"github.com/vwdilab/flashlight-grasshopper/grasshopper/mocks"
)

var subject *grasshopper.Server
var geckoboardService *mock_grasshopper.MockGeckoboardService

type MockedAppStatus struct {
	CommitID string `json:"CommitId"`
	Date     string `json:"Datum"`
	Stage    string `json:"Stage"`
	Status   string `json:"Status"`
}

type PublishResponse struct {
	Status string `json:"status"`
}

// func init() {
// 	server := initServer()
// 	fmt.Println(server)
// }

func testStartserver(t *testing.T) func() {
	// init mocks
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	geckoboardService = mock_grasshopper.NewMockGeckoboardService(ctrl)

	// create subject
	subject = grasshopper.NewServer(geckoboardService)

	subject.Start(":2345")
	testWaitForServer()

	return func() {
		subject.Stop()
		testWaitForServer()
	}
}

func testWaitForServer() {
	if os.Getenv("IS_PIPELINE") != "" {
		time.Sleep(time.Second * 5)
	}
	time.Sleep(time.Millisecond * 200)
}

func Test_publishStatus_whenReceivingPostRequest(t *testing.T) {
	// given
	defer testStartserver(t)()

	appStatus := grasshopper.AppStatus{
		AppName:  "dummy",
		CommitID: "asdasda",
		Date:     "2018-03-01",
		Stage:    "Production",
		Status:   "down",
	}

	requestBody, _ := json.Marshal(appStatus)
	someMockResponse := http.Response{}
	someMockResponse.StatusCode = http.StatusOK

	geckoboardService.EXPECT().PublishStatus(appStatus).Return(&someMockResponse, nil).Times(1)

	// when
	req, _ := http.NewRequest("POST", "http://0.0.0.0:2345/publish", bytes.NewBuffer([]byte(requestBody)))
	resp, err := http.DefaultClient.Do(req)

	// then
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, 200)
}
func Test_return400_whenFieldIsIncorrect(t *testing.T) {
	// given
	defer testStartserver(t)()

	incorrectAppStatusBody := MockedAppStatus{
		CommitID: "asdasda",
		Date:     "2018-03-01",
		Stage:    "Production",
		Status:   "down",
	}

	requestBody, _ := json.Marshal(incorrectAppStatusBody)

	// when
	req, _ := http.NewRequest("POST", "http://0.0.0.0:2345/publish", bytes.NewBuffer([]byte(requestBody)))
	resp, err := http.DefaultClient.Do(req)

	// then
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

// This BufferCloser is required to mock a io.ReaderCloser in the response body
type BufferCloser struct {
	io.Reader
}

func (BufferCloser) Close() error { return nil }

func Test_return403_whenGeckoboardReturnsError(t *testing.T) {
	// given
	defer testStartserver(t)()

	appStatus := grasshopper.AppStatus{
		AppName:  "dummy",
		CommitID: "asdasda",
		Date:     "2018-03-01",
		Stage:    "Production",
		Status:   "down",
	}

	requestBody, _ := json.Marshal(appStatus)
	someMockResponse := http.Response{}
	someMockResponse.StatusCode = http.StatusForbidden
	someMockResponse.Body = BufferCloser{bytes.NewBufferString("Some Geckoboard Forbidden Error MSG")}

	someMockError := errors.New("Some random error")

	geckoboardService.EXPECT().PublishStatus(appStatus).Return(&someMockResponse, someMockError).Times(1)

	// when
	req, _ := http.NewRequest("POST", "http://0.0.0.0:2345/publish", bytes.NewBuffer([]byte(requestBody)))
	resp, _ := http.DefaultClient.Do(req)

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	responseBodyString := string(bodyBytes)

	// then
	assert.Equal(t, resp.StatusCode, 403)
	assert.Equal(t, "Failed to invoke geckoboard: Some Geckoboard Forbidden Error MSG\n", responseBodyString)
}

func Test_return500_whenGeckoboardDoNotRespond(t *testing.T) {
	// given
	defer testStartserver(t)()

	appStatus := grasshopper.AppStatus{
		AppName:  "dummy",
		CommitID: "asdasda",
		Date:     "2018-03-01",
		Stage:    "Production",
		Status:   "down",
	}

	requestBody, _ := json.Marshal(appStatus)
	someMockError := errors.New("Some random error")

	geckoboardService.EXPECT().PublishStatus(appStatus).Return(nil, someMockError).Times(1)

	// when
	req, _ := http.NewRequest("POST", "http://0.0.0.0:2345/publish", bytes.NewBuffer([]byte(requestBody)))
	resp, _ := http.DefaultClient.Do(req)

	// then
	assert.Equal(t, resp.StatusCode, 500)
}
