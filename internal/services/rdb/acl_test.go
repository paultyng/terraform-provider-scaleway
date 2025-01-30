package rdb_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	rdbchecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/rdb/testfuncs"
)

func TestAccACL_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	instanceName := "rdb-acl-basic"
	latestEngineVersion := rdbchecks.GetLatestEngineVersion(tt, postgreSQLEngineName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      rdbchecks.IsInstanceDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "scaleway_instance_ip" "front1_ip" {
					}

					resource "scaleway_instance_ip" "front2_ip" {
					}

					resource scaleway_rdb_instance main {
						name = "%s"
						node_type = "db-dev-s"
						engine = %q
						is_ha_cluster = false
					}

					resource "scaleway_rdb_acl" "main" {
					  instance_id = scaleway_rdb_instance.main.id

					  acl_rules {
						ip = "${scaleway_instance_ip.front1_ip.address}/32"
					  }

					  acl_rules {
						ip = "${scaleway_instance_ip.front2_ip.address}/32"
					  }
					}
				`, instanceName, latestEngineVersion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("scaleway_rdb_acl.main", "acl_rules.0.description"),
					resource.TestCheckResourceAttrSet("scaleway_rdb_acl.main", "acl_rules.1.description"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource scaleway_rdb_instance main {
						name = "%s"
						node_type = "db-dev-s"
						engine = %q
						is_ha_cluster = false
					}

					resource scaleway_rdb_acl main {
						instance_id = scaleway_rdb_instance.main.id
						acl_rules {
							ip = "1.2.3.4/32"
							description = "foo"
						}

						acl_rules {
							ip = "4.5.6.7/32"
							description = "bar"
						}
					}`, instanceName, latestEngineVersion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaleway_rdb_acl.main", "acl_rules.0.ip", "1.2.3.4/32"),
					resource.TestCheckResourceAttr("scaleway_rdb_acl.main", "acl_rules.0.description", "foo"),
					resource.TestCheckResourceAttr("scaleway_rdb_acl.main", "acl_rules.1.ip", "4.5.6.7/32"),
					resource.TestCheckResourceAttr("scaleway_rdb_acl.main", "acl_rules.1.description", "bar"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource scaleway_rdb_instance main {
						name = "%s"
						node_type = "db-dev-s"
						engine = %q
						is_ha_cluster = false
					}

					resource scaleway_rdb_acl main {
						instance_id = scaleway_rdb_instance.main.id
						acl_rules {
							ip = "1.2.3.4/32"
							description = "foo"
						}

						acl_rules {
							ip = "9.0.0.0/16"
							description = "baz"
						}
					}`, instanceName, latestEngineVersion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaleway_rdb_acl.main", "acl_rules.1.ip", "9.0.0.0/16"),
					resource.TestCheckResourceAttr("scaleway_rdb_acl.main", "acl_rules.1.description", "baz"),
					resource.TestCheckResourceAttr("scaleway_rdb_acl.main", "acl_rules.0.ip", "1.2.3.4/32"),
					resource.TestCheckResourceAttr("scaleway_rdb_acl.main", "acl_rules.0.description", "foo"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource scaleway_rdb_instance main {
						name = "%s"
						node_type = "db-dev-s"
						engine = %q
						is_ha_cluster = false
					}`, instanceName, latestEngineVersion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaleway_rdb_instance.main", "name", "rdb-acl-basic"),
				),
			},
		},
	})
}
