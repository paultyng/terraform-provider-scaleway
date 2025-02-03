package k8s_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	k8sSDK "github.com/scaleway/scaleway-sdk-go/api/k8s/v1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/k8s"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/vpc"
	vpcchecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/vpc/testfuncs"
)

func testAccK8SClusterGetLatestK8SVersion(tt *acctest.TestTools) string {
	api := k8sSDK.NewAPI(tt.Meta.ScwClient())
	versions, err := api.ListVersions(&k8sSDK.ListVersionsRequest{})
	if err != nil {
		tt.T.Fatalf("Could not get latestK8SVersion: %s", err)
	}
	if len(versions.Versions) > 1 {
		latestK8SVersion := versions.Versions[0].Name
		return latestK8SVersion
	}
	return ""
}

func testAccK8SClusterGetLatestK8SVersionMinor(tt *acctest.TestTools) string {
	api := k8sSDK.NewAPI(tt.Meta.ScwClient())
	versions, err := api.ListVersions(&k8sSDK.ListVersionsRequest{})
	if err != nil {
		tt.T.Fatalf("Could not get latestK8SVersion: %s", err)
	}
	if len(versions.Versions) > 1 {
		latestK8SVersion := versions.Versions[0].Name
		latestK8SVersionMinor, _ := k8s.GetMinorVersionFromFull(latestK8SVersion)
		return latestK8SVersionMinor
	}
	return ""
}

func testAccK8SClusterGetPreviousK8SVersion(tt *acctest.TestTools) string {
	api := k8sSDK.NewAPI(tt.Meta.ScwClient())
	versions, err := api.ListVersions(&k8sSDK.ListVersionsRequest{})
	if err != nil {
		tt.T.Fatalf("Could not get latestK8SVersion: %s", err)
	}
	if len(versions.Versions) > 1 {
		previousK8SVersion := versions.Versions[1].Name
		return previousK8SVersion
	}
	return ""
}

func testAccK8SClusterGetPreviousK8SVersionMinor(tt *acctest.TestTools) string {
	api := k8sSDK.NewAPI(tt.Meta.ScwClient())
	versions, err := api.ListVersions(&k8sSDK.ListVersionsRequest{})
	if err != nil {
		tt.T.Fatalf("Could not get latestK8SVersion: %s", err)
	}
	if len(versions.Versions) > 1 {
		previousK8SVersion := versions.Versions[1].Name
		previousK8SVersionMinor, _ := k8s.GetMinorVersionFromFull(previousK8SVersion)
		return previousK8SVersionMinor
	}
	return ""
}

func TestAccCluster_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	latestK8SVersion := testAccK8SClusterGetLatestK8SVersion(tt)
	previousK8SVersion := testAccK8SClusterGetPreviousK8SVersion(tt)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckK8SClusterDestroy(tt),
			vpcchecks.CheckPrivateNetworkDestroy(tt),
		),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "scaleway_vpc_private_network" "minimal" {
  name       = "test-minimal"
}
resource "scaleway_k8s_cluster" "minimal" {
	cni = "calico"
	version = "%s"
	name = "test-minimal"
	tags = [ "terraform-test", "scaleway_k8s_cluster", "minimal" ]
	delete_additional_resources = false
	private_network_id = scaleway_vpc_private_network.minimal.id
	description = "terraform basic test cluster"
}`, previousK8SVersion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8SClusterExists(tt, "scaleway_k8s_cluster.minimal"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.minimal", "version", previousK8SVersion),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.minimal", "cni", "calico"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.minimal", "status", k8sSDK.ClusterStatusPoolRequired.String()),
					resource.TestCheckResourceAttrSet("scaleway_k8s_cluster.minimal", "kubeconfig.0.config_file"),
					resource.TestCheckResourceAttrSet("scaleway_k8s_cluster.minimal", "kubeconfig.0.host"),
					resource.TestCheckResourceAttrSet("scaleway_k8s_cluster.minimal", "kubeconfig.0.cluster_ca_certificate"),
					resource.TestCheckResourceAttrSet("scaleway_k8s_cluster.minimal", "kubeconfig.0.token"),
					resource.TestCheckResourceAttrSet("scaleway_k8s_cluster.minimal", "apiserver_url"),
					resource.TestCheckResourceAttrSet("scaleway_k8s_cluster.minimal", "wildcard_dns"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.minimal", "tags.0", "terraform-test"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.minimal", "tags.1", "scaleway_k8s_cluster"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.minimal", "tags.2", "minimal"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.minimal", "description", "terraform basic test cluster"),
				),
			},
			{
				Config: fmt.Sprintf(`
resource "scaleway_vpc_private_network" "minimal" {
  name       = "test-minimal"
}
resource "scaleway_k8s_cluster" "minimal" {
	cni = "calico"
	version = "%s"
	name = "test-minimal"
	tags = [ "terraform-test", "scaleway_k8s_cluster", "minimal" ]
	delete_additional_resources = false
	private_network_id = scaleway_vpc_private_network.minimal.id
}`, latestK8SVersion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8SClusterExists(tt, "scaleway_k8s_cluster.minimal"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.minimal", "version", latestK8SVersion),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.minimal", "cni", "calico"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.minimal", "status", k8sSDK.ClusterStatusPoolRequired.String()),
					resource.TestCheckResourceAttrSet("scaleway_k8s_cluster.minimal", "kubeconfig.0.config_file"),
					resource.TestCheckResourceAttrSet("scaleway_k8s_cluster.minimal", "kubeconfig.0.host"),
					resource.TestCheckResourceAttrSet("scaleway_k8s_cluster.minimal", "kubeconfig.0.cluster_ca_certificate"),
					resource.TestCheckResourceAttrSet("scaleway_k8s_cluster.minimal", "kubeconfig.0.token"),
					resource.TestCheckResourceAttrSet("scaleway_k8s_cluster.minimal", "apiserver_url"),
					resource.TestCheckResourceAttrSet("scaleway_k8s_cluster.minimal", "wildcard_dns"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.minimal", "tags.0", "terraform-test"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.minimal", "tags.1", "scaleway_k8s_cluster"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.minimal", "tags.2", "minimal"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.minimal", "description", ""),
				),
			},
		},
	})
}

func TestAccCluster_Autoscaling(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	latestK8SVersion := testAccK8SClusterGetLatestK8SVersion(tt)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckK8SClusterDestroy(tt),
			vpcchecks.CheckPrivateNetworkDestroy(tt),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckK8SClusterConfigAutoscaler(latestK8SVersion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8SClusterExists(tt, "scaleway_k8s_cluster.autoscaler"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.autoscaler", "version", latestK8SVersion),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.autoscaler", "cni", "calico"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.autoscaler", "status", k8sSDK.ClusterStatusPoolRequired.String()),
					resource.TestCheckResourceAttrSet("scaleway_k8s_cluster.autoscaler", "kubeconfig.0.config_file"),
					resource.TestCheckResourceAttrSet("scaleway_k8s_cluster.autoscaler", "kubeconfig.0.host"),
					resource.TestCheckResourceAttrSet("scaleway_k8s_cluster.autoscaler", "kubeconfig.0.cluster_ca_certificate"),
					resource.TestCheckResourceAttrSet("scaleway_k8s_cluster.autoscaler", "kubeconfig.0.token"),
					resource.TestCheckResourceAttrSet("scaleway_k8s_cluster.autoscaler", "apiserver_url"),
					resource.TestCheckResourceAttrSet("scaleway_k8s_cluster.autoscaler", "wildcard_dns"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.autoscaler", "autoscaler_config.0.disable_scale_down", "true"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.autoscaler", "autoscaler_config.0.scale_down_delay_after_add", "20m"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.autoscaler", "autoscaler_config.0.scale_down_unneeded_time", "20m"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.autoscaler", "autoscaler_config.0.estimator", "binpacking"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.autoscaler", "autoscaler_config.0.expander", "most_pods"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.autoscaler", "autoscaler_config.0.ignore_daemonsets_utilization", "true"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.autoscaler", "autoscaler_config.0.balance_similar_node_groups", "true"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.autoscaler", "autoscaler_config.0.expendable_pods_priority_cutoff", "10"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.autoscaler", "autoscaler_config.0.scale_down_utilization_threshold", "0.77"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.autoscaler", "autoscaler_config.0.max_graceful_termination_sec", "1337"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.autoscaler", "tags.0", "terraform-test"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.autoscaler", "tags.1", "scaleway_k8s_cluster"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.autoscaler", "tags.2", "autoscaler-config"),
				),
			},
			{
				Config: testAccCheckK8SClusterConfigAutoscalerChange(latestK8SVersion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8SClusterExists(tt, "scaleway_k8s_cluster.autoscaler"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.autoscaler", "version", latestK8SVersion),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.autoscaler", "cni", "calico"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.autoscaler", "status", k8sSDK.ClusterStatusPoolRequired.String()),
					resource.TestCheckResourceAttrSet("scaleway_k8s_cluster.autoscaler", "kubeconfig.0.config_file"),
					resource.TestCheckResourceAttrSet("scaleway_k8s_cluster.autoscaler", "kubeconfig.0.host"),
					resource.TestCheckResourceAttrSet("scaleway_k8s_cluster.autoscaler", "kubeconfig.0.cluster_ca_certificate"),
					resource.TestCheckResourceAttrSet("scaleway_k8s_cluster.autoscaler", "kubeconfig.0.token"),
					resource.TestCheckResourceAttrSet("scaleway_k8s_cluster.autoscaler", "apiserver_url"),
					resource.TestCheckResourceAttrSet("scaleway_k8s_cluster.autoscaler", "wildcard_dns"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.autoscaler", "autoscaler_config.0.disable_scale_down", "false"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.autoscaler", "autoscaler_config.0.scale_down_delay_after_add", "20m"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.autoscaler", "autoscaler_config.0.scale_down_unneeded_time", "5m"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.autoscaler", "autoscaler_config.0.estimator", "binpacking"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.autoscaler", "autoscaler_config.0.expander", "most_pods"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.autoscaler", "autoscaler_config.0.ignore_daemonsets_utilization", "false"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.autoscaler", "autoscaler_config.0.balance_similar_node_groups", "false"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.autoscaler", "autoscaler_config.0.expendable_pods_priority_cutoff", "0"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.autoscaler", "autoscaler_config.0.scale_down_utilization_threshold", "0.33"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.autoscaler", "autoscaler_config.0.max_graceful_termination_sec", "2664"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.autoscaler", "tags.0", "terraform-test"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.autoscaler", "tags.1", "scaleway_k8s_cluster"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.autoscaler", "tags.2", "autoscaler-config"),
				),
			},
		},
	})
}

func TestAccCluster_OIDC(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	latestK8SVersion := testAccK8SClusterGetLatestK8SVersion(tt)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckK8SClusterDestroy(tt),
			vpcchecks.CheckPrivateNetworkDestroy(tt),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckK8SClusterConfigOIDC(latestK8SVersion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8SClusterExists(tt, "scaleway_k8s_cluster.oidc"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.oidc", "version", latestK8SVersion),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.oidc", "cni", "cilium"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.oidc", "status", k8sSDK.ClusterStatusPoolRequired.String()),
					resource.TestCheckResourceAttrSet("scaleway_k8s_cluster.oidc", "kubeconfig.0.config_file"),
					resource.TestCheckResourceAttrSet("scaleway_k8s_cluster.oidc", "kubeconfig.0.host"),
					resource.TestCheckResourceAttrSet("scaleway_k8s_cluster.oidc", "kubeconfig.0.cluster_ca_certificate"),
					resource.TestCheckResourceAttrSet("scaleway_k8s_cluster.oidc", "kubeconfig.0.token"),
					resource.TestCheckResourceAttrSet("scaleway_k8s_cluster.oidc", "apiserver_url"),
					resource.TestCheckResourceAttrSet("scaleway_k8s_cluster.oidc", "wildcard_dns"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.oidc", "open_id_connect_config.0.issuer_url", "https://accounts.google.com"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.oidc", "open_id_connect_config.0.client_id", "my-super-id"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.oidc", "open_id_connect_config.0.username_claim", "mario"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.oidc", "open_id_connect_config.0.groups_prefix", "pouf"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.oidc", "open_id_connect_config.0.groups_claim.0", "k8s"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.oidc", "open_id_connect_config.0.groups_claim.1", "admin"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.oidc", "tags.0", "terraform-test"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.oidc", "tags.1", "scaleway_k8s_cluster"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.oidc", "tags.2", "oidc-config"),
				),
			},
			{
				Config: testAccCheckK8SClusterConfigOIDCChange(latestK8SVersion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8SClusterExists(tt, "scaleway_k8s_cluster.oidc"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.oidc", "version", latestK8SVersion),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.oidc", "cni", "cilium"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.oidc", "status", k8sSDK.ClusterStatusPoolRequired.String()),
					resource.TestCheckResourceAttrSet("scaleway_k8s_cluster.oidc", "kubeconfig.0.config_file"),
					resource.TestCheckResourceAttrSet("scaleway_k8s_cluster.oidc", "kubeconfig.0.host"),
					resource.TestCheckResourceAttrSet("scaleway_k8s_cluster.oidc", "kubeconfig.0.cluster_ca_certificate"),
					resource.TestCheckResourceAttrSet("scaleway_k8s_cluster.oidc", "kubeconfig.0.token"),
					resource.TestCheckResourceAttrSet("scaleway_k8s_cluster.oidc", "apiserver_url"),
					resource.TestCheckResourceAttrSet("scaleway_k8s_cluster.oidc", "wildcard_dns"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.oidc", "open_id_connect_config.0.issuer_url", "https://gitlab.com"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.oidc", "open_id_connect_config.0.client_id", "my-even-more-awesome-id"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.oidc", "open_id_connect_config.0.username_claim", "luigi"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.oidc", "open_id_connect_config.0.username_prefix", "boo"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.oidc", "open_id_connect_config.0.groups_prefix", ""),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.oidc", "tags.0", "terraform-test"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.oidc", "tags.1", "scaleway_k8s_cluster"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.oidc", "tags.2", "oidc-config"),
				),
			},
		},
	})
}

func TestAccCluster_AutoUpgrade(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	latestK8SVersion := testAccK8SClusterGetLatestK8SVersion(tt)
	latestK8SVersionMinor := testAccK8SClusterGetLatestK8SVersionMinor(tt)
	previousK8SVersion := testAccK8SClusterGetPreviousK8SVersion(tt)
	previousK8SVersionMinor := testAccK8SClusterGetPreviousK8SVersionMinor(tt)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckK8SClusterDestroy(tt),
			vpcchecks.CheckPrivateNetworkDestroy(tt),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckK8SClusterAutoUpgrade(false, "any", 0, previousK8SVersion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8SClusterExists(tt, "scaleway_k8s_cluster.auto_upgrade"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.auto_upgrade", "version", previousK8SVersion),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.auto_upgrade", "auto_upgrade.0.enable", "false"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.auto_upgrade", "auto_upgrade.0.maintenance_window_day", "any"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.auto_upgrade", "auto_upgrade.0.maintenance_window_start_hour", "0"),
				),
			},
			{
				Config: testAccCheckK8SClusterAutoUpgrade(true, "any", 0, previousK8SVersionMinor),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8SClusterExists(tt, "scaleway_k8s_cluster.auto_upgrade"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.auto_upgrade", "version", previousK8SVersionMinor),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.auto_upgrade", "auto_upgrade.0.enable", "true"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.auto_upgrade", "auto_upgrade.0.maintenance_window_day", "any"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.auto_upgrade", "auto_upgrade.0.maintenance_window_start_hour", "0"),
				),
			},
			{
				Config: testAccCheckK8SClusterAutoUpgrade(true, "any", 0, latestK8SVersionMinor),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8SClusterExists(tt, "scaleway_k8s_cluster.auto_upgrade"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.auto_upgrade", "version", latestK8SVersionMinor),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.auto_upgrade", "auto_upgrade.0.enable", "true"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.auto_upgrade", "auto_upgrade.0.maintenance_window_day", "any"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.auto_upgrade", "auto_upgrade.0.maintenance_window_start_hour", "0"),
				),
			},
			{
				Config: testAccCheckK8SClusterAutoUpgrade(false, "any", 0, latestK8SVersion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8SClusterExists(tt, "scaleway_k8s_cluster.auto_upgrade"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.auto_upgrade", "version", latestK8SVersion),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.auto_upgrade", "auto_upgrade.0.enable", "false"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.auto_upgrade", "auto_upgrade.0.maintenance_window_day", "any"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.auto_upgrade", "auto_upgrade.0.maintenance_window_start_hour", "0"),
				),
			},
			{
				Config: testAccCheckK8SClusterAutoUpgrade(true, "tuesday", 3, latestK8SVersionMinor),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8SClusterExists(tt, "scaleway_k8s_cluster.auto_upgrade"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.auto_upgrade", "version", latestK8SVersionMinor),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.auto_upgrade", "auto_upgrade.0.enable", "true"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.auto_upgrade", "auto_upgrade.0.maintenance_window_day", "tuesday"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.auto_upgrade", "auto_upgrade.0.maintenance_window_start_hour", "3"),
				),
			},
			{
				Config: testAccCheckK8SClusterAutoUpgrade(true, "any", 0, latestK8SVersionMinor),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8SClusterExists(tt, "scaleway_k8s_cluster.auto_upgrade"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.auto_upgrade", "version", latestK8SVersionMinor),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.auto_upgrade", "auto_upgrade.0.enable", "true"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.auto_upgrade", "auto_upgrade.0.maintenance_window_day", "any"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.auto_upgrade", "auto_upgrade.0.maintenance_window_start_hour", "0"),
				),
			},
		},
	})
}

func TestAccCluster_PrivateNetwork(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	latestK8SVersion := testAccK8SClusterGetLatestK8SVersion(tt)

	clusterID := ""

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckK8SClusterDestroy(tt),
			vpcchecks.CheckPrivateNetworkDestroy(tt),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckK8SClusterConfigPrivateNetworkLinked(latestK8SVersion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8SClusterExists(tt, "scaleway_k8s_cluster.private_network"),
					vpcchecks.IsPrivateNetworkPresent(tt, "scaleway_vpc_private_network.private_network"),
					testAccCheckK8sClusterPrivateNetworkID(tt, "scaleway_k8s_cluster.private_network", "scaleway_vpc_private_network.private_network"),
					acctest.CheckResourceIDPersisted("scaleway_k8s_cluster.private_network", &clusterID),
				),
			},
			{
				Config: testAccCheckK8SClusterConfigPrivateNetworkChange(latestK8SVersion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8SClusterExists(tt, "scaleway_k8s_cluster.private_network"),
					vpcchecks.IsPrivateNetworkPresent(tt, "scaleway_vpc_private_network.private_network"),
					vpcchecks.IsPrivateNetworkPresent(tt, "scaleway_vpc_private_network.private_network_2"),
					testAccCheckK8sClusterPrivateNetworkID(tt, "scaleway_k8s_cluster.private_network", "scaleway_vpc_private_network.private_network_2"),
					acctest.CheckResourceIDChanged("scaleway_k8s_cluster.private_network", &clusterID),
				),
			},
		},
	})
}

func TestAccCluster_Multicloud(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	latestK8SVersion := testAccK8SClusterGetLatestK8SVersion(tt)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckK8SClusterDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckK8SClusterMulticloud(latestK8SVersion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8SClusterExists(tt, "scaleway_k8s_cluster.multicloud"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.multicloud", "version", latestK8SVersion),
				),
			},
		},
	})
}

func TestAccCluster_TypeChange(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	latestK8SVersion := testAccK8SClusterGetLatestK8SVersion(tt)

	clusterID := ""

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckK8SClusterDestroy(tt),
			vpcchecks.CheckPrivateNetworkDestroy(tt),
		),
		Steps: []resource.TestStep{
			{
				// 1 : Start with a mutualized Kapsule cluster
				Config: testAccCheckK8SClusterTypeChange("kapsule", "cilium", latestK8SVersion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8SClusterExists(tt, "scaleway_k8s_cluster.type-change"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.type-change", "type", "kapsule"),
					acctest.CheckResourceIDPersisted("scaleway_k8s_cluster.type-change", &clusterID),
				),
			},
			{
				// 2 : Upgrade to a dedicated Kapsule --> should migrate
				Config: testAccCheckK8SClusterTypeChange("kapsule-dedicated-4", "cilium", latestK8SVersion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8SClusterExists(tt, "scaleway_k8s_cluster.type-change"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.type-change", "type", "kapsule-dedicated-4"),
					acctest.CheckResourceIDPersisted("scaleway_k8s_cluster.type-change", &clusterID),
				),
			},
			{
				// 3 : Upgrade to an even bigger dedicated Kapsule --> should migrate
				Config: testAccCheckK8SClusterTypeChange("kapsule-dedicated-8", "cilium", latestK8SVersion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8SClusterExists(tt, "scaleway_k8s_cluster.type-change"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.type-change", "type", "kapsule-dedicated-8"),
					acctest.CheckResourceIDPersisted("scaleway_k8s_cluster.type-change", &clusterID),
				),
			},
			{
				// 4 : Change to a dedicated Kosmos --> should recreate
				Config: testAccCheckK8SClusterTypeChange("multicloud-dedicated-4", "kilo", latestK8SVersion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8SClusterExists(tt, "scaleway_k8s_cluster.type-change"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.type-change", "type", "multicloud-dedicated-4"),
					acctest.CheckResourceIDChanged("scaleway_k8s_cluster.type-change", &clusterID),
				),
			},
			{
				// 5 : Upgrade to a bigger dedicated Kosmos --> should migrate
				Config: testAccCheckK8SClusterTypeChange("multicloud-dedicated-8", "kilo", latestK8SVersion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8SClusterExists(tt, "scaleway_k8s_cluster.type-change"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.type-change", "type", "multicloud-dedicated-8"),
					acctest.CheckResourceIDPersisted("scaleway_k8s_cluster.type-change", &clusterID),
				),
			},
			{
				// 6 : Downgrade to a smaller dedicated Kosmos --> should recreate
				Config: testAccCheckK8SClusterTypeChange("multicloud-dedicated-4", "kilo", latestK8SVersion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8SClusterExists(tt, "scaleway_k8s_cluster.type-change"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.type-change", "type", "multicloud-dedicated-4"),
					acctest.CheckResourceIDChanged("scaleway_k8s_cluster.type-change", &clusterID),
				),
			},
			{
				// 7 : Downgrade to a mutualized Kosmos --> should recreate
				Config: testAccCheckK8SClusterTypeChange("multicloud", "kilo", latestK8SVersion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8SClusterExists(tt, "scaleway_k8s_cluster.type-change"),
					resource.TestCheckResourceAttr("scaleway_k8s_cluster.type-change", "type", "multicloud"),
					acctest.CheckResourceIDChanged("scaleway_k8s_cluster.type-change", &clusterID),
				),
			},
		},
	})
}

func testAccCheckK8SClusterDestroy(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_k8s_cluster" {
				continue
			}

			k8sAPI, region, clusterID, err := k8s.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}

			_, err = k8sAPI.GetCluster(&k8sSDK.GetClusterRequest{
				Region:    region,
				ClusterID: clusterID,
			})

			// If no error resource still exist
			if err == nil {
				return fmt.Errorf("cluster (%s) still exists", rs.Primary.ID)
			}

			// Unexpected api error we return it
			if !httperrors.Is404(err) {
				return err
			}
		}
		return nil
	}
}

func testAccCheckK8SClusterExists(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		k8sAPI, region, clusterID, err := k8s.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = k8sAPI.GetCluster(&k8sSDK.GetClusterRequest{
			Region:    region,
			ClusterID: clusterID,
		})
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckK8sClusterPrivateNetworkID(tt *acctest.TestTools, clusterName, pnName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[clusterName]
		if !ok {
			return fmt.Errorf("resource not found: %s", clusterName)
		}

		k8sAPI, region, clusterID, err := k8s.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		cluster, err := k8sAPI.GetCluster(&k8sSDK.GetClusterRequest{
			Region:    region,
			ClusterID: clusterID,
		})
		if err != nil {
			return err
		}

		clusterPNID := cluster.PrivateNetworkID

		rs, ok = s.RootModule().Resources[pnName]
		if !ok {
			return fmt.Errorf("resource not found: %s", pnName)
		}

		_, _, pnID, err := vpc.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		if clusterPNID == nil {
			return fmt.Errorf("expected %s private_network_id to be %s, got nil", clusterName, pnID)
		}
		if *clusterPNID != pnID {
			return fmt.Errorf("expected %s private_network_id to be %s, got %s", clusterName, pnID, *clusterPNID)
		}

		return nil
	}
}

func testAccCheckK8SClusterConfigAutoscaler(version string) string {
	return fmt.Sprintf(`
resource "scaleway_vpc_private_network" "autoscaler" {
  name       = "test-autoscaler"
  region 	 = "nl-ams"
}
resource "scaleway_k8s_cluster" "autoscaler" {
	cni = "calico"
	version = "%s"
	name = "test-autoscaler-01"
	region = "nl-ams"
	autoscaler_config {
		disable_scale_down = true
		scale_down_delay_after_add = "20m"
		scale_down_unneeded_time = "20m"
		estimator = "binpacking"
		expander = "most_pods"
		ignore_daemonsets_utilization = true
		balance_similar_node_groups = true
		expendable_pods_priority_cutoff = 10
		scale_down_utilization_threshold = 0.77
		max_graceful_termination_sec = 1337
	}
	tags = [ "terraform-test", "scaleway_k8s_cluster", "autoscaler-config" ]
	delete_additional_resources = false
	private_network_id = scaleway_vpc_private_network.autoscaler.id
}`, version)
}

func testAccCheckK8SClusterConfigAutoscalerChange(version string) string {
	return fmt.Sprintf(`
resource "scaleway_vpc_private_network" "autoscaler" {
  name       = "test-autoscaler"
  region 	 = "nl-ams"
}
resource "scaleway_k8s_cluster" "autoscaler" {
	cni = "calico"
	version = "%s"
	name = "test-autoscaler-02"
	region = "nl-ams"
	autoscaler_config {
		disable_scale_down = false
		scale_down_delay_after_add = "20m"
		scale_down_unneeded_time = "5m"
		estimator = "binpacking"
		expander = "most_pods"
		expendable_pods_priority_cutoff = 0
		scale_down_utilization_threshold = 0.33
		max_graceful_termination_sec = 2664
	}
	tags = [ "terraform-test", "scaleway_k8s_cluster", "autoscaler-config" ]
	delete_additional_resources = false
	private_network_id = scaleway_vpc_private_network.autoscaler.id
}`, version)
}

func testAccCheckK8SClusterConfigOIDC(version string) string {
	return fmt.Sprintf(`
resource "scaleway_vpc_private_network" "oidc" {
  name       = "test-oidc"
}
resource "scaleway_k8s_cluster" "oidc" {
	cni = "cilium"
	version = "%s"
	name = "test-oidc"
	open_id_connect_config {
		issuer_url = "https://accounts.google.com"
		client_id = "my-super-id"
		username_claim = "mario"
		groups_claim = [ "k8s", "admin" ]
		groups_prefix = "pouf"
	}
	tags = [ "terraform-test", "scaleway_k8s_cluster", "oidc-config" ]
	delete_additional_resources = false
	private_network_id = scaleway_vpc_private_network.oidc.id
}
`, version)
}

func testAccCheckK8SClusterConfigOIDCChange(version string) string {
	return fmt.Sprintf(`
resource "scaleway_vpc_private_network" "oidc" {
  name       = "test-oidc"
}
resource "scaleway_k8s_cluster" "oidc" {
	cni = "cilium"
	version = "%s"
	name = "test-oidc"
	open_id_connect_config {
		issuer_url = "https://gitlab.com"
		client_id = "my-even-more-awesome-id"
		username_claim = "luigi"
		groups_claim = [ ]
		username_prefix = "boo"
	}
	tags = [ "terraform-test", "scaleway_k8s_cluster", "oidc-config" ]
	delete_additional_resources = false
	private_network_id = scaleway_vpc_private_network.oidc.id
}
`, version)
}

func testAccCheckK8SClusterAutoUpgrade(enable bool, day string, hour uint64, version string) string {
	return fmt.Sprintf(`
resource "scaleway_vpc_private_network" "auto_upgrade" {
  name       = "test-auto-upgrade"
}
resource "scaleway_k8s_cluster" "auto_upgrade" {
	cni = "calico"
	version = "%s"
	name = "test-auto-upgrade"
	auto_upgrade {
	    enable = %t
		maintenance_window_start_hour = %d
		maintenance_window_day = "%s"
	}
	tags = [ "terraform-test", "scaleway_k8s_cluster", "auto_upgrade" ]
	delete_additional_resources = false
	private_network_id = scaleway_vpc_private_network.auto_upgrade.id
}`, version, enable, hour, day)
}

func testAccCheckK8SClusterConfigPrivateNetworkLinked(version string) string {
	return fmt.Sprintf(`
resource "scaleway_vpc_private_network" "private_network" {
  name       = "k8s-private-network"
}
resource "scaleway_k8s_cluster" "private_network" {
	cni = "calico"
	version = "%s"
	name = "k8s-private-network-cluster"
	tags = [ "terraform-test", "scaleway_k8s_cluster", "private_network" ]
	delete_additional_resources = false
	depends_on = [scaleway_vpc_private_network.private_network]
    private_network_id = scaleway_vpc_private_network.private_network.id
}`, version)
}

func testAccCheckK8SClusterConfigPrivateNetworkChange(version string) string {
	return fmt.Sprintf(`
resource "scaleway_vpc_private_network" "private_network" {
  name       = "k8s-private-network"
}
resource "scaleway_vpc_private_network" "private_network_2" {
  name       = "other-private-network"
}
resource "scaleway_k8s_cluster" "private_network" {
	cni = "calico"
	version = "%s"
	name = "k8s-private-network-cluster"
	tags = [ "terraform-test", "scaleway_k8s_cluster", "private_network" ]
	delete_additional_resources = false
	depends_on = [scaleway_vpc_private_network.private_network_2]
    private_network_id = scaleway_vpc_private_network.private_network_2.id
}`, version)
}

func testAccCheckK8SClusterMulticloud(version string) string {
	return fmt.Sprintf(`
resource "scaleway_k8s_cluster" "multicloud" {
	name = "test-multicloud"
	version = "%s"
	cni = "kilo"
	type = "multicloud"
	delete_additional_resources = false
}

resource "scaleway_k8s_pool" "multicloud" {
	cluster_id = "${scaleway_k8s_cluster.multicloud.id}"
	name = "test-multicloud"
	node_type = "DEV1-M"
	size = 1
}
`, version)
}

func testAccCheckK8SClusterTypeChange(clusterType, cni, version string) string {
	config := ""
	isKapsule := strings.HasPrefix(clusterType, "kapsule")
	if isKapsule {
		config = `
resource "scaleway_vpc_private_network" "type-change" {
  name       = "test-type-change"
}`
	}

	config += fmt.Sprintf(`
resource "scaleway_k8s_cluster" "type-change" {
	type = "%s"
	cni = "%s"
	version = "%s"
	name = "test-type-change"
	tags = [ "terraform-test", "scaleway_k8s_cluster", "type-change" ]
	delete_additional_resources = false`, clusterType, cni, version)

	if isKapsule {
		config += "\nprivate_network_id = scaleway_vpc_private_network.type-change.id\n}"
	} else {
		config += "\n}"
	}

	return config
}
