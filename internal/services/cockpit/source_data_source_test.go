package cockpit_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

func TestAccCockpitSource_DataSource_ByID(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isSourceDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_account_project" "project" {
						name = "tf_tests_cockpit_datasource_by_id"
				  	}

					resource "scaleway_cockpit_source" "main" {
					  project_id = scaleway_account_project.project.id
					  name       = "source-by-id"
					  type       = "metrics"
					  retention_days = 30
					}

					data "scaleway_cockpit_source" "by_id" {
					  id = scaleway_cockpit_source.main.id
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.scaleway_cockpit_source.by_id", "id", "scaleway_cockpit_source.main", "id"),
					resource.TestCheckResourceAttr("data.scaleway_cockpit_source.by_id", "name", "source-by-id"),
					resource.TestCheckResourceAttr("data.scaleway_cockpit_source.by_id", "type", "metrics"),
				),
			},
		},
	})
}

func TestAccCockpitSource_DataSource_ByName(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isSourceDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_account_project" "project" {
						name = "tf_tests_cockpit_datasource_by_name"
				  	}

					resource "scaleway_cockpit_source" "main" {
					  project_id = scaleway_account_project.project.id
					  name       = "source-by-name"
					  type       = "logs"
					  retention_days = 30
					}

					data "scaleway_cockpit_source" "by_name" {
					  project_id = scaleway_account_project.project.id
					  name       = "source-by-name"
				      depends_on = [scaleway_cockpit_source.main]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.scaleway_cockpit_source.by_name", "id", "scaleway_cockpit_source.main", "id"),
					resource.TestCheckResourceAttr("data.scaleway_cockpit_source.by_name", "name", "source-by-name"),
					resource.TestCheckResourceAttr("data.scaleway_cockpit_source.by_name", "type", "logs"),
				),
			},
		},
	})
}

func TestAccCockpitSource_DataSource_Defaults(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	orgID, orgIDExists := tt.Meta.ScwClient().GetDefaultOrganizationID()
	if !orgIDExists {
		orgID = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`

					data scaleway_account_project "by_name" {
						name = "default"
						organization_id = "%s"
					}
					

					data "scaleway_cockpit_source" "default_metrics" {
					  	project_id = data.scaleway_account_project.by_name.id
	  					type       = "metrics"
						origin     = "scaleway"	
					}

					data "scaleway_cockpit_source" "default_logs" {
					  	project_id = data.scaleway_account_project.by_name.id
						type       = "logs"
						origin     = "scaleway"
					}
				`, orgID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.scaleway_cockpit_source.default_metrics", "name", "Scaleway Metrics"),
					resource.TestCheckResourceAttr("data.scaleway_cockpit_source.default_metrics", "type", "metrics"),
					resource.TestCheckResourceAttrSet("data.scaleway_cockpit_source.default_metrics", "url"),
					resource.TestCheckResourceAttr("data.scaleway_cockpit_source.default_logs", "name", "Scaleway Logs"),
					resource.TestCheckResourceAttr("data.scaleway_cockpit_source.default_logs", "type", "logs"),
					resource.TestCheckResourceAttrSet("data.scaleway_cockpit_source.default_logs", "url"),
				),
			},
		},
	})
}
