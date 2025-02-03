package lb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	lbchecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/lb/testfuncs"
)

func TestAccDataSourceIPs_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      lbchecks.IsIPDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_lb_ip ip1 {
					}
				`,
			},
			{
				Config: `
					resource scaleway_lb_ip ip1 {
					}
					resource scaleway_lb_ip ip2 {
					}

					data "scaleway_lb_ips" "ips_by_cidr_range" {
						ip_cidr_range = "0.0.0.0/0"
						depends_on = [scaleway_lb_ip.ip1, scaleway_lb_ip.ip2]
					}
					data "scaleway_lb_ips" "ips_by_cidr_range_other_zone" {
						ip_cidr_range  = "0.0.0.0/0"
						zone = "fr-par-2"
						depends_on = [scaleway_lb_ip.ip1, scaleway_lb_ip.ip2]
					}
					`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.scaleway_lb_ips.ips_by_cidr_range", "ips.0.id"),
					resource.TestCheckResourceAttrSet("data.scaleway_lb_ips.ips_by_cidr_range", "ips.0.ip_address"),
					resource.TestCheckResourceAttrSet("data.scaleway_lb_ips.ips_by_cidr_range", "ips.1.id"),
					resource.TestCheckResourceAttrSet("data.scaleway_lb_ips.ips_by_cidr_range", "ips.1.ip_address"),

					resource.TestCheckNoResourceAttr("data.scaleway_lb_ips.ips_by_cidr_range_other_zone", "ips.0.id"),
				),
			},
		},
	})
}

func TestAccDataSourceIPs_WithType(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      lbchecks.IsIPDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_lb_ip" "ip1" {
					  zone = "nl-ams-1"
					}					
				`,
			},
			{
				Config: `
					resource "scaleway_lb_ip" "ip1" {
					  zone = "nl-ams-1"
					}

					resource "scaleway_lb_ip" "ip2" {
					  is_ipv6 = true
					  zone    = "nl-ams-1"
					}					
				`,
			},
			{
				Config: `
					resource "scaleway_lb_ip" "ip1" {
					  zone = "nl-ams-1"
					}
					
					resource "scaleway_lb_ip" "ip2" {
					  is_ipv6 = true
					  zone    = "nl-ams-1"
					}
					
					resource "scaleway_lb_ip" "ip3" {
					  zone = "nl-ams-1"
					}
					
					data "scaleway_lb_ips" "ips_by_type" {
					  ip_type    = "ipv4"
					  zone       = "nl-ams-1"
					  depends_on = [scaleway_lb_ip.ip1, scaleway_lb_ip.ip2, scaleway_lb_ip.ip3]
					}					
					`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.scaleway_lb_ips.ips_by_type", "ips.#", "2"),
					resource.TestCheckResourceAttrSet("data.scaleway_lb_ips.ips_by_type", "ips.0.id"),
					resource.TestCheckResourceAttrSet("data.scaleway_lb_ips.ips_by_type", "ips.0.ip_address"),
					resource.TestCheckResourceAttrSet("data.scaleway_lb_ips.ips_by_type", "ips.1.id"),
					resource.TestCheckResourceAttrSet("data.scaleway_lb_ips.ips_by_type", "ips.1.ip_address"),
				),
			},
		},
	})
}

func TestAccDataSourceIPs_WithTags(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      lbchecks.IsIPDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_lb_ip ip1 {
					  tags = [ "ipv4", "ip" ]
					}
				`,
			},
			{
				Config: `
					resource scaleway_lb_ip ip1 {
					  tags = [ "ipv4", "ip" ]
					}
					resource scaleway_lb_ip ip2 {
					  tags = [ "ipv4", "ip" ]
					}
				`,
			},
			{
				Config: `
					resource scaleway_lb_ip ip1 {
					  tags = [ "ipv4", "ip" ]
					}
					resource scaleway_lb_ip ip2 {
					  tags = [ "ipv4", "ip" ]
					}
					resource scaleway_lb_ip ip3 {
					  tags = [ "other", "tags" ]
					}

					data "scaleway_lb_ips" "ips_by_tags" {
					    tags = [ "ipv4", "ip" ]
						depends_on = [scaleway_lb_ip.ip1, scaleway_lb_ip.ip2, scaleway_lb_ip.ip3]
					}
					`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.scaleway_lb_ips.ips_by_tags", "ips.#", "2"),
					resource.TestCheckResourceAttrSet("data.scaleway_lb_ips.ips_by_tags", "ips.0.id"),
					resource.TestCheckResourceAttrSet("data.scaleway_lb_ips.ips_by_tags", "ips.0.ip_address"),
					resource.TestCheckResourceAttrSet("data.scaleway_lb_ips.ips_by_tags", "ips.1.id"),
					resource.TestCheckResourceAttrSet("data.scaleway_lb_ips.ips_by_tags", "ips.1.ip_address"),
				),
			},
		},
	})
}
