package sdbtestfuncs

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	sdbSDK "github.com/scaleway/scaleway-sdk-go/api/serverless_sqldb/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/logging"
)

func AddTestSweepers() {
	resource.AddTestSweepers("scaleway_sdb_sql_database", &resource.Sweeper{
		Name: "scaleway_sdb_sql_database",
		F:    testSweepServerlessSQLDBDatabase,
	})
}

func testSweepServerlessSQLDBDatabase(_ string) error {
	return acctest.SweepRegions((&sdbSDK.API{}).Regions(), func(scwClient *scw.Client, region scw.Region) error {
		sdbAPI := sdbSDK.NewAPI(scwClient)
		logging.L.Debugf("sweeper: destroying the serverless sql database in (%s)", region)
		listServerlessSQLDBDatabases, err := sdbAPI.ListDatabases(
			&sdbSDK.ListDatabasesRequest{
				Region: region,
			}, scw.WithAllPages())
		if err != nil {
			return fmt.Errorf("error listing database in (%s) in sweeper: %s", region, err)
		}

		for _, database := range listServerlessSQLDBDatabases.Databases {
			_, err := sdbAPI.DeleteDatabase(&sdbSDK.DeleteDatabaseRequest{
				DatabaseID: database.ID,
				Region:     region,
			})
			if err != nil {
				logging.L.Debugf("sweeper: error (%s)", err)

				return fmt.Errorf("error deleting database in sweeper: %s", err)
			}
		}

		return nil
	})
}
