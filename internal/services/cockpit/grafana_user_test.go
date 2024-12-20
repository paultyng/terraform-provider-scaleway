package cockpit_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	cockpitSDK "github.com/scaleway/scaleway-sdk-go/api/cockpit/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/cockpit"
)

func TestAccGrafanaUser_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	projectName := "tf_tests_cockpit_grafana_user_basic"
	grafanaTestUsername := "testuserbasic"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isGrafanaUserDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "scaleway_account_project" "project" {
						name = "%[1]s"
					}

					resource scaleway_cockpit_grafana_user main {
						project_id = scaleway_account_project.project.id
						login = "%[2]s"
						role = "editor"
					}
				`, projectName, grafanaTestUsername),
				Check: resource.ComposeTestCheckFunc(
					isGrafanaUserPresent(tt, "scaleway_cockpit_grafana_user.main"),
					resource.TestCheckResourceAttrPair("scaleway_cockpit_grafana_user.main", "project_id", "scaleway_account_project.project", "id"),
					resource.TestCheckResourceAttr("scaleway_cockpit_grafana_user.main", "login", grafanaTestUsername),
					resource.TestCheckResourceAttr("scaleway_cockpit_grafana_user.main", "role", "editor"),
					resource.TestCheckResourceAttrSet("scaleway_cockpit_grafana_user.main", "password"),
					resource.TestCheckResourceAttrSet("scaleway_cockpit_grafana_user.main", "grafana_url"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "scaleway_account_project" "project" {
						name = "%[1]s"
					}
				`, projectName),
			},
		},
	})
}

func TestAccGrafanaUser_Update(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	projectName := "tf_tests_cockpit_grafana_user_update"
	grafanaTestUsername := "testuserupdate"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isGrafanaUserDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "scaleway_account_project" "project" {
						name = "%[1]s"
				  	}

					resource scaleway_cockpit_grafana_user main {
						project_id = scaleway_account_project.project.id
						login = "%[2]s"
						role = "editor"
					}
				`, projectName, grafanaTestUsername),
				Check: resource.ComposeTestCheckFunc(
					isGrafanaUserPresent(tt, "scaleway_cockpit_grafana_user.main"),
					resource.TestCheckResourceAttrPair("scaleway_cockpit_grafana_user.main", "project_id", "scaleway_account_project.project", "id"),
					resource.TestCheckResourceAttr("scaleway_cockpit_grafana_user.main", "login", grafanaTestUsername),
					resource.TestCheckResourceAttr("scaleway_cockpit_grafana_user.main", "role", "editor"),
					resource.TestCheckResourceAttrSet("scaleway_cockpit_grafana_user.main", "password"),
					resource.TestCheckResourceAttrSet("scaleway_cockpit_grafana_user.main", "grafana_url"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "scaleway_account_project" "project" {
						name = "%[1]s"
					}

					resource scaleway_cockpit_grafana_user main {
						project_id = scaleway_account_project.project.id
						login = "%[2]s"
						role = "viewer"
					}
				`, projectName, grafanaTestUsername),
				Check: resource.ComposeTestCheckFunc(
					isGrafanaUserPresent(tt, "scaleway_cockpit_grafana_user.main"),
					resource.TestCheckResourceAttrPair("scaleway_cockpit_grafana_user.main", "project_id", "scaleway_account_project.project", "id"),
					resource.TestCheckResourceAttr("scaleway_cockpit_grafana_user.main", "login", grafanaTestUsername),
					resource.TestCheckResourceAttr("scaleway_cockpit_grafana_user.main", "role", "viewer"),
					resource.TestCheckResourceAttrSet("scaleway_cockpit_grafana_user.main", "password"),
					resource.TestCheckResourceAttrSet("scaleway_cockpit_grafana_user.main", "grafana_url"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "scaleway_account_project" "project" {
						name = "%[1]s"
					}
				`, projectName),
			},
		},
	})
}

func isGrafanaUserPresent(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource cockpit grafana user not found: %s", n)
		}

		api, projectID, grafanaUserID, err := cockpit.NewAPIGrafanaUserID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		res, err := api.ListGrafanaUsers(&cockpitSDK.GlobalAPIListGrafanaUsersRequest{
			ProjectID: projectID,
		}, scw.WithAllPages())
		if err != nil {
			return err
		}

		var grafanaUser *cockpitSDK.GrafanaUser
		for _, user := range res.GrafanaUsers {
			if user.ID == grafanaUserID {
				grafanaUser = user
				break
			}
		}

		if grafanaUser == nil {
			return fmt.Errorf("cockpit grafana user (%d) (project %s) not found", grafanaUserID, projectID)
		}

		return nil
	}
}

func isGrafanaUserDestroyed(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_cockpit_grafana_user" {
				continue
			}

			api, projectID, grafanaUserID, err := cockpit.NewAPIGrafanaUserID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}

			err = api.DeleteGrafanaUser(&cockpitSDK.GlobalAPIDeleteGrafanaUserRequest{
				ProjectID:     projectID,
				GrafanaUserID: grafanaUserID,
			})
			if err == nil {
				return fmt.Errorf("cockpit grafana user (%s) still exists", rs.Primary.ID)
			}

			if !httperrors.Is404(err) && !httperrors.Is403(err) {
				return err
			}
		}

		return nil
	}
}
