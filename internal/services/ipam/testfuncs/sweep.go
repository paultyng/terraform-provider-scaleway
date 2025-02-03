package ipamtestfuncs

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	ipamSDK "github.com/scaleway/scaleway-sdk-go/api/ipam/v1"
	"github.com/scaleway/scaleway-sdk-go/api/ipam/v1/sweepers"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

func AddTestSweepers() {
	resource.AddTestSweepers("scaleway_ipam_ip", &resource.Sweeper{
		Name: "scaleway_ipam_ip",
		F:    testSweepIPAMIP,
	})
}

func testSweepIPAMIP(_ string) error {
	return acctest.SweepRegions((&ipamSDK.API{}).Regions(), sweepers.SweepIP)
}
