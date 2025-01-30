package ipam_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	ipamchecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/ipam/testfuncs"
)

func TestAccDataSourceIPAMIPs_Basic(t *testing.T) {
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
					  tags = ["terraform-test", "data_scaleway_ipam_ips"]
					}
					
					resource "scaleway_ipam_ip" "ip02" {
					  source {
						private_network_id = scaleway_vpc_private_network.pn01.id
					  }
					  tags = ["terraform-test", "data_scaleway_ipam_ips"]
					}
				`,
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
					  tags = ["terraform-test", "data_scaleway_ipam_ips"]
					}
					
					resource "scaleway_ipam_ip" "ip02" {
					  source {
						private_network_id = scaleway_vpc_private_network.pn01.id
					  }
					  tags = ["terraform-test", "data_scaleway_ipam_ips"]
					}
					
					data "scaleway_ipam_ips" "by_tag" {
					  tags = ["terraform-test", "data_scaleway_ipam_ips"]
					}
	
					data "scaleway_ipam_ips" "by_tag_other_zone" {
					  tags = ["terraform-test", "data_scaleway_ipam_ips"]
					  zonal = "fr-par-2"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.scaleway_ipam_ips.by_tag", "ips.0.id"),
					resource.TestCheckResourceAttrSet("data.scaleway_ipam_ips.by_tag", "ips.0.address"),
					resource.TestCheckResourceAttrSet("data.scaleway_ipam_ips.by_tag", "ips.1.id"),
					resource.TestCheckResourceAttrSet("data.scaleway_ipam_ips.by_tag", "ips.1.address"),

					resource.TestCheckNoResourceAttr("data.scaleway_ipam_ips.by_tag_other_zone", "ips.0.id"),
				),
			},
		},
	})
}

func TestAccDataSourceIPAMIPs_RedisCluster(t *testing.T) {
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
					
					resource "scaleway_redis_cluster" "redis01" {
					  name         = "test_redis_endpoints"
					  version      = "7.0.5"
					  node_type    = "RED1-XS"
					  user_name    = "my_initial_user"
					  password     = "thiZ_is_v&ry_s3cret"
					  cluster_size = 3
					  private_network {
						id = scaleway_vpc_private_network.pn01.id
					  }
					}
				`,
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
					
					resource "scaleway_redis_cluster" "redis01" {
					  name         = "test_redis_endpoints"
					  version      = "7.0.5"
					  node_type    = "RED1-XS"
					  user_name    = "my_initial_user"
					  password     = "thiZ_is_v&ry_s3cret"
					  cluster_size = 3
					  private_network {
						id = scaleway_vpc_private_network.pn01.id
					  }
					}
					
					data "scaleway_ipam_ips" "by_resource_and_type" {
					  type = "ipv4"
					  resource {
						id   = scaleway_redis_cluster.redis01.id
						type = "redis_cluster"
					  }
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.scaleway_ipam_ips.by_resource_and_type", "ips.#", "3"),
					resource.TestCheckResourceAttrSet("data.scaleway_ipam_ips.by_resource_and_type", "ips.0.id"),
					resource.TestCheckResourceAttrSet("data.scaleway_ipam_ips.by_resource_and_type", "ips.0.address"),
					resource.TestCheckResourceAttrSet("data.scaleway_ipam_ips.by_resource_and_type", "ips.1.id"),
					resource.TestCheckResourceAttrSet("data.scaleway_ipam_ips.by_resource_and_type", "ips.1.address"),
					resource.TestCheckResourceAttrSet("data.scaleway_ipam_ips.by_resource_and_type", "ips.2.id"),
					resource.TestCheckResourceAttrSet("data.scaleway_ipam_ips.by_resource_and_type", "ips.2.address"),
				),
			},
		},
	})
}
