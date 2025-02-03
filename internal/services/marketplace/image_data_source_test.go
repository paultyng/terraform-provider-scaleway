package marketplace_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	instancechecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/instance/testfuncs"
)

func TestAccDataSourceMarketplaceImage_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					data "scaleway_marketplace_image" "test1" {
						label = "ubuntu_focal"
					}
					`,
				Check: resource.ComposeTestCheckFunc(
					instancechecks.DoesImageExists(tt, "data.scaleway_marketplace_image.test1"),
					resource.TestCheckResourceAttr("data.scaleway_marketplace_image.test1", "label", "ubuntu_focal"),
				),
			},
		},
	})
}
