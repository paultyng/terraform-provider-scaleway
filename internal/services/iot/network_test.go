package iot_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	iotSDK "github.com/scaleway/scaleway-sdk-go/api/iot/v1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/iot"
)

func TestAccNetwork_Minimal(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		// Destruction is done via the hub destruction.
		CheckDestroy: isHubDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
						resource "scaleway_iot_hub" "minimal" {
							name         = "minimal"
							product_plan = "plan_shared"
						}

						resource "scaleway_iot_network" "default" {
							name   = "default"
							hub_id = scaleway_iot_hub.minimal.id
							type   = "rest"
						}
						`,
				Check: resource.ComposeTestCheckFunc(
					isHubPresent(tt, "scaleway_iot_hub.minimal"),
					isNetworkPresent(tt, "scaleway_iot_network.default"),
					resource.TestCheckResourceAttrSet("scaleway_iot_network.default", "id"),
					resource.TestCheckResourceAttrSet("scaleway_iot_network.default", "hub_id"),
					resource.TestCheckResourceAttr("scaleway_iot_network.default", "name", "default"),
					resource.TestCheckResourceAttr("scaleway_iot_network.default", "type", "rest"),
					resource.TestCheckResourceAttrSet("scaleway_iot_network.default", "endpoint"),
					resource.TestCheckResourceAttrSet("scaleway_iot_network.default", "secret"),
					resource.TestCheckResourceAttrSet("scaleway_iot_network.default", "created_at"),
				),
			},
		},
	})
}

func TestAccNetwork_RESTWithTopicPrefix(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		// Destruction is done via the hub destruction.
		CheckDestroy: isHubDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
						resource "scaleway_iot_hub" "minimal" {
							name         = "minimal"
							product_plan = "plan_shared"
						}

						resource "scaleway_iot_network" "default" {
							name         = "default"
							hub_id       = scaleway_iot_hub.minimal.id
							type         = "rest"
							topic_prefix = "foo/bar"
						}
						`,
				Check: resource.ComposeTestCheckFunc(
					isHubPresent(tt, "scaleway_iot_hub.minimal"),
					isNetworkPresent(tt, "scaleway_iot_network.default"),
					resource.TestCheckResourceAttrSet("scaleway_iot_network.default", "id"),
					resource.TestCheckResourceAttrSet("scaleway_iot_network.default", "hub_id"),
					resource.TestCheckResourceAttr("scaleway_iot_network.default", "name", "default"),
					resource.TestCheckResourceAttr("scaleway_iot_network.default", "type", "rest"),
					resource.TestCheckResourceAttr("scaleway_iot_network.default", "topic_prefix", "foo/bar"),
					resource.TestCheckResourceAttrSet("scaleway_iot_network.default", "endpoint"),
					resource.TestCheckResourceAttrSet("scaleway_iot_network.default", "secret"),
					resource.TestCheckResourceAttrSet("scaleway_iot_network.default", "created_at"),
				),
			},
		},
	})
}

func TestAccNetwork_Sigfox(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		// Destruction is done via the hub destruction.
		CheckDestroy: isHubDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
						resource "scaleway_iot_hub" "minimal" {
							name         = "minimal"
							product_plan = "plan_shared"
						}

						resource "scaleway_iot_network" "default" {
							name   = "default"
							hub_id = scaleway_iot_hub.minimal.id
							type   = "sigfox"
						}
						`,
				Check: resource.ComposeTestCheckFunc(
					isHubPresent(tt, "scaleway_iot_hub.minimal"),
					isNetworkPresent(tt, "scaleway_iot_network.default"),
					resource.TestCheckResourceAttrSet("scaleway_iot_network.default", "id"),
					resource.TestCheckResourceAttrSet("scaleway_iot_network.default", "hub_id"),
					resource.TestCheckResourceAttr("scaleway_iot_network.default", "name", "default"),
					resource.TestCheckResourceAttr("scaleway_iot_network.default", "type", "sigfox"),
					resource.TestCheckResourceAttrSet("scaleway_iot_network.default", "endpoint"),
					resource.TestCheckResourceAttrSet("scaleway_iot_network.default", "secret"),
					resource.TestCheckResourceAttrSet("scaleway_iot_network.default", "created_at"),
				),
			},
		},
	})
}

func isNetworkPresent(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		iotAPI, region, networkID, err := iot.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = iotAPI.GetNetwork(&iotSDK.GetNetworkRequest{
			Region:    region,
			NetworkID: networkID,
		})
		if err != nil {
			return err
		}

		return nil
	}
}
