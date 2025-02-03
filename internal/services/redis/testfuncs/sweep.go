package redistestfuncs

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	redisSDK "github.com/scaleway/scaleway-sdk-go/api/redis/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/logging"
)

func AddTestSweepers() {
	resource.AddTestSweepers("scaleway_redis_cluster", &resource.Sweeper{
		Name: "scaleway_redis_cluster",
		F:    testSweepRedisCluster,
	})
}

func testSweepRedisCluster(_ string) error {
	return acctest.SweepZones(scw.AllZones, func(scwClient *scw.Client, zone scw.Zone) error {
		redisAPI := redisSDK.NewAPI(scwClient)
		logging.L.Debugf("sweeper: destroying the redis cluster in (%s)", zone)
		listClusters, err := redisAPI.ListClusters(&redisSDK.ListClustersRequest{
			Zone: zone,
		}, scw.WithAllPages())
		if err != nil {
			return fmt.Errorf("error listing redis clusters in (%s) in sweeper: %w", zone, err)
		}

		for _, cluster := range listClusters.Clusters {
			_, err := redisAPI.DeleteCluster(&redisSDK.DeleteClusterRequest{
				Zone:      zone,
				ClusterID: cluster.ID,
			})
			if err != nil {
				return fmt.Errorf("error deleting redis cluster in sweeper: %w", err)
			}
		}

		return nil
	})
}
