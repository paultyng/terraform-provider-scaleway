package applesilicontestfuncs

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	applesiliconSDK "github.com/scaleway/scaleway-sdk-go/api/applesilicon/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/logging"
)

func AddTestSweepers() {
	resource.AddTestSweepers("scaleway_apple_silicon_instance", &resource.Sweeper{
		Name: "scaleway_apple_silicon",
		F:    testSweepAppleSiliconServer,
	})
}

func testSweepAppleSiliconServer(_ string) error {
	return acctest.SweepZones([]scw.Zone{scw.ZoneFrPar1}, func(scwClient *scw.Client, zone scw.Zone) error {
		asAPI := applesiliconSDK.NewAPI(scwClient)
		logging.L.Debugf("sweeper: destroying the apple silicon instance in (%s)", zone)
		listServers, err := asAPI.ListServers(&applesiliconSDK.ListServersRequest{Zone: zone}, scw.WithAllPages())
		if err != nil {
			return fmt.Errorf("error listing apple silicon servers in (%s) in sweeper: %s", zone, err)
		}

		for _, server := range listServers.Servers {
			errDelete := asAPI.DeleteServer(&applesiliconSDK.DeleteServerRequest{
				ServerID: server.ID,
				Zone:     zone,
			})
			if errDelete != nil {
				return fmt.Errorf("error deleting apple silicon server in sweeper: %s", err)
			}
		}

		return nil
	})
}
