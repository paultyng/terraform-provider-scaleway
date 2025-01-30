package mnq_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

func TestAccDataSourceSQS_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			isSQSDestroyed(tt),
		),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_account_project main {
						name = "tf_tests_ds_mnq_sqs_basic"
					}

					resource scaleway_mnq_sqs main {
						project_id = scaleway_account_project.main.id
					}

					data scaleway_mnq_sqs main {
						project_id = scaleway_mnq_sqs.main.project_id
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isSQSPresent(tt, "scaleway_mnq_sqs.main"),

					resource.TestCheckResourceAttrPair("scaleway_mnq_sqs.main", "id", "data.scaleway_mnq_sqs.main", "id"),
				),
			},
		},
	})
}
