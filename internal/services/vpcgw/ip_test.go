package vpcgw_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	vpcgwSDK "github.com/scaleway/scaleway-sdk-go/api/vpcgw/v1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/vpcgw"
	vpcgwchecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/vpcgw/testfuncs"
)

func TestAccVPCPublicGatewayIP_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      vpcgwchecks.IsIPDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_vpc_public_gateway_ip main {}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCPublicGatewayIPExists(tt, "scaleway_vpc_public_gateway_ip.main"),
					resource.TestCheckResourceAttrSet("scaleway_vpc_public_gateway_ip.main", "address"),
					resource.TestCheckResourceAttrSet("scaleway_vpc_public_gateway_ip.main", "created_at"),
					resource.TestCheckResourceAttrSet("scaleway_vpc_public_gateway_ip.main", "updated_at"),
				),
			},
		},
	})
}

func TestAccVPCPublicGatewayIP_WithZone(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      vpcgwchecks.IsIPDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_vpc_public_gateway_ip main {}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCPublicGatewayIPExists(tt, "scaleway_vpc_public_gateway_ip.main"),
					resource.TestCheckResourceAttr("scaleway_vpc_public_gateway_ip.main", "zone", "fr-par-1"),
				),
			},
			{
				Config: `
					resource scaleway_vpc_public_gateway_ip main {
						zone = "nl-ams-1"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCPublicGatewayIPExists(tt, "scaleway_vpc_public_gateway_ip.main"),
					resource.TestCheckResourceAttr("scaleway_vpc_public_gateway_ip.main", "zone", "nl-ams-1"),
				),
			},
		},
	})
}

func TestAccVPCPublicGatewayIP_WithTags(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      vpcgwchecks.IsIPDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_vpc_public_gateway_ip main {}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCPublicGatewayIPExists(tt, "scaleway_vpc_public_gateway_ip.main"),
					resource.TestCheckNoResourceAttr("scaleway_vpc_public_gateway_ip.main", "tags"),
				),
			},
			{
				Config: `
					resource scaleway_vpc_public_gateway_ip main {
						tags = ["foo", "bar"]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCPublicGatewayIPExists(tt, "scaleway_vpc_public_gateway_ip.main"),
					resource.TestCheckResourceAttr("scaleway_vpc_public_gateway_ip.main", "tags.0", "foo"),
					resource.TestCheckResourceAttr("scaleway_vpc_public_gateway_ip.main", "tags.1", "bar"),
				),
			},
		},
	})
}

func testAccCheckVPCPublicGatewayIPExists(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		api, zone, ID, err := vpcgw.NewAPIWithZoneAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = api.GetIP(&vpcgwSDK.GetIPRequest{
			IPID: ID,
			Zone: zone,
		})
		if err != nil {
			return err
		}

		return nil
	}
}
