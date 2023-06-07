package scaleway

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	v1 "github.com/scaleway/scaleway-sdk-go/api/vpc/v1"
	v2 "github.com/scaleway/scaleway-sdk-go/api/vpc/v2"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func init() {
	resource.AddTestSweepers("scaleway_vpc_private_network", &resource.Sweeper{
		Name: "scaleway_vpc_private_network",
		F:    testSweepVPCPrivateNetwork,
	})
}

func testSweepVPCPrivateNetwork(_ string) error {
	err := sweepZones(scw.AllZones, func(scwClient *scw.Client, zone scw.Zone) error {
		vpcAPI := v1.NewAPI(scwClient)
		l.Debugf("sweeper: destroying the private network in (%s)", zone)

		listPNResponse, err := vpcAPI.ListPrivateNetworks(&v1.ListPrivateNetworksRequest{
			Zone: zone,
		}, scw.WithAllPages())
		if err != nil {
			return fmt.Errorf("error listing private network in sweeper: %s", err)
		}

		for _, pn := range listPNResponse.PrivateNetworks {
			err := vpcAPI.DeletePrivateNetwork(&v1.DeletePrivateNetworkRequest{
				Zone:             zone,
				PrivateNetworkID: pn.ID,
			})
			if err != nil {
				return fmt.Errorf("error deleting private network in sweeper: %s", err)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	err = sweepRegions(scw.AllRegions, func(scwClient *scw.Client, region scw.Region) error {
		vpcAPI := v2.NewAPI(scwClient)

		l.Debugf("sweeper: destroying the private network in (%s)", region)

		listPNResponse, err := vpcAPI.ListPrivateNetworks(&v2.ListPrivateNetworksRequest{
			Region: region,
		}, scw.WithAllPages())
		if err != nil {
			return fmt.Errorf("error listing private network in sweeper: %s", err)
		}

		for _, pn := range listPNResponse.PrivateNetworks {
			err := vpcAPI.DeletePrivateNetwork(&v2.DeletePrivateNetworkRequest{
				Region:           region,
				PrivateNetworkID: pn.ID,
			})
			if err != nil {
				return fmt.Errorf("error deleting private network in sweeper: %s", err)
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func TestAccScalewayVPCPrivateNetwork_Basic(t *testing.T) {
	tt := NewTestTools(t)
	defer tt.Cleanup()
	privateNetworkName := "private-network-test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckScalewayVPCPrivateNetworkDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource scaleway_vpc_private_network pn01 {
						name = "%s"
					}
				`, privateNetworkName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayVPCPrivateNetworkExists(
						tt,
						"scaleway_vpc_private_network.pn01",
					),
					resource.TestCheckResourceAttr(
						"scaleway_vpc_private_network.pn01",
						"name",
						privateNetworkName,
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
					testAccCheckScalewayVPCPrivateNetworkExists(
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

func TestAccScalewayVPCPrivateNetwork_DefaultName(t *testing.T) {
	tt := NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckScalewayVPCPrivateNetworkDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `resource scaleway_vpc_private_network main {}`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayVPCPrivateNetworkExists(
						tt,
						"scaleway_vpc_private_network.main",
					),
					resource.TestCheckResourceAttrSet("scaleway_vpc_private_network.main", "name"),
				),
			},
		},
	})
}

func TestAccScalewayVPCPrivateNetwork_Subnets(t *testing.T) {
	tt := NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckScalewayVPCPrivateNetworkDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `resource scaleway_vpc_private_network test {}`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayVPCPrivateNetworkExists(
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
				Config: `resource scaleway_vpc_private_network test {
					ipv4_subnet {
						subnet = "172.16.32.0/22"
					}
					ipv6_subnets {
    				    subnet = "fd46:78ab:30b8:177c::/64"
					}
				}`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayVPCPrivateNetworkExists(
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
						"ipv6_subnets.0.subnet",
						"fd46:78ab:30b8:177c::/64",
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

func TestAccScalewayVPCPrivateNetwork_OneSubnet(t *testing.T) {
	tt := NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckScalewayVPCPrivateNetworkDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `resource scaleway_vpc_private_network test {
   					 ipv4_subnet {
						subnet = "172.16.32.0/22"
					 }
				}`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayVPCPrivateNetworkExists(
						tt,
						"scaleway_vpc_private_network.test",
					),
					resource.TestCheckResourceAttr(
						"scaleway_vpc_private_network.test",
						"ipv4_subnet.0.subnet",
						"172.16.32.0/22",
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

func TestAccScalewayVPCPrivateNetwork_Regional(t *testing.T) {
	tt := NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckScalewayVPCPrivateNetworkDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_vpc_private_network pn01 {
						name = "test-pn"
						is_regional = true
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayVPCPrivateNetworkExists(
						tt,
						"scaleway_vpc_private_network.pn01",
					),
					resource.TestCheckResourceAttrSet(
						"scaleway_vpc_private_network.pn01",
						"vpc_id"),
					resource.TestCheckResourceAttr(
						"scaleway_vpc_private_network.pn01",
						"name",
						"test-pn",
					),
				),
			},
			{
				Config: `
					resource scaleway_vpc_private_network pn01 {
						name = "test-pn"
						is_regional = true
						tags = ["tag0", "tag1"]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayVPCPrivateNetworkExists(
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

func TestAccScalewayVPCPrivateNetwork_RegionalWithTwoIPV6Subnets(t *testing.T) {
	tt := NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckScalewayVPCPrivateNetworkDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_vpc vpc01 {
						name = "test-vpc"
						tags = [ "terraform-test", "vpc", "update" ]
					}
					
					resource scaleway_vpc_private_network pn01 {
						name = "pn1"
						tags = ["tag0", "tag1"]
						vpc_id = scaleway_vpc.vpc01.id
						is_regional = true
						ipv4_subnet {
						  subnet = "192.168.0.0/24"
						}
						ipv6_subnets {
						  subnet = "fd46:78ab:30b8:177c::/64"
						}
						ipv6_subnets {
						  subnet = "fd46:78ab:30b8:c7df::/64"
						}
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayVPCPrivateNetworkExists(
						tt,
						"scaleway_vpc_private_network.pn01",
					),
					resource.TestCheckResourceAttr(
						"scaleway_vpc_private_network.pn01",
						"ipv4_subnet.0.subnet",
						"192.168.0.0/24",
					),
					resource.TestCheckResourceAttrSet(
						"scaleway_vpc_private_network.pn01",
						"ipv6_subnets.0.subnet",
					),
					resource.TestCheckTypeSetElemNestedAttrs(
						"scaleway_vpc_private_network.pn01", "ipv6_subnets.*", map[string]string{
							"subnet": "fd46:78ab:30b8:177c::/64",
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"scaleway_vpc_private_network.pn01", "ipv6_subnets.*", map[string]string{
							"subnet": "fd46:78ab:30b8:c7df::/64",
						}),
					resource.TestCheckResourceAttr(
						"scaleway_vpc_private_network.pn01",
						"ipv4_subnet.#",
						"1",
					),
					resource.TestCheckResourceAttr(
						"scaleway_vpc_private_network.pn01",
						"ipv6_subnets.#",
						"2",
					),
				),
			},
		},
	})
}

func testAccCheckScalewayVPCPrivateNetworkExists(tt *TestTools, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		if rs.Primary.Attributes["is_regional"] == "true" {
			vpcAPI, region, ID, err := vpcAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}

			_, err = vpcAPI.GetPrivateNetwork(&v2.GetPrivateNetworkRequest{
				PrivateNetworkID: ID,
				Region:           region,
			})
			if err != nil {
				return err
			}
		} else {
			vpcAPI, zone, ID, err := vpcAPIWithZoneAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}

			_, err = vpcAPI.GetPrivateNetwork(&v1.GetPrivateNetworkRequest{
				PrivateNetworkID: ID,
				Zone:             zone,
			})

			if err != nil {
				return err
			}
		}

		return nil
	}
}

func testAccCheckScalewayVPCPrivateNetworkDestroy(tt *TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_vpc_private_network" {
				continue
			}

			if rs.Primary.Attributes["is_regional"] == "true" {
				vpcAPI, region, ID, err := vpcAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
				if err != nil {
					return err
				}
				_, err = vpcAPI.GetPrivateNetwork(&v2.GetPrivateNetworkRequest{
					PrivateNetworkID: ID,
					Region:           region,
				})

				if err == nil {
					return fmt.Errorf(
						"VPC private network %s still exists",
						rs.Primary.ID,
					)
				}
				// Unexpected api error we return it
				if !is404Error(err) {
					return err
				}
			} else {
				vpcAPI, zone, ID, err := vpcAPIWithZoneAndID(tt.Meta, rs.Primary.ID)
				if err != nil {
					return err
				}

				_, err = vpcAPI.GetPrivateNetwork(&v1.GetPrivateNetworkRequest{
					PrivateNetworkID: ID,
					Zone:             zone,
				})

				if err == nil {
					return fmt.Errorf(
						"VPC private network %s still exists",
						rs.Primary.ID,
					)
				}

				// Unexpected api error we return it
				if !is404Error(err) {
					return err
				}
			}
		}

		return nil
	}
}
