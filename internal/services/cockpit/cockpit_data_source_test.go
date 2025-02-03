package cockpit_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

func TestAccDataSourceCockpit_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_account_project" "project" {
						name = "tf_tests_cockpit_project_premium"
					}
					
					resource "scaleway_cockpit_source" "metrics" {
						project_id = scaleway_account_project.project.id
						name       = "my-data-source-metrics"
						type       = "metrics"
						retention_days = 31
					}
					
					resource "scaleway_cockpit_source" "logs" {
						project_id = scaleway_account_project.project.id
						name       = "my-data-source-logs"
						type       = "logs"
						retention_days = 7
					}
					
					resource "scaleway_cockpit_source" "traces" {
						project_id = scaleway_account_project.project.id
						name       = "my-data-source-traces"
						type       = "traces"
						retention_days = 7
					}
					
					resource "scaleway_cockpit_alert_manager" "alert_manager" {
						project_id = scaleway_account_project.project.id
						enable_managed_alerts = true
					}

					
					resource "scaleway_cockpit" "main" {
						project_id = scaleway_account_project.project.id
						plan       = "free"
						depends_on = [
								scaleway_cockpit_source.metrics,
								scaleway_cockpit_source.logs,
								scaleway_cockpit_source.traces,
								scaleway_cockpit_alert_manager.alert_manager,
							]
					}

					data "scaleway_cockpit" "selected" {
						project_id = scaleway_cockpit.main.project_id
					}

					
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaleway_cockpit.main", "plan", "free"),
					resource.TestCheckResourceAttr("scaleway_cockpit.main", "plan_id", "free"),
					resource.TestCheckResourceAttrSet("scaleway_cockpit.main", "endpoints.0.metrics_url"),
					resource.TestCheckResourceAttrSet("scaleway_cockpit.main", "endpoints.0.logs_url"),
					resource.TestCheckResourceAttrSet("scaleway_cockpit.main", "endpoints.0.alertmanager_url"),
					resource.TestCheckResourceAttrSet("scaleway_cockpit.main", "endpoints.0.traces_url"),
					resource.TestCheckResourceAttrSet("scaleway_cockpit.main", "push_url.0.push_logs_url"),
					resource.TestCheckResourceAttrSet("scaleway_cockpit.main", "push_url.0.push_metrics_url"),
				),
			},
		},
	})
}
