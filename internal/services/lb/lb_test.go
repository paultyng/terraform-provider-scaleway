package lb_test

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	lbSDK "github.com/scaleway/scaleway-sdk-go/api/lb/v1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	instancechecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/instance/testfuncs"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/lb"
	lbchecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/lb/testfuncs"
	vpcchecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/vpc/testfuncs"
	vpcgwchecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/vpcgw/testfuncs"
)

func TestAccLB_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isLbDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_lb_ip main {
					}

					resource scaleway_lb main {
						ip_id = scaleway_lb_ip.main.id
						name = "test-lb-basic"
						description = "a description"
						type = "LB-S"
						tags = ["basic"]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isLbPresent(tt, "scaleway_lb.main"),
					acctest.CheckResourceAttrUUID("scaleway_lb.main", "id"),
					resource.TestCheckResourceAttr("scaleway_lb.main", "name", "test-lb-basic"),
					resource.TestCheckResourceAttr("scaleway_lb.main", "type", "LB-S"),
					resource.TestCheckResourceAttr("scaleway_lb.main", "tags.#", "1"),
					resource.TestCheckResourceAttr("scaleway_lb.main", "description", "a description"),
					resource.TestCheckResourceAttr("scaleway_lb.main", "ssl_compatibility_level", lbSDK.SSLCompatibilityLevelSslCompatibilityLevelIntermediate.String()),
				),
			},
			{
				Config: `
					resource scaleway_lb_ip main {
					}

					resource scaleway_lb main {
						ip_id = scaleway_lb_ip.main.id
						name = "test-lb-rename"
						description = "another description"
						type = "LB-S"
						tags = ["basic", "tag2"]
						ssl_compatibility_level = "ssl_compatibility_level_modern"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isLbPresent(tt, "scaleway_lb.main"),
					acctest.CheckResourceAttrUUID("scaleway_lb.main", "id"),
					resource.TestCheckResourceAttr("scaleway_lb.main", "name", "test-lb-rename"),
					resource.TestCheckResourceAttr("scaleway_lb.main", "tags.#", "2"),
					resource.TestCheckResourceAttr("scaleway_lb.main", "description", "another description"),
					resource.TestCheckResourceAttr("scaleway_lb.main", "ssl_compatibility_level", lbSDK.SSLCompatibilityLevelSslCompatibilityLevelModern.String()),
				),
			},
		},
	})
}

func TestAccLB_Private(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isLbDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_lb main {
						name = "test-lb-basic"
						description = "a description"
						type = "LB-S"
						tags = ["basic"]
						assign_flexible_ip = false
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isLbPresent(tt, "scaleway_lb.main"),
					acctest.CheckResourceAttrUUID("scaleway_lb.main", "id"),
					resource.TestCheckResourceAttr("scaleway_lb.main", "assign_flexible_ip", "false"),
					resource.TestCheckNoResourceAttr("scaleway_lb.main", "ip_address"),
					resource.TestCheckResourceAttr("scaleway_lb.main", "name", "test-lb-basic"),
					resource.TestCheckResourceAttr("scaleway_lb.main", "type", "LB-S"),
					resource.TestCheckResourceAttr("scaleway_lb.main", "tags.#", "1"),
					resource.TestCheckResourceAttr("scaleway_lb.main", "description", "a description"),
				),
			},
			{
				Config: `
					resource scaleway_lb main {
						name = "test-lb-rename"
						description = "another description"
						type = "LB-S"
						tags = ["basic", "tag2"]
						assign_flexible_ip = false
						ssl_compatibility_level = "ssl_compatibility_level_modern"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isLbPresent(tt, "scaleway_lb.main"),
					acctest.CheckResourceAttrUUID("scaleway_lb.main", "id"),
					resource.TestCheckResourceAttr("scaleway_lb.main", "assign_flexible_ip", "false"),
					resource.TestCheckNoResourceAttr("scaleway_lb.main", "ip_address"),
					resource.TestCheckResourceAttr("scaleway_lb.main", "name", "test-lb-rename"),
					resource.TestCheckResourceAttr("scaleway_lb.main", "tags.#", "2"),
					resource.TestCheckResourceAttr("scaleway_lb.main", "description", "another description"),
					resource.TestCheckResourceAttr("scaleway_lb.main", "ssl_compatibility_level", lbSDK.SSLCompatibilityLevelSslCompatibilityLevelModern.String()),
				),
			},
		},
	})
}

func TestAccLB_AssignedIPs(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isLbDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_lb main {
						name = "test-lb-assigned-ips"
						description = "a description"
						type = "LB-S"
						tags = ["basic"]
						assign_flexible_ip = true
					    assign_flexible_ipv6 = true
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isLbPresent(tt, "scaleway_lb.main"),
					resource.TestCheckResourceAttr("scaleway_lb.main", "assign_flexible_ip", "true"),
					resource.TestCheckResourceAttr("scaleway_lb.main", "assign_flexible_ipv6", "true"),
					acctest.CheckResourceAttrIPv4("scaleway_lb.main", "ip_address"),
					acctest.CheckResourceAttrIPv6("scaleway_lb.main", "ipv6_address"),
				),
			},
		},
	})
}

func TestAccLB_WithIPv6(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isLbDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_lb_ip" "v4" {
					}
				`,
			},
			{
				Config: `
					resource "scaleway_lb_ip" "v4" {
					}
					
					resource "scaleway_lb_ip" "v6" {
					  is_ipv6 = true
					}

					resource scaleway_lb main {
					    ip_ids = [scaleway_lb_ip.v4.id, scaleway_lb_ip.v6.id]
						name   = "test-lb-ip-ids"
						type   = "LB-S"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isLbPresent(tt, "scaleway_lb.main"),
					acctest.CheckResourceAttrIPv4("scaleway_lb.main", "ip_address"),
					acctest.CheckResourceAttrIPv6("scaleway_lb.main", "ipv6_address"),
				),
			},
		},
	})
}

func TestAccLB_UpdateToIPv6(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isLbDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_lb_ip" "v4" {
					}
					resource scaleway_lb main {
					    ip_ids = [scaleway_lb_ip.v4.id]
						name   = "test-lb-ip-ids"
						type   = "LB-S"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isLbPresent(tt, "scaleway_lb.main"),
					acctest.CheckResourceAttrIPv4("scaleway_lb.main", "ip_address"),
				),
			},
			{
				Config: `
					resource "scaleway_lb_ip" "v4" {
					}
					resource "scaleway_lb_ip" "v6" {
					  is_ipv6 = true
					}
					resource scaleway_lb main {
					    ip_ids = [scaleway_lb_ip.v4.id, scaleway_lb_ip.v6.id]
						name   = "test-lb-ip-ids"
						type   = "LB-S"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isLbPresent(tt, "scaleway_lb.main"),
					acctest.CheckResourceAttrIPv4("scaleway_lb.main", "ip_address"),
					acctest.CheckResourceAttrIPv6("scaleway_lb.main", "ipv6_address"),
				),
			},
		},
	})
}

func TestAccLB_Migrate(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	lbID := ""

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isLbDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					### IP for LB IP
					resource scaleway_lb_ip main {
					}

					resource scaleway_lb main {
						ip_id = scaleway_lb_ip.main.id
						name = "test-lb-migration"
						type = "LB-S"
						tags = ["basic", "tag2"]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isLbPresent(tt, "scaleway_lb.main"),
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["scaleway_lb.main"]
						if !ok {
							return fmt.Errorf("resource not found: %s", "scaleway_lb.main")
						}
						lbID = rs.Primary.ID
						return nil
					},
					resource.TestCheckResourceAttr("scaleway_lb.main", "type", "LB-S"),
					resource.TestCheckResourceAttr("scaleway_lb.main", "name", "test-lb-migration"),
					resource.TestCheckResourceAttr("scaleway_lb.main", "tags.#", "2"),
				),
			},
			{
				Config: `
					### IP for LB IP
					resource scaleway_lb_ip main {
					}

					resource scaleway_lb main {
						ip_id = scaleway_lb_ip.main.id
						name = "test-lb-migrate-lb-gp-m"
						type = "LB-GP-M"
						tags = ["migration"]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isLbPresent(tt, "scaleway_lb.main"),
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["scaleway_lb.main"]
						if !ok {
							return fmt.Errorf("resource not found: %s", "scaleway_lb.main")
						}
						if rs.Primary.ID != lbID {
							return errors.New("LB id has changed")
						}
						return nil
					},
					resource.TestCheckResourceAttr("scaleway_lb.main", "type", "LB-GP-M"),
					resource.TestCheckResourceAttr("scaleway_lb.main", "name", "test-lb-migrate-lb-gp-m"),
					resource.TestCheckResourceAttr("scaleway_lb.main", "tags.#", "1"),
				),
			},
			{
				Config: `
					### IP for LB IP
					resource scaleway_lb_ip main {
					}

					resource scaleway_lb main {
						ip_id = scaleway_lb_ip.main.id
						name = "test-lb-migrate-lb-gp-m"
						type = "LB-GP-M"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isLbPresent(tt, "scaleway_lb.main"),
					resource.TestCheckResourceAttr("scaleway_lb.main", "type", "LB-GP-M"),
					resource.TestCheckResourceAttr("scaleway_lb.main", "name", "test-lb-migrate-lb-gp-m"),
					resource.TestCheckResourceAttr("scaleway_lb.main", "tags.#", "0"),
				),
			},
			{
				Config: `
					### IP for LB IP
					resource scaleway_lb_ip main {
					}

					resource scaleway_lb main {
						ip_id = scaleway_lb_ip.main.id
						name = "test-lb-migrate-down"
						type = "LB-S"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isLbPresent(tt, "scaleway_lb.main"),
					resource.TestCheckResourceAttr("scaleway_lb.main", "type", "LB-S"),
					resource.TestCheckResourceAttr("scaleway_lb.main", "name", "test-lb-migrate-down"),
					resource.TestCheckResourceAttr("scaleway_lb.main", "tags.#", "0"),
				),
			},
		},
	})
}

func TestAccLB_WithPrivateNetworksOnDHCPConfig(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			instancechecks.IsServerDestroyed(tt),
			isLbDestroyed(tt),
			lbchecks.IsIPDestroyed(tt),
			vpcgwchecks.IsGatewayNetworkDestroyed(tt),
			vpcchecks.CheckPrivateNetworkDestroy(tt),
			vpcgwchecks.IsDHCPDestroyed(tt),
			vpcgwchecks.IsGatewayDestroyed(tt),
			vpcgwchecks.IsIPDestroyed(tt),
		),
		Steps: []resource.TestStep{
			{
				Config: `
				### IP for Public Gateway
				resource "scaleway_vpc_public_gateway_ip" "main" {
				}

				### The Public Gateway with the Attached IP
				resource "scaleway_vpc_public_gateway" "main" {
					name  = "tf-test-public-gw"
					type  = "VPC-GW-S"
					ip_id = scaleway_vpc_public_gateway_ip.main.id
				}
				
				### Scaleway Private Network
				resource "scaleway_vpc_private_network" "main" {
					name = "private network with a DHCP config"
				}
				
				### DHCP Space of VPC
				resource "scaleway_vpc_public_gateway_dhcp" "main" {
					subnet = "10.0.0.0/24"
				}
				
				### VPC Gateway Network
				resource "scaleway_vpc_gateway_network" "main" {
					gateway_id         = scaleway_vpc_public_gateway.main.id
					private_network_id = scaleway_vpc_private_network.main.id
					dhcp_id            = scaleway_vpc_public_gateway_dhcp.main.id
					cleanup_dhcp       = true
					enable_masquerade  = true
				}
				
				### Scaleway Instance
				resource "scaleway_instance_server" "main" {
					name        = "Scaleway Terraform Provider"
					type        = "DEV1-S"
					image       = "debian_bullseye"
					enable_ipv6 = false
				
					private_network {
						pn_id = scaleway_vpc_private_network.main.id
					}
				}

				### IP for LB IP
				resource scaleway_lb_ip ip01 {
				}
				
				resource scaleway_lb lb01 {
					ip_id = scaleway_lb_ip.ip01.id
					name = "test-lb-with-private-network-configs"
					type = "LB-S"
				
					private_network {
						private_network_id = scaleway_vpc_private_network.main.id
						dhcp_config = true
					}
				
					depends_on = [scaleway_vpc_public_gateway.main]
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					isLbPresent(tt, "scaleway_lb.lb01"),
					isIPPresent(tt, "scaleway_lb_ip.ip01"),
					resource.TestCheckResourceAttrSet("scaleway_vpc_private_network.main", "name"),
					resource.TestCheckResourceAttrPair(
						"scaleway_lb.lb01", "private_network.0.private_network_id",
						"scaleway_vpc_private_network.main", "id"),
					resource.TestCheckResourceAttrPair(
						"scaleway_instance_server.main", "private_network.0.pn_id",
						"scaleway_vpc_private_network.main", "id"),
					resource.TestCheckResourceAttr("scaleway_lb.lb01",
						"private_network.0.status", lbSDK.PrivateNetworkStatusReady.String()),
					resource.TestCheckResourceAttr("scaleway_lb.lb01",
						"private_network.0.dhcp_config", "true"),
					resource.TestCheckResourceAttr("scaleway_lb.lb01",
						"private_network.0.status", lbSDK.PrivateNetworkStatusReady.String()),
				),
			},
		},
	})
}

func TestAccLB_WithPrivateNetworksIPAMIDs(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			isLbDestroyed(tt),
			vpcchecks.CheckPrivateNetworkDestroy(tt),
		),
		Steps: []resource.TestStep{
			{
				Config: `
				resource "scaleway_vpc" "vpc01" {
				  name = "my vpc"
				}
				
				resource "scaleway_vpc_private_network" "pn01" {
				  vpc_id = scaleway_vpc.vpc01.id
				  ipv4_subnet {
					subnet = "172.16.32.0/22"
				  }
				}
				
				resource "scaleway_ipam_ip" "ip01" {
				  address = "172.16.32.7"
				  source {
					private_network_id = scaleway_vpc_private_network.pn01.id
				  }
				}
				
				resource scaleway_lb lb01 {
				  name = "test-lb-with-private-network-ipam"
				  type = "LB-S"
				
				  private_network {
				    private_network_id = scaleway_vpc_private_network.pn01.id
				    ipam_ids = [scaleway_ipam_ip.ip01.id]
				  }	
				}

				data "scaleway_ipam_ip" "by_name" {
				  resource {
					name = scaleway_lb.lb01.name
					type = "lb_server"
				  }
				  type = "ipv4"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					isLbPresent(tt, "scaleway_lb.lb01"),
					resource.TestCheckResourceAttrPair(
						"scaleway_lb.lb01", "private_network.0.private_network_id",
						"scaleway_vpc_private_network.pn01", "id"),
					resource.TestCheckResourceAttrPair(
						"scaleway_lb.lb01", "private_network.0.ipam_ids.0",
						"scaleway_ipam_ip.ip01", "id"),
					resource.TestCheckResourceAttrPair(
						"scaleway_ipam_ip.ip01", "address",
						"data.scaleway_ipam_ip.by_name", "address_cidr"),
				),
			},
		},
	})
}

func TestAccLB_WithoutPNConfig(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isLbDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_lb_ip" "ip01" {}

					resource "scaleway_vpc_private_network" "pn" {
						name = "pn-with-lb-static"
					}

					resource "scaleway_lb" "lb01" {
					    ip_id = scaleway_lb_ip.ip01.id
						name = "test-lb-with-pn-static-cidr"
						type = "LB-S"
						private_network {
							private_network_id = scaleway_vpc_private_network.pn.id
						}
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isLbPresent(tt, "scaleway_lb.lb01"),
					isIPPresent(tt, "scaleway_lb_ip.ip01"),
					resource.TestCheckResourceAttrPair(
						"scaleway_lb.lb01", "private_network.0.private_network_id",
						"scaleway_vpc_private_network.pn", "id"),
					resource.TestCheckResourceAttr("scaleway_lb.lb01",
						"private_network.0.status", lbSDK.PrivateNetworkStatusReady.String()),
					resource.TestCheckResourceAttr("scaleway_lb.lb01",
						"private_network.0.dhcp_config", "true"),
				),
			},
		},
	})
}

func TestAccLB_DifferentLocalityIPID(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isLbDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_lb_ip main {
						zone = "fr-par-2"
					}

					resource scaleway_lb main {
						ip_id = scaleway_lb_ip.main.id
						name = "test-lb-basic"
						type = "LB-S"
						zone = "fr-par-1"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isLbPresent(tt, "scaleway_lb.main"),
					acctest.CheckResourceAttrUUID("scaleway_lb.main", "id"),
					resource.TestCheckResourceAttr("scaleway_lb.main", "name", "test-lb-basic"),
					resource.TestCheckResourceAttr("scaleway_lb.main", "type", "LB-S"),
				),
				ExpectError: regexp.MustCompile("has different locality than the resource"),
			},
			{
				Config: `
					resource scaleway_lb_ip main {
						zone = "fr-par-2"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isIPPresent(tt, "scaleway_lb_ip.main"),
				),
			},
			{
				Config: `
					resource scaleway_lb_ip main {
						zone = "fr-par-2"
					}

					resource scaleway_lb main {
						ip_id = scaleway_lb_ip.main.id
						name = "test-lb-basic"
						type = "LB-S"
						zone = "fr-par-1"
					}
				`,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("has different locality than the resource"),
			},
		},
	})
}

func isLbPresent(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		lbAPI, zone, ID, err := lb.NewAPIWithZoneAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = lbAPI.GetLB(&lbSDK.ZonedAPIGetLBRequest{
			LBID: ID,
			Zone: zone,
		})
		if err != nil {
			return err
		}

		return nil
	}
}

func isLbDestroyed(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_lb" {
				continue
			}

			lbAPI, zone, ID, err := lb.NewAPIWithZoneAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}

			_, err = lbAPI.GetLB(&lbSDK.ZonedAPIGetLBRequest{
				Zone: zone,
				LBID: ID,
			})

			// If no error resource still exist
			if err == nil {
				return fmt.Errorf("load Balancer (%s) still exists", rs.Primary.ID)
			}

			// Unexpected api error we return it
			if !httperrors.Is404(err) {
				return err
			}
		}

		return nil
	}
}

func TestLbUpgradeV1SchemaUpgradeFunc(t *testing.T) {
	v0Schema := map[string]interface{}{
		"id": "fr-par/22c61530-834c-4ab4-aa71-aaaa2ac9d45a",
	}
	v1Schema := map[string]interface{}{
		"id": "fr-par-1/22c61530-834c-4ab4-aa71-aaaa2ac9d45a",
	}

	actual, err := lb.UpgradeStateV1Func(context.Background(), v0Schema, nil)
	if err != nil {
		t.Fatalf("error migrating state: %s", err)
	}

	if !reflect.DeepEqual(v1Schema, actual) {
		t.Fatalf("\n\nexpected:\n\n%#v\n\ngot:\n\n%#v\n\n", v1Schema, actual)
	}
}
