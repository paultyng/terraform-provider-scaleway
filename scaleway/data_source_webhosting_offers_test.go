package scaleway

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccScalewayDataSourceWebhostingOffer_Basic(t *testing.T) {
	tt := NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					data "scaleway_webhosting_offers" "by_name" {
						name = "performance"
					}

					data "scaleway_webhosting_offers" "by_id" {
						offer_id = "de2426b4-a9e9-11ec-b909-0242ac120002"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.scaleway_webhosting_offers.by_id", "offers.0.product.0.name", "performance"),

					resource.TestCheckResourceAttr("data.scaleway_webhosting_offers.by_name", "offers.0.product.0.name", "performance"),
					resource.TestCheckResourceAttr("data.scaleway_webhosting_offers.by_name", "offers.0.product.0.option", "false"),
					resource.TestCheckResourceAttr("data.scaleway_webhosting_offers.by_name", "offers.0.product.0.email_accounts_quota", "10"),
					resource.TestCheckResourceAttr("data.scaleway_webhosting_offers.by_name", "offers.0.product.0.email_storage_quota", "5"),
					resource.TestCheckResourceAttr("data.scaleway_webhosting_offers.by_name", "offers.0.product.0.databases_quota", "-1"),
					resource.TestCheckResourceAttr("data.scaleway_webhosting_offers.by_name", "offers.0.product.0.hosting_storage_quota", "100"),
					resource.TestCheckResourceAttr("data.scaleway_webhosting_offers.by_name", "offers.0.product.0.support_included", "true"),
					resource.TestCheckResourceAttr("data.scaleway_webhosting_offers.by_name", "offers.0.product.0.v_cpu", "4"),
					resource.TestCheckResourceAttr("data.scaleway_webhosting_offers.by_name", "offers.0.product.0.ram", "2"),
				),
			},
		},
	})
}
