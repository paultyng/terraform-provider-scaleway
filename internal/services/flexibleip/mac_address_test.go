package flexibleip_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	flexibleipSDK "github.com/scaleway/scaleway-sdk-go/api/flexibleip/v1alpha1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/locality"
	baremetalchecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/baremetal/testfuncs"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/flexibleip"
)

func TestAccFlexibleIPMACAddress_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckFlexibleIPDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
						resource "scaleway_flexible_ip" "main" {}

						resource "scaleway_flexible_ip_mac_address" "main" {
						  flexible_ip_id = scaleway_flexible_ip.main.id
						  type = "kvm"
						}
					`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlexibleIPAttachedMACAddress(tt, "scaleway_flexible_ip.main", "scaleway_flexible_ip_mac_address.main"),
					resource.TestCheckResourceAttr("scaleway_flexible_ip_mac_address.main", "zone", "fr-par-1"),
					resource.TestCheckResourceAttr("scaleway_flexible_ip_mac_address.main", "type", "kvm"),
					resource.TestCheckResourceAttrSet("scaleway_flexible_ip_mac_address.main", "address"),
					resource.TestCheckResourceAttrSet("scaleway_flexible_ip_mac_address.main", "created_at"),
					resource.TestCheckResourceAttrSet("scaleway_flexible_ip_mac_address.main", "updated_at")),
			},
		},
	})
}

func TestAccFlexibleIPMACAddress_MoveToAnotherFlexibleIP(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckFlexibleIPDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
						resource "scaleway_flexible_ip" "ip01" {
						  description = "attached to ip01"
						}

						resource "scaleway_flexible_ip" "ip02" {
						  description = "not attached to ip02"
						}

						resource "scaleway_flexible_ip_mac_address" "main" {
						  flexible_ip_id = scaleway_flexible_ip.ip01.id
						  type = "kvm"
						}
					`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlexibleIPExists(tt, "scaleway_flexible_ip.ip01"),
					testAccCheckFlexibleIPExists(tt, "scaleway_flexible_ip.ip02"),
					testAccCheckFlexibleIPAttachedMACAddress(tt, "scaleway_flexible_ip.ip01", "scaleway_flexible_ip_mac_address.main"),
					resource.TestCheckNoResourceAttr("scaleway_flexible_ip.ip02", "mac_address"),
				),
			},
			{
				Config: `
						resource "scaleway_flexible_ip" "ip01" {
						  description = "not attached to ip01"
						}

						resource "scaleway_flexible_ip" "ip02" {
						  description = "attached to ip02"
						}

						resource "scaleway_flexible_ip_mac_address" "main" {
						  flexible_ip_id = scaleway_flexible_ip.ip02.id
						  type = "kvm"
						}
					`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlexibleIPExists(tt, "scaleway_flexible_ip.ip01"),
					testAccCheckFlexibleIPExists(tt, "scaleway_flexible_ip.ip02"),
					testAccCheckFlexibleIPAttachedMACAddress(tt, "scaleway_flexible_ip.ip02", "scaleway_flexible_ip_mac_address.main"),
					resource.TestCheckNoResourceAttr("scaleway_flexible_ip.ip01", "mac_address"),
				),
			},
		},
	})
}

func TestAccFlexibleIPMACAddress_DuplicateOnOtherFlexibleIPs(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckFlexibleIPDestroy(tt),
			baremetalchecks.CheckServerDestroy(tt),
		),
		Steps: []resource.TestStep{
			{
				Config: `
						data "scaleway_baremetal_offer" "my_offer" {
					      	zone = "fr-par-1"
							name = "EM-A115X-SSD"
				     	}

						resource "scaleway_baremetal_server" "base" {
						  name 			             = "TestAccScalewayBaremetalServer_WithoutInstallConfig"
						  offer     				 = data.scaleway_baremetal_offer.my_offer.offer_id
						  install_config_afterward   = true
						}
			
						resource "scaleway_flexible_ip" "ip01" {
                          server_id = scaleway_baremetal_server.base.id	
						  description = "attached to ip01"
						}

						resource "scaleway_flexible_ip" "ip02" {
                          server_id = scaleway_baremetal_server.base.id
						  description = "not attached to ip02"
						}

						resource "scaleway_flexible_ip" "ip03" {
                          server_id = scaleway_baremetal_server.base.id
						  description = "not attached to ip03"
						}

						resource "scaleway_flexible_ip_mac_address" "main" {
						  flexible_ip_id = scaleway_flexible_ip.ip01.id
						  type           = "kvm"
						}
					`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlexibleIPExists(tt, "scaleway_flexible_ip.ip01"),
					testAccCheckFlexibleIPExists(tt, "scaleway_flexible_ip.ip02"),
					testAccCheckFlexibleIPAttachedMACAddress(tt, "scaleway_flexible_ip.ip01", "scaleway_flexible_ip_mac_address.main"),
					resource.TestCheckNoResourceAttr("scaleway_flexible_ip.ip02", "mac_address"),
				),
			},
			{
				Config: `
						data "scaleway_baremetal_offer" "my_offer" {
					      	zone = "fr-par-1"
							name = "EM-A115X-SSD"
				     	}

						resource "scaleway_baremetal_server" "base" {
						  name 			             = "TestAccScalewayBaremetalServer_WithoutInstallConfig"
						  offer     				 = data.scaleway_baremetal_offer.my_offer.offer_id
						  install_config_afterward   = true
						}
			
						resource "scaleway_flexible_ip" "ip01" {
                          server_id = scaleway_baremetal_server.base.id	
						  description = "attached to ip01"
						}

						resource "scaleway_flexible_ip" "ip02" {
                          server_id = scaleway_baremetal_server.base.id
						  description = "attached to ip02"
						}

						resource "scaleway_flexible_ip" "ip03" {
                          server_id = scaleway_baremetal_server.base.id
						  description = "attached to ip03"
						}

						resource "scaleway_flexible_ip_mac_address" "main" {
						  flexible_ip_id = scaleway_flexible_ip.ip01.id
						  type = "kvm"
                          flexible_ip_ids_to_duplicate = [
							scaleway_flexible_ip.ip02.id,
							scaleway_flexible_ip.ip03.id
						  ]
						}
					`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlexibleIPExists(tt, "scaleway_flexible_ip.ip01"),
					testAccCheckFlexibleIPExists(tt, "scaleway_flexible_ip.ip02"),
					testAccCheckFlexibleIPAttachedMACAddress(tt, "scaleway_flexible_ip.ip01", "scaleway_flexible_ip_mac_address.main"),
					testAccCheckFlexibleIPAttachedMACAddress(tt, "scaleway_flexible_ip.ip02", "scaleway_flexible_ip_mac_address.main"),
					testAccCheckFlexibleIPAttachedMACAddress(tt, "scaleway_flexible_ip.ip03", "scaleway_flexible_ip_mac_address.main"),
				),
			},
		},
	})
}

func testAccCheckFlexibleIPAttachedMACAddress(tt *acctest.TestTools, fipResource, macResource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fipState, ok := s.RootModule().Resources[fipResource]
		if !ok {
			return fmt.Errorf("resource not found: %s", fipResource)
		}
		macState, ok := s.RootModule().Resources[macResource]
		if !ok {
			return fmt.Errorf("resource not found: %s", macResource)
		}

		fipAPI, zone, ID, err := flexibleip.NewAPIWithZoneAndID(tt.Meta, fipState.Primary.ID)
		if err != nil {
			return err
		}
		ip, err := fipAPI.GetFlexibleIP(&flexibleipSDK.GetFlexibleIPRequest{
			FipID: ID,
			Zone:  zone,
		})
		if err != nil {
			return err
		}

		if ip.MacAddress != nil && locality.ExpandID(ip.MacAddress.ID) != locality.ExpandID(macState.Primary.ID) {
			return fmt.Errorf("IDs should be the same in %s and %s: %v is different than %v", fipResource, macResource, ip.MacAddress.ID, macState.Primary.ID)
		}

		return nil
	}
}
