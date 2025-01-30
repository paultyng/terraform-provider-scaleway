package baremetaltestfuncs

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	baremetalSDK "github.com/scaleway/scaleway-sdk-go/api/baremetal/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/logging"
)

func AddTestSweepers() {
	resource.AddTestSweepers("scaleway_baremetal_server", &resource.Sweeper{
		Name: "scaleway_baremetal_server",
		F:    testSweepServer,
	})
}

func testSweepServer(_ string) error {
	return acctest.SweepZones([]scw.Zone{scw.ZoneFrPar2}, func(scwClient *scw.Client, zone scw.Zone) error {
		baremetalAPI := baremetalSDK.NewAPI(scwClient)
		logging.L.Debugf("sweeper: destroying the baremetal server in (%s)", zone)
		listServers, err := baremetalAPI.ListServers(&baremetalSDK.ListServersRequest{Zone: zone}, scw.WithAllPages())
		if err != nil {
			logging.L.Warningf("error listing servers in (%s) in sweeper: %s", zone, err)
			return nil
		}

		for _, server := range listServers.Servers {
			_, err := baremetalAPI.DeleteServer(&baremetalSDK.DeleteServerRequest{
				Zone:     zone,
				ServerID: server.ID,
			})
			if err != nil {
				return fmt.Errorf("error deleting server in sweeper: %s", err)
			}
		}

		return nil
	})
}
