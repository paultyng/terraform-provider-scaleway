package iot_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	iotSDK "github.com/scaleway/scaleway-sdk-go/api/iot/v1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/iot"
)

func TestAccHub_Minimal(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isHubDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
						resource "scaleway_iot_hub" "minimal" {
							name = "minimal"
							product_plan = "plan_shared"
						}`,
				Check: resource.ComposeTestCheckFunc(
					isHubPresent(tt, "scaleway_iot_hub.minimal"),
					resource.TestCheckResourceAttr("scaleway_iot_hub.minimal", "product_plan", "plan_shared"),
					resource.TestCheckResourceAttr("scaleway_iot_hub.minimal", "status", iotSDK.HubStatusReady.String()),
					resource.TestCheckResourceAttrSet("scaleway_iot_hub.minimal", "endpoint"),
					resource.TestCheckResourceAttr("scaleway_iot_hub.minimal", "device_count", "0"),
					resource.TestCheckResourceAttr("scaleway_iot_hub.minimal", "connected_device_count", "0"),
					// If the plan is shared, there is no MQTT CA, so it should be empty
					resource.TestCheckResourceAttr("scaleway_iot_hub.minimal", "mqtt_ca_url", ""),
					resource.TestCheckResourceAttr("scaleway_iot_hub.minimal", "mqtt_ca", ""),
				),
			},
			{
				Config: `
						resource "scaleway_iot_hub" "minimal" {
							name = "minimal"
							product_plan = "plan_shared"
							enabled = false
						}`,
				Check: resource.ComposeTestCheckFunc(
					isHubPresent(tt, "scaleway_iot_hub.minimal"),
					resource.TestCheckResourceAttr("scaleway_iot_hub.minimal", "status", iotSDK.HubStatusDisabled.String()),
					resource.TestCheckResourceAttrSet("scaleway_iot_hub.minimal", "endpoint"),
					resource.TestCheckResourceAttr("scaleway_iot_hub.minimal", "device_count", "0"),
					resource.TestCheckResourceAttr("scaleway_iot_hub.minimal", "connected_device_count", "0"),
					resource.TestCheckResourceAttr("scaleway_iot_hub.minimal", "enabled", "false"),
					resource.TestCheckResourceAttr("scaleway_iot_hub.minimal", "product_plan", "plan_shared"),
					resource.TestCheckResourceAttr("scaleway_iot_hub.minimal", "mqtt_ca_url", ""),
					resource.TestCheckResourceAttr("scaleway_iot_hub.minimal", "mqtt_ca", ""),
				),
			},
		},
	})
}

func TestAccHub_Dedicated(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isHubDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
						resource "scaleway_iot_hub" "minimal" {
							name = "minimal"
							product_plan = "plan_dedicated"
						}`,
				Check: resource.ComposeTestCheckFunc(
					isHubPresent(tt, "scaleway_iot_hub.minimal"),
					resource.TestCheckResourceAttr("scaleway_iot_hub.minimal", "product_plan", "plan_dedicated"),
					resource.TestCheckResourceAttr("scaleway_iot_hub.minimal", "status", iotSDK.HubStatusReady.String()),
					resource.TestCheckResourceAttrSet("scaleway_iot_hub.minimal", "endpoint"),
					resource.TestCheckResourceAttr("scaleway_iot_hub.minimal", "device_count", "0"),
					resource.TestCheckResourceAttr("scaleway_iot_hub.minimal", "connected_device_count", "0"),
					resource.TestCheckResourceAttr("scaleway_iot_hub.minimal", "mqtt_ca_url", "https://iot.s3.nl-ams.scw.cloud/certificates/fr-par/iot-hub-ca.pem"),
					resource.TestCheckResourceAttrSet("scaleway_iot_hub.minimal", "mqtt_ca"),
				),
			},
			{
				Config: `
						resource "scaleway_iot_hub" "minimal" {
							name = "minimal"
							product_plan = "plan_dedicated"
						}`,
				Check: resource.ComposeTestCheckFunc(
					isHubPresent(tt, "scaleway_iot_hub.minimal"),
					resource.TestCheckResourceAttr("scaleway_iot_hub.minimal", "product_plan", "plan_dedicated"),
					resource.TestCheckResourceAttr("scaleway_iot_hub.minimal", "status", iotSDK.HubStatusReady.String()),
					resource.TestCheckResourceAttrSet("scaleway_iot_hub.minimal", "endpoint"),
					resource.TestCheckResourceAttr("scaleway_iot_hub.minimal", "device_count", "0"),
					resource.TestCheckResourceAttr("scaleway_iot_hub.minimal", "connected_device_count", "0"),
					resource.TestCheckResourceAttr("scaleway_iot_hub.minimal", "mqtt_ca_url", "https://iot.s3.nl-ams.scw.cloud/certificates/fr-par/iot-hub-ca.pem"),
					resource.TestCheckResourceAttrSet("scaleway_iot_hub.minimal", "mqtt_ca"),
				),
			},
		},
	})
}

func isHubDestroyed(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_iot_hub" {
				continue
			}

			iotAPI, region, hubID, err := iot.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}

			_, err = iotAPI.GetHub(&iotSDK.GetHubRequest{
				Region: region,
				HubID:  hubID,
			})

			// If no error resource still exist
			if err == nil {
				return fmt.Errorf("hub (%s) still exists", rs.Primary.ID)
			}

			// Unexpected api error we return it
			if !httperrors.Is404(err) {
				return err
			}
		}
		return nil
	}
}

func isHubPresent(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		iotAPI, region, hubID, err := iot.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = iotAPI.GetHub(&iotSDK.GetHubRequest{
			Region: region,
			HubID:  hubID,
		})
		if err != nil {
			return err
		}

		return nil
	}
}
