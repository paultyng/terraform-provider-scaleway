package iam_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	iamSDK "github.com/scaleway/scaleway-sdk-go/api/iam/v1alpha1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/iam"
)

func TestAccApiKey_WithApplication(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckIamAPIKeyDestroy(tt),
			testAccCheckIamApplicationDestroy(tt),
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
					`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIamAPIKeyExists(tt, "scaleway_iam_api_key.main"),
					resource.TestCheckResourceAttrPair("scaleway_iam_api_key.main", "application_id", "scaleway_iam_application.main", "id"),
					resource.TestCheckResourceAttr("scaleway_iam_api_key.main", "description", "tf_tests_with_application"),
					resource.TestCheckResourceAttrSet("scaleway_iam_api_key.main", "secret_key"),
				),
			},
			{
				Config: `
						resource "scaleway_iam_application" "main" {
							name = "tf_tests_app_key_basic"
						}

						resource "scaleway_iam_api_key" "main" {
							application_id = scaleway_iam_application.main.id
							description = "tf_tests_with_application_changed"
						}
					`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIamAPIKeyExists(tt, "scaleway_iam_api_key.main"),
					resource.TestCheckResourceAttrPair("scaleway_iam_api_key.main", "application_id", "scaleway_iam_application.main", "id"),
					resource.TestCheckResourceAttr("scaleway_iam_api_key.main", "description", "tf_tests_with_application_changed"),
					resource.TestCheckResourceAttrSet("scaleway_iam_api_key.main", "secret_key"),
				),
			},
			{
				ResourceName:            "scaleway_iam_api_key.main",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"secret_key"},
			},
		},
	})
}

func TestAccApiKey_WithApplicationChange(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckIamAPIKeyDestroy(tt),
			testAccCheckIamApplicationDestroy(tt),
		),
		Steps: []resource.TestStep{
			{
				Config: `
						resource "scaleway_iam_application" "main" {
							name = "tf_tests_api_key_app_change"
						}

						resource "scaleway_iam_application" "main2" {
							name = "tf_tests_api_key_app_change2"
						}

						resource "scaleway_iam_api_key" "main" {
							application_id = scaleway_iam_application.main.id
							description = "tf_tests_with_application_change"
						}
					`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIamAPIKeyExists(tt, "scaleway_iam_api_key.main"),
					resource.TestCheckResourceAttrPair("scaleway_iam_api_key.main", "application_id", "scaleway_iam_application.main", "id"),
					resource.TestCheckResourceAttr("scaleway_iam_api_key.main", "description", "tf_tests_with_application_change"),
					resource.TestCheckResourceAttrSet("scaleway_iam_api_key.main", "secret_key"),
				),
			},
			{
				Config: `
						resource "scaleway_iam_application" "main" {
							name = "tf_tests_api_key_app_change"
						}

						resource "scaleway_iam_application" "main2" {
							name = "tf_tests_api_key_app_change2"
						}

						resource "scaleway_iam_api_key" "main" {
							application_id = scaleway_iam_application.main2.id
							description = "tf_tests_with_application_change"
						}
					`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIamAPIKeyExists(tt, "scaleway_iam_api_key.main"),
					resource.TestCheckResourceAttrPair("scaleway_iam_api_key.main", "application_id", "scaleway_iam_application.main2", "id"),
					resource.TestCheckResourceAttr("scaleway_iam_api_key.main", "description", "tf_tests_with_application_change"),
					resource.TestCheckResourceAttrSet("scaleway_iam_api_key.main", "secret_key"),
				),
			},
			{
				ResourceName:            "scaleway_iam_api_key.main",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"secret_key"},
			},
		},
	})
}

func TestAccApiKey_Expires(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckIamAPIKeyDestroy(tt),
			testAccCheckIamApplicationDestroy(tt),
		),
		Steps: []resource.TestStep{
			{
				Config: `
						resource "scaleway_iam_application" "main" {
							name = "tf_tests_app_expires_at"
						}

						resource "scaleway_iam_api_key" "main" {
							application_id = scaleway_iam_application.main.id
							description = "tf_tests_expires"
							expires_at = "2025-07-06T11:00:00+02:00"
						}
					`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIamAPIKeyExists(tt, "scaleway_iam_api_key.main"),
					resource.TestCheckResourceAttrPair("scaleway_iam_api_key.main", "application_id", "scaleway_iam_application.main", "id"),
					resource.TestCheckResourceAttr("scaleway_iam_api_key.main", "description", "tf_tests_expires"),
					resource.TestCheckResourceAttr("scaleway_iam_api_key.main", "expires_at", "2025-07-06T09:00:00Z"),
				),
			},
			{
				Config: `
						resource "scaleway_iam_application" "main" {
							name = "tf_tests_app_expires_at"
						}

						resource "scaleway_iam_api_key" "main" {
							application_id = scaleway_iam_application.main.id
							description = "tf_tests_expires"
							expires_at = "2025-07-06T09:00:00Z"
						}
					`,
				PlanOnly: true,
			},
		},
	})
}

func TestAccApiKey_NoUpdate(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckIamAPIKeyDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
						resource "scaleway_iam_application" "main" {
							name = "tf_tests_app_key_noupdate"
						}

						resource "scaleway_iam_api_key" "main" {
							application_id = scaleway_iam_application.main.id
							description = "tf_tests_no_update"
						}
					`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIamAPIKeyExists(tt, "scaleway_iam_api_key.main"),
					resource.TestCheckResourceAttrPair("scaleway_iam_api_key.main", "application_id", "scaleway_iam_application.main", "id"),
					resource.TestCheckResourceAttr("scaleway_iam_api_key.main", "description", "tf_tests_no_update"),
				),
			},
			{
				Config: `
						resource "scaleway_iam_application" "main" {
							name = "tf_tests_app_key_noupdate"
						}

						resource "scaleway_iam_api_key" "main" {
							application_id = scaleway_iam_application.main.id
							description = "tf_tests_no_update"
						}
					`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIamAPIKeyExists(tt, "scaleway_iam_api_key.main"),
					resource.TestCheckResourceAttrPair("scaleway_iam_api_key.main", "application_id", "scaleway_iam_application.main", "id"),
					resource.TestCheckResourceAttr("scaleway_iam_api_key.main", "description", "tf_tests_no_update"),
				),
			},
		},
	})
}

func testAccCheckIamAPIKeyExists(tt *acctest.TestTools, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("resource not found: %s", name)
		}

		iamAPI := iam.NewAPI(tt.Meta)

		_, err := iamAPI.GetAPIKey(&iamSDK.GetAPIKeyRequest{
			AccessKey: rs.Primary.ID,
		})
		if err != nil {
			return fmt.Errorf("could not find api key: %w", err)
		}

		return nil
	}
}

func testAccCheckIamAPIKeyDestroy(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "scaleway_iam_api_key" {
				continue
			}

			iamAPI := iam.NewAPI(tt.Meta)

			_, err := iamAPI.GetAPIKey(&iamSDK.GetAPIKeyRequest{
				AccessKey: rs.Primary.ID,
			})

			// If no error resource still exist
			if err == nil {
				return fmt.Errorf("resource %s(%s) still exist", rs.Type, rs.Primary.ID)
			}

			// Unexpected api error we return it
			if !httperrors.Is404(err) {
				return err
			}
		}

		return nil
	}
}
