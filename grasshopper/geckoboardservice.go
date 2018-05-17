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
	AppName  string `json:"app_name"`
	CommitID string `json:"commit_id"`
	Date     string `json:"date"`
	Stage    string `json:"stage"`
	Status   string `json:"status"`
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
	body, _ := json.Marshal(data)

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(body))
	req.Header.Add("Authorization", "API_KEY_BASE64")
	if err != nil {
		return nil, err
	}

	return http.DefaultClient.Do(req)
}
