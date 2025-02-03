package function_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	functionSDK "github.com/scaleway/scaleway-sdk-go/api/function/v1beta1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/function"
)

func TestAccFunctionCron_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckFunctionCronDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_function_namespace main {
						name = "tf-tests-function-cron-basic"
					}

					resource scaleway_function main {
						name = "tf-tests-cron-basic"
						namespace_id = scaleway_function_namespace.main.id
						runtime = "node20"
						privacy = "private"
						handler = "handler.handle"
					}

					resource scaleway_function_cron main {
						name = "tf-tests-cron-basic"
						function_id = scaleway_function.main.id
						schedule = "0 0 * * *"
						args = jsonencode({})
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFunctionCronExists(tt, "scaleway_function_cron.main"),
					resource.TestCheckResourceAttr("scaleway_function_cron.main", "schedule", "0 0 * * *"),
					resource.TestCheckResourceAttr("scaleway_function_cron.main", "name", "tf-tests-cron-basic"),
				),
			},
		},
	})
}

func TestAccFunctionCron_NameUpdate(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckFunctionCronDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_function_namespace main {
						name = "tf-tests-function-cron-name-update"
					}

					resource scaleway_function main {
						name = "tf-tests-function-cron-name-update"
						namespace_id = scaleway_function_namespace.main.id
						runtime = "node20"
						privacy = "private"
						handler = "handler.handle"
					}

					resource scaleway_function_cron main {
						name = "tf-tests-function-cron-name-update"
						function_id = scaleway_function.main.id
						schedule = "0 0 * * *"
						args = jsonencode({})
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFunctionCronExists(tt, "scaleway_function_cron.main"),
					resource.TestCheckResourceAttr("scaleway_function_cron.main", "schedule", "0 0 * * *"),
					resource.TestCheckResourceAttr("scaleway_function_cron.main", "name", "tf-tests-function-cron-name-update"),
				),
			},
			{
				Config: `
					resource scaleway_function_namespace main {
						name = "tf-tests-function-cron-name-update"
					}

					resource scaleway_function main {
						name = "tf-tests-function-cron-name-update"
						namespace_id = scaleway_function_namespace.main.id
						runtime = "node20"
						privacy = "private"
						handler = "handler.handle"
					}

					resource scaleway_function_cron main {
						name = "name-changed"
						function_id = scaleway_function.main.id
						schedule = "0 0 * * *"
						args = jsonencode({test = "scw"})
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFunctionCronExists(tt, "scaleway_function_cron.main"),
					resource.TestCheckResourceAttr("scaleway_function_cron.main", "name", "name-changed"),
				),
			},
		},
	})
}

func TestAccFunctionCron_WithArgs(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckFunctionCronDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_function_namespace main {
						name = "tf-tests-function-cron-with-args"
					}

					resource scaleway_function main {
						name = "tf-tests-cron-with-args"
						namespace_id = scaleway_function_namespace.main.id
						runtime = "node20"
						privacy = "private"
						handler = "handler.handle"
					}

					resource scaleway_function_cron main {
						name = "tf-tests-cron-with-args"
						function_id = scaleway_function.main.id
						schedule = "0 0 * * *"
						args = jsonencode({test = "scw"})
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFunctionCronExists(tt, "scaleway_function_cron.main"),
					resource.TestCheckResourceAttr("scaleway_function_cron.main", "schedule", "0 0 * * *"),
					resource.TestCheckResourceAttr("scaleway_function_cron.main", "args", "{\"test\":\"scw\"}"),
					resource.TestCheckResourceAttr("scaleway_function_cron.main", "name", "tf-tests-cron-with-args"),
				),
			},
		},
	})
}

func testAccCheckFunctionCronExists(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		api, region, id, err := function.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = api.GetCron(&functionSDK.GetCronRequest{
			CronID: id,
			Region: region,
		})
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckFunctionCronDestroy(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_function_cron" {
				continue
			}

			api, region, id, err := function.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}

			_, err = api.DeleteCron(&functionSDK.DeleteCronRequest{
				CronID: id,
				Region: region,
			})

			if err == nil {
				return fmt.Errorf("function cron (%s) still exists", rs.Primary.ID)
			}

			if !httperrors.Is404(err) {
				return err
			}
		}

		return nil
	}
}
