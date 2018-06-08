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
	GetDomainName(path string) (string, error)
	FetchCloudFoundryApps() ([]CloudFoundryApp, error)
}

type CloudFoundryApp struct {
	AppName string
	AppUrl  string
	Stage   string
}

type DefCloudFoundryFetcher struct {
	BaseUrl string
}

type CloudFoundryApps struct {
	App     []CloudFoundryEntities `json:"resources"`
	Results int                    `json:"total_results"`
}

type CloudFoundryAppRoute struct {
	Resources []CloudFoundryEntities `json:"resources"`
}

type CloudFoundryDomain struct {
	Entity CloudFoundryEntity `json:"entity"`
}

type CloudFoundryEntities struct {
	Entity CloudFoundryEntity `json:"entity"`
}

// CloudFoundryEntity is a multipurpose container, used to parse the cf response
type CloudFoundryEntity struct {
	RoutesURL   string            `json:"routes_url"`
	DomainURL   string            `json:"domain_url"`
	Host        string            `json:"host"`
	Name        string            `json:"name"`
	Environment map[string]string `json:"environment_json"`
}

func NewCloudFoundryFetcher() *DefCloudFoundryFetcher {
	baseURL := os.Getenv("SERVER_URL")
	return &DefCloudFoundryFetcher{BaseUrl: baseURL}
}

func (me *DefCloudFoundryFetcher) GetApps() (*CloudFoundryApps, error) {
	cloudFoundryApps := CloudFoundryApps{}
	err := getRequest("/v2/apps", &cloudFoundryApps)

	if err != nil {
		return nil, err
	}

	return &cloudFoundryApps, nil
}

func (me *DefCloudFoundryFetcher) FetchCloudFoundryApps() ([]CloudFoundryApp, error) {

	t := []CloudFoundryApp{
		CloudFoundryApp{
			AppName: "App1",
			AppUrl:  "asad",
			Stage:   "",
		},
		CloudFoundryApp{
			AppName: "App2",
			AppUrl:  "asad",
			Stage:   "",
		},
	}
	return t, nil
}

// GetAppByRoute does something
func (me *DefCloudFoundryFetcher) GetAppByRoute(app CloudFoundryEntities) (*CloudFoundryAppRoute, error) {
	route := app.Entity.RoutesURL

	cloudFoundryAppRoute := CloudFoundryAppRoute{}
	err := getRequest(me.BaseUrl+route, &cloudFoundryAppRoute)

	if err != nil {
		return nil, err
	}

	return &cloudFoundryAppRoute, nil
}

// GetDomainName does something
func (me *DefCloudFoundryFetcher) GetDomainName(path string) (string, error) {

	cloudFoundryDomain := CloudFoundryDomain{}
	err := getRequest(me.BaseUrl+path, &cloudFoundryDomain)

	if err != nil {
		return "", err
	}

	return cloudFoundryDomain.Entity.Name, nil
}

func getRequest(url string, structPointer interface{}) error {
	// baseUrl is https://api.sys.emea.vwapps.io

	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header["Authorization"] = []string{"bearer blabla"}

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
