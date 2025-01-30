package rdb_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	rdbSDK "github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/rdb"
	rdbchecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/rdb/testfuncs"
)

func TestAccPrivilege_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	instanceName := "TestAccScalewayRdbPrivilege_Basic"
	latestEngineVersion := rdbchecks.GetLatestEngineVersion(tt, postgreSQLEngineName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      rdbchecks.IsInstanceDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "scaleway_rdb_instance" "instance" {
					  name          = "%s"
					  node_type     = "db-dev-s"
					  engine        = %q
					  is_ha_cluster = false
					  tags          = ["terraform-test", "scaleway_rdb_user", "minimal"]
					}

					resource "scaleway_rdb_database" "db01" {
					  instance_id = scaleway_rdb_instance.instance.id
					  name        = "foo"
					}

					resource "scaleway_rdb_user" "foo1" {
					  instance_id = scaleway_rdb_instance.instance.id
					  name        = "user_01"
					  password    = "R34lP4sSw#Rd"
					  is_admin    = true
					}

					resource "scaleway_rdb_privilege" "priv_admin" {
					  instance_id   = scaleway_rdb_instance.instance.id
					  user_name     = scaleway_rdb_user.foo1.name
					  database_name = scaleway_rdb_database.db01.name
					  permission    = "all"
					}
					`, instanceName, latestEngineVersion),
				Check: resource.ComposeTestCheckFunc(
					isPrivilegePresent(tt, "scaleway_rdb_instance.instance", "scaleway_rdb_database.db01", "scaleway_rdb_user.foo1"),
					resource.TestCheckResourceAttr("scaleway_rdb_privilege.priv_admin", "permission", "all"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "scaleway_rdb_instance" "instance" {
					  name          = "%s"
					  node_type     = "db-dev-s"
					  engine        = %q
					  is_ha_cluster = false
					  tags          = ["terraform-test", "scaleway_rdb_user", "minimal"]
					}

					resource "scaleway_rdb_database" "db01" {
					  instance_id = scaleway_rdb_instance.instance.id
					  name        = "foo"
					}

					resource "scaleway_rdb_user" "foo1" {
					  instance_id = scaleway_rdb_instance.instance.id
					  name        = "user_01"
					  password    = "R34lP4sSw#Rd"
					  is_admin    = true
					}

					resource "scaleway_rdb_privilege" "priv_admin" {
					  instance_id   = scaleway_rdb_instance.instance.id
					  user_name     = scaleway_rdb_user.foo1.name
					  database_name = scaleway_rdb_database.db01.name
					  permission    = "all"
					}

					resource "scaleway_rdb_user" "foo2" {
					  instance_id = scaleway_rdb_instance.instance.id
					  name        = "user_02"
					  password    = "R34lP4sSw#Rd"
					}

					resource "scaleway_rdb_privilege" "priv_foo_02" {
					  instance_id   = scaleway_rdb_instance.instance.id
					  user_name     = scaleway_rdb_user.foo2.name
					  database_name = scaleway_rdb_database.db01.name
					  permission    = "readwrite"
					}
					`, instanceName, latestEngineVersion),
				Check: resource.ComposeTestCheckFunc(
					isPrivilegePresent(tt, "scaleway_rdb_instance.instance", "scaleway_rdb_database.db01", "scaleway_rdb_user.foo2"),
					resource.TestCheckResourceAttr("scaleway_rdb_privilege.priv_foo_02", "permission", "readwrite"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "scaleway_rdb_instance" "instance" {
					  name          = "%s"
					  node_type     = "db-dev-s"
					  engine        = %q
					  is_ha_cluster = false
					  tags          = ["terraform-test", "scaleway_rdb_user", "minimal"]
					}

					resource "scaleway_rdb_database" "db01" {
					  instance_id = scaleway_rdb_instance.instance.id
					  name        = "foo"
					}

					resource "scaleway_rdb_user" "foo1" {
					  instance_id = scaleway_rdb_instance.instance.id
					  name        = "user_01"
					  password    = "R34lP4sSw#Rd"
					  is_admin    = true
					}

					resource "scaleway_rdb_privilege" "priv_admin" {
					  instance_id   = scaleway_rdb_instance.instance.id
					  user_name     = scaleway_rdb_user.foo1.name
					  database_name = scaleway_rdb_database.db01.name
					  permission    = "all"
					}

					resource "scaleway_rdb_user" "foo2" {
					  instance_id = scaleway_rdb_instance.instance.id
					  name        = "user_02"
					  password    = "R34lP4sSw#Rd"
					}

					resource "scaleway_rdb_privilege" "priv_foo_02" {
					  instance_id   = scaleway_rdb_instance.instance.id
					  user_name     = scaleway_rdb_user.foo2.name
					  database_name = scaleway_rdb_database.db01.name
					  permission    = "readwrite"
					}

					resource "scaleway_rdb_user" "foo3" {
					  instance_id = scaleway_rdb_instance.instance.id
					  name        = "user_03"
					  password    = "R34lP4sSw#Rd"
					}

					resource "scaleway_rdb_privilege" "priv_foo_03" {
					  instance_id   = scaleway_rdb_instance.instance.id
					  user_name     = scaleway_rdb_user.foo3.name
					  database_name = scaleway_rdb_database.db01.name
					  permission    = "none"
					}
					`, instanceName, latestEngineVersion),
				Check: resource.ComposeTestCheckFunc(
					isPrivilegePresent(tt, "scaleway_rdb_instance.instance", "scaleway_rdb_database.db01", "scaleway_rdb_user.foo3"),
					resource.TestCheckResourceAttr("scaleway_rdb_privilege.priv_foo_03", "permission", "none"),
				),
			},
		},
	})
}

func isPrivilegePresent(tt *acctest.TestTools, instance string, database string, user string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		instanceResource, ok := state.RootModule().Resources[instance]
		if !ok {
			return fmt.Errorf("resource not found: %s", instance)
		}

		databaseResource, ok := state.RootModule().Resources[database]
		if !ok {
			return fmt.Errorf("resource database not found: %s", database)
		}

		userResource, ok := state.RootModule().Resources[user]
		if !ok {
			return fmt.Errorf("resource not found: %s", user)
		}

		rdbAPI, _, _, err := rdb.NewAPIWithRegionAndID(tt.Meta, instanceResource.Primary.ID)
		if err != nil {
			return err
		}

		_, _, databaseName, err := rdb.ResourceRdbDatabaseParseID(databaseResource.Primary.ID)
		if err != nil {
			return err
		}

		region, instanceID, userName, err := rdb.ResourceUserParseID(userResource.Primary.ID)
		if err != nil {
			return err
		}

		databases, err := rdbAPI.ListPrivileges(&rdbSDK.ListPrivilegesRequest{
			Region:       region,
			InstanceID:   instanceID,
			DatabaseName: &databaseName,
			UserName:     &userName,
		})
		if err != nil {
			return err
		}

		if len(databases.Privileges) != 1 {
			return errors.New("no privilege found")
		}

		return nil
	}
}
