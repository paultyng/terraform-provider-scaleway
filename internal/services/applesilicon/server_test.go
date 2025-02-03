package applesilicon_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	applesiliconSDK "github.com/scaleway/scaleway-sdk-go/api/applesilicon/v1alpha1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/applesilicon"
)

func TestAccServer_Basic(t *testing.T) {
	//t.Skip("Skipping AppleSilicon deletion before 24h")
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isServerDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_apple_silicon_server main {
						name = "test-m2"
						type = "M2-M"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isServerPresent(tt, "scaleway_apple_silicon_server.main"),
					resource.TestCheckResourceAttr("scaleway_apple_silicon_server.main", "name", "test-m2"),
					resource.TestCheckResourceAttr("scaleway_apple_silicon_server.main", "type", "M2-M"),
					// Computed
					resource.TestCheckResourceAttrSet("scaleway_apple_silicon_server.main", "ip"),
					resource.TestCheckResourceAttrSet("scaleway_apple_silicon_server.main", "vnc_url"),
					resource.TestCheckResourceAttrSet("scaleway_apple_silicon_server.main", "created_at"),
					resource.TestCheckResourceAttrSet("scaleway_apple_silicon_server.main", "deletable_at"),
				),
			},
		},
	})
}

func TestAccServer_EnableVPC(t *testing.T) {
	//t.Skip("Skipping AppleSilicon VPC not available")
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isServerDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_vpc" "vpc01" {
					  name = "TestAccServer_EnableVPC"
					}
					
					resource "scaleway_vpc_private_network" "pn01" {
					  name = "TestAccServer_EnableVPC"
					  ipv4_subnet {
						subnet = "172.16.64.0/22"
					  }
					  vpc_id = scaleway_vpc.vpc01.id
					}
					
					resource "scaleway_ipam_ip" "ip01" {
					  address = "172.16.64.7"
					  source {
						private_network_id = scaleway_vpc_private_network.pn01.id
					  }
					}

					resource "scaleway_ipam_ip" "ip02" {
					  address = "172.16.64.9"
					  source {
						private_network_id = scaleway_vpc_private_network.pn01.id
					  }
					}

					resource scaleway_apple_silicon_server main {
						name = "test-m2"
						type = "M2-M"
						enable_vpc = true
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isServerPresent(tt, "scaleway_apple_silicon_server.main"),
					resource.TestCheckResourceAttr("scaleway_apple_silicon_server.main", "name", "test-m2"),
					resource.TestCheckResourceAttr("scaleway_apple_silicon_server.main", "type", "M2-M"),
					// Computed
					resource.TestCheckResourceAttrSet("scaleway_apple_silicon_server.main", "ip"),
					resource.TestCheckResourceAttrSet("scaleway_apple_silicon_server.main", "vnc_url"),
					resource.TestCheckResourceAttrSet("scaleway_apple_silicon_server.main", "created_at"),
					resource.TestCheckResourceAttrSet("scaleway_apple_silicon_server.main", "deletable_at"),
					resource.TestCheckResourceAttr("scaleway_apple_silicon_server.main", "vpc_status", "vpc_enabled"),
				),
			},
			{
				Config: `
					resource scaleway_apple_silicon_server main {
						name = "test-m2"
						type = "M2-M"
						enable_vpc = false
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isServerPresent(tt, "scaleway_apple_silicon_server.main"),
					resource.TestCheckResourceAttr("scaleway_apple_silicon_server.main", "name", "test-m2"),
					resource.TestCheckResourceAttr("scaleway_apple_silicon_server.main", "type", "M2-M"),
					// Computed
					resource.TestCheckResourceAttrSet("scaleway_apple_silicon_server.main", "ip"),
					resource.TestCheckResourceAttrSet("scaleway_apple_silicon_server.main", "vnc_url"),
					resource.TestCheckResourceAttrSet("scaleway_apple_silicon_server.main", "created_at"),
					resource.TestCheckResourceAttrSet("scaleway_apple_silicon_server.main", "deletable_at"),
					resource.TestCheckResourceAttr("scaleway_apple_silicon_server.main", "vpc_status", "vpc_updating"),
					resource.TestCheckNoResourceAttr("scaleway_apple_silicon_server.main", "vpc_id"),
				),
			},
		},
	})
}

func isServerPresent(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		asAPI, zone, ID, err := applesilicon.NewAPIWithZoneAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = asAPI.GetServer(&applesiliconSDK.GetServerRequest{
			ServerID: ID,
			Zone:     zone,
		})
		if err != nil {
			return err
		}

		return nil
	}
}

func isServerDestroyed(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_apple_silicon_server" {
				continue
			}

			asAPI, zone, ID, err := applesilicon.NewAPIWithZoneAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}

			_, err = asAPI.GetServer(&applesiliconSDK.GetServerRequest{
				ServerID: ID,
				Zone:     zone,
			})

			// If no error resource still exist
			if err == nil {
				return fmt.Errorf("server (%s) still exists", rs.Primary.ID)
			}

			// Unexpected api error we return it
			if !httperrors.Is403(err) {
				return err
			}
		}

		return nil
	}
}
