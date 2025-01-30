package instance_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

func TestAccDataSourcePlacementGroup_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			isPlacementGroupDestroyed(tt),
		),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_instance_placement_group main {
  						name = "test-ds-instance-placement-group-basic"
					}

					data scaleway_instance_placement_group find_by_name {
						name = scaleway_instance_placement_group.main.name
					}

					data scaleway_instance_placement_group find_by_id {
						placement_group_id = scaleway_instance_placement_group.main.id
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isPlacementGroupPresent(tt, "scaleway_instance_placement_group.main"),

					resource.TestCheckResourceAttrPair("scaleway_instance_placement_group.main", "name", "data.scaleway_instance_placement_group.find_by_name", "name"),
					resource.TestCheckResourceAttrPair("scaleway_instance_placement_group.main", "name", "data.scaleway_instance_placement_group.find_by_id", "name"),
					resource.TestCheckResourceAttrPair("scaleway_instance_placement_group.main", "id", "data.scaleway_instance_placement_group.find_by_name", "id"),
					resource.TestCheckResourceAttrPair("scaleway_instance_placement_group.main", "id", "data.scaleway_instance_placement_group.find_by_id", "id"),
				),
			},
		},
	})
}
