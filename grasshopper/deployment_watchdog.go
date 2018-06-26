package grasshopper

import (
	"strings"
)

type DeploymentWatchdog interface {
	CheckApps() bool
}

type DefaultDeploymentWatchdog struct {
	cloudFoundryFetcher CloudFoundryFetcher
	newRelicFetcher     NewRelicFetcher
}

// New Deployment
func NewDeploymentWatchdog(cf CloudFoundryFetcher, nr NewRelicFetcher) DeploymentWatchdog {
	return &DefaultDeploymentWatchdog{
		cloudFoundryFetcher: cf,
		newRelicFetcher:     nr,
	}
}

func (me *DefaultDeploymentWatchdog) CheckApps() bool {
	cloudFoundryApps, _ := me.cloudFoundryFetcher.GetApps()
	newRelicApps, _ := me.newRelicFetcher.GetApps()

	for _, cloudFoundryApp := range cloudFoundryApps.App {
		for _, newRelicApp := range newRelicApps.Applications {
			if strings.Compare(cloudFoundryApp.Entity.Environment["NEW_RELIC_APP_NAME"], newRelicApp.Name) == 0 {
				if strings.Compare(cloudFoundryApp.Entity.Environment["COMMIT_ID"], newRelicApp.Revision) == 0 {
					return true
				}
				return false
			}
		}
	}

	return false
}
