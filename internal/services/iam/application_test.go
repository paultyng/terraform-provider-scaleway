package iam_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	iamSDK "github.com/scaleway/scaleway-sdk-go/api/iam/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/iam"
)

func init() {
	resource.AddTestSweepers("scaleway_iam_application", &resource.Sweeper{
		Name: "scaleway_iam_application",
		F:    testSweepIamApplication,
	})
}

func testSweepIamApplication(_ string) error {
	return acctest.Sweep(func(scwClient *scw.Client) error {
		api := iamSDK.NewAPI(scwClient)

		orgID, exists := scwClient.GetDefaultOrganizationID()
		if !exists {
			return errors.New("missing organizationID")
		}

		listApps, err := api.ListApplications(&iamSDK.ListApplicationsRequest{
			OrganizationID: orgID,
		})
		if err != nil {
			return fmt.Errorf("failed to list applications: %w", err)
		}
		for _, app := range listApps.Applications {
			if !acctest.IsTestResource(app.Name) {
				continue
			}

			err = api.DeleteApplication(&iamSDK.DeleteApplicationRequest{
				ApplicationID: app.ID,
			})
			if err != nil {
				return fmt.Errorf("failed to delete application: %w", err)
			}
		}
		return nil
	})
}

func TestAccApplication_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckIamApplicationDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
						resource "scaleway_iam_application" "main" {
							name = "tf_tests_app_basic"
							description = "a description"
							tags = ["tf_tests", "tests"]
						}
					`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIamApplicationExists(tt, "scaleway_iam_application.main"),
					resource.TestCheckResourceAttr("scaleway_iam_application.main", "name", "tf_tests_app_basic"),
					resource.TestCheckResourceAttr("scaleway_iam_application.main", "description", "a description"),
					resource.TestCheckResourceAttr("scaleway_iam_application.main", "tags.#", "2"),
					resource.TestCheckResourceAttr("scaleway_iam_application.main", "tags.0", "tf_tests"),
					resource.TestCheckResourceAttr("scaleway_iam_application.main", "tags.1", "tests"),
				),
			},
			{
				Config: `
						resource "scaleway_iam_application" "main" {
							name = "tf_tests_app_basic_rename"
							description = "another description"
							tags = []
						}
					`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIamApplicationExists(tt, "scaleway_iam_application.main"),
					resource.TestCheckResourceAttr("scaleway_iam_application.main", "name", "tf_tests_app_basic_rename"),
					resource.TestCheckResourceAttr("scaleway_iam_application.main", "description", "another description"),
					resource.TestCheckResourceAttr("scaleway_iam_application.main", "tags.#", "0"),
				),
			},
		},
	})
}

func TestAccApplication_NoUpdate(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckIamApplicationDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
						resource "scaleway_iam_application" "main" {
							name = "tf_tests_app_noupdate"
						}
					`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIamApplicationExists(tt, "scaleway_iam_application.main"),
				),
			},
			{
				Config: `
						resource "scaleway_iam_application" "main" {
							name = "tf_tests_app_noupdate"
						}
					`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIamApplicationExists(tt, "scaleway_iam_application.main"),
				),
			},
		},
	})
}

func testAccCheckIamApplicationExists(tt *acctest.TestTools, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("resource not found: %s", name)
		}

		iamAPI := iam.NewAPI(tt.Meta)

		_, err := iamAPI.GetApplication(&iamSDK.GetApplicationRequest{
			ApplicationID: rs.Primary.ID,
		})
		if err != nil {
			return fmt.Errorf("could not find application: %w", err)
		}

		return nil
	}
}

func testAccCheckIamApplicationDestroy(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "scaleway_iam_application" {
				continue
			}

			iamAPI := iam.NewAPI(tt.Meta)

			_, err := iamAPI.GetApplication(&iamSDK.GetApplicationRequest{
				ApplicationID: rs.Primary.ID,
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
