package cockpit_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	cockpitSDK "github.com/scaleway/scaleway-sdk-go/api/cockpit/v1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/cockpit"
)

func TestAccCockpitSource_Basic_metrics(t *testing.T) {
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
						name = "tf_tests_cockpit_datasource_basic"
				  	}

					resource "scaleway_cockpit_source" "main" {
					  project_id = scaleway_account_project.project.id
					  name       = "my-source"
					  type       = "metrics"
					  retention_days      = 31
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isSourcePresent(tt, "scaleway_cockpit_source.main"),
					resource.TestCheckResourceAttr("scaleway_cockpit_source.main", "name", "my-source"),
					resource.TestCheckResourceAttr("scaleway_cockpit_source.main", "type", "metrics"),
					resource.TestCheckResourceAttr("scaleway_cockpit_source.main", "region", "fr-par"),
					resource.TestCheckResourceAttr("scaleway_cockpit_source.main", "retention_days", "31"),
					resource.TestCheckResourceAttrSet("scaleway_cockpit_source.main", "url"),
					resource.TestCheckResourceAttrSet("scaleway_cockpit_source.main", "push_url"),
					resource.TestCheckResourceAttrSet("scaleway_cockpit_source.main", "origin"),
					resource.TestCheckResourceAttrSet("scaleway_cockpit_source.main", "created_at"),
					resource.TestCheckResourceAttrSet("scaleway_cockpit_source.main", "updated_at"),
					resource.TestCheckResourceAttrSet("scaleway_cockpit_source.main", "synchronized_with_grafana"),
					resource.TestCheckResourceAttrPair("scaleway_cockpit_source.main", "project_id", "scaleway_account_project.project", "id"),
				),
			},
		},
	})
}

func TestAccCockpitSource_Basic_logs(t *testing.T) {
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
						name = "tf_tests_cockpit_datasource_basic"
				  	}

					resource "scaleway_cockpit_source" "main" {
					  project_id = scaleway_account_project.project.id
					  name       = "my-source"
					  type       = "logs"
					  retention_days  = 31

					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isSourcePresent(tt, "scaleway_cockpit_source.main"),
					resource.TestCheckResourceAttr("scaleway_cockpit_source.main", "name", "my-source"),
					resource.TestCheckResourceAttr("scaleway_cockpit_source.main", "type", "logs"),
					resource.TestCheckResourceAttr("scaleway_cockpit_source.main", "region", "fr-par"),
					resource.TestCheckResourceAttrSet("scaleway_cockpit_source.main", "url"),
					resource.TestCheckResourceAttrSet("scaleway_cockpit_source.main", "push_url"),
					resource.TestCheckResourceAttrSet("scaleway_cockpit_source.main", "origin"),
					resource.TestCheckResourceAttrSet("scaleway_cockpit_source.main", "created_at"),
					resource.TestCheckResourceAttrSet("scaleway_cockpit_source.main", "updated_at"),
					resource.TestCheckResourceAttrSet("scaleway_cockpit_source.main", "synchronized_with_grafana"),
					resource.TestCheckResourceAttrPair("scaleway_cockpit_source.main", "project_id", "scaleway_account_project.project", "id"),
				),
			},
		},
	})
}

func TestAccCockpitSource_retention_days(t *testing.T) {
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
						name = "tf_tests_cockpit_datasource_basic"
				  	}

					resource "scaleway_cockpit_source" "main" {
					  project_id = scaleway_account_project.project.id
					  name       = "my-source"
					  type       = "logs"
					  retention_days       = 13

					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isSourcePresent(tt, "scaleway_cockpit_source.main"),
					resource.TestCheckResourceAttr("scaleway_cockpit_source.main", "name", "my-source"),
					resource.TestCheckResourceAttr("scaleway_cockpit_source.main", "type", "logs"),
					resource.TestCheckResourceAttr("scaleway_cockpit_source.main", "region", "fr-par"),
					resource.TestCheckResourceAttr("scaleway_cockpit_source.main", "retention_days", "13"),
					resource.TestCheckResourceAttrSet("scaleway_cockpit_source.main", "url"),
					resource.TestCheckResourceAttrSet("scaleway_cockpit_source.main", "push_url"),
					resource.TestCheckResourceAttrSet("scaleway_cockpit_source.main", "origin"),
					resource.TestCheckResourceAttrSet("scaleway_cockpit_source.main", "created_at"),
					resource.TestCheckResourceAttrSet("scaleway_cockpit_source.main", "updated_at"),
					resource.TestCheckResourceAttrSet("scaleway_cockpit_source.main", "synchronized_with_grafana"),
					resource.TestCheckResourceAttrPair("scaleway_cockpit_source.main", "project_id", "scaleway_account_project.project", "id"),
				),
			},
		},
	})
}

func TestAccCockpitSource_Update(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isSourceDestroyed(tt),
		Steps: []resource.TestStep{
			// Initial creation
			{
				Config: `
					resource "scaleway_account_project" "project" {
						name = "tf_tests_cockpit_source_update"
					}

					resource "scaleway_cockpit_source" "main" {
						project_id = scaleway_account_project.project.id
						name       = "initial-name"
						type       = "logs"
						retention_days = 10
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isSourcePresent(tt, "scaleway_cockpit_source.main"),
					resource.TestCheckResourceAttr("scaleway_cockpit_source.main", "name", "initial-name"),
					resource.TestCheckResourceAttr("scaleway_cockpit_source.main", "retention_days", "10"),
				),
			},
			{
				Config: `
					resource "scaleway_account_project" "project" {
						name = "tf_tests_cockpit_source_update"
					}

					resource "scaleway_cockpit_source" "main" {
						project_id = scaleway_account_project.project.id
						name       = "initial-name"
						type       = "logs"
						retention_days = 20
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isSourcePresent(tt, "scaleway_cockpit_source.main"),
					resource.TestCheckResourceAttr("scaleway_cockpit_source.main", "retention_days", "20"),
				),
			},
			{
				Config: `
					resource "scaleway_account_project" "project" {
						name = "tf_tests_cockpit_source_update"
					}

					resource "scaleway_cockpit_source" "main" {
						project_id = scaleway_account_project.project.id
						name       = "updated-name"
						type       = "logs"
						retention_days = 20
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isSourcePresent(tt, "scaleway_cockpit_source.main"),
					resource.TestCheckResourceAttr("scaleway_cockpit_source.main", "name", "updated-name"),
				),
			},
		},
	})
}

func isSourcePresent(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource cockpit source not found: %s", n)
		}

		api, region, ID, err := cockpit.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = api.GetDataSource(&cockpitSDK.RegionalAPIGetDataSourceRequest{
			Region:       region,
			DataSourceID: ID,
		})
		if err != nil {
			return err
		}

		return nil
	}
}

func isSourceDestroyed(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_cockpit_source" {
				continue
			}

			api, region, ID, err := cockpit.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}

			_, err = api.GetDataSource(&cockpitSDK.RegionalAPIGetDataSourceRequest{
				Region:       region,
				DataSourceID: ID,
			})

			if err == nil {
				return fmt.Errorf("cockpit source (%s) still exists", rs.Primary.ID)
			}

			if !httperrors.Is404(err) {
				return err
			}
		}

		return nil
	}
}
