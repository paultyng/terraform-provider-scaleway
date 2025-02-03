package applesilicontestfuncs

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	applesiliconSDK "github.com/scaleway/scaleway-sdk-go/api/applesilicon/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1/sweepers"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

func AddTestSweepers() {
	resource.AddTestSweepers("scaleway_apple_silicon_instance", &resource.Sweeper{
		Name: "scaleway_apple_silicon",
		F:    testSweepAppleSiliconServer,
	})
}

func testSweepAppleSiliconServer(_ string) error {
	return acctest.SweepZones((&applesiliconSDK.API{}).Zones(), sweepers.SweepServers)
}
