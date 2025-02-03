package vpctestfuncs

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	vpcSDK "github.com/scaleway/scaleway-sdk-go/api/vpc/v2"
	"github.com/scaleway/scaleway-sdk-go/api/vpc/v2/sweepers"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

func AddTestSweepers() {
	resource.AddTestSweepers("scaleway_vpc", &resource.Sweeper{
		Name:         "scaleway_vpc",
		F:            testSweepVPC,
		Dependencies: []string{"scaleway_vpc_private_network"},
	})

	resource.AddTestSweepers("scaleway_vpc_private_network", &resource.Sweeper{
		Name:         "scaleway_vpc_private_network",
		F:            testSweepVPCPrivateNetwork,
		Dependencies: []string{"scaleway_ipam_ip"},
	})

	resource.AddTestSweepers("scaleway_vpc_route", &resource.Sweeper{
		Name: "scaleway_vpc_route",
		F:    testSweepVPCRoute,
	})
}

func testSweepVPC(_ string) error {
	return acctest.SweepRegions((&vpcSDK.API{}).Regions(), sweepers.SweepVPC)
}

func testSweepVPCPrivateNetwork(_ string) error {
	return acctest.SweepRegions((&vpcSDK.API{}).Regions(), sweepers.SweepPrivateNetwork)
}

func testSweepVPCRoute(_ string) error {
	return acctest.SweepRegions((&vpcSDK.API{}).Regions(), sweepers.SweepRoute)
}
