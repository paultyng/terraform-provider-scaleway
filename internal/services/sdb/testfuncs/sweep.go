package sdbtestfuncs

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	sdbSDK "github.com/scaleway/scaleway-sdk-go/api/serverless_sqldb/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/api/serverless_sqldb/v1alpha1/sweepers"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

func AddTestSweepers() {
	resource.AddTestSweepers("scaleway_sdb_sql_database", &resource.Sweeper{
		Name: "scaleway_sdb_sql_database",
		F:    testSweepServerlessSQLDBDatabase,
	})
}

func testSweepServerlessSQLDBDatabase(_ string) error {
	return acctest.SweepRegions((&sdbSDK.API{}).Regions(), sweepers.SweepDatabase)
}
