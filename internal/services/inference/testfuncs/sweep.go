package inferencetestfuncs

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	inference "github.com/scaleway/scaleway-sdk-go/api/inference/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/logging"
)

func AddTestSweepers() {
	resource.AddTestSweepers("scaleway_instance_deployment", &resource.Sweeper{
		Name:         "scaleway_instance_deployment",
		Dependencies: nil,
		F:            testSweepDeployment,
	})
}

func testSweepDeployment(_ string) error {
	return acctest.SweepRegions((&inference.API{}).Regions(), func(scwClient *scw.Client, region scw.Region) error {
		inferenceAPI := inference.NewAPI(scwClient)
		logging.L.Debugf("sweeper: destroying the inference deployments in (%s)", region)
		listDeployments, err := inferenceAPI.ListDeployments(
			&inference.ListDeploymentsRequest{
				Region: region,
			}, scw.WithAllPages())
		if err != nil {
			return fmt.Errorf("error listing deployment in (%s) in sweeper: %s", region, err)
		}

		for _, deployment := range listDeployments.Deployments {
			_, err := inferenceAPI.DeleteDeployment(&inference.DeleteDeploymentRequest{
				DeploymentID: deployment.ID,
				Region:       region,
			})
			if err != nil {
				logging.L.Debugf("sweeper: error (%s)", err)

				return fmt.Errorf("error deleting deployment in sweeper: %s", err)
			}
		}

		return nil
	})
}
