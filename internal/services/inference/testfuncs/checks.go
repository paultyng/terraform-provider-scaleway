package inferencetestfuncs

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	inferenceSDK "github.com/scaleway/scaleway-sdk-go/api/inference/v1beta1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/inference"
)

func IsDeploymentDestroyed(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "scaleway_inference_deployment" {
				continue
			}

			inferenceAPI, region, ID, err := inference.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}

			deployment, err := inferenceAPI.GetDeployment(&inferenceSDK.GetDeploymentRequest{
				Region:       region,
				DeploymentID: ID,
			})

			if err == nil {
				return fmt.Errorf("deployment %s (%s) still exists", deployment.Name, deployment.ID)
			}

			if !httperrors.Is404(err) {
				return err
			}
		}
		return nil
	}
}
