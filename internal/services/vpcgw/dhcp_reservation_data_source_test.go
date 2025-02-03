package vpcgw_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	vpcgwchecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/vpcgw/testfuncs"
)

func TestAccDataSourceVPCPublicGatewayDHCPReservation_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	pnName := "TestAccScalewayDataSourceVPCPublicGatewayDHCPReservation_Basic"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      vpcgwchecks.IsDHCPDestroyed(tt),

		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource scaleway_vpc_private_network main {
						name = "%s"
					}
	
					resource "scaleway_instance_server" "main" {
						image = "ubuntu_focal"
						type  = "DEV1-S"
						zone = "fr-par-1"
	
						private_network {
							pn_id = scaleway_vpc_private_network.main.id
						}
					}
	
					resource scaleway_vpc_public_gateway_ip main {
					}
	
					resource scaleway_vpc_public_gateway_dhcp main {
						subnet = "192.168.1.0/24"
					}
	
					resource scaleway_vpc_public_gateway main {
						name = "foobar"
						type = "VPC-GW-S"
						ip_id = scaleway_vpc_public_gateway_ip.main.id
					}
	
					resource scaleway_vpc_gateway_network main {
						gateway_id = scaleway_vpc_public_gateway.main.id
						private_network_id = scaleway_vpc_private_network.main.id
						dhcp_id = scaleway_vpc_public_gateway_dhcp.main.id
						cleanup_dhcp = true
						enable_masquerade = true
						depends_on = [scaleway_vpc_public_gateway_ip.main, scaleway_vpc_private_network.main]
					}
				`, pnName),
			},
			{
				Config: fmt.Sprintf(`
					resource scaleway_vpc_private_network main {
						name = "%s"
					}
	
					resource "scaleway_instance_server" "main" {
						image = "ubuntu_focal"
						type  = "DEV1-S"
						zone = "fr-par-1"
	
						private_network {
							pn_id = scaleway_vpc_private_network.main.id
						}
					}
	
					resource scaleway_vpc_public_gateway_ip main {
					}
	
					resource scaleway_vpc_public_gateway_dhcp main {
						subnet = "192.168.1.0/24"
					}
	
					resource scaleway_vpc_public_gateway main {
						name = "foobar"
						type = "VPC-GW-S"
						ip_id = scaleway_vpc_public_gateway_ip.main.id
					}
	
					resource scaleway_vpc_gateway_network main {
						gateway_id = scaleway_vpc_public_gateway.main.id
						private_network_id = scaleway_vpc_private_network.main.id
						dhcp_id = scaleway_vpc_public_gateway_dhcp.main.id
						cleanup_dhcp = true
						enable_masquerade = true
						depends_on = [scaleway_vpc_public_gateway_ip.main, scaleway_vpc_private_network.main]
					}
	
					data "scaleway_vpc_public_gateway_dhcp_reservation" "by_mac_address_and_gw_network" {
						mac_address = "${scaleway_instance_server.main.private_network.0.mac_address}"
					    gateway_network_id = scaleway_vpc_gateway_network.main.id
						wait_for_dhcp = true
						depends_on = [scaleway_vpc_gateway_network.main, scaleway_vpc_public_gateway_dhcp.main, scaleway_vpc_private_network.main]
					}
				`, pnName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.scaleway_vpc_public_gateway_dhcp_reservation.by_mac_address_and_gw_network", "mac_address",
						"scaleway_instance_server.main", "private_network.0.mac_address"),
				),
			},
		},
	})
}

func TestAccDataSourceVPCPublicGatewayDHCPReservation_Static(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	pnName := "TestAccScalewayDataSourceVPCPublicGatewayDHCPReservation_Static"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      vpcgwchecks.IsDHCPDestroyed(tt),

		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource scaleway_vpc_private_network main {
						name = "%s"
					}

					resource "scaleway_instance_security_group" main {
						inbound_default_policy  = "drop"
						outbound_default_policy = "accept"
						
						inbound_rule {
							action = "accept"
							port   = "22"
						}
					}

					resource "scaleway_instance_server" "main" {
						image = "ubuntu_focal"
						type  = "DEV1-S"
						zone = "fr-par-1"

						security_group_id = scaleway_instance_security_group.main.id

						private_network {
							pn_id = scaleway_vpc_private_network.main.id
						}
					}
	
					resource scaleway_vpc_public_gateway_ip main {
					}
	
					resource scaleway_vpc_public_gateway_dhcp main {
						subnet = "192.168.1.0/24"
					}
	
					resource scaleway_vpc_public_gateway main {
						name = "foobar"
						type = "VPC-GW-S"
						ip_id = scaleway_vpc_public_gateway_ip.main.id
					}
	
					resource scaleway_vpc_gateway_network main {
						gateway_id = scaleway_vpc_public_gateway.main.id
						private_network_id = scaleway_vpc_private_network.main.id
						dhcp_id = scaleway_vpc_public_gateway_dhcp.main.id
						cleanup_dhcp = true
						enable_masquerade = true
						depends_on = [scaleway_vpc_public_gateway_ip.main, scaleway_vpc_private_network.main]
					}

					resource scaleway_vpc_public_gateway_dhcp_reservation main {
						gateway_network_id = scaleway_vpc_gateway_network.main.id
						mac_address = scaleway_instance_server.main.private_network.0.mac_address
						ip_address = "192.168.1.4"
					}

					### VPC PAT RULE
					resource "scaleway_vpc_public_gateway_pat_rule" "main" {
						gateway_id   = scaleway_vpc_public_gateway.main.id
						private_ip   = scaleway_vpc_public_gateway_dhcp_reservation.main.ip_address
						private_port = 22
						public_port  = 2222
						protocol     = "tcp"
					}

					data "scaleway_vpc_public_gateway_dhcp_reservation" "by_id" {
						reservation_id = "${scaleway_vpc_public_gateway_dhcp_reservation.main.id}"
					}
				`, pnName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.scaleway_vpc_public_gateway_dhcp_reservation.by_id", "mac_address",
						"scaleway_instance_server.main", "private_network.0.mac_address"),
				),
			},
		},
	})
}
