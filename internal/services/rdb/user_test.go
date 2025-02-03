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

func TestAccUser_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	instanceName := "TestAccScalewayRdbUser_Basic"
	latestEngineVersion := rdbchecks.GetLatestEngineVersion(tt, postgreSQLEngineName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      rdbchecks.IsInstanceDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource scaleway_rdb_instance main {
						name = "%s"
						node_type = "db-dev-s"
						engine = %q
						is_ha_cluster = false
						tags = [ "terraform-test", "scaleway_rdb_user", "minimal" ]
					}

					resource scaleway_rdb_user db_user {
						instance_id = scaleway_rdb_instance.main.id
						name = "foo"
						password = "R34lP4sSw#Rd"
						is_admin = true
					}`, instanceName, latestEngineVersion),
				Check: resource.ComposeTestCheckFunc(
					isUserPresent(tt, "scaleway_rdb_instance.main", "scaleway_rdb_user.db_user"),
					resource.TestCheckResourceAttr("scaleway_rdb_user.db_user", "name", "foo"),
					resource.TestCheckResourceAttr("scaleway_rdb_user.db_user", "is_admin", "true"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource scaleway_rdb_instance main {
						name = "%s"
						node_type = "db-dev-s"
						engine = %q
						is_ha_cluster = false
						tags = [ "terraform-test", "scaleway_rdb_user", "minimal" ]
					}

					resource scaleway_rdb_user db_user {
						instance_id = scaleway_rdb_instance.main.id
						name = "bar"
						password = "R34lP4sSw#Rd"
						is_admin = false
					}`, instanceName, latestEngineVersion),
				Check: resource.ComposeTestCheckFunc(
					isUserPresent(tt, "scaleway_rdb_instance.main", "scaleway_rdb_user.db_user"),
					resource.TestCheckResourceAttr("scaleway_rdb_user.db_user", "name", "bar"),
					resource.TestCheckResourceAttr("scaleway_rdb_user.db_user", "is_admin", "false"),
				),
			},
		},
	})
}

func isUserPresent(tt *acctest.TestTools, instance string, user string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		instanceResource, ok := state.RootModule().Resources[instance]
		if !ok {
			return fmt.Errorf("resource not found: %s", instance)
		}

		userResource, ok := state.RootModule().Resources[user]
		if !ok {
			return fmt.Errorf("resource not found: %s", user)
		}

		rdbAPI, _, _, err := rdb.NewAPIWithRegionAndID(tt.Meta, instanceResource.Primary.ID)
		if err != nil {
			return err
		}

		region, instanceID, userName, err := rdb.ResourceUserParseID(userResource.Primary.ID)
		if err != nil {
			return err
		}

		users, err := rdbAPI.ListUsers(&rdbSDK.ListUsersRequest{
			InstanceID: instanceID,
			Region:     region,
			Name:       &userName,
		})
		if err != nil {
			return err
		}

		if len(users.Users) != 1 {
			return errors.New("no user found")
		}

		return nil
	}
}
