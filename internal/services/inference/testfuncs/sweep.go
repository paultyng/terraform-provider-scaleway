package inferencetestfuncs

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	inference "github.com/scaleway/scaleway-sdk-go/api/inference/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/api/inference/v1beta1/sweepers"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

func AddTestSweepers() {
	resource.AddTestSweepers("scaleway_instance_deployment", &resource.Sweeper{
		Name:         "scaleway_instance_deployment",
		Dependencies: nil,
		F:            testSweepDeployment,
	})
}

func testSweepDeployment(_ string) error {
	return acctest.SweepRegions((&inference.API{}).Regions(), sweepers.SweepDeployment)
}
