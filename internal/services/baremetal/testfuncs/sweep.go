package baremetaltestfuncs

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	baremetalSDK "github.com/scaleway/scaleway-sdk-go/api/baremetal/v1"
	"github.com/scaleway/scaleway-sdk-go/api/baremetal/v1/sweepers"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

func AddTestSweepers() {
	resource.AddTestSweepers("scaleway_baremetal_server", &resource.Sweeper{
		Name: "scaleway_baremetal_server",
		F:    testSweepServer,
	})
}

func testSweepServer(_ string) error {
	return acctest.SweepZones((&baremetalSDK.API{}).Zones(), sweepers.SweepServers)
}
