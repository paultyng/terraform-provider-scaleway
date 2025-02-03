package blocktestfuncs

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	blockSDK "github.com/scaleway/scaleway-sdk-go/api/block/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/api/block/v1alpha1/sweepers"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

func AddTestSweepers() {
	resource.AddTestSweepers("scaleway_block_snapshot", &resource.Sweeper{
		Name: "scaleway_block_snapshot",
		F:    testSweepSnapshot,
	})
	resource.AddTestSweepers("scaleway_block_volume", &resource.Sweeper{
		Name: "scaleway_block_volume",
		F:    testSweepBlockVolume,
	})
}

func testSweepBlockVolume(_ string) error {
	return acctest.SweepZones((&blockSDK.API{}).Zones(), sweepers.SweepVolumes)
}

func testSweepSnapshot(_ string) error {
	return acctest.SweepZones((&blockSDK.API{}).Zones(), sweepers.SweepSnapshots)
}
