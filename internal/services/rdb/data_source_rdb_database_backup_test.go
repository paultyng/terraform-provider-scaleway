package rdb_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	rdbchecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/rdb/testfuncs"
)

func TestAccDataSourceDatabaseBackup_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	latestEngineVersion := rdbchecks.GetLatestEngineVersion(tt, postgreSQLEngineName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rdbchecks.IsInstanceDestroyed(tt),
			isBackupDestroyed(tt),
		),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "scaleway_rdb_instance" "server" {
						name      = "test-terraform"
						node_type = "db-dev-s"
						engine    = %q
					}
					resource "scaleway_rdb_database" "database" {
						name        = "test-terraform"
						instance_id = scaleway_rdb_instance.server.id
					}`, latestEngineVersion),
			},
			{
				Config: fmt.Sprintf(`
					resource "scaleway_rdb_instance" "server" {
						name      = "test-terraform"
						node_type = "db-dev-s"
						engine    = %q
					}
					resource "scaleway_rdb_database" "database" {
						name        = "test-terraform"
						instance_id = scaleway_rdb_instance.server.id
					}

					resource scaleway_rdb_database_backup backup {
						instance_id 	= scaleway_rdb_instance.server.id
  						database_name 	= scaleway_rdb_database.database.name
  						name 			= "test_backup_datasource"
					}

					data scaleway_rdb_database_backup find_by_name {
						name        = scaleway_rdb_database_backup.backup.name
					}

					data scaleway_rdb_database_backup find_by_name_and_instance {
						name        = scaleway_rdb_database_backup.backup.name
						instance_id = scaleway_rdb_instance.server.id
					}

					data scaleway_rdb_database_backup find_by_id {
						backup_id 	= scaleway_rdb_database_backup.backup.id
					}
				`, latestEngineVersion),
				Check: resource.ComposeTestCheckFunc(
					isDatabasePresent(tt, "scaleway_rdb_instance.server", "scaleway_rdb_database.database"),
					isBackupPresent(tt, "scaleway_rdb_database_backup.backup"),

					resource.TestCheckResourceAttr("data.scaleway_rdb_database_backup.find_by_name", "name", "test_backup_datasource"),
					resource.TestCheckResourceAttr("data.scaleway_rdb_database_backup.find_by_name_and_instance", "name", "test_backup_datasource"),
					resource.TestCheckResourceAttr("data.scaleway_rdb_database_backup.find_by_id", "name", "test_backup_datasource"),
				),
			},
		},
	})
}
