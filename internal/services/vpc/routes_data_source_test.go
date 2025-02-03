package vpc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	vpcchecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/vpc/testfuncs"
)

func TestAccDataSourceRoutes_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      vpcchecks.CheckPrivateNetworkDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_vpc" "vpc01" {
					  name           = "tf-vpc-route-01"
					  enable_routing = true
					}
					
					resource "scaleway_vpc_private_network" "pn01" {
					  name   = "tf-pn_route"
					  vpc_id = scaleway_vpc.vpc01.id
					}
					
					resource "scaleway_vpc_private_network" "pn02" {
					  name   = "tf-pn_route_2"
					  vpc_id = scaleway_vpc.vpc01.id
					}
				`,
			},
			{
				Config: `
					resource "scaleway_vpc" "vpc01" {
					  name           = "tf-vpc-route-01"
					  enable_routing = true
					}
					
					resource "scaleway_vpc_private_network" "pn01" {
					  name   = "tf-pn_route"
					  vpc_id = scaleway_vpc.vpc01.id
					}
					
					resource "scaleway_vpc_private_network" "pn02" {
					  name   = "tf-pn_route_2"
					  vpc_id = scaleway_vpc.vpc01.id
					}
					
					resource "scaleway_vpc_public_gateway" "pg01" {
					  name = "tf-gw-route"
					  type = "VPC-GW-S"
					}
					
					resource "scaleway_vpc_gateway_network" "gn01" {
					  gateway_id         = scaleway_vpc_public_gateway.pg01.id
					  private_network_id = scaleway_vpc_private_network.pn01.id
					  enable_masquerade  = true
					  ipam_config {
						push_default_route = true
					  }
					}
					
					resource "scaleway_vpc_gateway_network" "gn02" {
					  gateway_id         = scaleway_vpc_public_gateway.pg01.id
					  private_network_id = scaleway_vpc_private_network.pn02.id
					  enable_masquerade  = true
					  ipam_config {
						push_default_route = true
					  }
					}
					`,
			},
			{
				Config: `
					resource "scaleway_vpc" "vpc01" {
					  name           = "tf-vpc-route-01"
					  enable_routing = true
					}
					
					resource "scaleway_vpc_private_network" "pn01" {
					  name   = "tf-pn_route"
					  vpc_id = scaleway_vpc.vpc01.id
					}
					
					resource "scaleway_vpc_private_network" "pn02" {
					  name   = "tf-pn_route_2"
					  vpc_id = scaleway_vpc.vpc01.id
					}
					
					resource "scaleway_vpc_public_gateway" "pg01" {
					  name = "tf-gw-route"
					  type = "VPC-GW-S"
					}
					
					resource "scaleway_vpc_gateway_network" "gn01" {
					  gateway_id         = scaleway_vpc_public_gateway.pg01.id
					  private_network_id = scaleway_vpc_private_network.pn01.id
					  enable_masquerade  = true
					  ipam_config {
						push_default_route = true
					  }
					}
					
					resource "scaleway_vpc_gateway_network" "gn02" {
					  gateway_id         = scaleway_vpc_public_gateway.pg01.id
					  private_network_id = scaleway_vpc_private_network.pn02.id
					  enable_masquerade  = true
					  ipam_config {
						push_default_route = true
					  }
					}
					
					data "scaleway_vpc_routes" "routes_by_vpc_id" {
					  vpc_id = scaleway_vpc.vpc01.id
					}
					
					data "scaleway_vpc_routes" "routes_by_ipv6" {
					  vpc_id  = scaleway_vpc.vpc01.id
					  is_ipv6 = true
					}
					
					data "scaleway_vpc_routes" "routes_by_pn_id" {
					  vpc_id                     = scaleway_vpc.vpc01.id
					  nexthop_private_network_id = scaleway_vpc_private_network.pn01.id
					}
					`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.scaleway_vpc_routes.routes_by_vpc_id", "routes.#", "6"),
					resource.TestCheckResourceAttr("data.scaleway_vpc_routes.routes_by_ipv6", "routes.#", "2"),
					resource.TestCheckResourceAttr("data.scaleway_vpc_routes.routes_by_pn_id", "routes.#", "3"),
				),
			},
		},
	})
}
