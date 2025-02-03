package webhosting_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	webhostingSDK "github.com/scaleway/scaleway-sdk-go/api/webhosting/v1alpha1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/webhosting"
)

func TestAccWebhosting_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckWebhostingDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
				data "scaleway_webhosting_offer" "by_name" {
				  name = "lite"
				}

				resource "scaleway_webhosting" "main" {
				  offer_id     = data.scaleway_webhosting_offer.by_name.offer_id
				  email        = "hashicorp@scaleway.com"
				  domain       = "scaleway.com"
				  tags         = ["devtools", "provider", "terraform"]
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWebhostingExists(tt, "scaleway_webhosting.main"),
					resource.TestCheckResourceAttrPair("scaleway_webhosting.main", "offer_id", "data.scaleway_webhosting_offer.by_name", "offer_id"),
					resource.TestCheckResourceAttr("scaleway_webhosting.main", "email", "hashicorp@scaleway.com"),
					resource.TestCheckResourceAttr("scaleway_webhosting.main", "domain", "scaleway.com"),
					resource.TestCheckResourceAttr("scaleway_webhosting.main", "status", webhostingSDK.HostingStatusReady.String()),
					resource.TestCheckResourceAttr("scaleway_webhosting.main", "tags.0", "devtools"),
					resource.TestCheckResourceAttr("scaleway_webhosting.main", "tags.1", "provider"),
					resource.TestCheckResourceAttr("scaleway_webhosting.main", "tags.2", "terraform"),
					resource.TestCheckResourceAttrSet("scaleway_webhosting.main", "updated_at"),
					resource.TestCheckResourceAttrSet("scaleway_webhosting.main", "created_at"),
					acctest.CheckResourceAttrUUID("scaleway_webhosting.main", "id"),
				),
			},
		},
	})
}

func testAccCheckWebhostingExists(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		api, region, id, err := webhosting.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = api.GetHosting(&webhostingSDK.GetHostingRequest{
			HostingID: id,
			Region:    region,
		})
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckWebhostingDestroy(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_webhosting" {
				continue
			}

			api, region, id, err := webhosting.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}

			res, err := api.GetHosting(&webhostingSDK.GetHostingRequest{
				HostingID: id,
				Region:    region,
			})

			if err == nil && res.Status != webhostingSDK.HostingStatusUnknownStatus {
				return fmt.Errorf("hosting (%s) still exists", rs.Primary.ID)
			}

			if !httperrors.Is404(err) {
				return err
			}
		}

		return nil
	}
}
