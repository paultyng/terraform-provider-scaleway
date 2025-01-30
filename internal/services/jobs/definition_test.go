package jobs_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	jobsSDK "github.com/scaleway/scaleway-sdk-go/api/jobs/v1alpha1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/jobs"
)

func TestAccJobDefinition_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckJobDefinitionDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_job_definition main {
						name = "test-jobs-job-definition-basic"
						cpu_limit = 120
						memory_limit = 256
						image_uri = "docker.io/alpine:latest"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobDefinitionExists(tt, "scaleway_job_definition.main"),
					acctest.CheckResourceAttrUUID("scaleway_job_definition.main", "id"),
					resource.TestCheckResourceAttr("scaleway_job_definition.main", "name", "test-jobs-job-definition-basic"),
				),
			},
		},
	})
}

func TestAccJobDefinition_Timeout(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckJobDefinitionDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_job_definition main {
						name = "test-jobs-job-definition-timeout"
						cpu_limit = 120
						memory_limit = 256
						image_uri = "docker.io/alpine:latest"
						timeout = "20m"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobDefinitionExists(tt, "scaleway_job_definition.main"),
					acctest.CheckResourceAttrUUID("scaleway_job_definition.main", "id"),
					resource.TestCheckResourceAttr("scaleway_job_definition.main", "name", "test-jobs-job-definition-timeout"),
					resource.TestCheckResourceAttr("scaleway_job_definition.main", "timeout", "20m0s"),
				),
			},
			{
				Config: `
					resource scaleway_job_definition main {
						name = "test-jobs-job-definition-timeout"
						cpu_limit = 120
						memory_limit = 256
						image_uri = "docker.io/alpine:latest"
						timeout = "1h30m"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobDefinitionExists(tt, "scaleway_job_definition.main"),
					acctest.CheckResourceAttrUUID("scaleway_job_definition.main", "id"),
					resource.TestCheckResourceAttr("scaleway_job_definition.main", "name", "test-jobs-job-definition-timeout"),
					resource.TestCheckResourceAttr("scaleway_job_definition.main", "timeout", "1h30m0s"),
				),
			},
		},
	})
}

func TestAccJobDefinition_Cron(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckJobDefinitionDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_job_definition main {
						name = "test-jobs-job-definition-cron"
						cpu_limit = 120
						memory_limit = 256
						image_uri = "docker.io/alpine:latest"
						cron {
							schedule = "5 4 1 * *"
							timezone = "Europe/Paris"
						}
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobDefinitionExists(tt, "scaleway_job_definition.main"),
					acctest.CheckResourceAttrUUID("scaleway_job_definition.main", "id"),
					resource.TestCheckResourceAttr("scaleway_job_definition.main", "name", "test-jobs-job-definition-cron"),
					resource.TestCheckResourceAttr("scaleway_job_definition.main", "cron.#", "1"),
					resource.TestCheckResourceAttr("scaleway_job_definition.main", "cron.0.schedule", "5 4 1 * *"),
					resource.TestCheckResourceAttr("scaleway_job_definition.main", "cron.0.timezone", "Europe/Paris"),
				),
			},
			{
				Config: `
					resource scaleway_job_definition main {
						name = "test-jobs-job-definition-cron"
						cpu_limit = 120
						memory_limit = 256
						image_uri = "docker.io/alpine:latest"
						cron {
							schedule = "5 5 * * *"
							timezone = "America/Jamaica"
						}
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobDefinitionExists(tt, "scaleway_job_definition.main"),
					acctest.CheckResourceAttrUUID("scaleway_job_definition.main", "id"),
					resource.TestCheckResourceAttr("scaleway_job_definition.main", "name", "test-jobs-job-definition-cron"),
					resource.TestCheckResourceAttr("scaleway_job_definition.main", "cron.#", "1"),
					resource.TestCheckResourceAttr("scaleway_job_definition.main", "cron.0.schedule", "5 5 * * *"),
					resource.TestCheckResourceAttr("scaleway_job_definition.main", "cron.0.timezone", "America/Jamaica"),
				),
			},
			{
				Config: `
					resource scaleway_job_definition main {
						name = "test-jobs-job-definition-cron"
						cpu_limit = 120
						memory_limit = 256
						image_uri = "docker.io/alpine:latest"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobDefinitionExists(tt, "scaleway_job_definition.main"),
					acctest.CheckResourceAttrUUID("scaleway_job_definition.main", "id"),
					resource.TestCheckResourceAttr("scaleway_job_definition.main", "name", "test-jobs-job-definition-cron"),
					resource.TestCheckResourceAttr("scaleway_job_definition.main", "cron.#", "0"),
				),
			},
		},
	})
}

func testAccCheckJobDefinitionExists(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		api, region, id, err := jobs.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = api.GetJobDefinition(&jobsSDK.GetJobDefinitionRequest{
			JobDefinitionID: id,
			Region:          region,
		})
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckJobDefinitionDestroy(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_job_definition" {
				continue
			}

			api, region, id, err := jobs.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}

			err = api.DeleteJobDefinition(&jobsSDK.DeleteJobDefinitionRequest{
				JobDefinitionID: id,
				Region:          region,
			})

			if err == nil {
				return fmt.Errorf("jobs jobdefinition (%s) still exists", rs.Primary.ID)
			}

			if !httperrors.Is404(err) {
				return err
			}
		}

		return nil
	}
}
