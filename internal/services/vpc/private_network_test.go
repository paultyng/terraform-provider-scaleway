package vpc_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	vpcchecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/vpc/testfuncs"
)

func TestAccVPCPrivateNetwork_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	privateNetworkName := "private-network-test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      vpcchecks.CheckPrivateNetworkDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource scaleway_vpc_private_network pn01 {
						name = "%s"
					}
				`, privateNetworkName),
				Check: resource.ComposeTestCheckFunc(
					vpcchecks.IsPrivateNetworkPresent(
						tt,
						"scaleway_vpc_private_network.pn01",
					),
					resource.TestCheckResourceAttrSet(
						"scaleway_vpc_private_network.pn01",
						"vpc_id"),
					resource.TestCheckResourceAttr(
						"scaleway_vpc_private_network.pn01",
						"name",
						privateNetworkName,
					),
					resource.TestCheckResourceAttr(
						"scaleway_vpc_private_network.pn01",
						"region",
						"fr-par",
					),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource scaleway_vpc_private_network pn01 {
						name = "%s"
						tags = ["tag0", "tag1"]
					}
				`, privateNetworkName),
				Check: resource.ComposeTestCheckFunc(
					vpcchecks.IsPrivateNetworkPresent(
						tt,
						"scaleway_vpc_private_network.pn01",
					),
					resource.TestCheckResourceAttr(
						"scaleway_vpc_private_network.pn01",
						"tags.0",
						"tag0",
					),
					resource.TestCheckResourceAttr(
						"scaleway_vpc_private_network.pn01",
						"tags.1",
						"tag1",
					),
				),
			},
		},
	})
}

func TestAccVPCPrivateNetwork_DefaultName(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      vpcchecks.CheckPrivateNetworkDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `resource scaleway_vpc_private_network main {}`,
				Check: resource.ComposeTestCheckFunc(
					vpcchecks.IsPrivateNetworkPresent(
						tt,
						"scaleway_vpc_private_network.main",
					),
					resource.TestCheckResourceAttrSet("scaleway_vpc_private_network.main", "name"),
				),
			},
		},
	})
}

func TestAccVPCPrivateNetwork_Subnets(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      vpcchecks.CheckPrivateNetworkDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_vpc vpc01 {
						name = "my vpc"
					}

					resource scaleway_vpc_private_network test {
						vpc_id = scaleway_vpc.vpc01.id
					}`,
				Check: resource.ComposeTestCheckFunc(
					vpcchecks.IsPrivateNetworkPresent(
						tt,
						"scaleway_vpc_private_network.test",
					),
					resource.TestCheckResourceAttr(
						"scaleway_vpc_private_network.test",
						"ipv4_subnet.#",
						"1",
					),
					resource.TestCheckResourceAttr(
						"scaleway_vpc_private_network.test",
						"ipv6_subnets.#",
						"1",
					),
				),
			},
			{
				Config: `
					resource scaleway_vpc vpc01 {
						name = "my vpc"
					}
					
					resource scaleway_vpc_private_network test {
						ipv4_subnet {
							subnet = "172.16.32.0/22"
						}
						ipv6_subnets {
							subnet = "fd46:78ab:30b8:177c::/64"
						}
						vpc_id = scaleway_vpc.vpc01.id
					}`,
				Check: resource.ComposeTestCheckFunc(
					vpcchecks.IsPrivateNetworkPresent(
						tt,
						"scaleway_vpc_private_network.test",
					),
					resource.TestCheckResourceAttr(
						"scaleway_vpc_private_network.test",
						"ipv4_subnet.0.subnet",
						"172.16.32.0/22",
					),
					resource.TestCheckResourceAttr(
						"scaleway_vpc_private_network.test",
						"ipv4_subnet.0.address",
						"172.16.32.0",
					),
					resource.TestCheckResourceAttr(
						"scaleway_vpc_private_network.test",
						"ipv4_subnet.0.subnet_mask",
						"255.255.252.0",
					),
					resource.TestCheckResourceAttr(
						"scaleway_vpc_private_network.test",
						"ipv4_subnet.0.prefix_length",
						"22",
					),
					resource.TestCheckResourceAttr(
						"scaleway_vpc_private_network.test",
						"ipv6_subnets.0.subnet",
						"fd46:78ab:30b8:177c::/64",
					),
					resource.TestCheckResourceAttr(
						"scaleway_vpc_private_network.test",
						"ipv6_subnets.0.address",
						"fd46:78ab:30b8:177c::",
					),
					resource.TestCheckResourceAttr(
						"scaleway_vpc_private_network.test",
						"ipv6_subnets.0.prefix_length",
						"64",
					),
					resource.TestCheckResourceAttr(
						"scaleway_vpc_private_network.test",
						"ipv4_subnet.#",
						"1",
					),
					resource.TestCheckResourceAttr(
						"scaleway_vpc_private_network.test",
						"ipv6_subnets.#",
						"1",
					),
				),
			},
		},
	})
}

func TestAccVPCPrivateNetwork_OneSubnet(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      vpcchecks.CheckPrivateNetworkDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_vpc vpc01 {
						name = "my vpc"
					}

					resource scaleway_vpc_private_network test {
						ipv4_subnet {
							subnet = "172.16.64.0/22"
						}
						vpc_id = scaleway_vpc.vpc01.id
					}`,
				Check: resource.ComposeTestCheckFunc(
					vpcchecks.IsPrivateNetworkPresent(
						tt,
						"scaleway_vpc_private_network.test",
					),
					resource.TestCheckResourceAttr(
						"scaleway_vpc_private_network.test",
						"ipv4_subnet.0.subnet",
						"172.16.64.0/22",
					),
					resource.TestCheckResourceAttrSet(
						"scaleway_vpc_private_network.test",
						"ipv6_subnets.0.subnet",
					),
					resource.TestCheckResourceAttr(
						"scaleway_vpc_private_network.test",
						"ipv4_subnet.#",
						"1",
					),
					resource.TestCheckResourceAttr(
						"scaleway_vpc_private_network.test",
						"ipv6_subnets.#",
						"1",
					),
				),
			},
		},
	})
}
