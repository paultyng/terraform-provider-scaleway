package lb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	lbchecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/lb/testfuncs"
)

func TestAccDataSourceIP_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      lbchecks.IsIPDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_lb_ip" "test" {
					}
					
					resource "scaleway_lb_ip" "test_par_2" {
					  zone = "fr-par-2"
					}
					
					data "scaleway_lb_ip" "test" {
					  ip_address = scaleway_lb_ip.test.ip_address
					}
					
					data "scaleway_lb_ip" "test2" {
					  ip_id = scaleway_lb_ip.test.id
					}
					
					data "scaleway_lb_ip" "test_another_zone" {
					  ip_address = scaleway_lb_ip.test_par_2.ip_address
					  zone       = "fr-par-2"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isIPPresent(tt, "data.scaleway_lb_ip.test"),
					resource.TestCheckResourceAttrPair("data.scaleway_lb_ip.test", "ip_address", "scaleway_lb_ip.test", "ip_address"),
					resource.TestCheckResourceAttrPair("data.scaleway_lb_ip.test2", "ip_address", "scaleway_lb_ip.test", "ip_address"),
					resource.TestCheckResourceAttrPair("data.scaleway_lb_ip.test_another_zone", "ip_address", "scaleway_lb_ip.test_par_2", "ip_address"),
				),
			},
		},
	})
}
