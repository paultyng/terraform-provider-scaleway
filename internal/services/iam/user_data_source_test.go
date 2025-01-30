package iam_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	iamSDK "github.com/scaleway/scaleway-sdk-go/api/iam/v1alpha1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/iam"
)

func TestAccDataSourceUser_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					data "scaleway_iam_user" "by_id" {
					  user_id = "b6360d4f-831c-45a8-889e-0b65ed079e63"
					}

					data "scaleway_iam_user" "by_email" {
					  email = "hashicorp@scaleway.com"
					  organization_id = "105bdce1-64c0-48ab-899d-868455867ecf"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIamUserExists(tt, "data.scaleway_iam_user.by_id"),
					testAccCheckIamUserExists(tt, "data.scaleway_iam_user.by_email"),

					resource.TestCheckResourceAttrSet("data.scaleway_iam_user.by_id", "user_id"),
					resource.TestCheckResourceAttrSet("data.scaleway_iam_user.by_id", "email"),

					resource.TestCheckResourceAttrSet("data.scaleway_iam_user.by_email", "user_id"),
					resource.TestCheckResourceAttrSet("data.scaleway_iam_user.by_email", "email"),
				),
			},
		},
	})
}

func testAccCheckIamUserExists(tt *acctest.TestTools, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("resource not found: %s", name)
		}

		iamAPI := iam.NewAPI(tt.Meta)

		_, err := iamAPI.GetUser(&iamSDK.GetUserRequest{
			UserID: rs.Primary.ID,
		})
		if err != nil {
			return fmt.Errorf("could not find user: %w", err)
		}

		return nil
	}
}
