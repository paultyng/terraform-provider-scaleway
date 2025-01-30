package block_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	blockSDK "github.com/scaleway/scaleway-sdk-go/api/block/v1alpha1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/block"
)

func TestAccSnapshot_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isSnapshotDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_block_volume main {
						iops = 5000
						size_in_gb = 10
					}

					resource scaleway_block_snapshot main {
						name = "test-block-snapshot-basic"
						volume_id = scaleway_block_volume.main.id
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isSnapshotPresent(tt, "scaleway_block_snapshot.main"),
					acctest.CheckResourceAttrUUID("scaleway_block_snapshot.main", "id"),
					resource.TestCheckResourceAttr("scaleway_block_snapshot.main", "name", "test-block-snapshot-basic"),
				),
			},
		},
	})
}

func isSnapshotPresent(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		api, zone, id, err := block.NewAPIWithZoneAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = api.GetSnapshot(&blockSDK.GetSnapshotRequest{
			SnapshotID: id,
			Zone:       zone,
		})
		if err != nil {
			return err
		}

		return nil
	}
}

func isSnapshotDestroyed(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_block_snapshot" {
				continue
			}

			api, zone, id, err := block.NewAPIWithZoneAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}

			err = api.DeleteSnapshot(&blockSDK.DeleteSnapshotRequest{
				SnapshotID: id,
				Zone:       zone,
			})

			if err == nil {
				return fmt.Errorf("block snapshot (%s) still exists", rs.Primary.ID)
			}

			if !httperrors.Is404(err) {
				return err
			}
		}

		return nil
	}
}
