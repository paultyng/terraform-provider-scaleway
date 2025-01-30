package k8s_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	vpcchecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/vpc/testfuncs"
)

func TestAccDataSourceCluster_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	clusterName := "tf-cluster"
	version := testAccK8SClusterGetLatestK8SVersion(tt)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckK8SPoolDestroy(tt, "scaleway_k8s_pool.default"),
			testAccCheckK8SClusterDestroy(tt),
			vpcchecks.CheckPrivateNetworkDestroy(tt),
		),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "scaleway_vpc_private_network" "main" {
						name = "test-data-source-cluster"
					}

					resource "scaleway_k8s_cluster" "main" {
					  	name 	= "%s"
						version = "%s"
						cni     = "cilium"
					  	tags    = [ "terraform-test", "data_scaleway_k8s_cluster", "basic" ]
						delete_additional_resources = false
						private_network_id = scaleway_vpc_private_network.main.id
					}

					resource "scaleway_k8s_pool" "default" {
						name = "default"
						cluster_id = "${scaleway_k8s_cluster.main.id}"
						node_type = "pro2_xxs"
						autohealing = true
						autoscaling = true
						size = 1
					}
					
					data "scaleway_k8s_cluster" "prod" {
					  	name = "${scaleway_k8s_cluster.main.name}"
					}
					
					data "scaleway_k8s_cluster" "stg" {
					  	cluster_id = "${scaleway_k8s_cluster.main.id}"
					}`, clusterName, version),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8SClusterExists(tt, "data.scaleway_k8s_cluster.prod"),
					resource.TestCheckResourceAttr("data.scaleway_k8s_cluster.prod", "name", clusterName),
					testAccCheckK8SClusterExists(tt, "data.scaleway_k8s_cluster.stg"),
					resource.TestCheckResourceAttr("data.scaleway_k8s_cluster.stg", "name", clusterName),
				),
			},
		},
	})
}
