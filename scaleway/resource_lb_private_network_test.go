package scaleway

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/scaleway/scaleway-sdk-go/api/lb/v1"
)

func TestAccScalewayLbPrivateNetwork_Basic(t *testing.T) {
	tt := NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckScalewayLbPrivateNetworkDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_lb_ip ip01 {}
					resource scaleway_lb lb01 {
						ip_id = scaleway_lb_ip.ip01.id
						name = "test-lb"
						type = "lb-s"
					}
					resource scaleway_vpc_private_network pn01 {
						name = "test-pn"
					}
					resource scaleway_lb_private_network lb01pn01 {
						lb_id = scaleway_lb.lb01.id
						private_network_id = scaleway_vpc_private_network.pn01.id
						static_config = ["172.16.0.100", "172.16.0.101"]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayLbPrivateNetworkExists(tt, "scaleway_lb_private_network.lb01pn01"),
					resource.TestCheckResourceAttr(
						"scaleway_vpc_private_network.lb01pn01",
						"static_config",
						"[\"172.16.0.100\", \"172.16.0.101\"]",
					),
					resource.TestCheckResourceAttr("scaleway_vpc_private_network.lb01pn01", "dhcp_config", ""),
				),
			},
			{
				Config: `
					resource scaleway_lb_ip ip01 {}
					resource scaleway_lb lb01 {
						ip_id = scaleway_lb_ip.ip01.id
						name = "test-lb"
						type = "lb-s"
					}
					resource scaleway_vpc_private_network pn01 {
						name = "test-pn"
					}
					resource scaleway_lb_private_network lb01pn01 {
						lb_id = scaleway_lb.lb01.id
						private_network_id = scaleway_vpc_private_network.pn01.id
						dhcp_config = true
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayLbPrivateNetworkExists(tt, "scaleway_lb_private_network.lb01pn01"),
					resource.TestCheckResourceAttr("scaleway_vpc_private_network.lb01pn01", "static_config", ""),
					resource.TestCheckResourceAttr("scaleway_vpc_private_network.lb01pn01", "dhcp_config", "true"),
				),
			},
		},
	})
}

func getLbPrivateNetwork(tt *TestTools, rs *terraform.ResourceState) (*lb.PrivateNetwork, error) {
	lbID := rs.Primary.Attributes["load_balancer_id"]
	pnID := rs.Primary.Attributes["private_network_id"]

	lbAPI, zone, pnID, err := lbAPIWithZoneAndID(tt.Meta, pnID)
	if err != nil {
		return nil, err
	}

	_, lbID, err = parseZonedID(lbID)
	if err != nil {
		return nil, fmt.Errorf("invalid resource: %s", err)
	}

	listPN, err := lbAPI.ListLBPrivateNetworks(&lb.ZonedAPIListLBPrivateNetworksRequest{
		LBID: lbID,
		Zone: zone,
	})
	if err != nil {
		return nil, err
	}

	for _, pn := range listPN.PrivateNetwork {
		if pn.PrivateNetworkID == pnID {
			return pn, nil
		}
	}

	return nil, fmt.Errorf("private network %s not found", pnID)
}

func testAccCheckScalewayLbPrivateNetworkExists(tt *TestTools, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		pn, err := getLbPrivateNetwork(tt, rs)
		if err != nil {
			return err
		}
		if pn == nil {
			return fmt.Errorf("resource not found: %s", rs.Primary.Attributes["private_network_id"])
		}

		return nil
	}
}

func testAccCheckScalewayLbPrivateNetworkDestroy(tt *TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_vpc_private_network" {
				continue
			}

			pn, err := getLbPrivateNetwork(tt, rs)
			if err != nil {
				return err
			}

			if pn != nil {
				fmt.Errorf("LB PN (%s) still exists", rs.Primary.Attributes["private_network_id"])
			}
		}

		return nil
	}
}
