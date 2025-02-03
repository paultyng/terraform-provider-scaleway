package vpc_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	vpcSDK "github.com/scaleway/scaleway-sdk-go/api/vpc/v2"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/vpc"
)

func TestAccVPCRoute_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isRouteDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_vpc" "vpc01" {
					  name = "tf-vpc-vpn"
					}
					
					resource "scaleway_vpc_private_network" "pn01" {
					  name = "tf-pn-vpn"
					  ipv4_subnet {
						subnet = "172.16.64.0/22"
					  }
					  vpc_id = scaleway_vpc.vpc01.id
					}
					
					resource "scaleway_instance_server" "server01" {
					  name  = "tf-server-vpn"
					  type  = "PLAY2-MICRO"
					  image = "openvpn"
					}
					
					resource "scaleway_instance_private_nic" "pnic01" {
					  private_network_id = scaleway_vpc_private_network.pn01.id
					  server_id          = scaleway_instance_server.server01.id
					}
					
					resource "scaleway_vpc_route" "rt01" {
					  vpc_id              = scaleway_vpc.vpc01.id
					  description         = "tf-route-vpn"
					  tags                = ["tf", "route"]
					  destination         = "10.0.0.0/24"
					  nexthop_resource_id = scaleway_instance_private_nic.pnic01.id
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isRoutePresent(tt, "scaleway_vpc_route.rt01"),
					resource.TestCheckResourceAttrPair("scaleway_vpc_route.rt01", "vpc_id", "scaleway_vpc.vpc01", "id"),
					resource.TestCheckResourceAttr("scaleway_vpc_route.rt01", "destination", "10.0.0.0/24"),
					resource.TestCheckResourceAttr("scaleway_vpc_route.rt01", "description", "tf-route-vpn"),
					resource.TestCheckResourceAttr("scaleway_vpc_route.rt01", "tags.#", "2"),
					resource.TestCheckResourceAttr("scaleway_vpc_route.rt01", "tags.0", "tf"),
					resource.TestCheckResourceAttr("scaleway_vpc_route.rt01", "tags.1", "route"),
					resource.TestCheckResourceAttr("scaleway_vpc_route.rt01", "region", "fr-par"),
				),
			},
			{
				Config: `
					resource "scaleway_vpc" "vpc01" {
					  name = "tf-vpc-vpn"
					}
					
					resource "scaleway_vpc_private_network" "pn01" {
					  name = "tf-pn-vpn"
					  ipv4_subnet {
						subnet = "172.16.64.0/22"
					  }
					  vpc_id = scaleway_vpc.vpc01.id
					}
					
					resource "scaleway_instance_server" "server01" {
					  name  = "tf-server-vpn"
					  type  = "PLAY2-MICRO"
					  image = "openvpn"
					}
					
					resource "scaleway_instance_private_nic" "pnic01" {
					  private_network_id = scaleway_vpc_private_network.pn01.id
					  server_id          = scaleway_instance_server.server01.id
					}

					resource "scaleway_instance_server" "server02" {
					  name  = "tf-server-vpn-2"
					  type  = "PLAY2-MICRO"
					  image = "openvpn"
					}
					
					resource "scaleway_instance_private_nic" "pnic02" {
					  private_network_id = scaleway_vpc_private_network.pn01.id
					  server_id          = scaleway_instance_server.server02.id
					}
					
					resource "scaleway_vpc_route" "rt01" {
					  vpc_id              = scaleway_vpc.vpc01.id
					  description         = "tf-route-vpn-updated"
					  tags                = ["tf", "route", "updated"]
					  destination         = "10.0.0.0/24"
					  nexthop_resource_id = scaleway_instance_private_nic.pnic02.id
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isRoutePresent(tt, "scaleway_vpc_route.rt01"),
					resource.TestCheckResourceAttrPair("scaleway_vpc_route.rt01", "vpc_id", "scaleway_vpc.vpc01", "id"),
					resource.TestCheckResourceAttr("scaleway_vpc_route.rt01", "destination", "10.0.0.0/24"),
					resource.TestCheckResourceAttr("scaleway_vpc_route.rt01", "description", "tf-route-vpn-updated"),
					resource.TestCheckResourceAttr("scaleway_vpc_route.rt01", "tags.#", "3"),
					resource.TestCheckResourceAttr("scaleway_vpc_route.rt01", "tags.0", "tf"),
					resource.TestCheckResourceAttr("scaleway_vpc_route.rt01", "tags.1", "route"),
					resource.TestCheckResourceAttr("scaleway_vpc_route.rt01", "tags.2", "updated"),
					resource.TestCheckResourceAttr("scaleway_vpc_route.rt01", "region", "fr-par"),
				),
			},
		},
	})
}

func isRoutePresent(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		vpcAPI, region, ID, err := vpc.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = vpcAPI.GetRoute(&vpcSDK.GetRouteRequest{
			RouteID: ID,
			Region:  region,
		})
		if err != nil {
			return err
		}

		return nil
	}
}

func isRouteDestroyed(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_vpc_route" {
				continue
			}

			vpcAPI, region, ID, err := vpc.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}

			_, err = vpcAPI.GetRoute(&vpcSDK.GetRouteRequest{
				RouteID: ID,
				Region:  region,
			})

			if err == nil {
				return fmt.Errorf("route (%s) still exists", rs.Primary.ID)
			}

			if !httperrors.Is404(err) {
				return err
			}
		}

		return nil
	}
}
