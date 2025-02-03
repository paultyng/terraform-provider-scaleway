package webhosting_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

func TestAccDataSourceOffer_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					data "scaleway_webhosting_offer" "by_name" {
						name = "professional"
					}

					data "scaleway_webhosting_offer" "by_id" {
						offer_id = "de2426b4-a9e9-11ec-b909-0242ac120002"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.scaleway_webhosting_offer.by_id", "id"),
					resource.TestCheckResourceAttr("data.scaleway_webhosting_offer.by_id", "name", "professional"),

					resource.TestCheckResourceAttrSet("data.scaleway_webhosting_offer.by_name", "id"),
					resource.TestCheckResourceAttr("data.scaleway_webhosting_offer.by_name", "name", "professional"),
					resource.TestCheckResourceAttr("data.scaleway_webhosting_offer.by_name", "product.0.option", "false"),
					resource.TestCheckResourceAttr("data.scaleway_webhosting_offer.by_name", "product.0.email_accounts_quota", "10"),
					resource.TestCheckResourceAttr("data.scaleway_webhosting_offer.by_name", "product.0.email_storage_quota", "5"),
					resource.TestCheckResourceAttr("data.scaleway_webhosting_offer.by_name", "product.0.databases_quota", "-1"),
					resource.TestCheckResourceAttr("data.scaleway_webhosting_offer.by_name", "product.0.hosting_storage_quota", "100"),
					resource.TestCheckResourceAttr("data.scaleway_webhosting_offer.by_name", "product.0.support_included", "true"),
					resource.TestCheckResourceAttr("data.scaleway_webhosting_offer.by_name", "product.0.v_cpu", "4"),
					resource.TestCheckResourceAttr("data.scaleway_webhosting_offer.by_name", "product.0.ram", "2"),
					resource.TestCheckResourceAttr("data.scaleway_webhosting_offer.by_name", "price", "€ 18.99"),
				),
			},
		},
	})
}
