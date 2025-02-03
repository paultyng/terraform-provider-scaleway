package cockpittestfuncs

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/scaleway/scaleway-sdk-go/api/cockpit/v1/sweepers"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

func AddTestSweepers() {
	resource.AddTestSweepers("scaleway_cockpit_grafana_user", &resource.Sweeper{
		Name: "scaleway_cockpit_grafana_user",
		F:    testSweepCockpitGrafanaUser,
	})
	resource.AddTestSweepers("scaleway_cockpit_token", &resource.Sweeper{
		Name: "scaleway_cockpit_token",
		F:    testSweepCockpitToken,
	})
	resource.AddTestSweepers("scaleway_cockpit_source", &resource.Sweeper{
		Name: "scaleway_cockpit_source",
		F:    testSweepCockpitSource,
	})
}

func testSweepCockpitToken(_ string) error {
	return acctest.Sweep(sweepers.SweepToken)
}

func testSweepCockpitGrafanaUser(_ string) error {
	return acctest.Sweep(sweepers.SweepGrafanaUser)
}

func testSweepCockpitSource(_ string) error {
	return acctest.SweepRegions(scw.AllRegions, sweepers.SweepSource)
}
