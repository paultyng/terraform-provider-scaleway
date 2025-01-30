package iam_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	iamSDK "github.com/scaleway/scaleway-sdk-go/api/iam/v1alpha1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/iam"
)

func TestAccGroupMembership_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckIamGroupDestroy(tt),
			testAccCheckIamApplicationDestroy(tt),
		), Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_iam_group main {
						name = "tf-tests-iam-group-membership-basic"
						external_membership = true
					}

					resource scaleway_iam_application main {
						name = "tf-tests-iam-group-membership-basic"
					}

					resource scaleway_iam_group_membership main {
						group_id = scaleway_iam_group.main.id
						application_id = scaleway_iam_application.main.id
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIamGroupMembershipApplicationInGroup(tt, "scaleway_iam_group_membership.main", "scaleway_iam_application.main"),
					acctest.CheckResourceAttrUUID("scaleway_iam_group_membership.main", "id"),
				),
			},
			{
				Config: `
					resource scaleway_iam_group main {
						name = "tf-tests-iam-group-membership-basic"
						external_membership = true
					}

					resource scaleway_iam_application main {
						name = "tf-tests-iam-group-membership-basic"
					}

					resource scaleway_iam_group_membership main {
						group_id = scaleway_iam_group.main.id
						application_id = scaleway_iam_application.main.id
					}

					resource scaleway_iam_group_membership import {
						group_id = scaleway_iam_group.main.id
						application_id = scaleway_iam_application.main.id
					}
				`,
				ImportState:  true,
				ResourceName: "scaleway_iam_group_membership.import",
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					groupID := state.RootModule().Resources["scaleway_iam_group.main"].Primary.ID
					applicationID := state.RootModule().Resources["scaleway_iam_application.main"].Primary.ID

					return iam.GroupMembershipID(groupID, nil, &applicationID), nil
				},
				ImportStatePersist: true,
			},
			{
				Config: `
					resource scaleway_iam_group main {
						name = "tf-tests-iam-group-membership-basic"
						external_membership = true
					}

					resource scaleway_iam_application main {
						name = "tf-tests-iam-group-membership-basic"
					}

					resource scaleway_iam_group_membership main {
						group_id = scaleway_iam_group.main.id
						application_id = scaleway_iam_application.main.id
					}

					resource scaleway_iam_group_membership import {
						group_id = scaleway_iam_group.main.id
						application_id = scaleway_iam_application.main.id
					}
				`,
				PlanOnly: true,
			},
		},
	})
}

func TestAccGroupMembership_User(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckIamGroupDestroy(tt),
			testAccCheckIamApplicationDestroy(tt),
		), Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_iam_group main {
						name = "tf-tests-iam-group-membership-user"
						external_membership = true
					}

					data "scaleway_iam_user" "main" {
						user_id = "b6360d4f-831c-45a8-889e-0b65ed079e63"
					}

					resource scaleway_iam_group_membership main {
						group_id = scaleway_iam_group.main.id
						user_id = data.scaleway_iam_user.main.id
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIamGroupMembershipUserInGroup(tt, "scaleway_iam_group_membership.main", "data.scaleway_iam_user.main"),
					acctest.CheckResourceAttrUUID("scaleway_iam_group_membership.main", "id"),
				),
			},
		},
	})
}

func testAccCheckIamGroupMembershipApplicationInGroup(tt *acctest.TestTools, n string, appN string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		appRS, ok := state.RootModule().Resources[appN]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		expectedApplicationID := appRS.Primary.ID

		api := iam.NewAPI(tt.Meta)
		groupID, _, applicationID, err := iam.ExpandGroupMembershipID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if applicationID != expectedApplicationID {
			return fmt.Errorf("group membership id does not contain expected application id, expected %s, got %s", expectedApplicationID, applicationID)
		}

		group, err := api.GetGroup(&iamSDK.GetGroupRequest{
			GroupID: groupID,
		})
		if err != nil {
			return err
		}

		foundInGroup := false

		for _, groupApplicationID := range group.ApplicationIDs {
			if groupApplicationID == applicationID {
				foundInGroup = true
			}
		}

		if !foundInGroup {
			return errors.New("application not found in group")
		}

		return nil
	}
}

func testAccCheckIamGroupMembershipUserInGroup(tt *acctest.TestTools, n string, appN string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		appRS, ok := state.RootModule().Resources[appN]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		expectedUserID := appRS.Primary.ID

		api := iam.NewAPI(tt.Meta)
		groupID, userID, _, err := iam.ExpandGroupMembershipID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if userID != expectedUserID {
			return fmt.Errorf("group membership id does not contain expected user id, expected %s, got %s", expectedUserID, userID)
		}

		group, err := api.GetGroup(&iamSDK.GetGroupRequest{
			GroupID: groupID,
		})
		if err != nil {
			return err
		}

		foundInGroup := false

		for _, groupUserID := range group.UserIDs {
			if groupUserID == userID {
				foundInGroup = true
			}
		}

		if !foundInGroup {
			return errors.New("user not found in group")
		}

		return nil
	}
}
