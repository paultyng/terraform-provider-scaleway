package vpcgwtestfuncs

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	vpcgwSDK "github.com/scaleway/scaleway-sdk-go/api/vpcgw/v1"
	"github.com/scaleway/scaleway-sdk-go/api/vpcgw/v1/sweepers"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

func AddTestSweepers() {
	resource.AddTestSweepers("scaleway_vpc_public_gateway_dhcp", &resource.Sweeper{
		Name: "scaleway_vpc_public_gateway_dhcp",
		F:    testSweepVPCPublicGatewayDHCP,
	})
	resource.AddTestSweepers("scaleway_vpc_public_gateway_ip", &resource.Sweeper{
		Name: "scaleway_vpc_public_gateway_ip",
		F:    testSweepVPCPublicGatewayIP,
	})
	resource.AddTestSweepers("scaleway_gateway_network", &resource.Sweeper{
		Name: "scaleway_gateway_network",
		F:    testSweepVPCGatewayNetwork,
	})
	resource.AddTestSweepers("scaleway_vpc_public_gateway", &resource.Sweeper{
		Name: "scaleway_vpc_public_gateway",
		F:    testSweepVPCPublicGateway,
	})
}

func testSweepVPCPublicGateway(_ string) error {
	return acctest.SweepZones((&vpcgwSDK.API{}).Zones(), sweepers.SweepVPCPublicGateway)
}

func testSweepVPCGatewayNetwork(_ string) error {
	return acctest.SweepZones((&vpcgwSDK.API{}).Zones(), sweepers.SweepGatewayNetworks)
}

func testSweepVPCPublicGatewayIP(_ string) error {
	return acctest.SweepZones((&vpcgwSDK.API{}).Zones(), sweepers.SweepVPCPublicGatewayIP)
}

func testSweepVPCPublicGatewayDHCP(_ string) error {
	return acctest.SweepZones((&vpcgwSDK.API{}).Zones(), sweepers.SweepVPCPublicGatewayDHCP)
}
