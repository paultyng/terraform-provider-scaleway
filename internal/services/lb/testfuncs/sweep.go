package lbtestfuncs

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	lbSDK "github.com/scaleway/scaleway-sdk-go/api/lb/v1"
	"github.com/scaleway/scaleway-sdk-go/api/lb/v1/sweepers"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
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
	return acctest.SweepZones((&lbSDK.ZonedAPI{}).Zones(), sweepers.SweepLB)
}

func testSweepIP(_ string) error {
	return acctest.SweepZones((&lbSDK.ZonedAPI{}).Zones(), sweepers.SweepIP)
}
