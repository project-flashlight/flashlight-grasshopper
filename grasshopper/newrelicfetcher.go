package grasshopper

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// NewRelicFetcher fetches applications from the NewRelic API
type NewRelicFetcher interface {
	GetApps() (*NewRelicApps, error)
}

type NewRelicApps struct {
	Applications []NewRelicApp `json:"applications"`
}

type NewRelicApp struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Revision string `json:"revision"`
}

type defaultNewRelicFetcher struct {
	ApiKey  string
	BaseUrl string
}

func NewNewRelicFetcher(baseURL string, apiKey string) NewRelicFetcher {
	return &defaultNewRelicFetcher{
		BaseUrl: baseURL,
		ApiKey:  apiKey,
	}
}

func (me *defaultNewRelicFetcher) GetApps() (*NewRelicApps, error) {
	path := "/v2/applications.json"
	newRelicApps := NewRelicApps{}

	req, _ := http.NewRequest(http.MethodGet, me.BaseUrl+path, nil)
	req.Header["X-Api-Key"] = []string{me.ApiKey}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	unmarshalError := json.Unmarshal(body, &newRelicApps)
	if unmarshalError != nil {
		return nil, unmarshalError
	}

	return &newRelicApps, nil
}
