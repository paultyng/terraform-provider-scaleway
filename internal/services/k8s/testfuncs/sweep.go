package k8stestfuncs

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	k8sSDK "github.com/scaleway/scaleway-sdk-go/api/k8s/v1"
	"github.com/scaleway/scaleway-sdk-go/api/k8s/v1/sweepers"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

func AddTestSweepers() {
	resource.AddTestSweepers("scaleway_k8s_cluster", &resource.Sweeper{
		Name: "scaleway_k8s_cluster",
		F:    testSweepK8SCluster,
	})
}

func testSweepK8SCluster(_ string) error {
	return acctest.SweepRegions((&k8sSDK.API{}).Regions(), sweepers.SweepCluster)
}
