package vpcgw_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	vpcgwchecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/vpcgw/testfuncs"
)

func TestAccDataSourceVPCPublicGatewayIP_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      vpcgwchecks.IsIPDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_vpc_public_gateway_ip" "main" {
					}`,
			},
			{
				Config: `
					resource "scaleway_vpc_public_gateway_ip" "main" {
					}

					data "scaleway_vpc_public_gateway_ip" "ip_by_id" {
						ip_id = "${scaleway_vpc_public_gateway_ip.main.id}"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCPublicGatewayIPExists(tt, "scaleway_vpc_public_gateway_ip.main"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_vpc_public_gateway_ip.ip_by_id", "ip_id",
						"scaleway_vpc_public_gateway_ip.main", "id"),
				),
			},
		},
	})
}
