package iamtestfuncs

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/scaleway/scaleway-sdk-go/api/iam/v1alpha1/sweepers"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

func AddTestSweepers() {
	resource.AddTestSweepers("scaleway_iam_api_key", &resource.Sweeper{
		Name: "scaleway_iam_api_key",
		F:    testSweepIamAPIKey,
	})
	resource.AddTestSweepers("scaleway_iam_application", &resource.Sweeper{
		Name: "scaleway_iam_application",
		F:    testSweepIamApplication,
	})
	resource.AddTestSweepers("scaleway_iam_group", &resource.Sweeper{
		Name: "scaleway_iam_group",
		F:    testSweepIamGroup,
	})
	resource.AddTestSweepers("scaleway_iam_policy", &resource.Sweeper{
		Name: "scaleway_iam_policy",
		F:    testSweepIamPolicy,
	})
	resource.AddTestSweepers("scaleway_iam_ssh_key", &resource.Sweeper{
		Name: "scaleway_iam_ssh_key",
		F:    testSweepSSHKey,
	})
	resource.AddTestSweepers("scaleway_iam_user", &resource.Sweeper{
		Name: "scaleway_iam_user",
		F:    testSweepUser,
	})
}

func testSweepUser(_ string) error {
	return acctest.Sweep(sweepers.SweepUser)
}

func testSweepSSHKey(_ string) error {
	return acctest.Sweep(sweepers.SweepSSHKey)
}

func testSweepIamPolicy(_ string) error {
	return acctest.Sweep(sweepers.SweepPolicy)
}

func testSweepIamGroup(_ string) error {
	return acctest.Sweep(sweepers.SweepGroup)
}

func testSweepIamApplication(_ string) error {
	return acctest.Sweep(sweepers.SweepApplication)
}

func testSweepIamAPIKey(_ string) error {
	return acctest.Sweep(sweepers.SweepAPIKey)
}
