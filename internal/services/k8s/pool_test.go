package k8s_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	k8sSDK "github.com/scaleway/scaleway-sdk-go/api/k8s/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/k8s"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/vpc"
	vpcchecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/vpc/testfuncs"
	vpcgwchecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/vpcgw/testfuncs"
)

func TestAccPool_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	latestK8SVersion := testAccK8SClusterGetLatestK8SVersion(tt)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckK8SPoolDestroy(tt, "scaleway_k8s_pool.default"),
			testAccCheckK8SPoolDestroy(tt, "scaleway_k8s_pool.minimal"),
			testAccCheckK8SClusterDestroy(tt),
			vpcchecks.CheckPrivateNetworkDestroy(tt),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckK8SPoolConfigMinimal(latestK8SVersion, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8SClusterExists(tt, "scaleway_k8s_cluster.minimal"),
					testAccCheckK8SPoolExists(tt, "scaleway_k8s_pool.default"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.default", "node_type", "pro2_xxs"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.default", "size", "1"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.default", "autohealing", "true"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.default", "autoscaling", "true"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.default", "version", latestK8SVersion),
					resource.TestCheckResourceAttrSet("scaleway_k8s_pool.default", "id"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.default", "tags.0", "terraform-test"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.default", "tags.1", "scaleway_k8s_cluster"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.default", "tags.2", "default"),
					testAccCheckK8SPoolServersAreInPrivateNetwork(tt, "scaleway_k8s_cluster.minimal", "scaleway_k8s_pool.default", "scaleway_vpc_private_network.minimal"),
					resource.TestCheckResourceAttrSet("scaleway_k8s_pool.default", "private_ip.0.id"),
					resource.TestCheckResourceAttrSet("scaleway_k8s_pool.default", "private_ip.0.address"),
				),
			},
			{
				Config: testAccCheckK8SPoolConfigMinimal(latestK8SVersion, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8SClusterExists(tt, "scaleway_k8s_cluster.minimal"),
					testAccCheckK8SPoolExists(tt, "scaleway_k8s_pool.minimal"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.minimal", "node_type", "pro2_xxs"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.minimal", "size", "1"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.minimal", "autohealing", "true"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.minimal", "autoscaling", "true"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.minimal", "version", latestK8SVersion),
					resource.TestCheckResourceAttrSet("scaleway_k8s_pool.minimal", "id"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.minimal", "tags.0", "terraform-test"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.minimal", "tags.1", "scaleway_k8s_cluster"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.minimal", "tags.2", "minimal"),
					resource.TestCheckResourceAttrSet("scaleway_k8s_pool.minimal", "nodes.0.public_ip"), // Deprecated attributes
					testAccCheckK8SPoolServersAreInPrivateNetwork(tt, "scaleway_k8s_cluster.minimal", "scaleway_k8s_pool.default", "scaleway_vpc_private_network.minimal"),
					testAccCheckK8SPoolServersAreInPrivateNetwork(tt, "scaleway_k8s_cluster.minimal", "scaleway_k8s_pool.minimal", "scaleway_vpc_private_network.minimal"),
					resource.TestCheckResourceAttrSet("scaleway_k8s_pool.default", "private_ip.0.id"),
					resource.TestCheckResourceAttrSet("scaleway_k8s_pool.default", "private_ip.0.address"),
					resource.TestCheckResourceAttrSet("scaleway_k8s_pool.default", "private_ip.1.id"),
					resource.TestCheckResourceAttrSet("scaleway_k8s_pool.default", "private_ip.1.address"),
				),
			},
			{
				Config: testAccCheckK8SPoolConfigMinimal(latestK8SVersion, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8SClusterExists(tt, "scaleway_k8s_cluster.minimal"),
					testAccCheckK8SPoolExists(tt, "scaleway_k8s_pool.default"),
					testAccCheckK8SPoolDestroy(tt, "scaleway_k8s_pool.minimal"),
					testAccCheckK8SPoolServersAreInPrivateNetwork(tt, "scaleway_k8s_cluster.minimal", "scaleway_k8s_pool.default", "scaleway_vpc_private_network.minimal"),
				),
			},
		},
	})
}

func TestAccPool_Wait(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	latestK8SVersion := testAccK8SClusterGetLatestK8SVersion(tt)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckK8SPoolDestroy(tt, "scaleway_k8s_pool.default"),
			testAccCheckK8SPoolDestroy(tt, "scaleway_k8s_pool.minimal"),
			testAccCheckK8SClusterDestroy(tt),
			vpcchecks.CheckPrivateNetworkDestroy(tt),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckK8SPoolConfigWait(latestK8SVersion, false, 0),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8SClusterExists(tt, "scaleway_k8s_cluster.minimal"),
					testAccCheckK8SPoolExists(tt, "scaleway_k8s_pool.default"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.default", "size", "1"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.default", "min_size", "1"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.default", "max_size", "1"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.default", "status", k8sSDK.PoolStatusReady.String()),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.default", "nodes.0.status", k8sSDK.NodeStatusReady.String()),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.default", "wait_for_pool_ready", "true"),
				),
			},
			{
				Config: testAccCheckK8SPoolConfigWait(latestK8SVersion, true, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8SClusterExists(tt, "scaleway_k8s_cluster.minimal"),
					testAccCheckK8SPoolExists(tt, "scaleway_k8s_pool.minimal"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.minimal", "size", "1"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.minimal", "min_size", "1"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.minimal", "max_size", "1"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.minimal", "status", k8sSDK.PoolStatusReady.String()),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.minimal", "nodes.0.status", k8sSDK.NodeStatusReady.String()),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.minimal", "wait_for_pool_ready", "true"),
				),
			},
			{
				Config: testAccCheckK8SPoolConfigWait(latestK8SVersion, true, 2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8SClusterExists(tt, "scaleway_k8s_cluster.minimal"),
					testAccCheckK8SPoolExists(tt, "scaleway_k8s_pool.minimal"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.minimal", "size", "2"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.minimal", "min_size", "1"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.minimal", "max_size", "2"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.minimal", "status", k8sSDK.PoolStatusReady.String()),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.minimal", "nodes.0.status", k8sSDK.NodeStatusReady.String()),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.minimal", "nodes.1.status", k8sSDK.NodeStatusReady.String()), // check that the new node has the "ready" status
					resource.TestCheckResourceAttr("scaleway_k8s_pool.minimal", "wait_for_pool_ready", "true"),
				),
			},
			{
				Config: testAccCheckK8SPoolConfigWait(latestK8SVersion, true, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8SClusterExists(tt, "scaleway_k8s_cluster.minimal"),
					testAccCheckK8SPoolExists(tt, "scaleway_k8s_pool.minimal"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.minimal", "size", "1"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.minimal", "min_size", "1"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.minimal", "max_size", "1"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.minimal", "status", k8sSDK.PoolStatusReady.String()),
					testAccCheckK8SPoolNodesOneOfIsDeleting("scaleway_k8s_pool.minimal"),        // check that one of the nodes is deleting (nodes are not ordered)
					resource.TestCheckResourceAttr("scaleway_k8s_pool.minimal", "nodes.#", "2"), // the node that is deleting should still exist
					resource.TestCheckResourceAttr("scaleway_k8s_pool.minimal", "wait_for_pool_ready", "true"),
				),
			},
			{
				Config: testAccCheckK8SPoolConfigWait(latestK8SVersion, false, 0),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8SClusterExists(tt, "scaleway_k8s_cluster.minimal"),
					testAccCheckK8SPoolDestroy(tt, "scaleway_k8s_pool.minimal"),
				),
			},
		},
	})
}

func TestAccPool_PlacementGroup(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	latestK8SVersion := testAccK8SClusterGetLatestK8SVersion(tt)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckK8SPoolDestroy(tt, "scaleway_k8s_pool.placement_group"),
			testAccCheckK8SPoolDestroy(tt, "scaleway_k8s_pool.placement_group_2"),
			testAccCheckK8SClusterDestroy(tt),
			vpcchecks.CheckPrivateNetworkDestroy(tt),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckK8SPoolConfigPlacementGroup(latestK8SVersion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8SClusterExists(tt, "scaleway_k8s_cluster.placement_group"),
					testAccCheckK8SPoolExists(tt, "scaleway_k8s_pool.placement_group"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.placement_group", "node_type", "pro2_xxs"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.placement_group", "size", "1"),
					resource.TestCheckResourceAttrSet("scaleway_k8s_pool.placement_group", "id"),
					resource.TestCheckResourceAttrSet("scaleway_k8s_pool.placement_group", "placement_group_id"),
				),
			},
			{
				Config: testAccCheckK8SPoolConfigPlacementGroupWithCustomZone(latestK8SVersion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8SClusterExists(tt, "scaleway_k8s_cluster.placement_group_2"),
					testAccCheckK8SPoolExists(tt, "scaleway_k8s_pool.placement_group_2"),
					resource.TestCheckResourceAttrPair("scaleway_k8s_pool.placement_group_2", "placement_group_id", "scaleway_instance_placement_group.placement_group", "id"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.placement_group_2", "zone", "nl-ams-2"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.placement_group_2", "node_type", "pro2_xxs"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.placement_group_2", "size", "1"),
					resource.TestCheckResourceAttrSet("scaleway_k8s_pool.placement_group_2", "id"),
					resource.TestCheckResourceAttrSet("scaleway_k8s_pool.placement_group_2", "placement_group_id"),
				),
			},
			{
				Config: testAccCheckK8SPoolConfigPlacementGroupWithMultiZone(latestK8SVersion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8SClusterExists(tt, "scaleway_k8s_cluster.placement_group_2"),
					testAccCheckK8SPoolExists(tt, "scaleway_k8s_pool.placement_group_2"),
					resource.TestCheckResourceAttrPair("scaleway_k8s_pool.placement_group_2", "placement_group_id", "scaleway_instance_placement_group.placement_group", "id"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.placement_group_2", "zone", "nl-ams-1"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.placement_group_2", "node_type", "pro2_xxs"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.placement_group_2", "size", "1"),
					resource.TestCheckResourceAttrSet("scaleway_k8s_pool.placement_group_2", "id"),
					resource.TestCheckResourceAttrSet("scaleway_k8s_pool.placement_group_2", "placement_group_id"),
				),
			},
		},
	})
}

func TestAccPool_UpgradePolicy(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	latestK8SVersion := testAccK8SClusterGetLatestK8SVersion(tt)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckK8SPoolDestroy(tt, "scaleway_k8s_pool.upgrade_policy"),
			testAccCheckK8SClusterDestroy(tt),
			vpcchecks.CheckPrivateNetworkDestroy(tt),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckK8SPoolConfigUpgradePolicy(latestK8SVersion, 2, 3),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8SClusterExists(tt, "scaleway_k8s_cluster.upgrade_policy"),
					testAccCheckK8SPoolExists(tt, "scaleway_k8s_pool.upgrade_policy"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.upgrade_policy", "node_type", "pro2_xxs"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.upgrade_policy", "size", "1"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.upgrade_policy", "autohealing", "true"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.upgrade_policy", "autoscaling", "true"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.upgrade_policy", "version", latestK8SVersion),
					resource.TestCheckResourceAttrSet("scaleway_k8s_pool.upgrade_policy", "id"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.upgrade_policy", "tags.0", "terraform-test"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.upgrade_policy", "tags.1", "scaleway_k8s_cluster"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.upgrade_policy", "tags.2", "upgrade_policy"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.upgrade_policy", "upgrade_policy.0.max_surge", "2"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.upgrade_policy", "upgrade_policy.0.max_unavailable", "3"),
				),
			},
			{
				Config: testAccCheckK8SPoolConfigUpgradePolicy(latestK8SVersion, 0, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8SClusterExists(tt, "scaleway_k8s_cluster.upgrade_policy"),
					testAccCheckK8SPoolExists(tt, "scaleway_k8s_pool.upgrade_policy"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.upgrade_policy", "node_type", "pro2_xxs"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.upgrade_policy", "size", "1"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.upgrade_policy", "autohealing", "true"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.upgrade_policy", "autoscaling", "true"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.upgrade_policy", "version", latestK8SVersion),
					resource.TestCheckResourceAttrSet("scaleway_k8s_pool.upgrade_policy", "id"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.upgrade_policy", "tags.0", "terraform-test"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.upgrade_policy", "tags.1", "scaleway_k8s_cluster"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.upgrade_policy", "tags.2", "upgrade_policy"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.upgrade_policy", "upgrade_policy.0.max_surge", "0"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.upgrade_policy", "upgrade_policy.0.max_unavailable", "1"),
				),
			},
		},
	})
}

func TestAccPool_KubeletArgs(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	latestK8SVersion := testAccK8SClusterGetLatestK8SVersion(tt)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckK8SPoolDestroy(tt, "scaleway_k8s_pool.kubelet_args"),
			testAccCheckK8SClusterDestroy(tt),
			vpcchecks.CheckPrivateNetworkDestroy(tt),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckK8SPoolConfigKubeletArgs(latestK8SVersion, 1337),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8SClusterExists(tt, "scaleway_k8s_cluster.kubelet_args"),
					testAccCheckK8SPoolExists(tt, "scaleway_k8s_pool.kubelet_args"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.kubelet_args", "version", latestK8SVersion),
					resource.TestCheckResourceAttrSet("scaleway_k8s_pool.kubelet_args", "id"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.kubelet_args", "kubelet_args.maxPods", "1337"),
				),
			},
			{
				Config: testAccCheckK8SPoolConfigKubeletArgs(latestK8SVersion, 50),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8SClusterExists(tt, "scaleway_k8s_cluster.kubelet_args"),
					testAccCheckK8SPoolExists(tt, "scaleway_k8s_pool.kubelet_args"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.kubelet_args", "version", latestK8SVersion),
					resource.TestCheckResourceAttrSet("scaleway_k8s_pool.kubelet_args", "id"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.kubelet_args", "kubelet_args.maxPods", "50"),
				),
			},
		},
	})
}

func TestAccPool_Zone(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	latestK8SVersion := testAccK8SClusterGetLatestK8SVersion(tt)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckK8SPoolDestroy(tt, "scaleway_k8s_pool.zone"),
			testAccCheckK8SClusterDestroy(tt),
			vpcchecks.CheckPrivateNetworkDestroy(tt),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckK8SPoolConfigZone(latestK8SVersion, "fr-par-2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8SClusterExists(tt, "scaleway_k8s_cluster.zone"),
					testAccCheckK8SPoolExists(tt, "scaleway_k8s_pool.zone"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.zone", "version", latestK8SVersion),
					resource.TestCheckResourceAttrSet("scaleway_k8s_pool.zone", "id"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.zone", "zone", "fr-par-2"),
				),
			},
		},
	})
}

func TestAccPool_Size(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	latestK8SVersionMinor := testAccK8SClusterGetLatestK8SVersionMinor(tt)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckK8SPoolDestroy(tt, "scaleway_k8s_pool.pool"),
			testAccCheckK8SClusterDestroy(tt),
			vpcchecks.CheckPrivateNetworkDestroy(tt),
		),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				resource "scaleway_vpc_private_network" "test-pool-size" {
					name = "test-pool-size"
				}
				resource "scaleway_k8s_cluster" "test-pool-size" {
				  name    = "test-pool-size"
				  version = "%s"
				  cni     = "cilium"
				  delete_additional_resources = false
				  private_network_id = scaleway_vpc_private_network.test-pool-size.id
				  auto_upgrade {
				    enable = true
				    maintenance_window_start_hour = 12
				    maintenance_window_day = "monday"
				  }
				}
				
				resource "scaleway_k8s_pool" "pool" {
				  cluster_id          = scaleway_k8s_cluster.test-pool-size.id
				  name                = "pool"
				  node_type           = "pro2_xxs"
				  size                = 1
				  autoscaling         = false
				  autohealing         = true
				  wait_for_pool_ready = true
				}`, latestK8SVersionMinor),
			},
			{
				Config: fmt.Sprintf(`
				resource "scaleway_vpc_private_network" "test-pool-size" {
					name = "test-pool-size"
				}
				resource "scaleway_k8s_cluster" "test-pool-size" {
				  name    = "test-pool-size"
				  version = "%s"
				  cni     = "cilium"
				  delete_additional_resources = false
				  private_network_id = scaleway_vpc_private_network.test-pool-size.id
				  auto_upgrade {
				    enable = true
				    maintenance_window_start_hour = 12
				    maintenance_window_day = "monday"
				  }
				}
				
				resource "scaleway_k8s_pool" "pool" {
				  cluster_id          = scaleway_k8s_cluster.test-pool-size.id
				  name                = "pool"
				  node_type           = "pro2_xxs"
				  size                = 2
				  autoscaling         = false
				  autohealing         = true
				  wait_for_pool_ready = true
				}`, latestK8SVersionMinor),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccPool_PublicIPDisabled(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	latestK8SVersion := testAccK8SClusterGetLatestK8SVersion(tt)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckK8SPoolDestroy(tt, "scaleway_k8s_pool.public_ip"),
			testAccCheckK8SClusterDestroy(tt),
			vpcgwchecks.IsGatewayNetworkDestroyed(tt),
			vpcgwchecks.IsDHCPDestroyed(tt),
			vpcgwchecks.IsGatewayDestroyed(tt),
			vpcchecks.CheckPrivateNetworkDestroy(tt),
		),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				resource "scaleway_vpc_private_network" "public_ip" {
				  name       = "test-k8s-public-ip"
				}
			
				resource "scaleway_k8s_cluster" "public_ip" {
				  name = "test-k8s-public-ip"
				  version = "%s"
				  cni     = "cilium"
				  private_network_id = scaleway_vpc_private_network.public_ip.id
				  tags = [ "terraform-test", "scaleway_k8s_cluster", "public_ip" ]
				  delete_additional_resources = false
				  depends_on = [scaleway_vpc_private_network.public_ip]
				}
			
				resource "scaleway_k8s_pool" "public_ip" {
				  cluster_id          = scaleway_k8s_cluster.public_ip.id
				  name                = "test-k8s-public-ip"
				  node_type           = "pro2_xxs"
				  size                = 1
				  autoscaling         = false
				  autohealing         = true
				  wait_for_pool_ready = true
				}`, latestK8SVersion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8SClusterExists(tt, "scaleway_k8s_cluster.public_ip"),
					vpcchecks.IsPrivateNetworkPresent(tt, "scaleway_vpc_private_network.public_ip"),
					testAccCheckK8SPoolExists(tt, "scaleway_k8s_pool.public_ip"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.public_ip", "public_ip_disabled", "false"),
					testAccCheckK8SPoolPublicIP(tt, "scaleway_k8s_cluster.public_ip", "scaleway_k8s_pool.public_ip", false),
				),
			},
			{
				Config: fmt.Sprintf(`
				resource "scaleway_vpc_private_network" "public_ip" {
				  name       = "test-k8s-public-ip"
				}
				resource "scaleway_vpc_public_gateway" "public_ip" {
   				  name = "test-k8s-public-ip"
    			  type = "VPC-GW-S"
				}
				resource "scaleway_vpc_public_gateway_dhcp" "public_ip" {
				  subnet = "192.168.0.0/22"
				  push_default_route = true
				}
				resource "scaleway_vpc_gateway_network" "public_ip" {
				  gateway_id = scaleway_vpc_public_gateway.public_ip.id
				  private_network_id = scaleway_vpc_private_network.public_ip.id
				  dhcp_id = scaleway_vpc_public_gateway_dhcp.public_ip.id
				}

				resource "scaleway_k8s_cluster" "public_ip" {
				  name = "test-k8s-public-ip"
				  version = "%s"
				  cni     = "cilium"
				  private_network_id = scaleway_vpc_private_network.public_ip.id
				  tags = [ "terraform-test", "scaleway_k8s_cluster", "public_ip" ]
				  delete_additional_resources = false
				  depends_on = [
					scaleway_vpc_private_network.public_ip,
					scaleway_vpc_gateway_network.public_ip,
				  ]
				}

				resource "scaleway_k8s_pool" "public_ip" {
				  cluster_id          = scaleway_k8s_cluster.public_ip.id
				  name                = "test-k8s-public-ip"
				  node_type           = "pro2_xxs"
				  size                = 1
				  autoscaling         = false
				  autohealing         = true
				  wait_for_pool_ready = true
				  public_ip_disabled  = true
				}`, latestK8SVersion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8SClusterExists(tt, "scaleway_k8s_cluster.public_ip"),
					vpcchecks.IsPrivateNetworkPresent(tt, "scaleway_vpc_private_network.public_ip"),
					testAccCheckK8SPoolExists(tt, "scaleway_k8s_pool.public_ip"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.public_ip", "public_ip_disabled", "true"),
					testAccCheckK8SPoolPublicIP(tt, "scaleway_k8s_cluster.public_ip", "scaleway_k8s_pool.public_ip", true),
				),
			},
		},
	})
}

func testAccCheckK8SPoolServersAreInPrivateNetwork(tt *acctest.TestTools, clusterTFName, poolTFName, pnTFName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[clusterTFName]
		if !ok {
			return fmt.Errorf("resource not found: %s", clusterTFName)
		}
		k8sAPI, region, clusterID, err := k8s.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		rs, ok = s.RootModule().Resources[poolTFName]
		if !ok {
			return fmt.Errorf("resource not found: %s", poolTFName)
		}
		_, _, poolID, err := k8s.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		rs, ok = s.RootModule().Resources[pnTFName]
		if !ok {
			return fmt.Errorf("resource not found: %s", pnTFName)
		}
		_, _, pnID, err := vpc.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		nodes, err := k8sAPI.ListNodes(&k8sSDK.ListNodesRequest{
			Region:    region,
			PoolID:    &poolID,
			ClusterID: clusterID,
		})
		if err != nil {
			return err
		}

		instanceAPI := instance.NewAPI(tt.Meta.ScwClient())

		for _, node := range nodes.Nodes {
			providerIDSplit := strings.SplitN(node.ProviderID, "/", 5)
			// node.ProviderID is of the form scaleway://instance/<zone>/<id>
			if len(providerIDSplit) < 5 {
				return fmt.Errorf("unexpected format for ProviderID in node %s", node.ID)
			}

			server, err := instanceAPI.GetServer(&instance.GetServerRequest{
				Zone:     scw.Zone(providerIDSplit[3]),
				ServerID: providerIDSplit[4],
			})
			if err != nil {
				return err
			}

			pnfound := false
			for _, privateNic := range server.Server.PrivateNics {
				if privateNic.PrivateNetworkID == pnID {
					pnfound = true
				}
			}
			if pnfound == false {
				return fmt.Errorf("node %s is not in linked to private network %s", node.ID, pnID)
			}
		}

		return nil
	}
}

func testAccCheckK8SPoolPublicIP(tt *acctest.TestTools, clusterTFName, poolTFName string, disabled bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[clusterTFName]
		if !ok {
			return fmt.Errorf("resource not found: %s", clusterTFName)
		}
		k8sAPI, region, clusterID, err := k8s.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		rs, ok = s.RootModule().Resources[poolTFName]
		if !ok {
			return fmt.Errorf("resource not found: %s", poolTFName)
		}
		_, _, poolID, err := k8s.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		nodes, err := k8sAPI.ListNodes(&k8sSDK.ListNodesRequest{
			Region:    region,
			PoolID:    &poolID,
			ClusterID: clusterID,
		})
		if err != nil {
			return err
		}

		instanceAPI := instance.NewAPI(tt.Meta.ScwClient())

		for _, node := range nodes.Nodes {
			providerIDSplit := strings.SplitN(node.ProviderID, "/", 5)
			// node.ProviderID is of the form scaleway://instance/<zone>/<id>
			if len(providerIDSplit) < 5 {
				return fmt.Errorf("unexpected format for ProviderID in node %s", node.ID)
			}

			server, err := instanceAPI.GetServer(&instance.GetServerRequest{
				Zone:     scw.Zone(providerIDSplit[3]),
				ServerID: providerIDSplit[4],
			})
			if err != nil {
				return err
			}

			if disabled == true && server.Server.PublicIPs != nil && len(server.Server.PublicIPs) > 0 {
				return errors.New("found node with public IP when none was expected")
			}
			if disabled == false && len(server.Server.PublicIPs) == 0 {
				return errors.New("found node with no public IP when one was expected")
			}
		}

		return nil
	}
}

func testAccCheckK8SPoolDestroy(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return nil
		}

		k8sAPI, region, poolID, err := k8s.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = k8sAPI.GetPool(&k8sSDK.GetPoolRequest{
			Region: region,
			PoolID: poolID,
		})
		// If no error resource still exist
		if err == nil {
			return fmt.Errorf("pool (%s) still exists", rs.Primary.ID)
		}
		// Unexpected api error we return it
		if !httperrors.Is404(err) {
			return err
		}

		return nil
	}
}

func testAccCheckK8SPoolExists(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		k8sAPI, region, poolID, err := k8s.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = k8sAPI.GetPool(&k8sSDK.GetPoolRequest{
			Region: region,
			PoolID: poolID,
		})
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckK8SPoolConfigMinimal(version string, otherPool bool) string {
	pool := ""
	if otherPool {
		pool += `
resource "scaleway_k8s_pool" "minimal" {
    name = "test-pool-minimal-2"
	cluster_id = "${scaleway_k8s_cluster.minimal.id}"
	node_type = "pro2_xxs"
	autohealing = true
	autoscaling = true
	size = 1
	tags = [ "terraform-test", "scaleway_k8s_cluster", "minimal" ]
}`
	}

	return fmt.Sprintf(`
resource "scaleway_k8s_pool" "default" {
    name = "test-pool-minimal"
	cluster_id = "${scaleway_k8s_cluster.minimal.id}"
	node_type = "pro2_xxs"
	autohealing = true
	autoscaling = true
	size = 1
	tags = [ "terraform-test", "scaleway_k8s_cluster", "default" ]
}

resource "scaleway_vpc_private_network" "minimal" {
	name = "test-pool-minimal"
}

resource "scaleway_k8s_cluster" "minimal" {
    name = "test-pool-minimal"
	cni = "calico"
	version = "%s"
	tags = [ "terraform-test", "scaleway_k8s_cluster", "minimal" ]
	delete_additional_resources = false
	private_network_id = scaleway_vpc_private_network.minimal.id
}
%s`, version, pool)
}

func testAccCheckK8SPoolConfigWait(version string, otherPool bool, otherPoolSize int) string {
	pool := ""
	if otherPool {
		pool += fmt.Sprintf(`
resource "scaleway_k8s_pool" "minimal" {
    name = "test-pool-wait-2"
	cluster_id = scaleway_k8s_cluster.minimal.id
	node_type = "pro2_xxs"
	size = %d
	min_size = 1
	max_size = %d

	wait_for_pool_ready = true
}`, otherPoolSize, otherPoolSize)
	}

	return fmt.Sprintf(`
resource "scaleway_k8s_pool" "default" {
    name = "test-pool-wait"
	cluster_id = scaleway_k8s_cluster.minimal.id
	node_type = "pro2_xxs"
	size = 1
	min_size = 1
	max_size = 1
	wait_for_pool_ready = true
}

resource "scaleway_vpc_private_network" "minimal" {
	name = "test-pool-wait"
}

resource "scaleway_k8s_cluster" "minimal" {
    name = "test-pool-wait"
	cni = "calico"
	version = "%s"
	tags = [ "terraform-test", "scaleway_k8s_cluster", "minimal" ]
	delete_additional_resources = false
	private_network_id = scaleway_vpc_private_network.minimal.id
}
%s`, version, pool)
}

func testAccCheckK8SPoolConfigPlacementGroup(version string) string {
	return fmt.Sprintf(`
resource "scaleway_instance_placement_group" "placement_group" {
  name        = "test-pool-placement-group"
  policy_type = "max_availability"
  policy_mode = "optional"
}

resource "scaleway_k8s_pool" "placement_group" {
    name               = "test-pool-placement-group"
    cluster_id         = scaleway_k8s_cluster.placement_group.id
    node_type          = "pro2_xxs"
    placement_group_id = scaleway_instance_placement_group.placement_group.id
    size               = 1
}

resource "scaleway_vpc_private_network" "placement_group" {
	name = "test-pool-placement-group"
}

resource "scaleway_k8s_cluster" "placement_group" {
  name	  = "test-pool-placement-group"
  cni	  = "calico"
  version = "%s"
  tags	  = [ "terraform-test", "scaleway_k8s_cluster", "placement_group" ]
  delete_additional_resources = false
  private_network_id = scaleway_vpc_private_network.placement_group.id
}`, version)
}

func testAccCheckK8SPoolConfigPlacementGroupWithCustomZone(version string) string {
	return fmt.Sprintf(`
resource "scaleway_instance_placement_group" "placement_group" {
  name        = "test-pool-placement-group"
  policy_type = "max_availability"
  policy_mode = "optional"
  zone        = "nl-ams-2"
}

resource "scaleway_k8s_pool" "placement_group_2" {
  name               = "test-pool-placement-group-2"
  cluster_id         = scaleway_k8s_cluster.placement_group_2.id
  node_type          = "pro2_xxs"
  placement_group_id = scaleway_instance_placement_group.placement_group.id
  size               = 1
  region             = scaleway_k8s_cluster.placement_group_2.region
  zone               = scaleway_instance_placement_group.placement_group.zone
}

resource "scaleway_vpc_private_network" "placement_group" {
	name = "test-pool-placement-group"
	region = "nl-ams"
}

resource "scaleway_k8s_cluster" "placement_group_2" {
  name	  = "test-pool-placement-group-2"
  cni	  = "calico"
  version = "%s"
  tags	  = [ "terraform-test", "scaleway_k8s_cluster", "placement_group" ]
  region  = "nl-ams"
  delete_additional_resources = false
  private_network_id = scaleway_vpc_private_network.placement_group.id
}`, version)
}

func testAccCheckK8SPoolConfigPlacementGroupWithMultiZone(version string) string {
	return fmt.Sprintf(`
resource "scaleway_instance_placement_group" "placement_group" {
  name        = "test-pool-placement-group"
  policy_type = "max_availability"
  policy_mode = "optional"
  zone        = "nl-ams-1"
}

resource "scaleway_k8s_pool" "placement_group_2" {
  name               = "test-pool-placement-group-2"
  cluster_id         = scaleway_k8s_cluster.placement_group_2.id
  node_type          = "pro2_xxs"
  placement_group_id = scaleway_instance_placement_group.placement_group.id
  size               = 1
  region             = scaleway_k8s_cluster.placement_group_2.region
  zone               = scaleway_instance_placement_group.placement_group.zone
}

resource "scaleway_k8s_cluster" "placement_group_2" {
  name		= "test-pool-placement-group-2"
  cni		= "kilo"
  version	= "%s"
  tags		= [ "terraform-test", "scaleway_k8s_cluster", "placement_group" ]
  region	= "fr-par"
  type		= "multicloud"
  delete_additional_resources = false
}`, version)
}

func testAccCheckK8SPoolConfigUpgradePolicy(version string, maxSurge, maxUnavailable int) string {
	return fmt.Sprintf(`
resource "scaleway_k8s_pool" "upgrade_policy" {
    name = "test-pool-upgrade-policy"
	cluster_id = "${scaleway_k8s_cluster.upgrade_policy.id}"
	node_type = "pro2_xxs"
	autohealing = true
	autoscaling = true
	size = 1
	tags = [ "terraform-test", "scaleway_k8s_cluster", "upgrade_policy" ]
	upgrade_policy {
		max_surge = %d
		max_unavailable = %d
	}
}

resource "scaleway_vpc_private_network" "upgrade_policy" {
	name = "test-pool-upgrade-policy"
}

resource "scaleway_k8s_cluster" "upgrade_policy" {
    name = "test-pool-upgrade-policy"
	cni = "cilium"
	version = "%s"
	tags = [ "terraform-test", "scaleway_k8s_cluster", "upgrade_policy" ]
	delete_additional_resources = false
	private_network_id = scaleway_vpc_private_network.upgrade_policy.id
}`, maxSurge, maxUnavailable, version)
}

func testAccCheckK8SPoolConfigKubeletArgs(version string, maxPods int) string {
	return fmt.Sprintf(`
resource "scaleway_k8s_pool" "kubelet_args" {
    name = "test-pool-kubelet-args"
	cluster_id = "${scaleway_k8s_cluster.kubelet_args.id}"
	node_type = "pro2_xxs"
	autohealing = true
	autoscaling = true
	size = 1
	tags = [ "terraform-test", "scaleway_k8s_cluster", "kubelet_args" ]
	kubelet_args = {
		maxPods = %d
	}
}

resource "scaleway_vpc_private_network" "kubelet_args" {
	name = "test-pool-kubelet-args"
}

resource "scaleway_k8s_cluster" "kubelet_args" {
    name = "test-pool-kubelet-args"
	cni = "cilium"
	version = "%s"
	tags = [ "terraform-test", "scaleway_k8s_cluster", "kubelet_args" ]
	delete_additional_resources = false
	private_network_id = scaleway_vpc_private_network.kubelet_args.id
}`, maxPods, version)
}

func testAccCheckK8SPoolConfigZone(version string, zone string) string {
	return fmt.Sprintf(`
resource "scaleway_k8s_pool" "zone" {
    name = "test-pool-zone"
	cluster_id = "${scaleway_k8s_cluster.zone.id}"
	node_type = "pro2_xxs"
	autohealing = true
	autoscaling = true
	size = 1
	tags = [ "terraform-test", "scaleway_k8s_cluster", "zone" ]
	zone = "%s"
}

resource "scaleway_vpc_private_network" "zone" {
	name = "test-pool-zone"
}

resource "scaleway_k8s_cluster" "zone" {
    name = "test-pool-zone"
	cni = "cilium"
	version = "%s"
	tags = [ "terraform-test", "scaleway_k8s_cluster", "zone" ]
	delete_additional_resources = false
	private_network_id = scaleway_vpc_private_network.zone.id
}`, zone, version)
}

func testAccCheckK8SPoolNodesOneOfIsDeleting(name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("resource not found: %s", name)
		}
		nodesZeroStatus, ok := rs.Primary.Attributes["nodes.0.status"]
		if !ok {
			return errors.New("attribute \"nodes.0.status\" was not set")
		}
		nodesOneStatus, ok := rs.Primary.Attributes["nodes.1.status"]
		if !ok {
			return errors.New("attribute \"nodes.1.status\" was not set")
		}
		if nodesZeroStatus == "ready" && nodesOneStatus == "deleting" ||
			nodesZeroStatus == "deleting" && nodesOneStatus == "ready" {
			return nil
		}
		return fmt.Errorf("nodes status were not as expected: got %q for nodes.0 and %q for nodes.1", nodesZeroStatus, nodesOneStatus)
	}
}
