package mongodb_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	mongodbSDK "github.com/scaleway/scaleway-sdk-go/api/mongodb/v1alpha1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/mongodb"
)

func TestAccMongoDBSnapshot_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isSnapshotDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_mongodb_instance" "main" {
						name       = "test-mongodb-instance"
						version    = "7.0.12"
						node_type  = "MGDB-PLAY2-NANO"
						node_number = 1
						user_name  = "my_initial_user"
						password   = "thiZ_is_v&ry_s3cret"
					}

					resource "scaleway_mongodb_snapshot" "main" {
						instance_id = scaleway_mongodb_instance.main.id
						name        = "test-snapshot"
						expires_at  = "2025-12-31T23:59:59Z"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaleway_mongodb_snapshot.main", "name", "test-snapshot"),
				),
			},
		},
	})
}

func TestAccMongoDBSnapshot_Update(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isSnapshotDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_mongodb_instance" "main" {
						name       = "test-mongodb-instance"
						version    = "7.0.12"
						node_type  = "MGDB-PLAY2-NANO"
						node_number = 1
						user_name  = "my_initial_user"
						password   = "thiZ_is_v&ry_s3cret"
					}

					resource "scaleway_mongodb_snapshot" "main" {
						instance_id = scaleway_mongodb_instance.main.id
						name        = "test-snapshot"
						expires_at  = "2025-12-31T23:59:59Z"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaleway_mongodb_snapshot.main", "expires_at", "2025-12-31T23:59:59Z"),
				),
			},
			{
				Config: `
					resource "scaleway_mongodb_instance" "main" {
						name       = "test-mongodb-instance"
						version    = "7.0.12"
						node_type  = "MGDB-PLAY2-NANO"
						node_number = 1
						user_name  = "my_initial_user"
						password   = "thiZ_is_v&ry_s3cret"
					}

					resource "scaleway_mongodb_snapshot" "main" {
						instance_id = scaleway_mongodb_instance.main.id
						name        = "updated-snapshot"
						expires_at  = "2025-09-20T23:59:59Z"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaleway_mongodb_snapshot.main", "name", "updated-snapshot"),
					resource.TestCheckResourceAttr("scaleway_mongodb_snapshot.main", "expires_at", "2025-09-20T23:59:59Z"),
				),
			},
		},
	})
}

func isSnapshotDestroyed(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_mongodb_snapshot" {
				continue
			}

			mongodbAPI, region, ID, err := mongodb.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}

			_, err = mongodbAPI.GetSnapshot(&mongodbSDK.GetSnapshotRequest{
				SnapshotID: ID,
				Region:     region,
			})
			if err == nil {
				return fmt.Errorf("instance (%s) still exists", rs.Primary.ID)
			}
			if !httperrors.Is404(err) {
				return err
			}
		}
		return nil
	}
}
