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

type CloudFoundryAppRoute struct {
	Resources []CloudFoundryAppRouteResource `json:"resources"`
}

type CloudFoundryAppRouteResource struct {
	Entity CloudFoundryAppRouteEntity `json:"entity"`
}

type CloudFoundryAppRouteEntity struct {
	DomainUrl string `json:"domain_url"`
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
func (me *DefCloudFoundryFetcher) GetAppByRoute(app CloudFoundryApp) (*CloudFoundryAppRoute, error) {
	route := app.Entity.RoutesUrl

	cloudFoundryAppRoute := CloudFoundryAppRoute{}
	err := getRequest(route, &cloudFoundryAppRoute)

	if err != nil {
		return nil, err
	}

	return &cloudFoundryAppRoute, nil
}

func getRequest(path string, structPointer interface{}) error {
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
