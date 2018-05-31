package grasshopper

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

type CloudFoundryFetcher interface {
	GetApps() (*CloudFoundryApps, error)
	GetAppByRoute(app CloudFoundryEntities) (*CloudFoundryAppRoute, error)
}

type DefCloudFoundryFetcher struct{}

type CloudFoundryApps struct {
	App     []CloudFoundryEntities `json:"resources"`
	Results int                    `json:"total_results"`
}

type CloudFoundryAppRoute struct {
	Resources []CloudFoundryEntities `json:"resources"`
}

type CloudFoundryEntities struct {
	Entity CloudFoundryEntity `json:"entity"`
}

// CloudFoundryEntity is a multipurpose container, used to parse the cf response
type CloudFoundryEntity struct {
	RoutesURL string `json:"routes_url"`
	DomainURL string `json:"domain_url"`
	Host      string `json:"host"`
}

func NewCloudFoundryFetcher() *DefCloudFoundryFetcher {
	return &DefCloudFoundryFetcher{}
}

func (me *DefCloudFoundryFetcher) GetApps() (*CloudFoundryApps, error) {
	cloudFoundryApps := CloudFoundryApps{}
	err := getRequest("/v2/apps", &cloudFoundryApps)

	if err != nil {
		return nil, err
	}

	return &cloudFoundryApps, nil
}

// GetAppByRoute does something
func (me *DefCloudFoundryFetcher) GetAppByRoute(app CloudFoundryEntities) (*CloudFoundryAppRoute, error) {
	route := app.Entity.RoutesURL

	cloudFoundryAppRoute := CloudFoundryAppRoute{}
	err := getRequest(route, &cloudFoundryAppRoute)

	if err != nil {
		return nil, err
	}

	return &cloudFoundryAppRoute, nil
}

func getRequest(path string, structPointer interface{}) error {
	// baseUrl is https://api.sys.emea.vwapps.io
	baseURL := os.Getenv("SERVER_URL")
	req, _ := http.NewRequest(http.MethodGet, baseURL+path, nil)

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
