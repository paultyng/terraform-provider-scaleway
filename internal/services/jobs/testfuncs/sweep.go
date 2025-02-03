package jobstestfuncs

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	jobsSDK "github.com/scaleway/scaleway-sdk-go/api/jobs/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/api/jobs/v1alpha1/sweepers"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

func AddTestSweepers() {
	resource.AddTestSweepers("scaleway_job_definition", &resource.Sweeper{
		Name: "scaleway_job_definition",
		F:    testSweepJobDefinition,
	})
}

func testSweepJobDefinition(_ string) error {
	return acctest.SweepRegions((&jobsSDK.API{}).Regions(), sweepers.SweepJobDefinition)
}
