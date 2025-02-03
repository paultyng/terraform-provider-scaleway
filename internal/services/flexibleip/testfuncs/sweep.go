package flexibleiptestfuncs

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	flexibleipSDK "github.com/scaleway/scaleway-sdk-go/api/flexibleip/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/api/flexibleip/v1alpha1/sweepers"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

func AddTestSweepers() {
	resource.AddTestSweepers("scaleway_flexible_ip", &resource.Sweeper{
		Name: "scaleway_flexible_ip",
		F:    testSweepFlexibleIP,
	})
}

func testSweepFlexibleIP(_ string) error {
	return acctest.SweepZones((&flexibleipSDK.API{}).Zones(), sweepers.SweepFlexibleIP)
}
