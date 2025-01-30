package iam_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

func TestAccDataSourceApiKey_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckIamAPIKeyDestroy(tt),
		),
		Steps: []resource.TestStep{
			{
				Config: `
						resource "scaleway_iam_application" "main" {
							name = "tf_tests_app_key_basic"
						}

						resource "scaleway_iam_api_key" "main" {
							application_id = scaleway_iam_application.main.id
							description = "tf_tests_with_application"
						}

						data "scaleway_iam_api_key" "main" {
							access_key = scaleway_iam_api_key.main.id
						}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIamAPIKeyExists(tt, "scaleway_iam_api_key.main"),
					resource.TestCheckResourceAttrPair("data.scaleway_iam_api_key.main", "access_key", "scaleway_iam_api_key.main", "id"),
					resource.TestCheckResourceAttrPair("data.scaleway_iam_api_key.main", "id", "scaleway_iam_api_key.main", "id"),
					resource.TestCheckResourceAttrPair("data.scaleway_iam_api_key.main", "application_id", "scaleway_iam_application.main", "id"),
				),
			},
		},
	})
}
