package container_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	containerSDK "github.com/scaleway/scaleway-sdk-go/api/container/v1beta1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/container"
)

func TestAccCron_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isCronDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_container_namespace main {
					}

					resource scaleway_container main {
						name = "my-container-with-cron-tf"
						namespace_id = scaleway_container_namespace.main.id
					}

					resource scaleway_container_cron main {
						name = "tf-tests-container-cron-basic"
						container_id = scaleway_container.main.id
						schedule = "5 4 * * *" #cron at 04:05
						args = jsonencode({test = "scw"})
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isCronPresent(tt, "scaleway_container_cron.main"),
					resource.TestCheckResourceAttr("scaleway_container_cron.main", "schedule", "5 4 * * *"),
					resource.TestCheckResourceAttr("scaleway_container_cron.main", "args", "{\"test\":\"scw\"}"),
					resource.TestCheckResourceAttr("scaleway_container_cron.main", "name", "tf-tests-container-cron-basic"),
				),
			},
			{
				Config: `
					resource scaleway_container_namespace main {
					}

					resource scaleway_container main {
						name = "my-container-with-cron-tf"
						namespace_id = scaleway_container_namespace.main.id
					}

					resource scaleway_container_cron main {
						name = "tf-tests-container-cron-basic-changed"
						container_id = scaleway_container.main.id
						schedule = "5 4 * * *" #cron at 04:05
						args = jsonencode({test = "scw"})
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isCronPresent(tt, "scaleway_container_cron.main"),
					resource.TestCheckResourceAttr("scaleway_container_cron.main", "schedule", "5 4 * * *"),
					resource.TestCheckResourceAttr("scaleway_container_cron.main", "args", "{\"test\":\"scw\"}"),
					resource.TestCheckResourceAttr("scaleway_container_cron.main", "name", "tf-tests-container-cron-basic-changed"),
				),
			},
		},
	})
}

func TestAccCron_WithMultiArgs(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isCronDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_container_namespace main {
					}

					resource scaleway_container main {
						name = "my-container-with-cron-tf"
						namespace_id = scaleway_container_namespace.main.id
					}

					resource scaleway_container_cron main {
						container_id = scaleway_container.main.id
						schedule = "5 4 1 * *" #cron at 04:05 on day-of-month 1
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
					isCronPresent(tt, "scaleway_container_cron.main"),
					resource.TestCheckResourceAttr("scaleway_container_cron.main", "schedule", "5 4 1 * *"),
					resource.TestCheckResourceAttr("scaleway_container_cron.main", "args", "{\"address\":{\"city\":\"Paris\",\"country\":\"FR\"},\"age\":23,\"firstName\":\"John\",\"isAlive\":true,\"lastName\":\"Smith\"}"),
				),
			},
			{
				Config: `
					resource scaleway_container_namespace main {
					}

					resource scaleway_container main {
					name = "my-container-with-cron-tf"
						namespace_id = scaleway_container_namespace.main.id
					}

					resource scaleway_container_cron main {
						container_id = scaleway_container.main.id
						schedule = "5 4 * * 1" #cron at 04:05
						args = jsonencode({test = "scw"})
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaleway_container_cron.main", "schedule", "5 4 * * 1"),
					resource.TestCheckResourceAttr("scaleway_container_cron.main", "args", "{\"test\":\"scw\"}"),
				),
			},
		},
	})
}

func isCronPresent(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource container cron not found: %s", n)
		}

		api, region, id, err := container.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = api.GetCron(&containerSDK.GetCronRequest{
			CronID: id,
			Region: region,
		})
		if err != nil {
			return err
		}

		return nil
	}
}

func isCronDestroyed(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_container_cron" {
				continue
			}

			api, region, id, err := container.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}

			_, err = api.DeleteCron(&containerSDK.DeleteCronRequest{
				CronID: id,
				Region: region,
			})

			if err == nil {
				return fmt.Errorf("containerSDK cron (%s) still exists", rs.Primary.ID)
			}

			if !httperrors.Is404(err) {
				return err
			}
		}

		return nil
	}
}
