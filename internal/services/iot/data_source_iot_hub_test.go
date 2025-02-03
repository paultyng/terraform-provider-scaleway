package iot_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

func TestAccDataSourceHub_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isHubDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_iot_hub" "test" {
						name = "test_iot_hub_datasource"
						product_plan = "plan_shared"
					}

					data "scaleway_iot_hub" "by_name" {
						name = scaleway_iot_hub.test.name
					}

					data "scaleway_iot_hub" "by_id" {
						hub_id = scaleway_iot_hub.test.id
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isHubPresent(tt, "scaleway_iot_hub.test"),

					resource.TestCheckResourceAttr("data.scaleway_iot_hub.by_name", "name", "test_iot_hub_datasource"),
					resource.TestCheckResourceAttrSet("data.scaleway_iot_hub.by_name", "id"),

					resource.TestCheckResourceAttr("data.scaleway_iot_hub.by_id", "name", "test_iot_hub_datasource"),
					resource.TestCheckResourceAttrSet("data.scaleway_iot_hub.by_id", "id"),
				),
			},
		},
	})
}
