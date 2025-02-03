package secrettestfuncs

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	secretSDK "github.com/scaleway/scaleway-sdk-go/api/secret/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/api/secret/v1beta1/sweepers"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

func AddTestSweepers() {
	resource.AddTestSweepers("scaleway_secret", &resource.Sweeper{
		Name: "scaleway_secret",
		F:    testSweepSecret,
	})
}

func testSweepSecret(_ string) error {
	return acctest.SweepRegions((&secretSDK.API{}).Regions(), sweepers.SweepSecret)
}
