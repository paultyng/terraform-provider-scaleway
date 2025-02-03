package rdbtestfuncs

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	rdbSDK "github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
	"github.com/scaleway/scaleway-sdk-go/api/rdb/v1/sweepers"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

func AddTestSweepers() {
	resource.AddTestSweepers("scaleway_rdb_instance", &resource.Sweeper{
		Name: "scaleway_rdb_instance",
		F:    testSweepInstance,
	})
}

func testSweepInstance(_ string) error {
	return acctest.SweepRegions((&rdbSDK.API{}).Regions(), sweepers.SweepInstance)
}
