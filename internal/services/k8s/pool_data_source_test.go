package k8s_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	vpcchecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/vpc/testfuncs"
)

func TestAccDataSourcePool_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	clusterName := "tf-cluster-pool"
	poolName := "tf-pool"
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
						name = "test-data-source-pool"
					}

					resource "scaleway_k8s_cluster" "main" {
					  	name 	= "%s"
						version = "%s"
						cni     = "cilium"
					  	tags    = [ "terraform-test", "data_scaleway_k8s_pool", "basic" ]
						delete_additional_resources = false
						private_network_id = scaleway_vpc_private_network.main.id
					}
					
					resource "scaleway_k8s_pool" "default" {
						cluster_id = "${scaleway_k8s_cluster.main.id}"
						name = "%s"
						node_type = "dev1_m"
						size = 1
						tags = [ "terraform-test", "data_scaleway_k8s_pool", "basic" ]
					}`, clusterName, version, poolName),
			},
			{
				Config: fmt.Sprintf(`
					resource "scaleway_vpc_private_network" "main" {
						name = "test-data-source-pool"
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
						cluster_id = "${scaleway_k8s_cluster.main.id}"
						name = "%s"
						node_type = "dev1_m"
						size = 1
						tags = [ "terraform-test", "data_scaleway_k8s_pool", "basic" ]
					}
					
					data "scaleway_k8s_pool" "prod" {
					  	name = "${scaleway_k8s_pool.default.name}"
						cluster_id = "${scaleway_k8s_cluster.main.id}"
					}
					
					data "scaleway_k8s_pool" "stg" {
					  	pool_id = "${scaleway_k8s_pool.default.id}"
					}`, clusterName, version, poolName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8SPoolExists(tt, "data.scaleway_k8s_pool.prod"),
					resource.TestCheckResourceAttr("data.scaleway_k8s_pool.prod", "name", poolName),
					testAccCheckK8SPoolExists(tt, "data.scaleway_k8s_pool.stg"),
					resource.TestCheckResourceAttr("data.scaleway_k8s_pool.stg", "name", poolName),
					resource.TestCheckResourceAttrSet("data.scaleway_k8s_pool.stg", "nodes.0.public_ip"), // Deprecated attributes
				),
			},
		},
	})
}
