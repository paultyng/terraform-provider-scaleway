package iamtestfuncs

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	iamSDK "github.com/scaleway/scaleway-sdk-go/api/iam/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/logging"
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
	return acctest.Sweep(func(scwClient *scw.Client) error {
		api := iamSDK.NewAPI(scwClient)

		orgID, exists := scwClient.GetDefaultOrganizationID()
		if !exists {
			return errors.New("missing organizationID")
		}

		listUsers, err := api.ListUsers(&iamSDK.ListUsersRequest{
			OrganizationID: &orgID,
		})
		if err != nil {
			return fmt.Errorf("failed to list users: %w", err)
		}
		for _, user := range listUsers.Users {
			if !acctest.IsTestResource(user.Email) {
				continue
			}
			err = api.DeleteUser(&iamSDK.DeleteUserRequest{
				UserID: user.ID,
			})
			if err != nil {
				return fmt.Errorf("failed to delete user: %w", err)
			}
		}
		return nil
	})
}

func testSweepSSHKey(_ string) error {
	return acctest.Sweep(func(scwClient *scw.Client) error {
		iamAPI := iamSDK.NewAPI(scwClient)

		logging.L.Debugf("sweeper: destroying the SSH keys")

		listSSHKeys, err := iamAPI.ListSSHKeys(&iamSDK.ListSSHKeysRequest{}, scw.WithAllPages())
		if err != nil {
			return fmt.Errorf("error listing SSH keys in sweeper: %s", err)
		}

		for _, sshKey := range listSSHKeys.SSHKeys {
			if !acctest.IsTestResource(sshKey.Name) {
				continue
			}
			err := iamAPI.DeleteSSHKey(&iamSDK.DeleteSSHKeyRequest{
				SSHKeyID: sshKey.ID,
			})
			if err != nil {
				return fmt.Errorf("error deleting SSH key in sweeper: %s", err)
			}
		}

		return nil
	})
}

func testSweepIamPolicy(_ string) error {
	return acctest.Sweep(func(scwClient *scw.Client) error {
		api := iamSDK.NewAPI(scwClient)

		orgID, exists := scwClient.GetDefaultOrganizationID()
		if !exists {
			return errors.New("missing organizationID")
		}

		listPols, err := api.ListPolicies(&iamSDK.ListPoliciesRequest{
			OrganizationID: orgID,
		})
		if err != nil {
			return fmt.Errorf("failed to list policies: %w", err)
		}
		for _, pol := range listPols.Policies {
			if !acctest.IsTestResource(pol.Name) {
				continue
			}
			err = api.DeletePolicy(&iamSDK.DeletePolicyRequest{
				PolicyID: pol.ID,
			})
			if err != nil {
				return fmt.Errorf("failed to delete policy: %w", err)
			}
		}
		return nil
	})
}

func testSweepIamGroup(_ string) error {
	return acctest.Sweep(func(scwClient *scw.Client) error {
		api := iamSDK.NewAPI(scwClient)

		orgID, exists := scwClient.GetDefaultOrganizationID()
		if !exists {
			return errors.New("missing organizationID")
		}

		listApps, err := api.ListGroups(&iamSDK.ListGroupsRequest{
			OrganizationID: orgID,
		})
		if err != nil {
			return fmt.Errorf("failed to list groups: %w", err)
		}
		for _, group := range listApps.Groups {
			if !acctest.IsTestResource(group.Name) {
				continue
			}
			err = api.DeleteGroup(&iamSDK.DeleteGroupRequest{
				GroupID: group.ID,
			})
			if err != nil {
				return fmt.Errorf("failed to delete group: %w", err)
			}
		}
		return nil
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

func testSweepIamAPIKey(_ string) error {
	return acctest.Sweep(func(scwClient *scw.Client) error {
		api := iamSDK.NewAPI(scwClient)

		logging.L.Debugf("sweeper: destroying the api keys")

		orgID, exists := scwClient.GetDefaultOrganizationID()
		if !exists {
			return errors.New("missing organizationID")
		}

		listAPIKeys, err := api.ListAPIKeys(&iamSDK.ListAPIKeysRequest{
			OrganizationID: &orgID,
		}, scw.WithAllPages())
		if err != nil {
			return fmt.Errorf("failed to list api keys: %w", err)
		}
		for _, key := range listAPIKeys.APIKeys {
			if !acctest.IsTestResource(key.Description) {
				continue
			}
			err = api.DeleteAPIKey(&iamSDK.DeleteAPIKeyRequest{
				AccessKey: key.AccessKey,
			})
			if err != nil {
				return fmt.Errorf("failed to delete api key: %w", err)
			}
		}
		return nil
	})
}
