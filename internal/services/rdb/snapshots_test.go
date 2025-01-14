package rdb_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	rdbSDK "github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/rdb"
	rdbchecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/rdb/testfuncs"
)

func TestAccRdbSnapshot_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	latestEngineVersion := rdbchecks.GetLatestEngineVersion(tt, postgreSQLEngineName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      IsSnapshotDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "scaleway_rdb_instance" "main" {
						name           = "test-rdb-basic"
						node_type      = "db-dev-s"
						engine         = %q
						is_ha_cluster  = false
						disable_backup = true
						user_name      = "my_initial_user"
						password       = "thiZ_is_v&ry_s3cret"
						tags           = ["terraform-test", "scaleway_rdb_instance", "minimal"]
						volume_type    = "bssd"
						volume_size_in_gb = 10
					}

					resource "scaleway_rdb_snapshot" "test" {
						name        = "test-snapshot"
						instance_id = scaleway_rdb_instance.main.id
						depends_on  = [scaleway_rdb_instance.main]
					}
				`, latestEngineVersion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaleway_rdb_snapshot.test", "name", "test-snapshot"),
					resource.TestCheckResourceAttrSet("scaleway_rdb_snapshot.test", "instance_id"),
					resource.TestCheckResourceAttrSet("scaleway_rdb_snapshot.test", "created_at"),
				),
			},
		},
	})
}

func TestAccRdbSnapshot_Update(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	latestEngineVersion := rdbchecks.GetLatestEngineVersion(tt, postgreSQLEngineName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      IsSnapshotDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "scaleway_rdb_instance" "main" {
						name           = "test-rdb-update"
						node_type      = "db-dev-s"
						engine         = %q
						is_ha_cluster  = false
						disable_backup = true
						user_name      = "my_initial_user"
						password       = "thiZ_is_v&ry_s3cret"
						tags           = ["terraform-test", "scaleway_rdb_instance", "minimal"]
						volume_type    = "bssd"
						volume_size_in_gb = 10
					}

					resource "scaleway_rdb_snapshot" "test" {
						name        = "initial-snapshot"
						instance_id = scaleway_rdb_instance.main.id
						depends_on  = [scaleway_rdb_instance.main]
					}
				`, latestEngineVersion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaleway_rdb_snapshot.test", "name", "initial-snapshot"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "scaleway_rdb_instance" "main" {
						name           = "test-rdb-update"
						node_type      = "db-dev-s"
						engine         = %q
						is_ha_cluster  = false
						disable_backup = true
						user_name      = "my_initial_user"
						password       = "thiZ_is_v&ry_s3cret"
						tags           = ["terraform-test", "scaleway_rdb_instance", "minimal"]
						volume_type    = "bssd"
						volume_size_in_gb = 10
					}

					resource "scaleway_rdb_snapshot" "test" {
						name        = "updated-snapshot"
						instance_id = scaleway_rdb_instance.main.id
						depends_on  = [scaleway_rdb_instance.main]
					}
				`, latestEngineVersion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaleway_rdb_snapshot.test", "name", "updated-snapshot"),
				),
			},
		},
	})
}

func IsSnapshotDestroyed(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_rdb_snapshot" {
				continue
			}

			rdbAPI, region, ID, err := rdb.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}

			_, err = rdbAPI.GetSnapshot(&rdbSDK.GetSnapshotRequest{
				SnapshotID: ID,
				Region:     region,
			})

			if err == nil {
				return fmt.Errorf("snapshot (%s) still exists", rs.Primary.ID)
			}

			if !httperrors.Is404(err) {
				return err
			}
		}

		return nil
	}
}
