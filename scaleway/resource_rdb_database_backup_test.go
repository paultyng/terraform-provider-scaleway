package scaleway

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("scaleway_rdb_database_backup", &resource.Sweeper{
		Name: "scaleway_rdb_database_backup",
		F:    testSweepRDBInstance,
	})
}

func TestAccScalewayRdbDatabaseBackup_Basic(t *testing.T) {
	tt := NewTestTools(t)
	defer tt.Cleanup()
	instanceName := "TestAccScalewayRdbUser_Basic"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckScalewayRdbInstanceDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource scaleway_rdb_instance main {
						name = "%s"
						node_type = "db-dev-s"
						engine = "PostgreSQL-12"
						is_ha_cluster = false
						tags = [ "terraform-test", "scaleway_rdb_user", "minimal" ]
					}

					resource scaleway_rdb_user db_user {
						instance_id = scaleway_rdb_instance.main.id
						name = "foo"
						password = "R34lP4sSw#Rd"
						is_admin = true
					}`, instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdbUserExists(tt, "scaleway_rdb_instance.main", "scaleway_rdb_user.db_user"),
					resource.TestCheckResourceAttr("scaleway_rdb_user.db_user", "name", "foo"),
					resource.TestCheckResourceAttr("scaleway_rdb_user.db_user", "is_admin", "true"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource scaleway_rdb_instance main {
						name = "%s"
						node_type = "db-dev-s"
						engine = "PostgreSQL-12"
						is_ha_cluster = false
						tags = [ "terraform-test", "scaleway_rdb_user", "minimal" ]
					}`, instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdbUserExists(tt, "scaleway_rdb_instance.main", "scaleway_rdb_user.db_user"),
					resource.TestCheckResourceAttr("scaleway_rdb_user.db_user", "name", "bar"),
					resource.TestCheckResourceAttr("scaleway_rdb_user.db_user", "is_admin", "false"),
				),
			},
		},
	})
}

//func testAccCheckRdbDatabaseBackupExists(tt *TestTools, instance string, databaseBackupName string) resource.TestCheckFunc {
//	return func(state *terraform.State) error {
//		instanceResource, ok := state.RootModule().Resources[instance]
//		if !ok {
//			return fmt.Errorf("resource not found: %s", instance)
//		}
//
//		userResource, ok := state.RootModule().Resources[databaseBackupName]
//		if !ok {
//			return fmt.Errorf("resource not found: %s", databaseBackupName)
//		}
//
//		rdbAPI, region, _, err := rdbAPIWithRegionAndID(tt.Meta, instanceResource.Primary.ID)
//		if err != nil {
//			return err
//		}
//
//		instanceID, databaseBackupName, err := resourceScalewayRdbDatabaseBackupParseID(userResource.Primary.ID)
//		if err != nil {
//			return err
//		}
//
//		databaseBackups, err := rdbAPI.ListDatabaseBackups(&rdb.ListDatabaseBackupsRequest{
//			Region:     region,
//			Name:       &databaseBackupName,
//			OrderBy:    "",
//			InstanceID: scw.StringPtr(instanceID),
//		})
//		if err != nil {
//			return err
//		}
//
//		if len(databaseBackups.DatabaseBackups) != 1 {
//			return fmt.Errorf("no database backup found")
//		}
//
//		return nil
//	}
//}
