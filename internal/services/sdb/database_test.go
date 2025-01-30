package sdb_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	sdbSDK "github.com/scaleway/scaleway-sdk-go/api/serverless_sqldb/v1alpha1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/sdb"
)

func TestAccServerlessSQLDBDatabase_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckServerlessSQLDBDatabaseDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_sdb_sql_database main {
						name = "test-sdb-sql-database-basic"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerlessSQLDBDatabaseExists(tt, "scaleway_sdb_sql_database.main"),
					acctest.CheckResourceAttrUUID("scaleway_sdb_sql_database.main", "id"),
					resource.TestCheckResourceAttr("scaleway_sdb_sql_database.main", "name", "test-sdb-sql-database-basic"),
					resource.TestCheckResourceAttr("scaleway_sdb_sql_database.main", "min_cpu", "0"),
					resource.TestCheckResourceAttr("scaleway_sdb_sql_database.main", "max_cpu", "15"),
					resource.TestCheckResourceAttrSet("scaleway_sdb_sql_database.main", "endpoint"),
				),
			},
			{
				Config: `
					resource scaleway_sdb_sql_database main {
						name = "test-sdb-sql-database-basic"
						max_cpu = 6
						min_cpu = 2
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerlessSQLDBDatabaseExists(tt, "scaleway_sdb_sql_database.main"),
					acctest.CheckResourceAttrUUID("scaleway_sdb_sql_database.main", "id"),
					resource.TestCheckResourceAttr("scaleway_sdb_sql_database.main", "min_cpu", "2"),
					resource.TestCheckResourceAttr("scaleway_sdb_sql_database.main", "max_cpu", "6"),
				),
			},
			{
				Config: `
					resource scaleway_sdb_sql_database main {
						name = "test-sdb-sql-database-basic"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerlessSQLDBDatabaseExists(tt, "scaleway_sdb_sql_database.main"),
					acctest.CheckResourceAttrUUID("scaleway_sdb_sql_database.main", "id"),
					resource.TestCheckResourceAttr("scaleway_sdb_sql_database.main", "min_cpu", "0"),
					resource.TestCheckResourceAttr("scaleway_sdb_sql_database.main", "max_cpu", "15"),
				),
			},
			{ // Should ForceNew, this is testing creation with cpu values
				Config: `
					resource scaleway_sdb_sql_database main {
						name = "test-sdb-sql-database-basic-rename"
						min_cpu = 4
						max_cpu = 8
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerlessSQLDBDatabaseExists(tt, "scaleway_sdb_sql_database.main"),
					acctest.CheckResourceAttrUUID("scaleway_sdb_sql_database.main", "id"),
					resource.TestCheckResourceAttr("scaleway_sdb_sql_database.main", "min_cpu", "4"),
					resource.TestCheckResourceAttr("scaleway_sdb_sql_database.main", "max_cpu", "8"),
				),
			},
		},
	})
}

func testAccCheckServerlessSQLDBDatabaseExists(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		api, region, id, err := sdb.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = api.GetDatabase(&sdbSDK.GetDatabaseRequest{
			DatabaseID: id,
			Region:     region,
		})
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckServerlessSQLDBDatabaseDestroy(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_sdb_sql_database" {
				continue
			}

			api, region, id, err := sdb.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}

			_, err = api.DeleteDatabase(&sdbSDK.DeleteDatabaseRequest{
				DatabaseID: id,
				Region:     region,
			})

			if err == nil {
				return fmt.Errorf("serverless_sql database (%s) still exists", rs.Primary.ID)
			}

			if !httperrors.Is403(err) {
				return err
			}
		}

		return nil
	}
}
