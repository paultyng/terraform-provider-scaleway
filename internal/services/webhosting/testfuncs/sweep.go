package webhostingtestfuncs

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	webhostingSDK "github.com/scaleway/scaleway-sdk-go/api/webhosting/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/logging"
)

func AddTestSweepers() {
	resource.AddTestSweepers("scaleway_webhosting", &resource.Sweeper{
		Name: "scaleway_webhosting",
		F:    testSweepWebhosting,
	})
}

func testSweepWebhosting(_ string) error {
	return acctest.SweepRegions(scw.AllRegions, func(scwClient *scw.Client, region scw.Region) error {
		webhsotingAPI := webhostingSDK.NewAPI(scwClient)

		logging.L.Debugf("sweeper: deleting the hostings in (%s)", region)

		listHostings, err := webhsotingAPI.ListHostings(&webhostingSDK.ListHostingsRequest{Region: region}, scw.WithAllPages())
		if err != nil {
			return fmt.Errorf("error listing hostings in (%s) in sweeper: %s", region, err)
		}

		for _, hosting := range listHostings.Hostings {
			_, err := webhsotingAPI.DeleteHosting(&webhostingSDK.DeleteHostingRequest{
				HostingID: hosting.ID,
				Region:    region,
			})
			if err != nil {
				logging.L.Debugf("sweeper: error (%s)", err)

				return fmt.Errorf("error deleting hosting in sweeper: %s", err)
			}
		}

		return nil
	})
}
