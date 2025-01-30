package rdb_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	rdbchecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/rdb/testfuncs"
)

func TestAccDataSourceACL_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	instanceName := "TestAccScalewayDataSourceRDBAcl_Basic"
	latestEngineVersion := rdbchecks.GetLatestEngineVersion(tt, postgreSQLEngineName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      rdbchecks.IsInstanceDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "scaleway_rdb_instance" "main" {
						name = "%s"
						node_type = "db-dev-s"
						engine = %q
						is_ha_cluster = false
					}

					resource "scaleway_rdb_acl" "main" {
						instance_id = scaleway_rdb_instance.main.id
						acl_rules {
							ip = "1.2.3.4/32"
							description = "foo"
						}

						acl_rules {
							ip = "4.5.6.7/32"
							description = "bar"
						}
					}
					`, instanceName, latestEngineVersion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaleway_rdb_acl.main", "acl_rules.0.ip", "1.2.3.4/32"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "scaleway_rdb_instance" "main" {
						name = "%s"
						node_type = "db-dev-s"
						engine = %q
						is_ha_cluster = false
					}

					resource "scaleway_rdb_acl" "main" {
						instance_id = scaleway_rdb_instance.main.id
						acl_rules {
							ip = "1.2.3.4/32"
							description = "foo"
						}

						acl_rules {
							ip = "4.5.6.7/32"
							description = "bar"
						}
					}

					data "scaleway_rdb_acl" "maindata" {
						instance_id = scaleway_rdb_instance.main.id
					}`, instanceName, latestEngineVersion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaleway_rdb_acl.main", "acl_rules.0.ip", "1.2.3.4/32"),
					resource.TestCheckResourceAttr("data.scaleway_rdb_acl.maindata", "acl_rules.0.ip", "1.2.3.4/32"),
				),
			},
		},
	})
}
