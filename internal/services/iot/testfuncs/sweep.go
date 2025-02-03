package iottestfuncs

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	iotSDK "github.com/scaleway/scaleway-sdk-go/api/iot/v1"
	"github.com/scaleway/scaleway-sdk-go/api/iot/v1/sweepers"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

func AddTestSweepers() {
	resource.AddTestSweepers("scaleway_iot_hub", &resource.Sweeper{
		Name: "scaleway_iot_hub",
		F:    testSweepHub,
	})
}

func testSweepHub(_ string) error {
	return acctest.SweepRegions((&iotSDK.API{}).Regions(), sweepers.SweepHub)
}
