package redis_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

func TestAccDataSourceCluster_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	latestRedisVersion := getLatestVersion(tt)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isClusterDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "scaleway_redis_cluster" "test" {
    					name = "test_redis_datasource_terraform"
    					version = "%s"
    					node_type = "RED1-micro"
    					user_name = "my_initial_user"
    					password = "thiZ_is_v&ry_s3cret"
					}

					data "scaleway_redis_cluster" "test" {
						name = scaleway_redis_cluster.test.name
					}

					data "scaleway_redis_cluster" "test2" {
						cluster_id = scaleway_redis_cluster.test.id
					}
				`, latestRedisVersion),
				Check: resource.ComposeTestCheckFunc(
					isClusterPresent(tt, "scaleway_redis_cluster.test"),

					resource.TestCheckResourceAttr("data.scaleway_redis_cluster.test", "name", "test_redis_datasource_terraform"),
					resource.TestCheckResourceAttrSet("data.scaleway_redis_cluster.test", "id"),

					resource.TestCheckResourceAttr("data.scaleway_redis_cluster.test2", "name", "test_redis_datasource_terraform"),
					resource.TestCheckResourceAttrSet("data.scaleway_redis_cluster.test2", "id"),
				),
			},
		},
	})
}
