package rdbtestfuncs

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	rdbSDK "github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/rdb"
)

func IsInstanceDestroyed(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_rdb_instance" {
				continue
			}

			rdbAPI, region, ID, err := rdb.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}

			_, err = rdbAPI.GetInstance(&rdbSDK.GetInstanceRequest{
				InstanceID: ID,
				Region:     region,
			})

			// If no error resource still exist
			if err == nil {
				return fmt.Errorf("instance (%s) still exists", rs.Primary.ID)
			}

			// Unexpected api error we return it
			if !httperrors.Is404(err) {
				return err
			}
		}

		return nil
	}
}

func GetLatestEngineVersion(tt *acctest.TestTools, engineName string) string {
	api := rdbSDK.NewAPI(tt.Meta.ScwClient())
	engines, err := api.ListDatabaseEngines(&rdbSDK.ListDatabaseEnginesRequest{})
	if err != nil {
		tt.T.Fatalf("Could not get latest engine version: %s", err)
	}

	latestEngineVersion := ""
	for _, engine := range engines.Engines {
		if engine.Name == engineName {
			if len(engine.Versions) > 0 {
				latestEngineVersion = engine.Versions[0].Name
				break
			}
		}
	}
	return latestEngineVersion
}
