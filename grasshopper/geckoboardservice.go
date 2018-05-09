package grasshopper

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

type GeckoboardService interface {
	PublishStatus(appStatus AppStatus) (*http.Response, error)
}

type AppStatus struct {
	AppName  string `json:"AppName"`
	CommitID string `json:"CommitId"`
	Date     string `json:"Date"`
	Stage    string `json:"Stage"`
	Status   string `json:"Status"`
}

type DataSet struct {
	Data []AppStatus `json:"data"`
}

type DefGeckoboardService struct{}

func NewGeckoboardService() *DefGeckoboardService {
	return &DefGeckoboardService{}
}

func (me *DefGeckoboardService) PublishStatus(appStatus AppStatus) (*http.Response, error) {
	url := os.Getenv("VARIABLE")
	// add variable https://api.geckoboard.com/datasets/flashlight.services
	endpoint := url + "/data"

	data := DataSet{[]AppStatus{appStatus}}
	b, _ := json.Marshal(data)
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(b))
	req.Header.Add("Authorization", "API_KEY_BASE64")

	if err != nil {
		return nil, err
	}

	return http.DefaultClient.Do(req)
}
