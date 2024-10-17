package rdb_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	rdbSDK "github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/rdb"
	rdbchecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/rdb/testfuncs"
	vpcchecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/vpc/testfuncs"
)

const (
	mySQLEngineName      = "MySQL"
	postgreSQLEngineName = "PostgreSQL"
)

func TestAccInstance_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	latestEngineVersion := rdbchecks.GetLatestEngineVersion(tt, postgreSQLEngineName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      rdbchecks.IsInstanceDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource scaleway_rdb_instance main {
						name = "test-rdb-basic"
						node_type = "db-dev-s"
						engine = %q
						is_ha_cluster = false
						disable_backup = true
						user_name = "my_initial_user"
						password = "thiZ_is_v&ry_s3cret"
						tags = [ "terraform-test", "scaleway_rdb_instance", "minimal" ]
					}
				`, latestEngineVersion),
				Check: resource.ComposeTestCheckFunc(
					isInstancePresent(tt, "scaleway_rdb_instance.main"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "name", "test-rdb-basic"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "node_type", "db-dev-s"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "engine", latestEngineVersion),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "is_ha_cluster", "false"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "disable_backup", "true"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "user_name", "my_initial_user"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "password", "thiZ_is_v&ry_s3cret"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "tags.0", "terraform-test"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "tags.1", "scaleway_rdb_instance"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "tags.2", "minimal"),
					resource.TestCheckResourceAttrSet("scaleway_rdb_instance.main", "endpoint_ip"),   // Deprecated attribute, might be deleted later
					resource.TestCheckResourceAttrSet("scaleway_rdb_instance.main", "endpoint_port"), // Deprecated attribute, might be deleted later
					resource.TestCheckResourceAttrSet("scaleway_rdb_instance.main", "certificate"),
					resource.TestCheckResourceAttrSet("scaleway_rdb_instance.main", "load_balancer.0.ip"),
					resource.TestCheckResourceAttrSet("scaleway_rdb_instance.main", "load_balancer.0.endpoint_id"),
					resource.TestCheckResourceAttrSet("scaleway_rdb_instance.main", "load_balancer.0.port"),
					resource.TestCheckResourceAttrSet("scaleway_rdb_instance.main", "logs_policy.0.max_age_retention"),
					resource.TestCheckResourceAttrSet("scaleway_rdb_instance.main", "logs_policy.0.total_disk_retention"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource scaleway_rdb_instance main {
						name = "test-rdb-basic"
						node_type = "db-dev-s"
						engine = %q
						is_ha_cluster = false
						disable_backup = true
						user_name = "my_initial_user"
						password = "thiZ_is_v&ry_s3cret"
						tags = [ "terraform-change-tag" ]
					}
				`, latestEngineVersion),
				Check: resource.ComposeTestCheckFunc(
					isInstancePresent(tt, "scaleway_rdb_instance.main"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "name", "test-rdb-basic"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "node_type", "db-dev-s"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "engine", latestEngineVersion),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "is_ha_cluster", "false"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "disable_backup", "true"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "user_name", "my_initial_user"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "password", "thiZ_is_v&ry_s3cret"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "tags.0", "terraform-change-tag"),
					resource.TestCheckResourceAttrSet("scaleway_rdb_instance.main", "endpoint_ip"),   // Deprecated attribute, might be deleted later
					resource.TestCheckResourceAttrSet("scaleway_rdb_instance.main", "endpoint_port"), // Deprecated attribute, might be deleted later
					resource.TestCheckResourceAttrSet("scaleway_rdb_instance.main", "certificate"),
					resource.TestCheckResourceAttrSet("scaleway_rdb_instance.main", "load_balancer.0.ip"),
					resource.TestCheckResourceAttrSet("scaleway_rdb_instance.main", "load_balancer.0.endpoint_id"),
					resource.TestCheckResourceAttrSet("scaleway_rdb_instance.main", "load_balancer.0.port"),
				),
			},
		},
	})
}

func TestAccInstance_WithCluster(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	latestEngineVersion := rdbchecks.GetLatestEngineVersion(tt, postgreSQLEngineName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      rdbchecks.IsInstanceDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource scaleway_rdb_instance main {
						name = "test-rdb-with-cluster"
						node_type = "db-dev-m"
						engine = %q
						is_ha_cluster = true
						disable_backup = false
						user_name = "my_initial_user"
						password = "thiZ_is_v&ry_s8cret"
						tags = [ "terraform-test", "scaleway_rdb_instance", "minimal" ]
					}
				`, latestEngineVersion),
				Check: resource.ComposeTestCheckFunc(
					isInstancePresent(tt, "scaleway_rdb_instance.main"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "name", "test-rdb-with-cluster"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "node_type", "db-dev-m"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "engine", latestEngineVersion),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "is_ha_cluster", "true"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "disable_backup", "false"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "user_name", "my_initial_user"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "password", "thiZ_is_v&ry_s8cret"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "tags.0", "terraform-test"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "tags.1", "scaleway_rdb_instance"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "tags.2", "minimal"),
				),
			},
		},
	})
}

func TestAccInstance_LogsPolicy(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	latestEngineVersion := rdbchecks.GetLatestEngineVersion(tt, postgreSQLEngineName)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      rdbchecks.IsInstanceDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource scaleway_rdb_instance main {
						name = "test-rdb-log-policy"
						node_type = "db-dev-s"
						engine = %q
						is_ha_cluster = true
						disable_backup = false
						user_name = "my_initial_user"
						password = "thiZ_is_v&ry_s8cret"
						tags = [ "terraform-test", "scaleway_rdb_instance", "minimal" ]
  						logs_policy {
							max_age_retention    = 30
							total_disk_retention = 100000000
					  }
					}
				`, latestEngineVersion),
				Check: resource.ComposeTestCheckFunc(
					isInstancePresent(tt, "scaleway_rdb_instance.main"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "logs_policy.0.max_age_retention", "30"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "logs_policy.0.total_disk_retention", "100000000"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource scaleway_rdb_instance main {
						name = "test-rdb-log-policy"
						node_type = "db-dev-s"
						engine = %q
						is_ha_cluster = true
						disable_backup = false
						user_name = "my_initial_user"
						password = "thiZ_is_v&ry_s8cret"
						tags = [ "terraform-test", "scaleway_rdb_instance", "minimal" ]
  						logs_policy {
							max_age_retention    = 10
							total_disk_retention = 200000000
					  }
					}
				`, latestEngineVersion),
				Check: resource.ComposeTestCheckFunc(
					isInstancePresent(tt, "scaleway_rdb_instance.main"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "logs_policy.0.max_age_retention", "10"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "logs_policy.0.total_disk_retention", "200000000"),
				),
			},
		},
	})
}

func TestAccInstance_Settings(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	latestEngineVersion := rdbchecks.GetLatestEngineVersion(tt, postgreSQLEngineName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      rdbchecks.IsInstanceDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource scaleway_rdb_instance main {
						name = "test-rdb-instance-settings"
						node_type = "db-dev-s"
						disable_backup = true
						engine = %q
						user_name = "my_initial_user"
						password = "thiZ_is_v&ry_s3cret"
						settings = {
							work_mem = "4"
							max_connections = "200"
							effective_cache_size = "1300"
							maintenance_work_mem = "150"
							max_parallel_workers = "2"
							max_parallel_workers_per_gather = "2"
						}
					}
				`, latestEngineVersion),
				Check: resource.ComposeTestCheckFunc(
					isInstancePresent(tt, "scaleway_rdb_instance.main"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "settings.work_mem", "4"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "settings.max_connections", "200"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "settings.effective_cache_size", "1300"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "settings.maintenance_work_mem", "150"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "settings.max_parallel_workers", "2"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "settings.max_parallel_workers_per_gather", "2"),
				),
			},
		},
	})
}

func TestAccInstance_InitSettings(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	latestEngineVersion := rdbchecks.GetLatestEngineVersion(tt, mySQLEngineName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      rdbchecks.IsInstanceDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource scaleway_rdb_instance main {
						name = "test-rdb-init-settings"
						node_type = "db-dev-s"
						disable_backup = true
						engine = %q
						user_name = "my_initial_user"
						password = "thiZ_is_v&ry_s3cret"
						init_settings = {
							"lower_case_table_names" = 1
						}
						settings = {
							"max_connections" = "350"
						}
					}
				`, latestEngineVersion),
				Check: resource.ComposeTestCheckFunc(
					isInstancePresent(tt, "scaleway_rdb_instance.main"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "init_settings.lower_case_table_names", "1"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "settings.max_connections", "350"),
				),
			},
		},
	})
}

func TestAccInstance_Capitalize(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	latestEngineVersion := rdbchecks.GetLatestEngineVersion(tt, postgreSQLEngineName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      rdbchecks.IsInstanceDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource scaleway_rdb_instance main {
						name = "test-rdb-instance-capitalize"
						node_type = "DB-DEV-S"
						engine = %q
						is_ha_cluster = false
						disable_backup = true
						user_name = "my_initial_user"
						password = "thiZ_is_v&ry_s3cret"
					}
				`, latestEngineVersion),
				Check: resource.ComposeTestCheckFunc(
					isInstancePresent(tt, "scaleway_rdb_instance.main"),
				),
			},
		},
	})
}

func TestAccInstance_PrivateNetwork(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	latestEngineVersion := rdbchecks.GetLatestEngineVersion(tt, postgreSQLEngineName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      rdbchecks.IsInstanceDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_vpc_private_network pn01 {
						name = "my_private_network"
						region= "nl-ams"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaleway_vpc_private_network.pn01", "name", "my_private_network"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource scaleway_vpc_private_network pn01 {
						name = "my_private_network"
						region= "nl-ams"
					}

					resource scaleway_rdb_instance main {
						name = "test-rdb-private-network"
						node_type = "db-dev-s"
						engine = %q
						is_ha_cluster = false
						disable_backup = true
						user_name = "my_initial_user"
						password = "thiZ_is_v&ry_s3cret"
						region= "nl-ams"
						tags = [ "terraform-test", "scaleway_rdb_instance", "volume", "rdb_pn" ]
						volume_type = "bssd"
						volume_size_in_gb = 10
						private_network {
							ip_net = "192.168.1.42/24"
							pn_id = "${scaleway_vpc_private_network.pn01.id}"
						}
					}
				`, latestEngineVersion),
				Check: resource.ComposeTestCheckFunc(
					isInstancePresent(tt, "scaleway_rdb_instance.main"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "private_network.#", "1"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "private_network.0.ip_net", "192.168.1.42/24"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource scaleway_vpc_private_network pn01 {
						name = "my_private_network_to_be_replaced"
						region= "nl-ams"
					}

					resource scaleway_vpc_private_network pn02 {
						name = "my_private_network"
						region= "nl-ams"
					}

					resource scaleway_rdb_instance main {
						name = "test-rdb-private-network"
						node_type = "db-dev-s"
						engine = %q
						is_ha_cluster = false
						disable_backup = true
						user_name = "my_initial_user"
						password = "thiZ_is_v&ry_s3cret"
						region= "nl-ams"
						tags = [ "terraform-test", "scaleway_rdb_instance", "volume", "rdb_pn" ]
						volume_type = "bssd"
						volume_size_in_gb = 10
						private_network {
							ip_net = "192.168.1.254/24"
							pn_id = "${scaleway_vpc_private_network.pn02.id}"
						}
					}
				`, latestEngineVersion),
				Check: resource.ComposeTestCheckFunc(
					isInstancePresent(tt, "scaleway_rdb_instance.main"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "private_network.#", "1"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "private_network.0.ip_net", "192.168.1.254/24"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource scaleway_vpc_private_network pn01 {
						name = "my_private_network_to_be_replaced"
						region= "nl-ams"
					}

					resource scaleway_vpc_private_network pn02 {
						name = "my_private_network"
						region= "nl-ams"
					}

					resource scaleway_vpc_public_gateway_dhcp main {
						subnet = "192.168.1.0/24"
						zone = "nl-ams-1"
					}

					resource scaleway_vpc_public_gateway_ip main {
						zone = "nl-ams-1"
					}

					resource scaleway_vpc_public_gateway main {
						name = "foobar"
						type = "VPC-GW-S"
						zone = "nl-ams-1"
						ip_id = scaleway_vpc_public_gateway_ip.main.id
					}

					resource scaleway_vpc_public_gateway_pat_rule main {
						gateway_id = scaleway_vpc_public_gateway.main.id
						private_ip = scaleway_vpc_public_gateway_dhcp.main.address
						private_port = scaleway_rdb_instance.main.private_network.0.port
						public_port = 42
						protocol = "both"
						zone = "nl-ams-1"
						depends_on = [scaleway_vpc_gateway_network.main, scaleway_vpc_private_network.pn02]
					}

					resource scaleway_vpc_gateway_network main {
						gateway_id = scaleway_vpc_public_gateway.main.id
						private_network_id = scaleway_vpc_private_network.pn02.id
						dhcp_id = scaleway_vpc_public_gateway_dhcp.main.id
						cleanup_dhcp = true
						enable_masquerade = true
						zone = "nl-ams-1"
						depends_on = [scaleway_vpc_public_gateway_ip.main, scaleway_vpc_private_network.pn02]
					}

					resource scaleway_rdb_instance main {
						name = "test-rdb-private-network"
						node_type = "db-dev-s"
						engine = %q
						is_ha_cluster = false
						disable_backup = true
						user_name = "my_initial_user"
						password = "thiZ_is_v&ry_s3cret"
						region= "nl-ams"
						tags = [ "terraform-test", "scaleway_rdb_instance", "volume", "rdb_pn" ]
						volume_type = "bssd"
						volume_size_in_gb = 10
						private_network {
							ip_net = "192.168.1.254/24" #pool high
							pn_id = "${scaleway_vpc_private_network.pn02.id}"
						}
					}
				`, latestEngineVersion),
				Check: resource.ComposeTestCheckFunc(
					isInstancePresent(tt, "scaleway_rdb_instance.main"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "private_network.0.ip_net", "192.168.1.254/24"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource scaleway_vpc_private_network pn02 {
						name = "my_private_network"
						region= "nl-ams"
					}

					resource scaleway_rdb_instance main {
						name = "test-rdb-private-network"
						node_type = "db-dev-s"
						engine = %q
						is_ha_cluster = false
						disable_backup = true
						user_name = "my_initial_user"
						password = "thiZ_is_v&ry_s3cret"
						region= "nl-ams"
						tags = [ "terraform-test", "scaleway_rdb_instance", "volume", "rdb_pn" ]
						volume_type = "bssd"
						volume_size_in_gb = 10
					}
				`, latestEngineVersion),
			},
		},
	})
}

func TestAccInstance_PrivateNetworkUpdate(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	latestEngineVersion := rdbchecks.GetLatestEngineVersion(tt, postgreSQLEngineName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rdbchecks.IsInstanceDestroyed(tt),
			vpcchecks.CheckPrivateNetworkDestroy(tt),
		),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource scaleway_vpc_private_network pn01 {
						name = "my_private_network"
					}

					resource scaleway_rdb_instance main {
						name = "test-rdb-private-network-update"
						node_type = "db-dev-s"
						engine = %q
						is_ha_cluster = false
						disable_backup = true
						user_name = "my_initial_user"
						password = "thiZ_is_v&ry_s3cret"
						tags = [ "terraform-test", "scaleway_rdb_instance", "volume", "rdb_pn" ]
						volume_type = "bssd"
						volume_size_in_gb = 10
						private_network {
							pn_id = "${scaleway_vpc_private_network.pn01.id}"
							enable_ipam = true
						}
					}
				`, latestEngineVersion),
				Check: resource.ComposeTestCheckFunc(
					vpcchecks.IsPrivateNetworkPresent(tt, "scaleway_vpc_private_network.pn01"),
					isInstancePresent(tt, "scaleway_rdb_instance.main"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "private_network.#", "1"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "private_network.0.enable_ipam", "true"),
					resource.TestCheckResourceAttrPair("scaleway_rdb_instance.main", "private_network.0.pn_id", "scaleway_vpc_private_network.pn01", "id"),
				),
			},
			// Change PN but keep ipam config
			{
				Config: fmt.Sprintf(`
					resource scaleway_vpc_private_network pn01 {
						name = "my_private_network"
					}
					resource scaleway_vpc_private_network pn02 {
						name = "my_second_private_network"
					}
			
					resource scaleway_rdb_instance main {
						name = "test-rdb-private-network-update"
						node_type = "db-dev-s"
						engine = %q
						is_ha_cluster = false
						disable_backup = true
						user_name = "my_initial_user"
						password = "thiZ_is_v&ry_s3cret"
						tags = [ "terraform-test", "scaleway_rdb_instance", "volume", "rdb_pn" ]
						volume_type = "bssd"
						volume_size_in_gb = 10
						private_network {
							pn_id = "${scaleway_vpc_private_network.pn02.id}"
							enable_ipam = true
						}
					}
				`, latestEngineVersion),
				Check: resource.ComposeTestCheckFunc(
					vpcchecks.IsPrivateNetworkPresent(tt, "scaleway_vpc_private_network.pn01"),
					vpcchecks.IsPrivateNetworkPresent(tt, "scaleway_vpc_private_network.pn02"),
					isInstancePresent(tt, "scaleway_rdb_instance.main"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "private_network.#", "1"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "private_network.0.enable_ipam", "true"),
					resource.TestCheckResourceAttrPair("scaleway_rdb_instance.main", "private_network.0.pn_id", "scaleway_vpc_private_network.pn02", "id"),
				),
			},
			// Keep PN but change ipam config -> static
			{
				Config: fmt.Sprintf(`
					resource scaleway_vpc_private_network pn01 {
						name = "my_private_network"
					}
					resource scaleway_vpc_private_network pn02 {
						name = "my_second_private_network"
					}
			
					locals {
						ip_address  = cidrhost(scaleway_vpc_private_network.pn02.ipv4_subnet.0.subnet, 4)
						cidr_prefix = split("/", scaleway_vpc_private_network.pn02.ipv4_subnet.0.subnet)[1]
					}
			
					resource scaleway_rdb_instance main {
						name = "test-rdb-private-network-update"
						node_type = "db-dev-s"
						engine = %q
						is_ha_cluster = false
						disable_backup = true
						user_name = "my_initial_user"
						password = "thiZ_is_v&ry_s3cret"
						tags = [ "terraform-test", "scaleway_rdb_instance", "volume", "rdb_pn" ]
						volume_type = "bssd"
						volume_size_in_gb = 10`, latestEngineVersion) + `
						private_network {
							ip_net = format("%s/%s", local.ip_address, local.cidr_prefix)
							pn_id  = scaleway_vpc_private_network.pn02.id
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					vpcchecks.IsPrivateNetworkPresent(tt, "scaleway_vpc_private_network.pn02"),
					isInstancePresent(tt, "scaleway_rdb_instance.main"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "private_network.#", "1"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "private_network.0.enable_ipam", "false"),
					resource.TestCheckResourceAttrPair("scaleway_rdb_instance.main", "private_network.0.pn_id", "scaleway_vpc_private_network.pn02", "id"),
				),
			},
			// Change PN but keep static config
			{
				Config: fmt.Sprintf(`
					resource scaleway_vpc_private_network pn01 {
						name = "my_private_network"
					}
					resource scaleway_vpc_private_network pn02 {
						name = "my_second_private_network"
					}

					locals {
						ip_address  = cidrhost(scaleway_vpc_private_network.pn01.ipv4_subnet.0.subnet, 4)
						cidr_prefix = split("/", scaleway_vpc_private_network.pn01.ipv4_subnet.0.subnet)[1]
					}

					resource scaleway_rdb_instance main {
						name = "test-rdb-private-network-update"
						node_type = "db-dev-s"
						engine = %q
						is_ha_cluster = false
						disable_backup = true
						user_name = "my_initial_user"
						password = "thiZ_is_v&ry_s3cret"
						tags = [ "terraform-test", "scaleway_rdb_instance", "volume", "rdb_pn" ]
						volume_type = "bssd"
						volume_size_in_gb = 10`, latestEngineVersion) + `
						private_network {
							ip_net = format("%s/%s", local.ip_address, local.cidr_prefix)
							pn_id  = scaleway_vpc_private_network.pn01.id
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					vpcchecks.IsPrivateNetworkPresent(tt, "scaleway_vpc_private_network.pn01"),
					vpcchecks.IsPrivateNetworkPresent(tt, "scaleway_vpc_private_network.pn02"),
					isInstancePresent(tt, "scaleway_rdb_instance.main"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "private_network.#", "1"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "private_network.0.enable_ipam", "false"),
					resource.TestCheckResourceAttrPair("scaleway_rdb_instance.main", "private_network.0.pn_id", "scaleway_vpc_private_network.pn01", "id"),
				),
			},
			// Keep PN but change static config -> ipam
			{
				Config: fmt.Sprintf(`
					resource scaleway_vpc_private_network pn01 {
						name = "my_private_network"
					}
			
					resource scaleway_rdb_instance main {
						name = "test-rdb-private-network-update"
						node_type = "db-dev-s"
						engine = %q
						is_ha_cluster = false
						disable_backup = true
						user_name = "my_initial_user"
						password = "thiZ_is_v&ry_s3cret"
						tags = [ "terraform-test", "scaleway_rdb_instance", "volume", "rdb_pn" ]
						volume_type = "bssd"
						volume_size_in_gb = 10`, latestEngineVersion) + `
						private_network {
							pn_id  = scaleway_vpc_private_network.pn01.id
							enable_ipam = true
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					vpcchecks.IsPrivateNetworkPresent(tt, "scaleway_vpc_private_network.pn01"),
					isInstancePresent(tt, "scaleway_rdb_instance.main"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "private_network.#", "1"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "private_network.0.enable_ipam", "true"),
					resource.TestCheckResourceAttrPair("scaleway_rdb_instance.main", "private_network.0.pn_id", "scaleway_vpc_private_network.pn01", "id"),
				),
			},
		},
	})
}

func TestAccInstance_PrivateNetwork_DHCP(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	latestEngineVersion := rdbchecks.GetLatestEngineVersion(tt, postgreSQLEngineName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      rdbchecks.IsInstanceDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource scaleway_vpc_private_network pn02 {
						name = "my_private_network"
						region= "nl-ams"
					}

					resource scaleway_vpc_public_gateway_dhcp main {
						subnet = "192.168.1.0/24"
						zone = "nl-ams-1"
					}

					resource scaleway_vpc_public_gateway_ip main {
						zone = "nl-ams-1"
					}

					resource scaleway_vpc_public_gateway main {
						name = "foobar"
						type = "VPC-GW-S"
						zone = "nl-ams-1"
						ip_id = scaleway_vpc_public_gateway_ip.main.id
					}

					resource scaleway_vpc_public_gateway_pat_rule main {
						gateway_id = scaleway_vpc_public_gateway.main.id
						private_ip = scaleway_vpc_public_gateway_dhcp.main.address
						private_port = scaleway_rdb_instance.main.private_network.0.port
						public_port = 42
						protocol = "both"
						zone = "nl-ams-1"
						depends_on = [scaleway_vpc_gateway_network.main, scaleway_vpc_private_network.pn02]
					}

					resource scaleway_vpc_gateway_network main {
						gateway_id = scaleway_vpc_public_gateway.main.id
						private_network_id = scaleway_vpc_private_network.pn02.id
						dhcp_id = scaleway_vpc_public_gateway_dhcp.main.id
						cleanup_dhcp = true
						enable_masquerade = true
						zone = "nl-ams-1"
						depends_on = [scaleway_vpc_public_gateway_ip.main, scaleway_vpc_private_network.pn02]
					}

					resource scaleway_rdb_instance main {
						name = "test-rdb-private-network-dhcp"
						node_type = "db-dev-s"
						engine = %q
						is_ha_cluster = false
						disable_backup = true
						user_name = "my_initial_user"
						password = "thiZ_is_v&ry_s3cret"
						region= "nl-ams"
						tags = [ "terraform-test", "scaleway_rdb_instance", "volume", "rdb_pn" ]
						volume_type = "bssd"
						volume_size_in_gb = 10
						private_network {
							ip_net = "192.168.1.254/24" #pool high
							pn_id = "${scaleway_vpc_private_network.pn02.id}"
						}
					}
				`, latestEngineVersion),
				Check: resource.ComposeTestCheckFunc(
					isInstancePresent(tt, "scaleway_rdb_instance.main"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "private_network.0.ip_net", "192.168.1.254/24"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource scaleway_vpc_private_network pn02 {
						name = "my_private_network"
						region= "nl-ams"
					}

					resource scaleway_rdb_instance main {
						name = "test-rdb-private-network-dhcp"
						node_type = "db-dev-s"
						engine = %q
						is_ha_cluster = false
						disable_backup = true
						user_name = "my_initial_user"
						password = "thiZ_is_v&ry_s3cret"
						region= "nl-ams"
						tags = [ "terraform-test", "scaleway_rdb_instance", "volume", "rdb_pn" ]
						volume_type = "bssd"
						volume_size_in_gb = 10
					}
				`, latestEngineVersion),
			},
		},
	})
}

func TestAccInstance_BackupSchedule(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	latestEngineVersion := rdbchecks.GetLatestEngineVersion(tt, postgreSQLEngineName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      rdbchecks.IsInstanceDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource scaleway_rdb_instance main {
						name                      = "test-rdb-instance-backup-schedule"
						node_type                 = "db-dev-s"
						engine                    = %q
						is_ha_cluster             = false
						disable_backup            = false
                        backup_schedule_frequency = 24
                        backup_schedule_retention = 7
						backup_same_region        = true
						user_name                 = "my_initial_user"
						password                  = "thiZ_is_v&ry_s3cret"
						region                    = "nl-ams"
						tags                      = [ "terraform-test", "scaleway_rdb_instance", "volume" ]
					}
				`, latestEngineVersion),
				Check: resource.ComposeTestCheckFunc(
					isInstancePresent(tt, "scaleway_rdb_instance.main"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "disable_backup", "false"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "backup_schedule_frequency", "24"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "backup_schedule_retention", "7"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "backup_same_region", "true"),
				),
			},
		},
	})
}

func TestAccInstance_Volume(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	latestEngineVersion := rdbchecks.GetLatestEngineVersion(tt, postgreSQLEngineName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      rdbchecks.IsInstanceDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource scaleway_rdb_instance main {
						name = "test-rdb-instance-volume"
						node_type = "db-dev-s"
						engine = %q
						is_ha_cluster = false
						disable_backup = true
						user_name = "my_initial_user"
						password = "thiZ_is_v&ry_s3cret"
						region= "nl-ams"
						tags = [ "terraform-test", "scaleway_rdb_instance", "volume" ]
					}
				`, latestEngineVersion),
				Check: resource.ComposeTestCheckFunc(
					isInstancePresent(tt, "scaleway_rdb_instance.main"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "volume_type", "lssd"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource scaleway_rdb_instance main {
						name = "test-rdb-instance-volume"
						node_type = "db-dev-s"
						engine = %q
						is_ha_cluster = false
						disable_backup = true
						user_name = "my_initial_user"
						password = "thiZ_is_v&ry_s3cret"
						region= "nl-ams"
						tags = [ "terraform-test", "scaleway_rdb_instance", "volume" ]
						volume_type = "bssd"
						volume_size_in_gb = 10
					}
				`, latestEngineVersion),
				Check: resource.ComposeTestCheckFunc(
					isInstancePresent(tt, "scaleway_rdb_instance.main"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "volume_type", "bssd"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "volume_size_in_gb", "10"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource scaleway_rdb_instance main {
						name = "test-rdb-instance-volume"
						node_type = "db-dev-m"
						engine = %q
						is_ha_cluster = false
						disable_backup = true
						user_name = "my_initial_user"
						password = "thiZ_is_v&ry_s3cret"
						region= "nl-ams"
						tags = [ "terraform-test", "scaleway_rdb_instance", "volume" ]
					}
				`, latestEngineVersion),
				Check: resource.ComposeTestCheckFunc(
					isInstancePresent(tt, "scaleway_rdb_instance.main"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "volume_type", "lssd"),
				),
			},
		},
	})
}

func TestAccInstance_SBSVolume(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	latestEngineVersion := rdbchecks.GetLatestEngineVersion(tt, postgreSQLEngineName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      rdbchecks.IsInstanceDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource scaleway_rdb_instance main {
						name = "test-rdb-instance-volume"
						node_type = "db-play2-pico"
						engine = %q
						is_ha_cluster = false
						disable_backup = true
						user_name = "my_initial_user"
						password = "thiZ_is_v&ry_s3cret"
						region= "nl-ams"
						tags = [ "terraform-test", "scaleway_rdb_instance", "sdb-volume" ]
						volume_type = "sbs_5k"
						volume_size_in_gb = 10
					}
				`, latestEngineVersion),
				Check: resource.ComposeTestCheckFunc(
					isInstancePresent(tt, "scaleway_rdb_instance.main"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "volume_type", "sbs_5k"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "volume_size_in_gb", "10"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource scaleway_rdb_instance main {
						name = "test-rdb-instance-volume"
						node_type = "db-play2-pico"
						engine = %q
						is_ha_cluster = false
						disable_backup = true
						user_name = "my_initial_user"
						password = "thiZ_is_v&ry_s3cret"
						region= "nl-ams"
						tags = [ "terraform-test", "scaleway_rdb_instance", "volume" ]
						volume_type = "sbs_5k"
						volume_size_in_gb = 20
					}
				`, latestEngineVersion),
				Check: resource.ComposeTestCheckFunc(
					isInstancePresent(tt, "scaleway_rdb_instance.main"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "volume_type", "sbs_5k"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "volume_size_in_gb", "20"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource scaleway_rdb_instance main {
						name = "test-rdb-instance-volume"
						node_type = "db-play2-pico"
						engine = %q
						is_ha_cluster = false
						disable_backup = true
						user_name = "my_initial_user"
						password = "thiZ_is_v&ry_s3cret"
						region= "nl-ams"
						tags = [ "terraform-test", "scaleway_rdb_instance", "volume" ]
						volume_type = "sbs_15k"
						volume_size_in_gb = 20
					}
				`, latestEngineVersion),
				Check: resource.ComposeTestCheckFunc(
					isInstancePresent(tt, "scaleway_rdb_instance.main"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "volume_type", "sbs_15k"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "volume_size_in_gb", "20"),
				),
			},
		},
	})
}

func TestAccInstance_ChangeVolumeType(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	latestEngineVersion := rdbchecks.GetLatestEngineVersion(tt, postgreSQLEngineName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      rdbchecks.IsInstanceDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource scaleway_rdb_instance main {
						name = "test-rdb-instance-volume"
						node_type = "db-play2-pico"
						engine = %q
						is_ha_cluster = false
						disable_backup = true
						user_name = "my_initial_user"
						password = "thiZ_is_v&ry_s3cret"
						region= "nl-ams"
						tags = [ "terraform-test", "scaleway_rdb_instance", "sdb-volume" ]
						volume_type = "bssd"
						volume_size_in_gb = 10
					}
				`, latestEngineVersion),
				Check: resource.ComposeTestCheckFunc(
					isInstancePresent(tt, "scaleway_rdb_instance.main"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "volume_type", "bssd"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "volume_size_in_gb", "10"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource scaleway_rdb_instance main {
						name = "test-rdb-instance-volume"
						node_type = "db-play2-pico"
						engine = %q
						is_ha_cluster = false
						disable_backup = true
						user_name = "my_initial_user"
						password = "thiZ_is_v&ry_s3cret"
						region= "nl-ams"
						tags = [ "terraform-test", "scaleway_rdb_instance", "volume" ]
						volume_type = "sbs_5k"
						volume_size_in_gb = 20
					}
				`, latestEngineVersion),
				Check: resource.ComposeTestCheckFunc(
					isInstancePresent(tt, "scaleway_rdb_instance.main"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "volume_type", "sbs_5k"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "volume_size_in_gb", "20"),
				),
			},
		},
	})
}

func TestAccInstance_ChangeNodeType(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	latestEngineVersion := rdbchecks.GetLatestEngineVersion(tt, postgreSQLEngineName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() { acctest.PreCheck(t) },

		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      rdbchecks.IsInstanceDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource scaleway_rdb_instance main {
						name = "test-rdb-instance-volume"
						node_type = "db-play2-pico"
						engine = %q
						is_ha_cluster = false
						disable_backup = true
						user_name = "my_initial_user"
						password = "thiZ_is_v&ry_s3cret"
						region= "nl-ams"
						tags = [ "terraform-test", "scaleway_rdb_instance" ]
						volume_type = "bssd"
						volume_size_in_gb = 10
					}
				`, latestEngineVersion),
				Check: resource.ComposeTestCheckFunc(
					isInstancePresent(tt, "scaleway_rdb_instance.main"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "node_type", "db-play2-pico"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource scaleway_rdb_instance main {
						name = "test-rdb-instance-volume"
						node_type = "db-play2-nano"
						engine = %q
						is_ha_cluster = false
						disable_backup = true
						user_name = "my_initial_user"
						password = "thiZ_is_v&ry_s3cret"
						region= "nl-ams"
						tags = [ "terraform-test", "scaleway_rdb_instance"]
						volume_type = "bssd"
						volume_size_in_gb = 10
					}
				`, latestEngineVersion),
				Check: resource.ComposeTestCheckFunc(
					isInstancePresent(tt, "scaleway_rdb_instance.main"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "node_type", "db-play2-nano"),
				),
			},
		},
	})
}

func TestAccInstance_Endpoints(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	latestEngineVersion := rdbchecks.GetLatestEngineVersion(tt, postgreSQLEngineName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      rdbchecks.IsInstanceDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "scaleway_vpc_private_network" "test_endpoints" {
						name = "test-rdb-endpoints"
					}

					resource "scaleway_rdb_instance" "test_endpoints" {
						name = "test-rdb-endpoints"
						node_type = "db-dev-s"
						engine = %q
						is_ha_cluster = false
						disable_backup = true
						user_name = "my_initial_user"
						password = "thiZ_is_v&ry_s3cret"
						tags = [ "terraform-test", "scaleway_rdb_instance", "test_endpoints" ]
						private_network {
							pn_id = scaleway_vpc_private_network.test_endpoints.id
							enable_ipam = true
						}
					}
				`, latestEngineVersion),
				Check: resource.ComposeTestCheckFunc(
					isInstancePresent(tt, "scaleway_rdb_instance.test_endpoints"),
					vpcchecks.IsPrivateNetworkPresent(tt, "scaleway_vpc_private_network.test_endpoints"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.test_endpoints", "private_network.#", "1"),
					resource.TestCheckResourceAttrPair("scaleway_rdb_instance.test_endpoints", "private_network.0.pn_id", "scaleway_vpc_private_network.test_endpoints", "id"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.test_endpoints", "private_network.0.enable_ipam", "true"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.test_endpoints", "load_balancer.#", "0"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.test_endpoints", "endpoint_ip", ""),    // Deprecated attribute, might be deleted later
					resource.TestCheckResourceAttr("scaleway_rdb_instance.test_endpoints", "endpoint_port", "0"), // Deprecated attribute, might be deleted later
					resource.TestCheckResourceAttrSet("scaleway_rdb_instance.test_endpoints", "private_ip.0.id"),
					resource.TestCheckResourceAttrSet("scaleway_rdb_instance.test_endpoints", "private_ip.0.address"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "scaleway_vpc_private_network" "test_endpoints" {
						name = "test-rdb-endpoints"
					}

					resource "scaleway_rdb_instance" "test_endpoints" {
						name = "test-rdb-endpoints"
						node_type = "db-dev-s"
						engine = %q
						is_ha_cluster = false
						disable_backup = true
						user_name = "my_initial_user"
						password = "thiZ_is_v&ry_s3cret"
						tags = [ "terraform-test", "scaleway_rdb_instance", "test_endpoints" ]
						private_network {
							pn_id = scaleway_vpc_private_network.test_endpoints.id
							enable_ipam = true
						}
						load_balancer {}
					}
				`, latestEngineVersion),
				Check: resource.ComposeTestCheckFunc(
					isInstancePresent(tt, "scaleway_rdb_instance.test_endpoints"),
					vpcchecks.IsPrivateNetworkPresent(tt, "scaleway_vpc_private_network.test_endpoints"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.test_endpoints", "load_balancer.#", "1"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.test_endpoints", "private_network.#", "1"),
					resource.TestCheckResourceAttrPair("scaleway_rdb_instance.test_endpoints", "private_network.0.pn_id", "scaleway_vpc_private_network.test_endpoints", "id"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.test_endpoints", "private_network.0.enable_ipam", "true"),
					resource.TestCheckResourceAttrSet("scaleway_rdb_instance.test_endpoints", "endpoint_ip"),   // Deprecated attribute, might be deleted later
					resource.TestCheckResourceAttrSet("scaleway_rdb_instance.test_endpoints", "endpoint_port"), // Deprecated attribute, might be deleted later
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "scaleway_vpc_private_network" "test_endpoints" {
						name = "test-rdb-endpoints"
					}

					resource "scaleway_rdb_instance" "test_endpoints" {
						name = "test-rdb-endpoints"
						node_type = "db-dev-s"
						engine = %q
						is_ha_cluster = false
						disable_backup = true
						user_name = "my_initial_user"
						password = "thiZ_is_v&ry_s3cret"
						tags = [ "terraform-test", "scaleway_rdb_instance", "test_endpoints" ]
					}
				`, latestEngineVersion),
				Check: resource.ComposeTestCheckFunc(
					isInstancePresent(tt, "scaleway_rdb_instance.test_endpoints"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.test_endpoints", "load_balancer.#", "1"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.test_endpoints", "private_network.#", "0"),
					resource.TestCheckResourceAttrSet("scaleway_rdb_instance.test_endpoints", "endpoint_ip"),   // Deprecated attribute, might be deleted later
					resource.TestCheckResourceAttrSet("scaleway_rdb_instance.test_endpoints", "endpoint_port"), // Deprecated attribute, might be deleted later
				),
			},
		},
	})
}

func TestAccInstance_EncryptionAtRest(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	latestEngineVersion := rdbchecks.GetLatestEngineVersion(tt, postgreSQLEngineName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      rdbchecks.IsInstanceDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource scaleway_rdb_instance main {
						name               = "test-rdb-encryption"
						node_type          = "db-dev-s"
						engine             = %q
						is_ha_cluster      = false
						disable_backup     = true
						user_name          = "my_initial_user"
						password           = "thiZ_is_v&ry_s3cret"
						encryption_at_rest = true
						tags               = [ "terraform-test", "scaleway_rdb_instance", "encryption" ]
					}
				`, latestEngineVersion),
				Check: resource.ComposeTestCheckFunc(
					isInstancePresent(tt, "scaleway_rdb_instance.main"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "name", "test-rdb-encryption"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "node_type", "db-dev-s"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "engine", latestEngineVersion),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "is_ha_cluster", "false"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "disable_backup", "true"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "user_name", "my_initial_user"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "password", "thiZ_is_v&ry_s3cret"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "encryption_at_rest", "true"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "tags.0", "terraform-test"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "tags.1", "scaleway_rdb_instance"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "tags.2", "encryption"),
				),
			},
		},
	})
}

func TestAccInstance_EncryptionAtRestFalse(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	latestEngineVersion := rdbchecks.GetLatestEngineVersion(tt, postgreSQLEngineName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      rdbchecks.IsInstanceDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource scaleway_rdb_instance main {
						name               = "test-rdb-no-encryption"
						node_type          = "db-dev-s"
						engine             = %q
						is_ha_cluster      = false
						disable_backup     = true
						user_name          = "my_initial_user_no_enc"
						password           = "thiZ_is_v&ry_s3cret"
						encryption_at_rest = false
						tags               = [ "terraform-test", "scaleway_rdb_instance", "no_encryption" ]
					}
				`, latestEngineVersion),
				Check: resource.ComposeTestCheckFunc(
					isInstancePresent(tt, "scaleway_rdb_instance.main"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "name", "test-rdb-no-encryption"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "node_type", "db-dev-s"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "engine", latestEngineVersion),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "is_ha_cluster", "false"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "disable_backup", "true"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "user_name", "my_initial_user_no_enc"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "password", "thiZ_is_v&ry_s3cret"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "encryption_at_rest", "false"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "tags.0", "terraform-test"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "tags.1", "scaleway_rdb_instance"),
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "tags.2", "no_encryption"),
				),
			},
		},
	})
}

func isInstancePresent(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		rdbAPI, region, ID, err := rdb.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = rdbAPI.GetInstance(&rdbSDK.GetInstanceRequest{
			InstanceID: ID,
			Region:     region,
		})
		if err != nil {
			return err
		}

		return nil
	}
}
