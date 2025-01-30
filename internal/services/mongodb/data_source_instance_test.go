package mongodb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

func TestAccDataSourceMongoDBInstance_ByName(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      IsInstanceDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_mongodb_instance" "test" {
						name        = "test-mongodb-instance-by-name"
						version     = "7.0.12"
						node_type   = "MGDB-PLAY2-NANO"
						node_number = 1
						user_name   = "my_initial_user"
						password    = "thiZ_is_v&ry_s3cret"
					}

					data "scaleway_mongodb_instance" "test_by_name" {
						name = scaleway_mongodb_instance.test.name
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.scaleway_mongodb_instance.test_by_name", "name", "test-mongodb-instance-by-name"),
					resource.TestCheckResourceAttr("data.scaleway_mongodb_instance.test_by_name", "version", "7.0.12"),
					resource.TestCheckResourceAttr("data.scaleway_mongodb_instance.test_by_name", "node_type", "mgdb-play2-nano"),
					resource.TestCheckResourceAttr("data.scaleway_mongodb_instance.test_by_name", "node_number", "1"),
				),
			},
		},
	})
}

func TestAccDataSourceMongoDBInstance_ByID(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      IsInstanceDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_mongodb_instance" "test" {
						name        = "test-mongodb-instance-id"
						version     = "7.0.12"
						node_type   = "MGDB-PLAY2-NANO"
						node_number = 1
						user_name   = "my_initial_user"
						password    = "thiZ_is_v&ry_s3cret"
					}

					data "scaleway_mongodb_instance" "test_by_id" {
						instance_id = scaleway_mongodb_instance.test.id
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.scaleway_mongodb_instance.test_by_id", "name", "test-mongodb-instance-id"),
					resource.TestCheckResourceAttr("data.scaleway_mongodb_instance.test_by_id", "version", "7.0.12"),
					resource.TestCheckResourceAttr("data.scaleway_mongodb_instance.test_by_id", "node_type", "mgdb-play2-nano"),
					resource.TestCheckResourceAttr("data.scaleway_mongodb_instance.test_by_id", "node_number", "1"),
				),
			},
		},
	})
}
