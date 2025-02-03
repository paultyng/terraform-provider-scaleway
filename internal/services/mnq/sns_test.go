package mnq_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	mnqSDK "github.com/scaleway/scaleway-sdk-go/api/mnq/v1beta1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/mnq"
)

func TestAccSNS_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isSNSDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_account_project main {
						name = "tf_tests_mnq_sns_basic"
					}

					resource scaleway_mnq_sns main {
						project_id = scaleway_account_project.main.id
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isSNSPresent(tt, "scaleway_mnq_sns.main"),
					acctest.CheckResourceAttrUUID("scaleway_mnq_sns.main", "id"),
					resource.TestCheckResourceAttrSet("scaleway_mnq_sns.main", "endpoint"),
				),
			},
		},
	})
}

func isSNSPresent(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		api, region, id, err := mnq.NewSNSAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		sns, err := api.GetSnsInfo(&mnqSDK.SnsAPIGetSnsInfoRequest{
			ProjectID: id,
			Region:    region,
		})

		if sns.Status != mnqSDK.SnsInfoStatusEnabled {
			return fmt.Errorf("sns status should be enabled, got: %s", sns.Status)
		}

		if err != nil {
			return err
		}

		return nil
	}
}

func isSNSDestroyed(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_mnq_sns" {
				continue
			}

			api, region, id, err := mnq.NewSNSAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}

			sns, err := api.DeactivateSns(&mnqSDK.SnsAPIDeactivateSnsRequest{
				ProjectID: id,
				Region:    region,
			})
			if err != nil {
				if httperrors.Is404(err) {
					return nil
				}
				return err
			}

			if sns.Status != mnqSDK.SnsInfoStatusDisabled {
				return fmt.Errorf("mnq sns (%s) should be disabled", rs.Primary.ID)
			}

			if !httperrors.Is404(err) {
				return err
			}
		}

		return nil
	}
}
