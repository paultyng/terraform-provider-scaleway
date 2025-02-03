package redistestfuncs

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	redisSDK "github.com/scaleway/scaleway-sdk-go/api/redis/v1"
	"github.com/scaleway/scaleway-sdk-go/api/redis/v1/sweepers"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

func AddTestSweepers() {
	resource.AddTestSweepers("scaleway_redis_cluster", &resource.Sweeper{
		Name: "scaleway_redis_cluster",
		F:    testSweepRedisCluster,
	})
}

func testSweepRedisCluster(_ string) error {
	return acctest.SweepZones((&redisSDK.API{}).Zones(), sweepers.SweepCluster)
}
