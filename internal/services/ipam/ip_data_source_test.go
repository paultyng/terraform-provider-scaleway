package ipam_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	instancechecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/instance/testfuncs"
	ipamchecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/ipam/testfuncs"
	rdbchecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/rdb/testfuncs"
)

func TestAccDataSourceIPAMIP_Instance(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      instancechecks.IsServerDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_vpc" "main" {
						name = "tf-tests-ipam-ip-datasource-instance"
					}

					resource "scaleway_vpc_private_network" "main" {
						vpc_id = scaleway_vpc.main.id
						name = "tf-tests-ipam-ip-datasource-instance"
					}

					resource "scaleway_instance_server" "main" {
						name  = "tf-tests-ipam-ip-datasource-instance"
						image = "ubuntu_jammy"
						type  = "PLAY2-MICRO"
						tags  = [ "terraform-test", "data_scaleway_instance_servers", "basic" ]
					}

					resource "scaleway_instance_private_nic" "main" {
						private_network_id = scaleway_vpc_private_network.main.id
						server_id = scaleway_instance_server.main.id
					}

					data "scaleway_ipam_ip" "by_mac" {
						mac_address = scaleway_instance_private_nic.main.mac_address
						type = "ipv4"
					}

					data "scaleway_ipam_ip" "by_id" {
						resource {
							id = scaleway_instance_private_nic.main.id
							type = "instance_private_nic"
						}
						type = "ipv4"
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.scaleway_ipam_ip.by_mac", "address"),
					resource.TestCheckResourceAttrSet("data.scaleway_ipam_ip.by_id", "address"),
					resource.TestCheckResourceAttrPair("data.scaleway_ipam_ip.by_mac", "address", "data.scaleway_ipam_ip.by_id", "address"),
				),
			},
		},
	})
}

func TestAccDataSourceIPAMIP_InstanceLB(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      instancechecks.IsServerDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_vpc" "main" {
						name = "tf-tests-ipam-ip-datasource-instance"
					}

					resource "scaleway_vpc_private_network" "main" {
						vpc_id = scaleway_vpc.main.id
						name = "tf-tests-ipam-ip-datasource-instance"
					}

					resource "scaleway_instance_server" "main" {
						name  = "tf-tests-ipam-ip-datasource-instance"
						image = "ubuntu_jammy"
						type  = "PLAY2-MICRO"
						tags  = [ "terraform-test", "data_scaleway_instance_servers", "basic" ]
					}

					resource "scaleway_instance_private_nic" "main" {
						private_network_id = scaleway_vpc_private_network.main.id
						server_id = scaleway_instance_server.main.id
					}

					data "scaleway_ipam_ip" "main" {
						mac_address = scaleway_instance_private_nic.main.mac_address
						type = "ipv4"
					}

					resource "scaleway_lb_ip" "main" {}

					resource "scaleway_lb" "main" {
						ip_id = scaleway_lb_ip.main.id
						type = "LB-S"
					}
					
					resource "scaleway_lb_backend" "main" {
						lb_id = scaleway_lb.main.id
						forward_protocol = "http"
						forward_port = "80"
						server_ips = [data.scaleway_ipam_ip.main.address]
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.scaleway_ipam_ip.main", "address"),
				),
			},
		},
	})
}

func TestAccDataSourceIPAMIP_RDB(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      rdbchecks.IsInstanceDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_vpc" "main" {
						name = "tf-tests-ipam-ip-datasource-rdb"
					}

					resource "scaleway_vpc_private_network" "main" {
						vpc_id = scaleway_vpc.main.id
						name = "tf-tests-ipam-ip-datasource-rdb"
						ipv4_subnet {
							subnet = "172.16.0.0/22"
						}
					}

					resource scaleway_rdb_instance main {
						name = "test-ipam-ip-rdb"
						node_type = "db-dev-s"
						engine = "PostgreSQL-14"
						is_ha_cluster = false
						disable_backup = true
						user_name = "my_initial_user"
						password = "thiZ_is_v&ry_s3cret"
						tags = [ "terraform-test", "scaleway_rdb_instance", "volume", "rdb_pn" ]
						volume_type = "bssd"
						volume_size_in_gb = 10
						private_network {
							pn_id = "${scaleway_vpc_private_network.main.id}"
							enable_ipam = true
						}
					}

					data "scaleway_ipam_ip" "main" {
						resource {
							name = scaleway_rdb_instance.main.name
							type = "rdb_instance"
						}
						type = "ipv4"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.scaleway_ipam_ip.main", "address"),
				),
			},
		},
	})
}

func TestAccDataSourceIPAMIP_ID(t *testing.T) {
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
					  address = "172.16.32.5"
					  source {
						private_network_id = scaleway_vpc_private_network.pn01.id
					  }
					}
					
					data "scaleway_ipam_ip" "by_id" {
					  ipam_ip_id = scaleway_ipam_ip.ip01.id
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.scaleway_ipam_ip.by_id", "address", "172.16.32.5"),
					resource.TestCheckResourceAttr("data.scaleway_ipam_ip.by_id", "address_cidr", "172.16.32.5/22"),
					resource.TestCheckResourceAttrPair("data.scaleway_ipam_ip.by_id", "ipam_ip_id", "scaleway_ipam_ip.ip01", "id"),
				),
			},
		},
	})
}
