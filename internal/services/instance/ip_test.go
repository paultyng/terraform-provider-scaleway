package instance_test

import (
	"fmt"
	"net"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	instanceSDK "github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/locality"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/logging"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/instance"
	instancechecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/instance/testfuncs"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

func init() {
	resource.AddTestSweepers("scaleway_instance_ip", &resource.Sweeper{
		Name: "scaleway_instance_ip",
		F:    testSweepIP,
	})
}

func testSweepIP(_ string) error {
	return acctest.SweepZones(scw.AllZones, func(scwClient *scw.Client, zone scw.Zone) error {
		instanceAPI := instanceSDK.NewAPI(scwClient)

		listIPs, err := instanceAPI.ListIPs(&instanceSDK.ListIPsRequest{Zone: zone}, scw.WithAllPages())
		if err != nil {
			logging.L.Warningf("error listing ips in (%s) in sweeper: %s", zone, err)
			return nil
		}

		for _, ip := range listIPs.IPs {
			err := instanceAPI.DeleteIP(&instanceSDK.DeleteIPRequest{
				IP:   ip.ID,
				Zone: zone,
			})
			if err != nil {
				return fmt.Errorf("error deleting ip in sweeper: %s", err)
			}
		}

		return nil
	})
}

func TestAccIP_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      instancechecks.IsIPDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
						resource "scaleway_instance_ip" "base" {}
						resource "scaleway_instance_ip" "scaleway" {}
					`,
				Check: resource.ComposeTestCheckFunc(
					instancechecks.CheckIPExists(tt, "scaleway_instance_ip.base"),
					instancechecks.CheckIPExists(tt, "scaleway_instance_ip.scaleway"),
				),
			},
		},
	})
}

func TestAccIP_WithZone(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      instancechecks.IsIPDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
						resource "scaleway_instance_ip" "base" {}
					`,
				Check: resource.ComposeTestCheckFunc(
					instancechecks.CheckIPExists(tt, "scaleway_instance_ip.base"),
					resource.TestCheckResourceAttr("scaleway_instance_ip.base", "zone", "fr-par-1"),
				),
			},
			{
				Config: `
						resource "scaleway_instance_ip" "base" {
							zone = "nl-ams-1"
						}
					`,
				Check: resource.ComposeTestCheckFunc(
					instancechecks.CheckIPExists(tt, "scaleway_instance_ip.base"),
					resource.TestCheckResourceAttr("scaleway_instance_ip.base", "zone", "nl-ams-1"),
				),
			},
		},
	})
}

func TestAccIP_Tags(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      instancechecks.IsIPDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
						resource "scaleway_instance_ip" "main" {}
					`,
				Check: resource.ComposeTestCheckFunc(
					instancechecks.CheckIPExists(tt, "scaleway_instance_ip.main"),
					resource.TestCheckNoResourceAttr("scaleway_instance_ip.main", "tags"),
				),
			},
			{
				Config: `
						resource "scaleway_instance_ip" "main" {
							tags = ["foo", "bar"]
						}
					`,
				Check: resource.ComposeTestCheckFunc(
					instancechecks.CheckIPExists(tt, "scaleway_instance_ip.main"),
					resource.TestCheckResourceAttr("scaleway_instance_ip.main", "tags.0", "foo"),
					resource.TestCheckResourceAttr("scaleway_instance_ip.main", "tags.1", "bar"),
				),
			},
		},
	})
}

func TestAccIP_RoutedMigrate(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      instancechecks.IsIPDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
						resource "scaleway_instance_ip" "main" {
						}
					`,
				Check: resource.ComposeTestCheckFunc(
					instancechecks.CheckIPExists(tt, "scaleway_instance_ip.main"),
					resource.TestCheckResourceAttr("scaleway_instance_ip.main", "type", "nat"),
				),
			},
			{
				Config: `
						resource "scaleway_instance_ip" "main" {
						}
						resource "scaleway_instance_ip" "copy" {
						}
					`,
				Check: resource.ComposeTestCheckFunc(
					instancechecks.CheckIPExists(tt, "scaleway_instance_ip.main"),
					instancechecks.CheckIPExists(tt, "scaleway_instance_ip.copy"),
					resource.TestCheckResourceAttr("scaleway_instance_ip.main", "type", "nat"),
					resource.TestCheckResourceAttr("scaleway_instance_ip.copy", "type", "nat"),
					resource.TestCheckResourceAttrPair("scaleway_instance_ip.main", "id", "scaleway_instance_ip.copy", "id"),
				),
				ResourceName: "scaleway_instance_ip.copy",
				ImportState:  true,
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					return state.RootModule().Resources["scaleway_instance_ip.main"].Primary.ID, nil
				},
				ImportStatePersist: true,
			},
			{
				Config: `
						resource "scaleway_instance_ip" "main" {
							type = "routed_ipv4"
						}
						resource "scaleway_instance_ip" "copy" {
						}
					`,
				Check: resource.ComposeTestCheckFunc(
					instancechecks.CheckIPExists(tt, "scaleway_instance_ip.main"),
					instancechecks.CheckIPExists(tt, "scaleway_instance_ip.copy"),
					resource.TestCheckResourceAttr("scaleway_instance_ip.main", "type", "routed_ipv4"),
					resource.TestCheckResourceAttrPair("scaleway_instance_ip.main", "id", "scaleway_instance_ip.copy", "id"),
				),
			},
			{
				// After the main IP migrated, we check that there is no ForceNew on the copy
				// This check that the ip is not deleted if the migration is done outside terraform
				RefreshState: true,
				Check: resource.ComposeTestCheckFunc(
					instancechecks.CheckIPExists(tt, "scaleway_instance_ip.main"),
					instancechecks.CheckIPExists(tt, "scaleway_instance_ip.copy"),
					resource.TestCheckResourceAttr("scaleway_instance_ip.main", "type", "routed_ipv4"),
					resource.TestCheckResourceAttr("scaleway_instance_ip.copy", "type", "routed_ipv4"),
					resource.TestCheckResourceAttrPair("scaleway_instance_ip.main", "id", "scaleway_instance_ip.copy", "id"),
				),
			},
		},
	})
}

func TestAccIP_RoutedDowngrade(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      instancechecks.IsIPDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
						resource "scaleway_instance_ip" "main" {
							type = "routed_ipv4"
						}
					`,
				Check: resource.ComposeTestCheckFunc(
					instancechecks.CheckIPExists(tt, "scaleway_instance_ip.main"),
					resource.TestCheckResourceAttr("scaleway_instance_ip.main", "type", "routed_ipv4"),
					isIPValid("scaleway_instance_ip.main", "address"),
					isIPCIDRValid("scaleway_instance_ip.main", "prefix"),
				),
			},
			{
				Config: `
						resource "scaleway_instance_ip" "main" {
							type = "nat"
						}
					`,
				Check: resource.ComposeTestCheckFunc(
					instancechecks.CheckIPExists(tt, "scaleway_instance_ip.main"),
					resource.TestCheckResourceAttr("scaleway_instance_ip.main", "type", "nat"),
					isIPValid("scaleway_instance_ip.main", "address"),
					isIPCIDRValid("scaleway_instance_ip.main", "prefix"),
				),
			},
		},
	})
}

func TestAccIP_RoutedIPV6(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      instancechecks.IsIPDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
						resource "scaleway_instance_ip" "main" {
							type = "routed_ipv6"
						}
					`,
				Check: resource.ComposeTestCheckFunc(
					instancechecks.CheckIPExists(tt, "scaleway_instance_ip.main"),
					resource.TestCheckResourceAttr("scaleway_instance_ip.main", "type", "routed_ipv6"),
					resource.TestCheckResourceAttrSet("scaleway_instance_ip.main", "address"),
					resource.TestCheckResourceAttrSet("scaleway_instance_ip.main", "prefix"),
					isIPValid("scaleway_instance_ip.main", "address"),
					isIPCIDRValid("scaleway_instance_ip.main", "prefix"),
				),
			},
		},
	})
}

func isIPCIDRValid(name string, key string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("resource not found: %s", name)
		}

		cidr, exists := rs.Primary.Attributes[key]
		if !exists {
			return fmt.Errorf("requested attribute %s[%q] does not exist", name, key)
		}
		_, _, err := net.ParseCIDR(cidr)
		if err != nil {
			return fmt.Errorf("invalid cidr (%s) in %s[%q]", cidr, name, key)
		}

		return nil
	}
}

func isIPValid(name string, key string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("resource not found: %s", name)
		}

		ip, exists := rs.Primary.Attributes[key]
		if !exists {
			return fmt.Errorf("requested attribute %s[%q] does not exist", name, key)
		}

		parsedIP := net.ParseIP(ip)
		if parsedIP == nil {
			return fmt.Errorf("invalid ip (%s) in %s[%q]", ip, name, key)
		}

		return nil
	}
}

func isIPAttachedToServer(tt *acctest.TestTools, ipResource, serverResource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ipState, ok := s.RootModule().Resources[ipResource]
		if !ok {
			return fmt.Errorf("resource not found: %s", ipResource)
		}
		serverState, ok := s.RootModule().Resources[serverResource]
		if !ok {
			return fmt.Errorf("resource not found: %s", serverResource)
		}

		instanceAPI, zone, ID, err := instance.NewAPIWithZoneAndID(tt.Meta, ipState.Primary.ID)
		if err != nil {
			return err
		}

		server, err := instanceAPI.GetServer(&instanceSDK.GetServerRequest{
			Zone:     zone,
			ServerID: locality.ExpandID(serverState.Primary.ID),
		})
		if err != nil {
			return err
		}

		ip, err := instanceAPI.GetIP(&instanceSDK.GetIPRequest{
			IP:   ID,
			Zone: zone,
		})
		if err != nil {
			return err
		}

		if server.Server.PublicIP.Address.String() != ip.IP.Address.String() {
			return fmt.Errorf("IPs should be the same in %s and %s: %v is different than %v", ipResource, serverResource, server.Server.PublicIP.Address, ip.IP.Address)
		}

		return nil
	}
}

func serverHasNoIPAssigned(tt *acctest.TestTools, serverResource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		serverState, ok := s.RootModule().Resources[serverResource]
		if !ok {
			return fmt.Errorf("resource not found: %s", serverResource)
		}

		instanceAPI, zone, ID, err := instance.NewAPIWithZoneAndID(tt.Meta, serverState.Primary.ID)
		if err != nil {
			return err
		}

		server, err := instanceAPI.GetServer(&instanceSDK.GetServerRequest{
			Zone:     zone,
			ServerID: ID,
		})
		if err != nil {
			return err
		}

		if server.Server.PublicIP != nil && !server.Server.PublicIP.Dynamic {
			return fmt.Errorf("no flexible IP should be assigned to %s", serverResource)
		}

		return nil
	}
}
