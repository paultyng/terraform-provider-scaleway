package az_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

func TestAccDataSourceAvailabilityZones_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					data scaleway_availability_zones main {
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.scaleway_availability_zones.main", "region", "fr-par"),
					resource.TestCheckResourceAttr(
						"data.scaleway_availability_zones.main", "zones.0", "fr-par-1"),
				),
			},
			{
				Config: `
					data scaleway_availability_zones main {
						region = "nl-ams"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.scaleway_availability_zones.main", "region", "nl-ams"),
					resource.TestCheckResourceAttr(
						"data.scaleway_availability_zones.main", "zones.0", "nl-ams-1"),
					resource.TestCheckResourceAttr(
						"data.scaleway_availability_zones.main", "zones.1", "nl-ams-2"),
				),
			},
		},
	})
}
