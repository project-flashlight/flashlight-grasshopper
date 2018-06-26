package grasshopper_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/vwdilab/mango/assert"

	"github.com/vwdilab/flashlight-grasshopper/grasshopper"
	"github.com/vwdilab/flashlight-grasshopper/mocks"
)

type mockedCloudfoundryFetcher struct {
}

func (me *mockedCloudfoundryFetcher) GetApps() (*grasshopper.CloudFoundryApps, error) {
	return nil, nil
}

func Test_shouldPublishNewDeploymentWhenCommitIDMismatch(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)

	nr := mock_grasshopper.NewMockNewRelicFetcher(ctrl)
	cf := mock_grasshopper.NewMockCloudFoundryFetcher(ctrl)

	cloudfoundryApp1 := grasshopper.CloudFoundryEntities{
		Entity: grasshopper.CloudFoundryEntity{
			Name: "Name1",
			Environment: map[string]string{
				"COMMIT_ID":          "commitId1",
				"NEW_RELIC_APP_NAME": "TestApp1",
			},
		},
	}

	cloudfoundryApp2 := grasshopper.CloudFoundryEntities{
		Entity: grasshopper.CloudFoundryEntity{
			Name: "Name2",
			Environment: map[string]string{
				"COMMIT_ID":          "commitId2",
				"NEW_RELIC_APP_NAME": "TestApp2",
			},
		},
	}

	cloudFoundryAppList := []grasshopper.CloudFoundryEntities{cloudfoundryApp2, cloudfoundryApp1}
	cloundFoundaryApps := grasshopper.CloudFoundryApps{
		App:     cloudFoundryAppList,
		Results: 1,
	}

	newRelicApp1 := grasshopper.NewRelicApp{
		Id:       1000,
		Name:     "TestApp1",
		Revision: "commitId1",
	}

	newRelicApp2 := grasshopper.NewRelicApp{
		Id:       2000,
		Name:     "TestApp2",
		Revision: "commitId2",
	}

	newRelicAppList := []grasshopper.NewRelicApp{newRelicApp1, newRelicApp2}

	newRelicApps := grasshopper.NewRelicApps{
		Applications: newRelicAppList,
	}

	cf.EXPECT().GetApps().Return(&cloundFoundaryApps, nil)
	nr.EXPECT().GetApps().Return(&newRelicApps, nil)

	subject := grasshopper.NewDeploymentWatchdog(cf, nr)

	// when
	comparisionResult := subject.CheckApps()

	// then
	assert.Equal(t, true, comparisionResult)
}
