package k8stestfuncs

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	k8sSDK "github.com/scaleway/scaleway-sdk-go/api/k8s/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/logging"
)

func AddTestSweepers() {
	resource.AddTestSweepers("scaleway_k8s_cluster", &resource.Sweeper{
		Name: "scaleway_k8s_cluster",
		F:    testSweepK8SCluster,
	})
}

func testSweepK8SCluster(_ string) error {
	return acctest.SweepRegions([]scw.Region{scw.RegionFrPar, scw.RegionNlAms}, func(scwClient *scw.Client, region scw.Region) error {
		k8sAPI := k8sSDK.NewAPI(scwClient)

		logging.L.Debugf("sweeper: destroying the k8s cluster in (%s)", region)
		listClusters, err := k8sAPI.ListClusters(&k8sSDK.ListClustersRequest{Region: region}, scw.WithAllPages())
		if err != nil {
			return fmt.Errorf("error listing clusters in (%s) in sweeper: %s", region, err)
		}

		for _, cluster := range listClusters.Clusters {
			// remove pools
			listPools, err := k8sAPI.ListPools(&k8sSDK.ListPoolsRequest{
				Region:    region,
				ClusterID: cluster.ID,
			}, scw.WithAllPages())
			if err != nil {
				return fmt.Errorf("error listing pool in (%s) in sweeper: %s", region, err)
			}

			for _, pool := range listPools.Pools {
				_, err := k8sAPI.DeletePool(&k8sSDK.DeletePoolRequest{
					Region: region,
					PoolID: pool.ID,
				})
				if err != nil {
					return fmt.Errorf("error deleting pool in sweeper: %s", err)
				}
			}
			_, err = k8sAPI.DeleteCluster(&k8sSDK.DeleteClusterRequest{
				Region:    region,
				ClusterID: cluster.ID,
			})
			if err != nil {
				return fmt.Errorf("error deleting cluster in sweeper: %s", err)
			}
		}

		return nil
	})
}
