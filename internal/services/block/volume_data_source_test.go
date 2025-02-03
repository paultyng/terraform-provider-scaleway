package block_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

func TestAccDataSourceVolume_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			isVolumeDestroyed(tt),
		),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_block_volume main {
						iops = 5000
						size_in_gb = 10
  						name = "test-ds-block-volume-basic"
					}

					data scaleway_block_volume find_by_name {
						name = scaleway_block_volume.main.name
					}

					data scaleway_block_volume find_by_id {
						volume_id = scaleway_block_volume.main.id
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isVolumePresent(tt, "scaleway_block_volume.main"),

					resource.TestCheckResourceAttrPair("scaleway_block_volume.main", "name", "data.scaleway_block_volume.find_by_name", "name"),
					resource.TestCheckResourceAttrPair("scaleway_block_volume.main", "name", "data.scaleway_block_volume.find_by_id", "name"),
					resource.TestCheckResourceAttrPair("scaleway_block_volume.main", "id", "data.scaleway_block_volume.find_by_name", "id"),
					resource.TestCheckResourceAttrPair("scaleway_block_volume.main", "id", "data.scaleway_block_volume.find_by_id", "id"),
				),
			},
		},
	})
}
