package grasshopper

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

type CloudFoundryFetcher interface {
	GetApps() (*CloudFoundryApps, error)
}

type DefCloudFoundryFetcher struct{}

type CloudFoundryApps struct {
	App     []CloudFoundryApp `json:"resources"`
	Results int               `json:"total_results"`
}
type CloudFoundryApp struct {
	Entity CloudFoundryEntity `json:"entity"`
}
type CloudFoundryEntity struct {
	RoutesUrl string `json:"routes_url"`
}

func NewCloudFoundryFetcher() *DefCloudFoundryFetcher {
	return &DefCloudFoundryFetcher{}
}

func (me *DefCloudFoundryFetcher) GetApps() (*CloudFoundryApps, error) {

	url := os.Getenv("SERVER_URL")
	req, _ := http.NewRequest(http.MethodGet, url+"/v2/apps", nil)

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	cloudFoundryApps := CloudFoundryApps{}
	unmarshalError := json.Unmarshal(body, &cloudFoundryApps)
	if unmarshalError != nil {
		return nil, unmarshalError
	}
	resp.Body.Close()

	return &cloudFoundryApps, nil
}
