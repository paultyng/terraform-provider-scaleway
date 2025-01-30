package cockpit_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	cockpitSDK "github.com/scaleway/scaleway-sdk-go/api/cockpit/v1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/cockpit"
)

func TestAccToken_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	projectName := "tf_tests_cockpit_token_basic"
	tokenName := projectName

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isTokenDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "scaleway_account_project" "project" {
						name = "%[1]s"
				  	}

					resource scaleway_cockpit_token main {
						project_id = scaleway_account_project.project.id
						name = "%[2]s"
						scopes {
							query_metrics = true
							write_logs = false
						}
					}
				`, projectName, tokenName),
				Check: resource.ComposeTestCheckFunc(
					isTokenPresent(tt, "scaleway_cockpit_token.main"),
					resource.TestCheckResourceAttrPair("scaleway_cockpit_token.main", "project_id", "scaleway_account_project.project", "id"),
					resource.TestCheckResourceAttrSet("scaleway_cockpit_token.main", "secret_key"),
					resource.TestCheckResourceAttr("scaleway_cockpit_token.main", "name", tokenName),
					resource.TestCheckResourceAttr("scaleway_cockpit_token.main", "scopes.0.query_metrics", "true"),
					resource.TestCheckResourceAttr("scaleway_cockpit_token.main", "scopes.0.write_metrics", "true"),
					resource.TestCheckResourceAttr("scaleway_cockpit_token.main", "scopes.0.setup_metrics_rules", "false"),
					resource.TestCheckResourceAttr("scaleway_cockpit_token.main", "scopes.0.query_logs", "false"),
					resource.TestCheckResourceAttr("scaleway_cockpit_token.main", "scopes.0.write_logs", "false"),
					resource.TestCheckResourceAttr("scaleway_cockpit_token.main", "scopes.0.setup_logs_rules", "false"),
					resource.TestCheckResourceAttr("scaleway_cockpit_token.main", "scopes.0.query_traces", "false"),
					resource.TestCheckResourceAttr("scaleway_cockpit_token.main", "scopes.0.write_traces", "false"),
				),
			},
		},
	})
}

func TestAccToken_NoScopes(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	projectName := "tf_tests_cockpit_token_no_scopes"
	tokenName := "tf_tests_cockpit_token_no_scopes"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isTokenDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "scaleway_account_project" "project" {
						name = "%[1]s"
				  	}

					resource scaleway_cockpit_token main {
						project_id = scaleway_account_project.project.id
						name = "%[2]s"
					}
				`, projectName, tokenName),
				Check: resource.ComposeTestCheckFunc(
					isTokenPresent(tt, "scaleway_cockpit_token.main"),
					resource.TestCheckResourceAttrPair("scaleway_cockpit_token.main", "project_id", "scaleway_account_project.project", "id"),
					resource.TestCheckResourceAttrSet("scaleway_cockpit_token.main", "secret_key"),
					resource.TestCheckResourceAttr("scaleway_cockpit_token.main", "name", tokenName),
					resource.TestCheckResourceAttr("scaleway_cockpit_token.main", "scopes.0.query_metrics", "false"),
					resource.TestCheckResourceAttr("scaleway_cockpit_token.main", "scopes.0.write_metrics", "true"),
					resource.TestCheckResourceAttr("scaleway_cockpit_token.main", "scopes.0.setup_metrics_rules", "false"),
					resource.TestCheckResourceAttr("scaleway_cockpit_token.main", "scopes.0.query_logs", "false"),
					resource.TestCheckResourceAttr("scaleway_cockpit_token.main", "scopes.0.write_logs", "true"),
					resource.TestCheckResourceAttr("scaleway_cockpit_token.main", "scopes.0.setup_logs_rules", "false"),
					resource.TestCheckResourceAttr("scaleway_cockpit_token.main", "scopes.0.query_traces", "false"),
					resource.TestCheckResourceAttr("scaleway_cockpit_token.main", "scopes.0.write_traces", "false"),
				),
			},
		},
	})
}

func TestAccToken_Update(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	projectName := "tf_tests_cockpit_token_update"
	tokenName := projectName

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isTokenDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "scaleway_account_project" "project" {
						name = "%[1]s"
				  	}

					resource scaleway_cockpit_token main {
						project_id = scaleway_account_project.project.id
						name = "%[2]s"
						scopes {
							query_metrics = true
							write_logs = false
						}
					}
				`, projectName, tokenName),
				Check: resource.ComposeTestCheckFunc(
					isTokenPresent(tt, "scaleway_cockpit_token.main"),
					resource.TestCheckResourceAttrPair("scaleway_cockpit_token.main", "project_id", "scaleway_account_project.project", "id"),
					resource.TestCheckResourceAttrSet("scaleway_cockpit_token.main", "secret_key"),
					resource.TestCheckResourceAttr("scaleway_cockpit_token.main", "name", tokenName),
					resource.TestCheckResourceAttr("scaleway_cockpit_token.main", "scopes.0.query_metrics", "true"),
					resource.TestCheckResourceAttr("scaleway_cockpit_token.main", "scopes.0.write_metrics", "true"),
					resource.TestCheckResourceAttr("scaleway_cockpit_token.main", "scopes.0.setup_metrics_rules", "false"),
					resource.TestCheckResourceAttr("scaleway_cockpit_token.main", "scopes.0.query_logs", "false"),
					resource.TestCheckResourceAttr("scaleway_cockpit_token.main", "scopes.0.write_logs", "false"),
					resource.TestCheckResourceAttr("scaleway_cockpit_token.main", "scopes.0.setup_logs_rules", "false"),
					resource.TestCheckResourceAttr("scaleway_cockpit_token.main", "scopes.0.query_traces", "false"),
					resource.TestCheckResourceAttr("scaleway_cockpit_token.main", "scopes.0.write_traces", "false"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "scaleway_account_project" "project" {
						name = "%[1]s"
					}

					resource scaleway_cockpit_token main {
						project_id = scaleway_account_project.project.id
						name = "%[2]s"
						scopes {
							query_metrics = true
							write_logs = true
						}
					}
				`, projectName, tokenName),
				Check: resource.ComposeTestCheckFunc(
					isTokenPresent(tt, "scaleway_cockpit_token.main"),
					resource.TestCheckResourceAttrPair("scaleway_cockpit_token.main", "project_id", "scaleway_account_project.project", "id"),
					resource.TestCheckResourceAttrSet("scaleway_cockpit_token.main", "secret_key"),
					resource.TestCheckResourceAttr("scaleway_cockpit_token.main", "name", tokenName),
					resource.TestCheckResourceAttr("scaleway_cockpit_token.main", "scopes.0.query_metrics", "true"),
					resource.TestCheckResourceAttr("scaleway_cockpit_token.main", "scopes.0.write_metrics", "true"),
					resource.TestCheckResourceAttr("scaleway_cockpit_token.main", "scopes.0.setup_metrics_rules", "false"),
					resource.TestCheckResourceAttr("scaleway_cockpit_token.main", "scopes.0.query_logs", "false"),
					resource.TestCheckResourceAttr("scaleway_cockpit_token.main", "scopes.0.write_logs", "true"),
					resource.TestCheckResourceAttr("scaleway_cockpit_token.main", "scopes.0.setup_logs_rules", "false"),
					resource.TestCheckResourceAttr("scaleway_cockpit_token.main", "scopes.0.query_traces", "false"),
					resource.TestCheckResourceAttr("scaleway_cockpit_token.main", "scopes.0.write_traces", "false"),
				),
			},
		},
	})
}

func isTokenPresent(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource cockpit token not found: %s", n)
		}

		api, region, ID, err := cockpit.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = api.GetToken(&cockpitSDK.RegionalAPIGetTokenRequest{
			TokenID: ID,
			Region:  region,
		})
		if err != nil {
			return err
		}

		return nil
	}
}

func isTokenDestroyed(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_cockpit_token" {
				continue
			}

			api, region, ID, err := cockpit.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}

			err = api.DeleteToken(&cockpitSDK.RegionalAPIDeleteTokenRequest{
				TokenID: ID,
				Region:  region,
			})
			if err == nil {
				return fmt.Errorf("cockpit token (%s) still exists", rs.Primary.ID)
			}

			// Currently the API returns a 403 error when we try to delete a token that does not exist
			if !httperrors.Is404(err) && !httperrors.Is403(err) {
				return err
			}
		}

		return nil
	}
}
