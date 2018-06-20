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

type CloudFoundryApp struct {
	AppName string
	AppUrl  string
	Stage   string
}

type DefCloudFoundryFetcher struct {
	AuthToken string
	BaseUrl   string
}

type CloudFoundryApps struct {
	App     []CloudFoundryEntities `json:"resources"`
	Results int                    `json:"total_results"`
}

type CloudFoundryEntities struct {
	Entity CloudFoundryEntity `json:"entity"`
}

// CloudFoundryEntity is a multipurpose container, used to parse the cf response
type CloudFoundryEntity struct {
	Name        string            `json:"name"`
	Environment map[string]string `json:"environment_json"`
}

func NewCloudFoundryFetcher() *DefCloudFoundryFetcher {
	baseURL := os.Getenv("SERVER_URL")
	authToken := os.Getenv("CLOUDFOUNDRY_TOKEN")
	return &DefCloudFoundryFetcher{
		AuthToken: authToken,
		BaseUrl:   baseURL,
	}
}

func (me *DefCloudFoundryFetcher) GetApps() (*CloudFoundryApps, error) {
	cloudFoundryApps := CloudFoundryApps{}
	err := getRequest(me.BaseUrl+"/v2/apps", me.AuthToken, &cloudFoundryApps)

	if err != nil {
		return nil, err
	}

	return &cloudFoundryApps, nil
}

func getRequest(url string, token string, structPointer interface{}) error {
	// baseUrl is https://api.sys.emea.vwapps.io
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header["Authorization"] = []string{token}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	unmarshalError := json.Unmarshal(body, &structPointer)
	if unmarshalError != nil {
		return unmarshalError
	}
	resp.Body.Close()

	return nil
}
