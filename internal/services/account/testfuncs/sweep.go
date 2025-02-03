package accounttestfuncs

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/scaleway/scaleway-sdk-go/api/account/v3/sweepers"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

func AddTestSweepers() {
	resource.AddTestSweepers("scaleway_account_project", &resource.Sweeper{
		Name: "scaleway_account_project",
		F:    testSweepAccountProject,
	})
}

func testSweepAccountProject(_ string) error {
	return acctest.Sweep(sweepers.SweepProjects)
}
