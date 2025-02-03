package mnq_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	mnqSDK "github.com/scaleway/scaleway-sdk-go/api/mnq/v1beta1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/mnq"
)

func TestAccSQS_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isSQSDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_account_project main {
						name = "tf_tests_mnq_sqs_basic"
					}

					resource scaleway_mnq_sqs main {
						project_id = scaleway_account_project.main.id
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isSQSPresent(tt, "scaleway_mnq_sqs.main"),
					acctest.CheckResourceAttrUUID("scaleway_mnq_sqs.main", "id"),
					resource.TestCheckResourceAttrSet("scaleway_mnq_sqs.main", "endpoint"),
				),
			},
		},
	})
}

func TestAccSQS_AlreadyActivated(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isSQSDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_account_project main {
						name = "tf_tests_mnq_sqs_already_activated"
					}

					resource scaleway_mnq_sqs main {
						project_id = scaleway_account_project.main.id
					}
				`,
			},
			{
				Config: `
					resource scaleway_account_project main {
						name = "tf_tests_mnq_sqs_already_activated"
					}

					resource scaleway_mnq_sqs main {
						project_id = scaleway_account_project.main.id
					}

					resource scaleway_mnq_sqs duplicated {
						project_id = scaleway_account_project.main.id
					}
				`,
				ExpectError: regexp.MustCompile(".*Conflict.*"),
			},
		},
	})
}

func isSQSPresent(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		api, region, id, err := mnq.NewSQSAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		sqs, err := api.GetSqsInfo(&mnqSDK.SqsAPIGetSqsInfoRequest{
			ProjectID: id,
			Region:    region,
		})
		if err != nil {
			return err
		}

		if sqs.Status != mnqSDK.SqsInfoStatusEnabled {
			return fmt.Errorf("sqs status should be enabled, got: %s", sqs.Status)
		}

		return nil
	}
}

func isSQSDestroyed(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_mnq_sqs" {
				continue
			}

			api, region, id, err := mnq.NewSQSAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}

			sqs, err := api.DeactivateSqs(&mnqSDK.SqsAPIDeactivateSqsRequest{
				ProjectID: id,
				Region:    region,
			})
			if err != nil {
				if httperrors.Is404(err) { // Project may have been deleted
					return nil
				}
				return err
			}

			if sqs.Status != mnqSDK.SqsInfoStatusDisabled {
				return fmt.Errorf("mnq sqs (%s) should be disabled", rs.Primary.ID)
			}

			if !httperrors.Is404(err) {
				return err
			}
		}

		return nil
	}
}
