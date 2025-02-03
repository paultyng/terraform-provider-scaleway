package secret_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/stretchr/testify/require"
)

func TestAccDataSourceSecret_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	ctx := context.Background()
	secretName := "scalewayDataSourceSecret"
	project, iamAPIKey, terminateFakeSideProject, err := acctest.CreateFakeIAMManager(tt)
	require.NoError(t, err)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: acctest.FakeSideProjectProviders(ctx, tt, project, iamAPIKey),
		CheckDestroy: resource.ComposeAggregateTestCheckFunc(
			func(_ *terraform.State) error {
				return terminateFakeSideProject()
			},
			testAccCheckSecretDestroy(tt),
		),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "scaleway_secret" "main" {
					  name        = "%[1]s"
					  description = "DataSourceSecret test description"
					  project_id  = "%[3]s"
					}
					
					data "scaleway_secret" "by_name" {
					  name            = scaleway_secret.main.name
					  organization_id = "%[2]s"
					  project_id      = "%[3]s"
					}
					
					data "scaleway_secret" "by_id" {
					  secret_id       = scaleway_secret.main.id
					  organization_id = "%[2]s"
					  project_id      = "%[3]s"
					}
				`, secretName, project.OrganizationID, project.ID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecretExists(tt, "data.scaleway_secret.by_name"),
					resource.TestCheckResourceAttr("data.scaleway_secret.by_name", "name", secretName),

					testAccCheckSecretExists(tt, "data.scaleway_secret.by_id"),
					resource.TestCheckResourceAttr("data.scaleway_secret.by_id", "name", secretName),
				),
			},
		},
	})
}

func TestAccDataSourceSecret_Path(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		CheckDestroy:      testAccCheckSecretDestroy(tt),
		ProviderFactories: tt.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_account_project" "project" {
						name = "tf-tests-secret-ds-path"
					}

					resource "scaleway_secret" "main" {
					  name = "test-secret-ds-path"
					  path = "/test-secret-ds-path-path"
					  project_id = scaleway_account_project.project.id
					}
					
					data "scaleway_secret" "by_name" {
					  name = scaleway_secret.main.name
					  path = "/test-secret-ds-path-path"
					  project_id = scaleway_account_project.project.id
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecretExists(tt, "data.scaleway_secret.by_name"),
					resource.TestCheckResourceAttr("data.scaleway_secret.by_name", "name", "test-secret-ds-path"),
				),
			},
		},
	})
}
