package flexibleip_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

func TestAccDataSourceFlexibleIP_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckFlexibleIPDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_flexible_ip" "main" {
					}
					
					data "scaleway_flexible_ip" "by_address" {
						ip_address = "${scaleway_flexible_ip.main.ip_address}"
					}
					
					data "scaleway_flexible_ip" "by_id" {
						flexible_ip_id = "${scaleway_flexible_ip.main.id}"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlexibleIPExists(tt, "scaleway_flexible_ip.main"),

					resource.TestCheckResourceAttrSet("scaleway_flexible_ip.main", "ip_address"),
					resource.TestCheckResourceAttrSet("data.scaleway_flexible_ip.by_address", "ip_address"),
					resource.TestCheckResourceAttrSet("data.scaleway_flexible_ip.by_id", "ip_address"),

					resource.TestCheckResourceAttrSet("scaleway_flexible_ip.main", "id"),
					resource.TestCheckResourceAttrSet("data.scaleway_flexible_ip.by_address", "id"),
					resource.TestCheckResourceAttrSet("data.scaleway_flexible_ip.by_id", "flexible_ip_id"),
				),
			},
		},
	})
}

func TestAccDataSourceFlexibleIP_Multiple(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckFlexibleIPDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_flexible_ip" "first" {
					}
					
					data "scaleway_flexible_ip" "by_address_first" {
						ip_address = "${scaleway_flexible_ip.first.ip_address}"
					}
					
					data "scaleway_flexible_ip" "by_id_first" {
						flexible_ip_id = "${scaleway_flexible_ip.first.id}"
					}

					resource "scaleway_flexible_ip" "second" {
					}
					
					data "scaleway_flexible_ip" "by_address_second" {
						ip_address = "${scaleway_flexible_ip.second.ip_address}"
					}
					
					data "scaleway_flexible_ip" "by_id_second" {
						flexible_ip_id = "${scaleway_flexible_ip.second.id}"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlexibleIPExists(tt, "scaleway_flexible_ip.first"),

					resource.TestCheckResourceAttrSet("scaleway_flexible_ip.first", "ip_address"),
					resource.TestCheckResourceAttrSet("data.scaleway_flexible_ip.by_address_first", "ip_address"),
					resource.TestCheckResourceAttrSet("data.scaleway_flexible_ip.by_id_first", "ip_address"),

					resource.TestCheckResourceAttrSet("scaleway_flexible_ip.first", "id"),
					resource.TestCheckResourceAttrSet("data.scaleway_flexible_ip.by_address_first", "id"),
					resource.TestCheckResourceAttrSet("data.scaleway_flexible_ip.by_id_first", "flexible_ip_id"),

					testAccCheckFlexibleIPExists(tt, "scaleway_flexible_ip.second"),

					resource.TestCheckResourceAttrSet("scaleway_flexible_ip.second", "ip_address"),
					resource.TestCheckResourceAttrSet("data.scaleway_flexible_ip.by_address_second", "ip_address"),
					resource.TestCheckResourceAttrSet("data.scaleway_flexible_ip.by_id_second", "ip_address"),

					resource.TestCheckResourceAttrSet("scaleway_flexible_ip.second", "id"),
					resource.TestCheckResourceAttrSet("data.scaleway_flexible_ip.by_address_second", "id"),
					resource.TestCheckResourceAttrSet("data.scaleway_flexible_ip.by_id_second", "flexible_ip_id"),
				),
			},
		},
	})
}
