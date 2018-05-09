package grasshopper

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

type GeckoboardService interface {
	PublishStatus(appStatus AppStatus) error
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

func (me *DefGeckoboardService) PublishStatus(appStatus AppStatus) error {
	url := os.Getenv("VARIABLE")
	// add variable https://api.geckoboard.com/datasets/flashlight.services
	endpoint := url + "/data"

	data := DataSet{[]AppStatus{appStatus}}
	b, _ := json.Marshal(data)
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(b))

	if err != nil {
		return err
	}

	_, errDo := http.DefaultClient.Do(req)

	return errDo
}
