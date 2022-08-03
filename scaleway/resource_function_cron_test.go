package scaleway

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	function "github.com/scaleway/scaleway-sdk-go/api/function/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func init() {
	resource.AddTestSweepers("scaleway_function_cron", &resource.Sweeper{
		Name: "scaleway_function_cron",
		F:    testSweepFunctionCron,
	})
}

func testSweepFunctionCron(_ string) error {
	return sweepRegions([]scw.Region{scw.RegionFrPar}, func(scwClient *scw.Client, region scw.Region) error {
		functionAPI := function.NewAPI(scwClient)
		l.Debugf("sweeper: destroying the function cron in (%s)", region)
		listCron, err := functionAPI.ListCrons(
			&function.ListCronsRequest{
				Region: region,
			}, scw.WithAllPages())
		if err != nil {
			return fmt.Errorf("error listing cron in (%s) in sweeper: %s", region, err)
		}

		for _, cron := range listCron.Crons {
			_, err := functionAPI.DeleteCron(&function.DeleteCronRequest{
				CronID: cron.ID,
				Region: region,
			})
			if err != nil {
				l.Debugf("sweeper: error (%s)", err)

				return fmt.Errorf("error deleting cron in sweeper: %s", err)
			}
		}

		return nil
	})
}

func TestAccScalewayFunctionCron_Basic(t *testing.T) {
	tt := NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckScalewayFunctionCronDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_function_namespace main {
						name = "test-cron"
					}

					resource scaleway_function main {
						name = "test-cron"
						namespace_id = scaleway_function_namespace.main.id
						runtime = "node14"
						privacy = "private"
						handler = "handler.handle"
					}

					resource scaleway_function_cron main {
						function_id = scaleway_function.main.id
						schedule = "0 0 * * *"
						args = jsonencode({})
					}

					resource scaleway_function_cron func {
						function_id = scaleway_function.main.id
						schedule = "0 1 * * *"
						args = jsonencode({})
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayFunctionCronExists(tt, "scaleway_function_cron.main"),
					testAccCheckScalewayFunctionCronExists(tt, "scaleway_function_cron.func"),
					resource.TestCheckResourceAttr("scaleway_function_cron.main", "schedule", "0 0 * * *"),
					resource.TestCheckResourceAttr("scaleway_function_cron.func", "schedule", "0 1 * * *"),
				),
			},
			{
				Config: `
					resource scaleway_function_namespace main {
						name = "test-cron"
					}

					resource scaleway_function main {
						name = "test-cron"
						namespace_id = scaleway_function_namespace.main.id
						runtime = "node14"
						privacy = "private"
						handler = "handler.handle"
					}

					resource scaleway_function_cron main {
						function_id = scaleway_function.main.id
						schedule = "1 1 * * *"
						args = jsonencode({})
					}

					resource scaleway_function_cron func {
						function_id = scaleway_function.main.id
						schedule = "1 2 * * *"
						args = jsonencode({})
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayFunctionCronExists(tt, "scaleway_function_cron.main"),
					testAccCheckScalewayFunctionCronExists(tt, "scaleway_function_cron.func"),
					resource.TestCheckResourceAttr("scaleway_function_cron.main", "schedule", "1 1 * * *"),
					resource.TestCheckResourceAttr("scaleway_function_cron.func", "schedule", "1 2 * * *"),
				),
			},
		},
	})
}

func TestAccScalewayFunctionCron_WithArgs(t *testing.T) {
	tt := NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckScalewayFunctionCronDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_function_namespace main {
						name = "test-cron"
					}

					resource scaleway_function main {
						name = "test-cron"
						namespace_id = scaleway_function_namespace.main.id
						runtime = "node14"
						privacy = "private"
						handler = "handler.handle"
					}

					resource scaleway_function_cron main {
						function_id = scaleway_function.main.id
						schedule = "0 0 * * *"
						args = jsonencode({test = "scw"})
					}

					resource scaleway_function_cron func {
						function_id = scaleway_function.main.id
						schedule = "0 1 * * *"
						args = jsonencode({my_var = "terraform"})
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayFunctionCronExists(tt, "scaleway_function_cron.main"),
					testAccCheckScalewayFunctionCronExists(tt, "scaleway_function_cron.func"),
					resource.TestCheckResourceAttr("scaleway_function_cron.main", "schedule", "0 0 * * *"),
					resource.TestCheckResourceAttr("scaleway_function_cron.main", "args", "{\"test\":\"scw\"}"),
					resource.TestCheckResourceAttr("scaleway_function_cron.func", "schedule", "0 1 * * *"),
					resource.TestCheckResourceAttr("scaleway_function_cron.func", "args", "{\"my_var\":\"terraform\"}"),
				),
			},
			{
				Config: `
					resource scaleway_function_namespace main {
						name = "test-cron"
					}

					resource scaleway_function main {
						name = "test-cron"
						namespace_id = scaleway_function_namespace.main.id
						runtime = "node14"
						privacy = "private"
						handler = "handler.handle"
					}

					resource scaleway_function_cron main {
						function_id = scaleway_function.main.id
						schedule = "1 1 * * *"
						args = jsonencode(
						{
							address   = {
								city    = "Paris"
								country = "FR"
							}
							age       = 23
							firstName = "John"
							isAlive   = true
							lastName  = "Smith"
						}
                		)
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayFunctionCronExists(tt, "scaleway_function_cron.main"),
					resource.TestCheckResourceAttr("scaleway_function_cron.main", "schedule", "1 1 * * *"),
					resource.TestCheckResourceAttr("scaleway_function_cron.main", "args", "{\"address\":{\"city\":\"Paris\",\"country\":\"FR\"},\"age\":23,\"firstName\":\"John\",\"isAlive\":true,\"lastName\":\"Smith\"}"),
				),
			},
		},
	})
}

func testAccCheckScalewayFunctionCronExists(tt *TestTools, n string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		api, region, id, err := functionAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return nil
		}

		_, err = api.GetCron(&function.GetCronRequest{
			CronID: id,
			Region: region,
		})

		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckScalewayFunctionCronDestroy(tt *TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_function_namespace" {
				continue
			}

			api, region, id, err := functionAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}

			_, err = api.DeleteNamespace(&function.DeleteNamespaceRequest{
				NamespaceID: id,
				Region:      region,
			})

			if err == nil {
				return fmt.Errorf("function namespace (%s) still exists", rs.Primary.ID)
			}

			if !is404Error(err) {
				return err
			}
		}

		return nil
	}
}
