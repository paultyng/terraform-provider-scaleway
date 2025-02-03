package lb_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	lbSDK "github.com/scaleway/scaleway-sdk-go/api/lb/v1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/lb"
	lbchecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/lb/testfuncs"
)

func TestAccIP_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      lbchecks.IsIPDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_lb_ip ipZone {
						zone = "nl-ams-1"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isIPPresent(tt, "scaleway_lb_ip.ipZone"),
					acctest.CheckResourceAttrIPv4("scaleway_lb_ip.ipZone", "ip_address"),
					resource.TestCheckResourceAttrSet("scaleway_lb_ip.ipZone", "reverse"),
					resource.TestCheckResourceAttr("scaleway_lb_ip.ipZone", "is_ipv6", "false"),
					resource.TestCheckResourceAttr("scaleway_lb_ip.ipZone", "zone", "nl-ams-1"),
				),
			},
			{
				Config: `
					resource scaleway_lb_ip ip01 {
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isIPPresent(tt, "scaleway_lb_ip.ip01"),
					acctest.CheckResourceAttrIPv4("scaleway_lb_ip.ip01", "ip_address"),
					resource.TestCheckResourceAttrSet("scaleway_lb_ip.ip01", "reverse"),
					resource.TestCheckResourceAttr("scaleway_lb_ip.ip01", "zone", "fr-par-1"),
				),
			},
			{
				Config: `
					resource scaleway_lb_ip ip01 {
						reverse = "myreverse.com"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isIPPresent(tt, "scaleway_lb_ip.ip01"),
					acctest.CheckResourceAttrIPv4("scaleway_lb_ip.ip01", "ip_address"),
					resource.TestCheckResourceAttr("scaleway_lb_ip.ip01", "reverse", "myreverse.com"),
				),
			},
			{
				Config: `
					resource scaleway_lb_ip ip01 {
						reverse = "myreverse.com"
					}

					resource scaleway_lb main {
					    ip_id = scaleway_lb_ip.ip01.id
						name = "test-lb-with-release-ip"
						type = "LB-S"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isLbPresent(tt, "scaleway_lb.main"),
					isIPPresent(tt, "scaleway_lb_ip.ip01"),
				),
			},
		},
	})
}

func TestAccIP_IPv6(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      lbchecks.IsIPDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_lb_ip ipv6 {
						is_ipv6 = true
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isIPPresent(tt, "scaleway_lb_ip.ipv6"),
					acctest.CheckResourceAttrIPv6("scaleway_lb_ip.ipv6", "ip_address"),
					resource.TestCheckResourceAttr("scaleway_lb_ip.ipv6", "is_ipv6", "true"),
				),
			},
		},
	})
}

func TestAccIP_WithTags(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      lbchecks.IsIPDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_lb_ip tags {
					  tags = [ "terraform-test", "lb", "ip" ]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isIPPresent(tt, "scaleway_lb_ip.tags"),
					resource.TestCheckResourceAttr("scaleway_lb_ip.tags", "tags.#", "3"),
					resource.TestCheckResourceAttr("scaleway_lb_ip.tags", "tags.0", "terraform-test"),
					resource.TestCheckResourceAttr("scaleway_lb_ip.tags", "tags.1", "lb"),
					resource.TestCheckResourceAttr("scaleway_lb_ip.tags", "tags.2", "ip")),
			},
			{
				Config: `
					resource scaleway_lb_ip tags {
					  tags = [ "terraform-test", "lb", "ip", "updated" ]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isIPPresent(tt, "scaleway_lb_ip.tags"),
					resource.TestCheckResourceAttr("scaleway_lb_ip.tags", "tags.#", "4"),
					resource.TestCheckResourceAttr("scaleway_lb_ip.tags", "tags.0", "terraform-test"),
					resource.TestCheckResourceAttr("scaleway_lb_ip.tags", "tags.1", "lb"),
					resource.TestCheckResourceAttr("scaleway_lb_ip.tags", "tags.2", "ip"),
					resource.TestCheckResourceAttr("scaleway_lb_ip.tags", "tags.3", "updated"),
				),
			},
		},
	})
}

func isIPPresent(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		lbAPI, zone, ID, err := lb.NewAPIWithZoneAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = lbAPI.GetIP(&lbSDK.ZonedAPIGetIPRequest{
			IPID: ID,
			Zone: zone,
		})
		if err != nil {
			return err
		}

		return nil
	}
}
