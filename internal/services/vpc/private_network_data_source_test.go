package vpc_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	vpcchecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/vpc/testfuncs"
)

func TestAccDataSourceVPCPrivateNetwork_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	pnName := "TestAccScalewayDataSourceVPCPrivateNetwork_Basic"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      vpcchecks.CheckPrivateNetworkDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "scaleway_vpc_private_network" "pn_test" {
					  name = "%s"
					}`, pnName),
			},
			{
				Config: fmt.Sprintf(`
					resource "scaleway_vpc_private_network" "pn_test" {
					  name = "%s"
					}

					data "scaleway_vpc_private_network" "pn_test_by_name" {
						name = "${scaleway_vpc_private_network.pn_test.name}"
					}

					data "scaleway_vpc_private_network" "pn_test_by_id" {
						private_network_id = "${scaleway_vpc_private_network.pn_test.id}"
					}
				`, pnName),
				Check: resource.ComposeTestCheckFunc(
					vpcchecks.IsPrivateNetworkPresent(tt, "scaleway_vpc_private_network.pn_test"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_vpc_private_network.pn_test_by_name", "name",
						"scaleway_vpc_private_network.pn_test", "name"),
					resource.TestCheckResourceAttr(
						"data.scaleway_vpc_private_network.pn_test_by_name",
						"ipv4_subnet.#", "1"),
					resource.TestCheckResourceAttr(
						"data.scaleway_vpc_private_network.pn_test_by_name",
						"ipv6_subnets.#", "1"),
					resource.TestCheckTypeSetElemAttrPair(
						"data.scaleway_vpc_private_network.pn_test_by_name", "ipv4_subnet.0.subnet",
						"scaleway_vpc_private_network.pn_test", "ipv4_subnet.0.subnet"),
					resource.TestCheckTypeSetElemAttrPair(
						"data.scaleway_vpc_private_network.pn_test_by_name", "ipv6_subnets.0.subnet",
						"scaleway_vpc_private_network.pn_test", "ipv6_subnets.0.subnet"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_vpc_private_network.pn_test_by_id", "private_network_id",
						"scaleway_vpc_private_network.pn_test", "id"),
					resource.TestCheckResourceAttr(
						"data.scaleway_vpc_private_network.pn_test_by_id",
						"ipv4_subnet.#", "1"),
					resource.TestCheckResourceAttr(
						"data.scaleway_vpc_private_network.pn_test_by_id",
						"ipv6_subnets.#", "1"),
					resource.TestCheckTypeSetElemAttrPair(
						"data.scaleway_vpc_private_network.pn_test_by_id", "ipv4_subnet.0.subnet",
						"scaleway_vpc_private_network.pn_test", "ipv4_subnet.0.subnet"),
					resource.TestCheckTypeSetElemAttrPair(
						"data.scaleway_vpc_private_network.pn_test_by_id", "ipv6_subnets.0.subnet",
						"scaleway_vpc_private_network.pn_test", "ipv6_subnets.0.subnet"),
				),
			},
		},
	})
}

func TestAccDataSourceVPCPrivateNetwork_VpcID(t *testing.T) {
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
					  name = "TestAccScalewayResourceVPC_Basic01"
					}

					resource "scaleway_vpc" "vpc02" {
					  name = "TestAccScalewayResourceVPC_Basic02"
					}
					
					resource "scaleway_vpc_private_network" "pn01" {
					  name = "TestAccScalewayResourceVPCPrivateNetwork_Basic"
					  vpc_id = scaleway_vpc.vpc01.id
					}

					resource "scaleway_vpc_private_network" "pn02" {
					  name = "TestAccScalewayResourceVPCPrivateNetwork_Basic"
					  vpc_id = scaleway_vpc.vpc02.id
					}
				`,
			},
			{
				Config: `
					resource "scaleway_vpc" "vpc01" {
					  name = "TestAccScalewayResourceVPC_Basic01"
					}

					resource "scaleway_vpc" "vpc02" {
					  name = "TestAccScalewayResourceVPC_Basic02"
					}
					
					resource "scaleway_vpc_private_network" "pn01" {
					  name = "TestAccScalewayResourceVPCPrivateNetwork_Basic"
					  vpc_id = scaleway_vpc.vpc01.id
					}

					resource "scaleway_vpc_private_network" "pn02" {
					  name = "TestAccScalewayResourceVPCPrivateNetwork_Basic"
					  vpc_id = scaleway_vpc.vpc02.id
					}

					data "scaleway_vpc_private_network" "by_vpc_id" {
						name = "${scaleway_vpc_private_network.pn01.name}"
						vpc_id = "${scaleway_vpc.vpc01.id}"
					}

					data "scaleway_vpc_private_network" "by_vpc_id_2" {
						name = "${scaleway_vpc_private_network.pn02.name}"
						vpc_id = "${scaleway_vpc.vpc02.id}"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					vpcchecks.IsPrivateNetworkPresent(tt, "scaleway_vpc_private_network.pn01"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_vpc_private_network.by_vpc_id", "name",
						"scaleway_vpc_private_network.pn01", "name"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_vpc_private_network.by_vpc_id", "vpc_id",
						"scaleway_vpc_private_network.pn01", "vpc_id"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_vpc_private_network.by_vpc_id", "vpc_id",
						"scaleway_vpc.vpc01", "id"),

					vpcchecks.IsPrivateNetworkPresent(tt, "scaleway_vpc_private_network.pn02"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_vpc_private_network.by_vpc_id_2", "name",
						"scaleway_vpc_private_network.pn02", "name"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_vpc_private_network.by_vpc_id_2", "vpc_id",
						"scaleway_vpc_private_network.pn02", "vpc_id"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_vpc_private_network.by_vpc_id_2", "vpc_id",
						"scaleway_vpc.vpc02", "id"),
				),
			},
		},
	})
}
