package registrytestfuncs

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	registrySDK "github.com/scaleway/scaleway-sdk-go/api/registry/v1"
	"github.com/scaleway/scaleway-sdk-go/api/registry/v1/sweepers"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

func AddTestSweepers() {
	resource.AddTestSweepers("scaleway_registry_namespace", &resource.Sweeper{
		Name: "scaleway_registry_namespace",
		F:    testSweepNamespace,
	})
}

func testSweepNamespace(_ string) error {
	return acctest.SweepRegions((&registrySDK.API{}).Regions(), sweepers.SweepNamespace)
}
