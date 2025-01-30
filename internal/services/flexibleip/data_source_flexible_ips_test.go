package flexibleip_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	lbchecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/lb/testfuncs"
)

func TestAccDataSourceFlexibleIPs_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      lbchecks.IsIPDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_flexible_ip" "first" {
						tags = [ "minimal", "first" ]
					}
				`,
			},
			{
				Config: `
					resource "scaleway_flexible_ip" "first" {
						tags = [ "minimal", "first" ]
					}
					
					resource "scaleway_flexible_ip" "second" {
						tags = [ "minimal", "second" ]
					}

					data "scaleway_flexible_ips" "fips_by_tags" {
						tags = [ "minimal" ]
						depends_on = [scaleway_flexible_ip.first, scaleway_flexible_ip.second]
					}

					data "scaleway_flexible_ips" "fips_by_tags_other_zone" {
						tags = [ "minimal" ]
						zone = "fr-par-2"
						depends_on = [scaleway_flexible_ip.first, scaleway_flexible_ip.second]
					}
					`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.scaleway_flexible_ips.fips_by_tags", "ips.0.id"),
					resource.TestCheckResourceAttrSet("data.scaleway_flexible_ips.fips_by_tags", "ips.0.ip_address"),
					resource.TestCheckResourceAttrSet("data.scaleway_flexible_ips.fips_by_tags", "ips.1.id"),
					resource.TestCheckResourceAttrSet("data.scaleway_flexible_ips.fips_by_tags", "ips.1.ip_address"),

					resource.TestCheckNoResourceAttr("data.scaleway_flexible_ips.fips_by_tags_other_zone", "ips.0.id"),
				),
			},
		},
	})
}

func TestAccDataSourceFlexibleIPs_WithBaremetalIDs(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      lbchecks.IsIPDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					data "scaleway_baremetal_offer" "my_offer" {
					  name = "EM-A115X-SSD"
					  zone = "fr-par-1"
					}

					resource "scaleway_baremetal_server" "base" {
					  name 			             = "TestAccScalewayBaremetalServer_WithoutInstallConfig"
					  offer     				 = data.scaleway_baremetal_offer.my_offer.offer_id
					  install_config_afterward   = true
					}

					resource "scaleway_flexible_ip" "first" {
						tags = [ "minimal", "first" ]
					}
				`,
			},
			{
				Config: `
					data "scaleway_baremetal_offer" "my_offer" {
					  name = "EM-A115X-SSD"
					  zone = "fr-par-1"
					}

					resource "scaleway_baremetal_server" "base" {
					  name 			             = "TestAccScalewayBaremetalServer_WithoutInstallConfig"
					  offer     				 = data.scaleway_baremetal_offer.my_offer.offer_id
					  install_config_afterward   = true
					}

					resource "scaleway_flexible_ip" "first" {
						tags = [ "minimal", "first" ]
					}
					
					resource "scaleway_flexible_ip" "second" {
						tags = [ "minimal", "second" ]
					}
					`,
			},
			{
				Config: `
					data "scaleway_baremetal_offer" "my_offer" {
					  name = "EM-A115X-SSD"
					  zone = "fr-par-1"
					}

					resource "scaleway_baremetal_server" "base" {
					  name 			             = "TestAccScalewayBaremetalServer_WithoutInstallConfig"
					  offer     				 = data.scaleway_baremetal_offer.my_offer.offer_id
					  install_config_afterward   = true
					}

					resource "scaleway_flexible_ip" "first" {
						server_id = scaleway_baremetal_server.base.id
						tags = [ "minimal", "first" ]
					}
					
					resource "scaleway_flexible_ip" "second" {
						server_id = scaleway_baremetal_server.base.id
						tags = [ "minimal", "second" ]
					}

					data "scaleway_flexible_ips" "fips_by_server_id" {
						server_ids = [scaleway_baremetal_server.base.id]
						depends_on = [scaleway_flexible_ip.first, scaleway_flexible_ip.second]
					}
					`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.scaleway_flexible_ips.fips_by_server_id", "ips.0.id"),
					resource.TestCheckResourceAttrSet("data.scaleway_flexible_ips.fips_by_server_id", "ips.0.ip_address"),
					resource.TestCheckResourceAttrSet("data.scaleway_flexible_ips.fips_by_server_id", "ips.1.id"),
					resource.TestCheckResourceAttrSet("data.scaleway_flexible_ips.fips_by_server_id", "ips.1.ip_address"),
				),
			},
		},
	})
}
