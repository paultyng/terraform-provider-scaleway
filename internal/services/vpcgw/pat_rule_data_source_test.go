package vpcgw_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

func TestAccDataSourceVPCPublicGatewayPATRule_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckVPCPublicGatewayPATRuleDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_vpc_public_gateway pg01 {
						type = "VPC-GW-S"
					}

					resource scaleway_vpc_public_gateway_dhcp dhcp01 {
						subnet = "192.168.1.0/24"
					}

					resource scaleway_vpc_private_network pn01 {
						name = "pn_test_network"
					}
				`,
			},
			{
				Config: `
					resource scaleway_vpc_public_gateway pg01 {
						type = "VPC-GW-S"
					}

					resource scaleway_vpc_public_gateway_dhcp dhcp01 {
						subnet = "192.168.1.0/24"
					}

					resource scaleway_vpc_private_network pn01 {
						name = "pn_test_network"
					}

					resource scaleway_vpc_gateway_network gn01 {
					    gateway_id = scaleway_vpc_public_gateway.pg01.id
					    private_network_id = scaleway_vpc_private_network.pn01.id
					    dhcp_id = scaleway_vpc_public_gateway_dhcp.dhcp01.id
						depends_on = [scaleway_vpc_private_network.pn01]
						cleanup_dhcp = true
						enable_masquerade = true
					}

					resource scaleway_vpc_public_gateway_pat_rule main {
						gateway_id = scaleway_vpc_public_gateway.pg01.id
						private_ip = scaleway_vpc_public_gateway_dhcp.dhcp01.address
						private_port = 42
						public_port = 42
						protocol = "both"
						depends_on = [scaleway_vpc_gateway_network.gn01, scaleway_vpc_private_network.pn01]
					}

					data "scaleway_vpc_public_gateway_pat_rule" "main" {
						pat_rule_id = "${scaleway_vpc_public_gateway_pat_rule.main.id}"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCPublicGatewayPATRuleExists(
						tt,
						"scaleway_vpc_public_gateway_pat_rule.main",
					),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_vpc_public_gateway_pat_rule.main", "gateway_id",
						"scaleway_vpc_public_gateway_pat_rule.main", "gateway_id"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_vpc_public_gateway_pat_rule.main", "private_ip",
						"scaleway_vpc_public_gateway_pat_rule.main", "private_ip"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_vpc_public_gateway_pat_rule.main", "created_at",
						"scaleway_vpc_public_gateway_pat_rule.main", "created_at"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_vpc_public_gateway_pat_rule.main", "updated_at",
						"scaleway_vpc_public_gateway_pat_rule.main", "updated_at"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_vpc_public_gateway_pat_rule.main", "protocol",
						"scaleway_vpc_public_gateway_pat_rule.main", "protocol"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_vpc_public_gateway_pat_rule.main", "public_port",
						"scaleway_vpc_public_gateway_pat_rule.main", "public_port"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_vpc_public_gateway_pat_rule.main", "private_port",
						"scaleway_vpc_public_gateway_pat_rule.main", "private_port"),
				),
			},
		},
	})
}
