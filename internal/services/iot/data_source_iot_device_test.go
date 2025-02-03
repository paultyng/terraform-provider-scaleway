package iot_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

func TestAccDataSourceDevice_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isHubDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_iot_device" "test" {
						name = "test_iot_device_datasource"
						hub_id = scaleway_iot_hub.test.id
					}

					resource "scaleway_iot_hub" "test" {
						name = "test_iot_device_datasource"
						product_plan = "plan_shared"
					}

					data "scaleway_iot_device" "by_name" {
						name = scaleway_iot_device.test.name
					}

					data "scaleway_iot_device" "by_name_and_hub" {
						name = scaleway_iot_device.test.name
						hub_id = scaleway_iot_hub.test.id
					}

					data "scaleway_iot_device" "by_id" {
						device_id = scaleway_iot_device.test.id
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isDevicePresent(tt, "scaleway_iot_device.test"),

					resource.TestCheckResourceAttr("data.scaleway_iot_device.by_name", "name", "test_iot_device_datasource"),
					resource.TestCheckResourceAttrSet("data.scaleway_iot_device.by_name", "id"),

					resource.TestCheckResourceAttr("data.scaleway_iot_device.by_name_and_hub", "name", "test_iot_device_datasource"),
					resource.TestCheckResourceAttrSet("data.scaleway_iot_device.by_name_and_hub", "id"),

					resource.TestCheckResourceAttr("data.scaleway_iot_device.by_id", "name", "test_iot_device_datasource"),
					resource.TestCheckResourceAttrSet("data.scaleway_iot_device.by_id", "id"),
				),
			},
		},
	})
}
