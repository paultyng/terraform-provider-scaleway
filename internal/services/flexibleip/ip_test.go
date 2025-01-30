package flexibleip_test

import (
	"fmt"
	"net"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	baremetalSDK "github.com/scaleway/scaleway-sdk-go/api/baremetal/v1"
	flexibleipSDK "github.com/scaleway/scaleway-sdk-go/api/flexibleip/v1alpha1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/locality"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/baremetal"
	baremetalchecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/baremetal/testfuncs"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/flexibleip"
)

const SSHKeyFlexibleIP = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIM7HUxRyQtB2rnlhQUcbDGCZcTJg7OvoznOiyC9W6IxH opensource@scaleway.com"

func TestAccFlexibleIP_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckFlexibleIPDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
						resource "scaleway_flexible_ip" "main" {}
					`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlexibleIPExists(tt, "scaleway_flexible_ip.main"),
				),
			},
			{
				ResourceName:            "scaleway_flexible_ip.main",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"is_ipv6"},
			},
		},
	})
}

func TestAccFlexibleIP_WithZone(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckFlexibleIPDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
						resource "scaleway_flexible_ip" "base" {}
					`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlexibleIPExists(tt, "scaleway_flexible_ip.base"),
					resource.TestCheckResourceAttr("scaleway_flexible_ip.base", "zone", "fr-par-1"),
				),
			},
			{
				Config: `
						resource "scaleway_flexible_ip" "base" {
							zone = "nl-ams-1"
						}
					`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlexibleIPExists(tt, "scaleway_flexible_ip.base"),
					resource.TestCheckResourceAttr("scaleway_flexible_ip.base", "zone", "nl-ams-1"),
				),
			},
		},
	})
}

func TestAccFlexibleIP_IPv6(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckFlexibleIPDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
						resource "scaleway_flexible_ip" "main" {
							is_ipv6 = true
						}
					`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlexibleIPExists(tt, "scaleway_flexible_ip.main"),
					resource.TestCheckResourceAttr("scaleway_flexible_ip.main", "is_ipv6", "true"),
					testAccCheckFlexibleIPIsIPv6(tt, "scaleway_flexible_ip.main"),
				),
			},
		},
	})
}

func TestAccFlexibleIP_CreateAndAttachToBaremetalServer(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	SSHKeyName := "TestAccScalewayFlexibleIP_CreateAndAttachToBaremetalServer"
	name := "TestAccScalewayFlexibleIP_CreateAndAttachToBaremetalServer"

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckFlexibleIPDestroy(tt),
			baremetalchecks.CheckServerDestroy(tt),
		),
		Steps: []resource.TestStep{
			{
				Config: `
						resource "scaleway_flexible_ip" "base" {
							zone = "fr-par-1"
						}
					`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlexibleIPExists(tt, "scaleway_flexible_ip.base"),
					resource.TestCheckResourceAttr("scaleway_flexible_ip.base", "zone", "fr-par-1"),
				),
			},
			{
				Config: fmt.Sprintf(`
						data "scaleway_baremetal_os" "by_id" {
							zone = "fr-par-1"
							name = "Ubuntu"
							version = "22.04 LTS (Jammy Jellyfish)"						
						}

						data "scaleway_baremetal_offer" "my_offer" {
							name = "EM-A115X-SSD"
					  		zone = "fr-par-1"
						}				

						resource "scaleway_iam_ssh_key" "main" {
							name 	   = "%s"
							public_key = "%s"
						}

						resource "scaleway_baremetal_server" "base" {
							name        = "%s"
							zone        = "fr-par-1"
							offer       = data.scaleway_baremetal_offer.my_offer.offer_id
							os          = data.scaleway_baremetal_os.by_id.os_id

							ssh_key_ids = [ scaleway_iam_ssh_key.main.id ]
						}
					`, SSHKeyName, SSHKeyFlexibleIP, name),
			},
			{
				Config: fmt.Sprintf(`
						data "scaleway_baremetal_os" "by_id" {
							zone = "fr-par-1"
							name = "Ubuntu"
							version = "22.04 LTS (Jammy Jellyfish)"						
						}

						data "scaleway_baremetal_offer" "my_offer" {
							name = "EM-A115X-SSD"
					  		zone = "fr-par-1"
						}				

						resource "scaleway_iam_ssh_key" "main" {
							name 	   = "%s"
							public_key = "%s"
						}

						resource "scaleway_baremetal_server" "base" {
							name        = "%s"
							zone        = "fr-par-1"
							offer       = data.scaleway_baremetal_offer.my_offer.offer_id
							os          = data.scaleway_baremetal_os.by_id.os_id

							ssh_key_ids = [ scaleway_iam_ssh_key.main.id ]
						}

						resource "scaleway_flexible_ip" "base" {
							server_id = scaleway_baremetal_server.base.id
							zone = "fr-par-1"
						}
					`, SSHKeyName, SSHKeyFlexibleIP, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlexibleIPExists(tt, "scaleway_flexible_ip.base"),
					testAccCheckFlexibleIPAttachedToBaremetalServer(tt, "scaleway_flexible_ip.base", "scaleway_baremetal_server.base"),
					resource.TestCheckResourceAttr("scaleway_flexible_ip.base", "zone", "fr-par-1"),
				),
			},
		},
	})
}

func TestAccFlexibleIP_AttachAndDetachFromBaremetalServer(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	SSHKeyName := "TestAccScalewayFlexibleIP_AttachAndDetachFromBaremetalServer"
	name := "TestAccScalewayFlexibleIP_AttachAndDetachFromBaremetalServer"
	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckFlexibleIPDestroy(tt),
			baremetalchecks.CheckServerDestroy(tt),
		),
		Steps: []resource.TestStep{
			{
				Config: `
						resource "scaleway_flexible_ip" "base" {
							zone = "fr-par-1"
						}
					`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlexibleIPExists(tt, "scaleway_flexible_ip.base"),
					resource.TestCheckResourceAttr("scaleway_flexible_ip.base", "zone", "fr-par-1"),
				),
			},
			{
				Config: fmt.Sprintf(`
						data "scaleway_baremetal_os" "by_id" {
							zone = "fr-par-1"
							name = "Ubuntu"
							version = "22.04 LTS (Jammy Jellyfish)"						
						}

						data "scaleway_baremetal_offer" "my_offer" {
							name = "EM-A115X-SSD"
					  		zone = "fr-par-1"
						}		

						resource "scaleway_iam_ssh_key" "main" {
							name 	   = "%s"
							public_key = "%s"
						}

						resource "scaleway_baremetal_server" "base" {
							name        = "%s"
							zone        = "fr-par-1"
							offer       = data.scaleway_baremetal_offer.my_offer.offer_id
							os          = data.scaleway_baremetal_os.by_id.os_id

							ssh_key_ids = [ scaleway_iam_ssh_key.main.id ]
						}
					`, SSHKeyName, SSHKeyFlexibleIP, name),
			},
			{
				Config: fmt.Sprintf(`
						data "scaleway_baremetal_os" "by_id" {
							zone = "fr-par-1"
							name = "Ubuntu"
							version = "22.04 LTS (Jammy Jellyfish)"						
						}

						data "scaleway_baremetal_offer" "my_offer" {
							name = "EM-A115X-SSD"
					  		zone = "fr-par-1"
						}		

						resource "scaleway_iam_ssh_key" "main" {
							name 	   = "%s"
							public_key = "%s"
						}

						resource "scaleway_baremetal_server" "base" {
							name        = "%s"
							zone        = "fr-par-1"
							offer       = data.scaleway_baremetal_offer.my_offer.offer_id
							os          = data.scaleway_baremetal_os.by_id.os_id

							ssh_key_ids = [ scaleway_iam_ssh_key.main.id ]
						}

						resource "scaleway_flexible_ip" "base" {
							server_id = scaleway_baremetal_server.base.id
							zone = "fr-par-1"
						}
					`, SSHKeyName, SSHKeyFlexibleIP, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlexibleIPExists(tt, "scaleway_flexible_ip.base"),
					testAccCheckFlexibleIPAttachedToBaremetalServer(tt, "scaleway_flexible_ip.base", "scaleway_baremetal_server.base"),
					resource.TestCheckResourceAttr("scaleway_flexible_ip.base", "zone", "fr-par-1"),
				),
			},
			{
				Config: `
						resource "scaleway_flexible_ip" "base" {
							zone = "fr-par-1"
						}
					`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlexibleIPExists(tt, "scaleway_flexible_ip.base"),
					resource.TestCheckResourceAttr("scaleway_flexible_ip.base", "zone", "fr-par-1"),
				),
			},
		},
	})
}

func testAccCheckFlexibleIPExists(tt *acctest.TestTools, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("resource not found: %s", name)
		}

		fipAPI, zone, ID, err := flexibleip.NewAPIWithZoneAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = fipAPI.GetFlexibleIP(&flexibleipSDK.GetFlexibleIPRequest{
			FipID: ID,
			Zone:  zone,
		})
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckFlexibleIPDestroy(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "scaleway_flexible_ip" {
				continue
			}

			fipAPI, zone, id, err := flexibleip.NewAPIWithZoneAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}

			_, err = fipAPI.GetFlexibleIP(&flexibleipSDK.GetFlexibleIPRequest{
				FipID: id,
				Zone:  zone,
			})

			// If no error resource still exist
			if err == nil {
				return fmt.Errorf("resource %s(%s) still exist", rs.Type, rs.Primary.ID)
			}

			// Unexpected api error we return it
			// We check for 403 because instance API return 403 for deleted IP
			if !httperrors.Is404(err) && !httperrors.Is403(err) {
				return err
			}
		}

		return nil
	}
}

func testAccCheckFlexibleIPAttachedToBaremetalServer(tt *acctest.TestTools, ipResource, serverResource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ipState, ok := s.RootModule().Resources[ipResource]
		if !ok {
			return fmt.Errorf("resource not found: %s", ipResource)
		}
		serverState, ok := s.RootModule().Resources[serverResource]
		if !ok {
			return fmt.Errorf("resource not found: %s", serverResource)
		}

		baremetalAPI, zoneID, err := baremetal.NewAPIWithZoneAndID(tt.Meta, ipState.Primary.ID)
		if err != nil {
			return err
		}

		server, err := baremetalAPI.GetServer(&baremetalSDK.GetServerRequest{
			Zone:     zoneID.Zone,
			ServerID: locality.ExpandID(serverState.Primary.ID),
		})
		if err != nil {
			return err
		}

		fipAPI, zone, ID, err := flexibleip.NewAPIWithZoneAndID(tt.Meta, ipState.Primary.ID)
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

		if ip.ServerID == nil || server.ID != *ip.ServerID {
			return fmt.Errorf("IDs should be the same in %s and %s: %v is different than %v", ipResource, serverResource, server.ID, ip.ServerID)
		}

		return nil
	}
}

func testAccCheckFlexibleIPIsIPv6(tt *acctest.TestTools, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		fipAPI, zone, ID, err := flexibleip.NewAPIWithZoneAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		flexibleIP, err := fipAPI.GetFlexibleIP(&flexibleipSDK.GetFlexibleIPRequest{
			Zone:  zone,
			FipID: ID,
		})
		if err != nil {
			return err
		}

		if len(flexibleIP.IPAddress.IP.To16()) != net.IPv6len {
			return fmt.Errorf("expected an IPv6 address but got: %s", flexibleIP.IPAddress.String())
		}

		return nil
	}
}
