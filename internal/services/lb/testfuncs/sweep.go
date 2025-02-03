package lbtestfuncs

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	lbSDK "github.com/scaleway/scaleway-sdk-go/api/lb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/logging"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/instance"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/lb"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/transport"
)

func AddTestSweepers() {
	resource.AddTestSweepers("scaleway_lb_ip", &resource.Sweeper{
		Name: "scaleway_lb_ip",
		F:    testSweepIP,
	})
	resource.AddTestSweepers("scaleway_lb", &resource.Sweeper{
		Name: "scaleway_lb",
		F:    testSweepLB,
	})
}

func testSweepLB(_ string) error {
	return acctest.SweepZones([]scw.Zone{scw.ZoneFrPar1, scw.ZoneNlAms1, scw.ZonePlWaw1}, func(scwClient *scw.Client, zone scw.Zone) error {
		lbAPI := lbSDK.NewZonedAPI(scwClient)

		logging.L.Debugf("sweeper: destroying the lbs in (%s)", zone)
		listLBs, err := lbAPI.ListLBs(&lbSDK.ZonedAPIListLBsRequest{
			Zone: zone,
		}, scw.WithAllPages())
		if err != nil {
			return fmt.Errorf("error listing lbs in (%s) in sweeper: %s", zone, err)
		}

		for _, l := range listLBs.LBs {
			retryInterval := lb.DefaultWaitLBRetryInterval

			if transport.DefaultWaitRetryInterval != nil {
				retryInterval = *transport.DefaultWaitRetryInterval
			}

			_, err := lbAPI.WaitForLbInstances(&lbSDK.ZonedAPIWaitForLBInstancesRequest{
				Zone:          zone,
				LBID:          l.ID,
				Timeout:       scw.TimeDurationPtr(instance.DefaultInstanceServerWaitTimeout),
				RetryInterval: &retryInterval,
			})
			if err != nil {
				return fmt.Errorf("error waiting for lb in sweeper: %s", err)
			}
			err = lbAPI.DeleteLB(&lbSDK.ZonedAPIDeleteLBRequest{
				LBID:      l.ID,
				ReleaseIP: true,
				Zone:      zone,
			})
			if err != nil {
				return fmt.Errorf("error deleting lb in sweeper: %s", err)
			}
		}

		return nil
	})
}

func testSweepIP(_ string) error {
	return acctest.SweepZones([]scw.Zone{scw.ZoneFrPar1, scw.ZoneNlAms1, scw.ZonePlWaw1}, func(scwClient *scw.Client, zone scw.Zone) error {
		lbAPI := lbSDK.NewZonedAPI(scwClient)

		logging.L.Debugf("sweeper: destroying the lb ips in zone (%s)", zone)
		listIPs, err := lbAPI.ListIPs(&lbSDK.ZonedAPIListIPsRequest{Zone: zone}, scw.WithAllPages())
		if err != nil {
			return fmt.Errorf("error listing lb ips in (%s) in sweeper: %s", zone, err)
		}

		for _, ip := range listIPs.IPs {
			if ip.LBID == nil {
				err := lbAPI.ReleaseIP(&lbSDK.ZonedAPIReleaseIPRequest{
					Zone: zone,
					IPID: ip.ID,
				})
				if err != nil {
					return fmt.Errorf("error deleting lb ip in sweeper: %s", err)
				}
			}
		}

		return nil
	})
}
