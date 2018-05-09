package grasshopper_test

import (
	"bytes"
	"encoding/json"
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
	geckoboardService.EXPECT().PublishStatus(appStatus).Return(nil).Times(1)

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
