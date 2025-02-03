package flexibleiptestfuncs

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	flexibleipSDK "github.com/scaleway/scaleway-sdk-go/api/flexibleip/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/logging"
)

func AddTestSweepers() {
	resource.AddTestSweepers("scaleway_flexible_ip", &resource.Sweeper{
		Name: "scaleway_flexible_ip",
		F:    testSweepFlexibleIP,
	})
}

func testSweepFlexibleIP(_ string) error {
	return acctest.SweepZones(scw.AllZones, func(scwClient *scw.Client, zone scw.Zone) error {
		fipAPI := flexibleipSDK.NewAPI(scwClient)

		listIPs, err := fipAPI.ListFlexibleIPs(&flexibleipSDK.ListFlexibleIPsRequest{Zone: zone}, scw.WithAllPages())
		if err != nil {
			logging.L.Warningf("error listing ips in (%s) in sweeper: %s", zone, err)
			return nil
		}

		for _, ip := range listIPs.FlexibleIPs {
			err := fipAPI.DeleteFlexibleIP(&flexibleipSDK.DeleteFlexibleIPRequest{
				FipID: ip.ID,
				Zone:  zone,
			})
			if err != nil {
				return fmt.Errorf("error deleting ip in sweeper: %s", err)
			}
		}

		return nil
	})
}
