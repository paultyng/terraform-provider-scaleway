package containertestfuncs

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	containerSDK "github.com/scaleway/scaleway-sdk-go/api/container/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/api/container/v1beta1/sweepers"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

func AddTestSweepers() {
	resource.AddTestSweepers("scaleway_container_namespace", &resource.Sweeper{
		Name:         "scaleway_container_namespace",
		F:            testSweepNamespace,
		Dependencies: []string{"scaleway_container"},
	})
	resource.AddTestSweepers("scaleway_container", &resource.Sweeper{
		Name: "scaleway_container",
		F:    testSweepContainer,
	})
	resource.AddTestSweepers("scaleway_container_trigger", &resource.Sweeper{
		Name: "scaleway_container_trigger",
		F:    testSweepTrigger,
	})
}

func testSweepTrigger(_ string) error {
	return acctest.SweepRegions((&containerSDK.API{}).Regions(), sweepers.SweepTrigger)
}

func testSweepContainer(_ string) error {
	return acctest.SweepRegions((&containerSDK.API{}).Regions(), sweepers.SweepContainer)
}

func testSweepNamespace(_ string) error {
	return acctest.SweepRegions((&containerSDK.API{}).Regions(), sweepers.SweepNamespace)
}
