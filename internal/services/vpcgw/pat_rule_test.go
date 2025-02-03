package vpcgw_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	vpcgwSDK "github.com/scaleway/scaleway-sdk-go/api/vpcgw/v1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/vpcgw"
)

func TestAccVPCPublicGatewayPATRule_Basic(t *testing.T) {
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
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCPublicGatewayPATRuleExists(
						tt,
						"scaleway_vpc_public_gateway_pat_rule.main",
					),
					resource.TestCheckResourceAttrSet("scaleway_vpc_public_gateway_pat_rule.main", "gateway_id"),
					resource.TestCheckResourceAttrSet("scaleway_vpc_public_gateway_pat_rule.main", "private_ip"),
					resource.TestCheckResourceAttrSet("scaleway_vpc_public_gateway_pat_rule.main", "created_at"),
					resource.TestCheckResourceAttrSet("scaleway_vpc_public_gateway_pat_rule.main", "updated_at"),
					resource.TestCheckResourceAttrSet("scaleway_vpc_public_gateway_pat_rule.main", "protocol"),
					resource.TestCheckResourceAttr("scaleway_vpc_public_gateway_pat_rule.main", "protocol", "both"),
					resource.TestCheckResourceAttr("scaleway_vpc_public_gateway_pat_rule.main", "public_port", "42"),
					resource.TestCheckResourceAttr("scaleway_vpc_public_gateway_pat_rule.main", "private_port", "42"),
				),
			},
		},
	})
}

func TestAccVPCPublicGatewayPATRule_WithInstance(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckVPCPublicGatewayPATRuleDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					### Scaleway Private Network
					resource "scaleway_vpc_private_network" "main" {
					  name = "My Private Network"
					}
					
					### IP for Public Gateway
					resource "scaleway_vpc_public_gateway_ip" "main" {
					}
					
					### DHCP Space of VPC
					resource "scaleway_vpc_public_gateway_dhcp" "main" {
					  subnet = "10.0.0.0/24"
					}
					
					### The Public Gateway with the Attached IP
					resource "scaleway_vpc_public_gateway" "main" {
					  name  = "The Public Gateway"
					  type  = "VPC-GW-S"
					  ip_id = scaleway_vpc_public_gateway_ip.main.id
					}
					
					### VPC Gateway Network
					resource "scaleway_vpc_gateway_network" "main" {
					  gateway_id         = scaleway_vpc_public_gateway.main.id
					  private_network_id = scaleway_vpc_private_network.main.id
					  dhcp_id            = scaleway_vpc_public_gateway_dhcp.main.id
					  cleanup_dhcp       = true
					  enable_masquerade  = true
					  depends_on = [
						scaleway_vpc_public_gateway_ip.main,
						scaleway_vpc_private_network.main
					  ]
					}
					
					### Scaleway Instance
					resource "scaleway_instance_server" "main" {
					  name        = "Scaleway Instance"
					  type        = "DEV1-S"
					  image       = "debian_bullseye"
					  enable_ipv6 = false
					
					  private_network {
						pn_id = scaleway_vpc_private_network.main.id
					  }
					}
					
					### DHCP Reservation for Instance
					resource "scaleway_vpc_public_gateway_dhcp_reservation" "main" {
					  gateway_network_id = scaleway_vpc_gateway_network.main.id
					  mac_address        = scaleway_instance_server.main.private_network.0.mac_address
					  ip_address         = "10.0.0.3"
					}
					
					### PAT RULE for SSH Port
					resource "scaleway_vpc_public_gateway_pat_rule" "main" {
					  gateway_id   = scaleway_vpc_public_gateway.main.id
					  private_ip   = scaleway_vpc_public_gateway_dhcp_reservation.main.ip_address
					  private_port = 22
					  public_port  = 2023
					  protocol     = "tcp"
					  depends_on = [
						scaleway_vpc_gateway_network.main,
						scaleway_vpc_private_network.main
					  ]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCPublicGatewayDHCPExists(tt, "scaleway_vpc_public_gateway_dhcp.main"),
					testAccCheckVPCPublicGatewayPATRuleExists(tt, "scaleway_vpc_public_gateway_pat_rule.main"),
					resource.TestCheckResourceAttr("scaleway_vpc_public_gateway_dhcp.main", "subnet", "10.0.0.0/24"),
					resource.TestCheckResourceAttr("scaleway_vpc_public_gateway_pat_rule.main", "protocol", "tcp"),
					resource.TestCheckResourceAttr("scaleway_vpc_public_gateway_pat_rule.main", "public_port", "2023"),
					resource.TestCheckResourceAttr("scaleway_vpc_public_gateway_pat_rule.main", "private_port", "22"),
					resource.TestCheckResourceAttrPair(
						"scaleway_vpc_public_gateway_pat_rule.main", "private_ip",
						"scaleway_vpc_public_gateway_dhcp_reservation.main", "ip_address"),
					resource.TestCheckResourceAttrPair(
						"scaleway_vpc_public_gateway_pat_rule.main", "gateway_id",
						"scaleway_vpc_public_gateway.main", "id"),
				),
			},
		},
	})
}

func testAccCheckVPCPublicGatewayPATRuleExists(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		api, zone, ID, err := vpcgw.NewAPIWithZoneAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = api.GetPATRule(&vpcgwSDK.GetPATRuleRequest{
			PatRuleID: ID,
			Zone:      zone,
		})
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckVPCPublicGatewayPATRuleDestroy(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_vpc_public_gateway_pat_rules" {
				continue
			}

			api, zone, ID, err := vpcgw.NewAPIWithZoneAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}

			_, err = api.GetPATRule(&vpcgwSDK.GetPATRuleRequest{
				PatRuleID: ID,
				Zone:      zone,
			})

			if err == nil {
				return fmt.Errorf(
					"VPC public gateway pat rule %s still exists",
					rs.Primary.ID,
				)
			}

			// Unexpected api error we return it
			if !httperrors.Is404(err) {
				return err
			}
		}

		return nil
	}
}
