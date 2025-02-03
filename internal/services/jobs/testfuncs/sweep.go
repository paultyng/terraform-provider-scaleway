package jobstestfuncs

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	jobsSDK "github.com/scaleway/scaleway-sdk-go/api/jobs/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/logging"
)

func AddTestSweepers() {
	resource.AddTestSweepers("scaleway_job_definition", &resource.Sweeper{
		Name: "scaleway_job_definition",
		F:    testSweepJobDefinition,
	})
}

func testSweepJobDefinition(_ string) error {
	return acctest.SweepRegions((&jobsSDK.API{}).Regions(), func(scwClient *scw.Client, region scw.Region) error {
		jobsAPI := jobsSDK.NewAPI(scwClient)
		logging.L.Debugf("sweeper: destroying the jobs definitions in (%s)", region)
		listJobDefinitions, err := jobsAPI.ListJobDefinitions(
			&jobsSDK.ListJobDefinitionsRequest{
				Region: region,
			}, scw.WithAllPages())
		if err != nil {
			return fmt.Errorf("error listing definition in (%s) in sweeper: %s", region, err)
		}

		for _, definition := range listJobDefinitions.JobDefinitions {
			err := jobsAPI.DeleteJobDefinition(&jobsSDK.DeleteJobDefinitionRequest{
				JobDefinitionID: definition.ID,
				Region:          region,
			})
			if err != nil {
				logging.L.Debugf("sweeper: error (%s)", err)

				return fmt.Errorf("error deleting definition in sweeper: %s", err)
			}
		}

		return nil
	})
}
