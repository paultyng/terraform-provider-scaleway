package instance_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

func TestAccDataSourcePrivateNIC_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			isPrivateNICDestroyed(tt),
		),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_instance_server" "server" {
						name = "test-terraform-datasource-private-nic"
						type = "DEV1-S"
						image = "ubuntu_jammy"
					}
					resource "scaleway_vpc_private_network" "vpc" {}

					resource "scaleway_instance_private_nic" "nic" {
						server_id = scaleway_instance_server.server.id
						private_network_id = scaleway_vpc_private_network.vpc.id
						tags = ["test-terraform-datasource-private-nic"]
					}`,
			},
			{
				Config: `
					resource "scaleway_instance_server" "server" {
						name = "test-terraform-datasource-private-nic"
						type = "DEV1-S"
						image = "ubuntu_jammy"
					}
					resource "scaleway_vpc_private_network" "vpc" {}

					resource "scaleway_instance_private_nic" "nic" {
						server_id = scaleway_instance_server.server.id
						private_network_id = scaleway_vpc_private_network.vpc.id
						tags = ["test-terraform-datasource-private-nic"]
					}

					data scaleway_instance_private_nic find_by_nic_id {
						server_id = scaleway_instance_server.server.id
						private_nic_id = split("/", scaleway_instance_private_nic.nic.id)[2]
					}

					data scaleway_instance_private_nic find_by_vpc_id {
						server_id = scaleway_instance_server.server.id
						private_network_id = scaleway_vpc_private_network.vpc.id
					}

					data scaleway_instance_private_nic find_by_tags {
						server_id = scaleway_instance_server.server.id
						tags = ["test-terraform-datasource-private-nic"]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isPrivateNICPresent(tt, "scaleway_instance_private_nic.nic"),

					resource.TestCheckResourceAttrPair("scaleway_instance_private_nic.nic", "id", "data.scaleway_instance_private_nic.find_by_nic_id", "id"),
					resource.TestCheckResourceAttrPair("scaleway_instance_private_nic.nic", "id", "data.scaleway_instance_private_nic.find_by_vpc_id", "id"),
					resource.TestCheckResourceAttrPair("scaleway_instance_private_nic.nic", "id", "data.scaleway_instance_private_nic.find_by_tags", "id"),
					resource.TestCheckResourceAttrPair("data.scaleway_instance_private_nic.find_by_nic_id", "id", "data.scaleway_instance_private_nic.find_by_vpc_id", "id"),
					resource.TestCheckResourceAttrSet("data.scaleway_instance_private_nic.find_by_nic_id", "mac_address"),
					resource.TestCheckResourceAttrSet("data.scaleway_instance_private_nic.find_by_vpc_id", "mac_address"),
					resource.TestCheckResourceAttrSet("data.scaleway_instance_private_nic.find_by_tags", "mac_address"),
				),
			},
		},
	})
}
