package webhostingtestfuncs

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	webhostingSDK "github.com/scaleway/scaleway-sdk-go/api/webhosting/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/api/webhosting/v1alpha1/sweepers"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

func AddTestSweepers() {
	resource.AddTestSweepers("scaleway_webhosting", &resource.Sweeper{
		Name: "scaleway_webhosting",
		F:    testSweepWebhosting,
	})
}

func testSweepWebhosting(_ string) error {
	return acctest.SweepRegions((&webhostingSDK.API{}).Regions(), sweepers.SweepWebHosting)
}
