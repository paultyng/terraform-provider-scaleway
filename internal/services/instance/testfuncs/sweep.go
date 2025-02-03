package instancetestfuncs

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	instanceSDK "github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1/sweepers"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

func AddTestSweepers() {
	resource.AddTestSweepers("scaleway_instance_image", &resource.Sweeper{
		Name:         "scaleway_instance_image",
		Dependencies: []string{"scaleway_instance_server"},
		F:            testSweepImage,
	})
	resource.AddTestSweepers("scaleway_instance_ip", &resource.Sweeper{
		Name: "scaleway_instance_ip",
		F:    testSweepIP,
	})
	resource.AddTestSweepers("scaleway_instance_placement_group", &resource.Sweeper{
		Name: "scaleway_instance_placement_group",
		F:    testSweepPlacementGroup,
	})
	resource.AddTestSweepers("scaleway_instance_security_group", &resource.Sweeper{
		Name: "scaleway_instance_security_group",
		F:    testSweepSecurityGroup,
	})
	resource.AddTestSweepers("scaleway_instance_server", &resource.Sweeper{
		Name: "scaleway_instance_server",
		F:    testSweepServer,
	})
	resource.AddTestSweepers("scaleway_instance_snapshot", &resource.Sweeper{
		Name: "scaleway_instance_snapshot",
		F:    testSweepSnapshot,
	})
	resource.AddTestSweepers("scaleway_instance_volume", &resource.Sweeper{
		Name: "scaleway_instance_volume",
		F:    testSweepVolume,
	})
}

func testSweepVolume(_ string) error {
	return acctest.SweepZones((&instanceSDK.API{}).Zones(), sweepers.SweepVolumes)
}

func testSweepSnapshot(_ string) error {
	return acctest.SweepZones((&instanceSDK.API{}).Zones(), sweepers.SweepSnapshots)
}

func testSweepServer(_ string) error {
	return acctest.SweepZones((&instanceSDK.API{}).Zones(), sweepers.SweepServers)
}

func testSweepSecurityGroup(_ string) error {
	return acctest.SweepZones((&instanceSDK.API{}).Zones(), sweepers.SweepSecurityGroups)
}

func testSweepPlacementGroup(_ string) error {
	return acctest.SweepZones((&instanceSDK.API{}).Zones(), sweepers.SweepPlacementGroup)
}

func testSweepIP(_ string) error {
	return acctest.SweepZones((&instanceSDK.API{}).Zones(), sweepers.SweepIP)
}

func testSweepImage(_ string) error {
	return acctest.SweepZones((&instanceSDK.API{}).Zones(), sweepers.SweepImages)
}
