package ipamtestfuncs

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	ipamSDK "github.com/scaleway/scaleway-sdk-go/api/ipam/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/logging"
)

func AddTestSweepers() {
	resource.AddTestSweepers("scaleway_ipam_ip", &resource.Sweeper{
		Name: "scaleway_ipam_ip",
		F:    testSweepIPAMIP,
	})
}

func testSweepIPAMIP(_ string) error {
	return acctest.SweepRegions(scw.AllRegions, func(scwClient *scw.Client, region scw.Region) error {
		ipamAPI := ipamSDK.NewAPI(scwClient)

		logging.L.Debugf("sweeper: deleting the IPs in (%s)", region)

		listIPs, err := ipamAPI.ListIPs(&ipamSDK.ListIPsRequest{Region: region}, scw.WithAllPages())
		if err != nil {
			return fmt.Errorf("error listing ips in (%s) in sweeper: %s", region, err)
		}

		for _, v := range listIPs.IPs {
			err := ipamAPI.ReleaseIP(&ipamSDK.ReleaseIPRequest{
				IPID:   v.ID,
				Region: region,
			})
			if err != nil {
				logging.L.Debugf("sweeper: error (%s)", err)

				return fmt.Errorf("error releasing IP in sweeper: %s", err)
			}
		}

		return nil
	})
}
