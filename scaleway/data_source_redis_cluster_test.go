package scaleway

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccScalewayDataSourceRedisCluster_Basic(t *testing.T) {
	tt := NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckScalewayRedisClusterDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_redis_cluster" "test" {
    					name = "test_redis_datasource_terraform"
    					version = "6.2.6"
    					node_type = "MDB-BETA-M"
    					user_name = "my_initial_user"
    					password = "thiZ_is_v&ry_s3cret"
					}

					data "scaleway_redis_cluster" "test" {
						name = scaleway_redis_cluster.test.name
					}

					data "scaleway_redis_cluster" "test2" {
						cluster_id = scaleway_redis_cluster.test.id
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayRedisExists(tt, "scaleway_redis_cluster.test"),

					resource.TestCheckResourceAttr("data.scaleway_redis_cluster.test", "name", "test_redis_datasource_terraform"),
					resource.TestCheckResourceAttrSet("data.scaleway_redis_cluster.test", "id"),

					resource.TestCheckResourceAttr("data.scaleway_redis_cluster.test2", "name", "test_redis_datasource_terraform"),
					resource.TestCheckResourceAttrSet("data.scaleway_redis_cluster.test2", "id"),
				),
			},
		},
	})
}
