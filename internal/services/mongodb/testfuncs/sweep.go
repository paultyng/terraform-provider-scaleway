package mongodbtestfuncs

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	mongodb "github.com/scaleway/scaleway-sdk-go/api/mongodb/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/api/mongodb/v1alpha1/sweepers"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

func AddTestSweepers() {
	resource.AddTestSweepers("scaleway_mongodb_instance", &resource.Sweeper{
		Name: "scaleway_mongodb_instance",
		F:    testSweepMongodbInstance,
	})
}

func testSweepMongodbInstance(_ string) error {
	return acctest.SweepRegions((&mongodb.API{}).Regions(), sweepers.SweepInstances)
}
