package iottestfuncs

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	iotSDK "github.com/scaleway/scaleway-sdk-go/api/iot/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/logging"
)

func AddTestSweepers() {
	resource.AddTestSweepers("scaleway_iot_hub", &resource.Sweeper{
		Name: "scaleway_iot_hub",
		F:    testSweepHub,
	})
}

func testSweepHub(_ string) error {
	return acctest.SweepRegions(scw.AllRegions, func(scwClient *scw.Client, region scw.Region) error {
		iotAPI := iotSDK.NewAPI(scwClient)
		logging.L.Debugf("sweeper: destroying the iot hub in (%s)", region)
		listHubs, err := iotAPI.ListHubs(&iotSDK.ListHubsRequest{Region: region}, scw.WithAllPages())
		if err != nil {
			logging.L.Debugf("sweeper: destroying the iot hub in (%s)", region)
			return fmt.Errorf("error listing hubs in (%s) in sweeper: %s", region, err)
		}

		deleteDevices := true
		for _, hub := range listHubs.Hubs {
			err := iotAPI.DeleteHub(&iotSDK.DeleteHubRequest{
				HubID:         hub.ID,
				Region:        hub.Region,
				DeleteDevices: &deleteDevices,
			})
			if err != nil {
				return fmt.Errorf("error deleting hub in sweeper: %s", err)
			}
		}

		return nil
	})
}
