package accounttestfuncs

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	accountSDK "github.com/scaleway/scaleway-sdk-go/api/account/v3"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/logging"
)

func AddTestSweepers() {
	resource.AddTestSweepers("scaleway_account_project", &resource.Sweeper{
		Name: "scaleway_account_project",
		F:    testSweepAccountProject,
	})
}

func testSweepAccountProject(_ string) error {
	return acctest.Sweep(func(scwClient *scw.Client) error {
		accountAPI := accountSDK.NewProjectAPI(scwClient)

		logging.L.Debugf("sweeper: destroying the project")

		req := &accountSDK.ProjectAPIListProjectsRequest{}
		listProjects, err := accountAPI.ListProjects(req, scw.WithAllPages())
		if err != nil {
			return fmt.Errorf("failed to list projects: %w", err)
		}
		for _, project := range listProjects.Projects {
			// Do not delete default project
			if project.ID == req.OrganizationID || !acctest.IsTestResource(project.Name) {
				continue
			}
			err = accountAPI.DeleteProject(&accountSDK.ProjectAPIDeleteProjectRequest{
				ProjectID: project.ID,
			})
			if err != nil {
				return fmt.Errorf("failed to delete project: %w", err)
			}
		}
		return nil
	})
}
