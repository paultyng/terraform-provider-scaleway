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

func TestAccMongoDBInstance_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      IsInstanceDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_mongodb_instance main {
						name = "test-mongodb-basic1"
						version = "7.0.12"
						node_type = "MGDB-PLAY2-NANO"
						node_number = 1
						user_name = "my_initial_user"
						password = "thiZ_is_v&ry_s3cret"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isMongoDBInstancePresent(tt, "scaleway_mongodb_instance.main"),
					resource.TestCheckResourceAttr("scaleway_mongodb_instance.main", "name", "test-mongodb-basic1"),
					resource.TestCheckResourceAttr("scaleway_mongodb_instance.main", "node_type", "mgdb-play2-nano"),
					resource.TestCheckResourceAttr("scaleway_mongodb_instance.main", "version", "7.0.12"),
					resource.TestCheckResourceAttr("scaleway_mongodb_instance.main", "node_number", "1"),
					resource.TestCheckResourceAttr("scaleway_mongodb_instance.main", "user_name", "my_initial_user"),
					resource.TestCheckResourceAttr("scaleway_mongodb_instance.main", "password", "thiZ_is_v&ry_s3cret"),
				),
			},
		},
	})
}

func TestAccMongoDBInstance_VolumeUpdate(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      IsInstanceDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_mongodb_instance main {
						name = "test-mongodb-volume-update1"
						version = "7.0.12"
						node_type = "MGDB-PLAY2-NANO"
						node_number = "1"
						user_name = "my_initial_user"
						password = "thiZ_is_v&ry_s3cret"
						volume_size_in_gb = 5
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isMongoDBInstancePresent(tt, "scaleway_mongodb_instance.main"),
					resource.TestCheckResourceAttr("scaleway_mongodb_instance.main", "volume_size_in_gb", "5"),
				),
			},
			{
				Config: `
					resource scaleway_mongodb_instance main {
						name = "test-mongodb-volume-update1"
						version = "7.0.12"
						node_type = "MGDB-PLAY2-NANO"
						node_number = "1"
						user_name = "my_initial_user"
						password = "thiZ_is_v&ry_s3cret"
						volume_size_in_gb = 10
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isMongoDBInstancePresent(tt, "scaleway_mongodb_instance.main"),
					resource.TestCheckResourceAttr("scaleway_mongodb_instance.main", "volume_size_in_gb", "10"),
				),
			},
		},
	})
}

func TestAccMongoDBInstance_UpdateNameTagsUser(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      IsInstanceDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_mongodb_instance main {
						name = "test-mongodb-update-initial"
						version = "7.0.12"
						node_type = "MGDB-PLAY2-NANO"
						node_number = 1
						user_name = "user"
						password = "initial_password"
						tags = ["initial_tag1", "initial_tag2"]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isMongoDBInstancePresent(tt, "scaleway_mongodb_instance.main"),
					resource.TestCheckResourceAttr("scaleway_mongodb_instance.main", "name", "test-mongodb-update-initial"),
					resource.TestCheckResourceAttr("scaleway_mongodb_instance.main", "user_name", "user"),
					resource.TestCheckResourceAttr("scaleway_mongodb_instance.main", "tags.#", "2"),
					resource.TestCheckResourceAttr("scaleway_mongodb_instance.main", "tags.0", "initial_tag1"),
					resource.TestCheckResourceAttr("scaleway_mongodb_instance.main", "tags.1", "initial_tag2"),
				),
			},
			{
				Config: `
					resource scaleway_mongodb_instance main {
						name = "test-mongodb-update-final"
						version = "7.0.12"
						node_type = "MGDB-PLAY2-NANO"
						node_number = 1
						user_name = "user"
						password = "updated_password"
						tags = ["updated_tag1", "updated_tag2", "updated_tag3"]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isMongoDBInstancePresent(tt, "scaleway_mongodb_instance.main"),
					resource.TestCheckResourceAttr("scaleway_mongodb_instance.main", "name", "test-mongodb-update-final"),
					resource.TestCheckResourceAttr("scaleway_mongodb_instance.main", "tags.#", "3"),
					resource.TestCheckResourceAttr("scaleway_mongodb_instance.main", "tags.0", "updated_tag1"),
					resource.TestCheckResourceAttr("scaleway_mongodb_instance.main", "tags.1", "updated_tag2"),
					resource.TestCheckResourceAttr("scaleway_mongodb_instance.main", "tags.2", "updated_tag3"),
				),
			},
		},
	})
}

func TestAccMongoDBInstance_FromSnapshot(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      IsInstanceDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_mongodb_instance" "main" {
						name        = "test-mongodb-from-snapshot"
						version     = "7.0.12"
						node_type   = "MGDB-PLAY2-NANO"
						node_number = 1
						user_name   = "my_initial_user"
						password    = "thiZ_is_v&ry_s3cret"
					}

					resource "scaleway_mongodb_snapshot" "main_snapshot" {
						instance_id = scaleway_mongodb_instance.main.id
						name        = "test-snapshot"
						expires_at  = "2024-12-31T23:59:59Z"
 						depends_on = [
    						scaleway_mongodb_instance.main
  						]
					}
					
					resource "scaleway_mongodb_instance" "restored_instance" {
						snapshot_id = scaleway_mongodb_snapshot.main_snapshot.id
						name        = "restored-mongodb-from-snapshot"
						node_type   = "MGDB-PLAY2-NANO"
						node_number = 1
						depends_on = [
    						scaleway_mongodb_instance.main,
    						scaleway_mongodb_snapshot.main_snapshot
  						]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isMongoDBInstancePresent(tt, "scaleway_mongodb_instance.restored_instance"),
				),
			},
		},
	})
}

func isMongoDBInstancePresent(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		mongodbAPI, region, ID, err := mongodb.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = mongodbAPI.GetInstance(&mongodbSDK.GetInstanceRequest{
			InstanceID: ID,
			Region:     region,
		})
		if err != nil {
			return err
		}

		return nil
	}
}

func IsInstanceDestroyed(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_mongodb_instance" {
				continue
			}

			mongodbAPI, zone, ID, err := mongodb.NewAPIWithZoneAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}
			extractRegion, err := zone.Region()
			if err != nil {
				return err
			}
			_, err = mongodbAPI.GetInstance(&mongodbSDK.GetInstanceRequest{
				InstanceID: ID,
				Region:     extractRegion,
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
