package ipam_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	ipamSDK "github.com/scaleway/scaleway-sdk-go/api/ipam/v1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/ipam"
	ipamchecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/ipam/testfuncs"
)

func TestAccIPAMIP_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      ipamchecks.CheckIPDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_vpc vpc01 {
						name = "my vpc"
					}

					resource scaleway_vpc_private_network pn01 {
						vpc_id = scaleway_vpc.vpc01.id
						ipv4_subnet {
							subnet = "172.16.32.0/22"
						}
					}

					resource scaleway_ipam_ip ip01 {
					 source {
						private_network_id = scaleway_vpc_private_network.pn01.id
					  }
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPAMIPExists(tt, "scaleway_ipam_ip.ip01"),
					resource.TestCheckResourceAttr("scaleway_ipam_ip.ip01", "is_ipv6", "false"),
					acctest.CheckResourceRawIDMatches("scaleway_ipam_ip.ip01", "source.0.private_network_id", "scaleway_vpc_private_network.pn01", "id"),
					resource.TestCheckResourceAttrSet("scaleway_ipam_ip.ip01", "source.0.subnet_id"),
					resource.TestCheckResourceAttrSet("scaleway_ipam_ip.ip01", "address"),
					resource.TestCheckResourceAttrSet("scaleway_ipam_ip.ip01", "created_at"),
					resource.TestCheckResourceAttrSet("scaleway_ipam_ip.ip01", "updated_at"),
				),
			},
		},
	})
}

func TestAccIPAMIP_WithStandaloneAddress(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      ipamchecks.CheckIPDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_vpc vpc01 {
						name = "my vpc"
					}

					resource scaleway_vpc_private_network pn01 {
						vpc_id = scaleway_vpc.vpc01.id
						ipv4_subnet {
							subnet = "172.16.32.0/22"
						}
					}

					resource scaleway_ipam_ip ip01 {
					 address = "172.16.32.7"
					 source {
						private_network_id = scaleway_vpc_private_network.pn01.id
					  }
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPAMIPExists(tt, "scaleway_ipam_ip.ip01"),
					resource.TestCheckResourceAttr("scaleway_ipam_ip.ip01", "address", "172.16.32.7/22"),
				),
			},
		},
	})
}

func TestAccIPAMIP_WithTags(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      ipamchecks.CheckIPDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_vpc vpc01 {
						name = "my vpc"
					}

					resource scaleway_vpc_private_network pn01 {
						vpc_id = scaleway_vpc.vpc01.id
						ipv4_subnet {
							subnet = "172.16.32.0/22"
						}
					}

					resource scaleway_ipam_ip ip01 {
					 source {
						private_network_id = scaleway_vpc_private_network.pn01.id
					  }
					  tags = [ "terraform-test", "ipam" ]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPAMIPExists(tt, "scaleway_ipam_ip.ip01"),
					resource.TestCheckResourceAttr("scaleway_ipam_ip.ip01", "tags.#", "2"),
					resource.TestCheckResourceAttr("scaleway_ipam_ip.ip01", "tags.0", "terraform-test"),
					resource.TestCheckResourceAttr("scaleway_ipam_ip.ip01", "tags.1", "ipam"),
				),
			},
			{
				Config: `
					resource scaleway_vpc vpc01 {
						name = "my vpc"
					}

					resource scaleway_vpc_private_network pn01 {
						vpc_id = scaleway_vpc.vpc01.id
						ipv4_subnet {
							subnet = "172.16.32.0/22"
						}
					}

					resource scaleway_ipam_ip ip01 {
					 source {
						private_network_id = scaleway_vpc_private_network.pn01.id
					  }
					  tags = [ "terraform-test", "ipam", "updated" ]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPAMIPExists(tt, "scaleway_ipam_ip.ip01"),
					resource.TestCheckResourceAttr("scaleway_ipam_ip.ip01", "tags.#", "3"),
					resource.TestCheckResourceAttr("scaleway_ipam_ip.ip01", "tags.0", "terraform-test"),
					resource.TestCheckResourceAttr("scaleway_ipam_ip.ip01", "tags.1", "ipam"),
					resource.TestCheckResourceAttr("scaleway_ipam_ip.ip01", "tags.2", "updated"),
				),
			},
		},
	})
}

func TestAccIPAMIP_WithCustomResource(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      ipamchecks.CheckIPDestroy(tt),
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
					  source {
						private_network_id = scaleway_vpc_private_network.pn01.id
					  }
					  custom_resource {
						mac_address = "bc:24:11:74:d0:5a"
					  }
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPAMIPExists(tt, "scaleway_ipam_ip.ip01"),
					resource.TestCheckResourceAttr("scaleway_ipam_ip.ip01", "custom_resource.0.mac_address", "bc:24:11:74:d0:5a"),
				),
			},
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
					  source {
						private_network_id = scaleway_vpc_private_network.pn01.id
					  }
					  custom_resource {
						mac_address = "bc:24:11:74:d0:5b"
					  }
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPAMIPExists(tt, "scaleway_ipam_ip.ip01"),
					resource.TestCheckResourceAttr("scaleway_ipam_ip.ip01", "custom_resource.0.mac_address", "bc:24:11:74:d0:5b"),
				),
			},
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
					  source {
						private_network_id = scaleway_vpc_private_network.pn01.id
					  }
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPAMIPExists(tt, "scaleway_ipam_ip.ip01"),
					resource.TestCheckNoResourceAttr("scaleway_ipam_ip.ip01", "custom_resource.0.mac_address"),
				),
			},
		},
	})
}

func testAccCheckIPAMIPExists(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		ipamAPI, region, ID, err := ipam.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = ipamAPI.GetIP(&ipamSDK.GetIPRequest{
			IPID:   ID,
			Region: region,
		})
		if err != nil {
			return err
		}

		return nil
	}
}
