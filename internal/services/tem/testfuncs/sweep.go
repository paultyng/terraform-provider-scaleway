package temtestfuncs

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	temSDK "github.com/scaleway/scaleway-sdk-go/api/tem/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/api/tem/v1alpha1/sweepers"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

const skippedDomain = "test.scaleway-terraform.com"

func AddTestSweepers() {
	resource.AddTestSweepers("scaleway_tem_domain", &resource.Sweeper{
		Name: "scaleway_tem_domain",
		F:    testSweepDomain,
	})
}

func testSweepDomain(_ string) error {
	return acctest.SweepRegions((&temSDK.API{}).Regions(), func(scwClient *scw.Client, region scw.Region) error {
		return sweepers.SweepDomain(scwClient, region, skippedDomain)
	})
}
