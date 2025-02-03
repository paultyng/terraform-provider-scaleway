package mnqtestfuncs

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	accountSDK "github.com/scaleway/scaleway-sdk-go/api/account/v3"
	mnqSDK "github.com/scaleway/scaleway-sdk-go/api/mnq/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/logging"
)

func AddTestSweepers() {
	resource.AddTestSweepers("scaleway_mnq_nats_account", &resource.Sweeper{
		Name: "scaleway_mnq_nats_account",
		F:    testSweepNatsAccount,
	})
	resource.AddTestSweepers("scaleway_mnq_sns", &resource.Sweeper{
		Name: "scaleway_mnq_sns",
		F:    testSweepSNS,
	})
	resource.AddTestSweepers("scaleway_mnq_sns_credentials", &resource.Sweeper{
		Name: "scaleway_mnq_sns_credentials",
		F:    testSweepSNSCredentials,
	})
	resource.AddTestSweepers("scaleway_mnq_sqs", &resource.Sweeper{
		Name: "scaleway_mnq_sqs",
		F:    testSweepSQS,
	})
	resource.AddTestSweepers("scaleway_mnq_sqs_credentials", &resource.Sweeper{
		Name: "scaleway_mnq_sqs_credentials",
		F:    testSweepSQSCredentials,
	})
}

func testSweepSQSCredentials(_ string) error {
	return acctest.SweepRegions((&mnqSDK.SqsAPI{}).Regions(), func(scwClient *scw.Client, region scw.Region) error {
		mnqAPI := mnqSDK.NewSqsAPI(scwClient)
		logging.L.Debugf("sweeper: destroying the mnq sqs credentials in (%s)", region)
		listSqsCredentials, err := mnqAPI.ListSqsCredentials(
			&mnqSDK.SqsAPIListSqsCredentialsRequest{
				Region: region,
			}, scw.WithAllPages())
		if err != nil {
			return fmt.Errorf("error listing sqs credentials in (%s) in sweeper: %s", region, err)
		}

		for _, credentials := range listSqsCredentials.SqsCredentials {
			err := mnqAPI.DeleteSqsCredentials(&mnqSDK.SqsAPIDeleteSqsCredentialsRequest{
				SqsCredentialsID: credentials.ID,
				Region:           region,
			})
			if err != nil {
				logging.L.Debugf("sweeper: error (%s)", err)

				return fmt.Errorf("error deleting sqs credentials in sweeper: %s", err)
			}
		}

		return nil
	})
}

func testSweepSQS(_ string) error {
	return acctest.SweepRegions((&mnqSDK.SqsAPI{}).Regions(), func(scwClient *scw.Client, region scw.Region) error {
		accountAPI := accountSDK.NewProjectAPI(scwClient)
		mnqAPI := mnqSDK.NewSqsAPI(scwClient)

		logging.L.Debugf("sweeper: destroying the mnq sqss in (%s)", region)

		listProjects, err := accountAPI.ListProjects(&accountSDK.ProjectAPIListProjectsRequest{}, scw.WithAllPages())
		if err != nil {
			return fmt.Errorf("failed to list projects: %w", err)
		}
		for _, project := range listProjects.Projects {
			if !strings.HasPrefix(project.Name, "tf_tests") {
				continue
			}

			_, err := mnqAPI.DeactivateSqs(&mnqSDK.SqsAPIDeactivateSqsRequest{
				Region:    region,
				ProjectID: project.ID,
			})
			if err != nil {
				logging.L.Debugf("sweeper: error (%s)", err)
				return err
			}
		}

		return nil
	})
}

func testSweepSNSCredentials(_ string) error {
	return acctest.SweepRegions((&mnqSDK.SnsAPI{}).Regions(), func(scwClient *scw.Client, region scw.Region) error {
		mnqAPI := mnqSDK.NewSnsAPI(scwClient)
		logging.L.Debugf("sweeper: destroying the mnq sns credentials in (%s)", region)
		listSnsCredentials, err := mnqAPI.ListSnsCredentials(
			&mnqSDK.SnsAPIListSnsCredentialsRequest{
				Region: region,
			}, scw.WithAllPages())
		if err != nil {
			return fmt.Errorf("error listing sns credentials in (%s) in sweeper: %s", region, err)
		}

		for _, credentials := range listSnsCredentials.SnsCredentials {
			err := mnqAPI.DeleteSnsCredentials(&mnqSDK.SnsAPIDeleteSnsCredentialsRequest{
				SnsCredentialsID: credentials.ID,
				Region:           region,
			})
			if err != nil {
				logging.L.Debugf("sweeper: error (%s)", err)

				return fmt.Errorf("error deleting sns credentials in sweeper: %s", err)
			}
		}

		return nil
	})
}

func testSweepSNS(_ string) error {
	return acctest.SweepRegions((&mnqSDK.SnsAPI{}).Regions(), func(scwClient *scw.Client, region scw.Region) error {
		accountAPI := accountSDK.NewProjectAPI(scwClient)
		mnqAPI := mnqSDK.NewSnsAPI(scwClient)

		logging.L.Debugf("sweeper: destroying the mnq sns in (%s)", region)

		listProjects, err := accountAPI.ListProjects(&accountSDK.ProjectAPIListProjectsRequest{}, scw.WithAllPages())
		if err != nil {
			return fmt.Errorf("failed to list projects: %w", err)
		}
		for _, project := range listProjects.Projects {
			if !strings.HasPrefix(project.Name, "tf_tests") {
				continue
			}

			_, err := mnqAPI.DeactivateSns(&mnqSDK.SnsAPIDeactivateSnsRequest{
				Region:    region,
				ProjectID: project.ID,
			})
			if err != nil {
				logging.L.Debugf("sweeper: error (%s)", err)
				return err
			}
		}

		return nil
	})
}

func testSweepNatsAccount(_ string) error {
	return acctest.SweepRegions((&mnqSDK.NatsAPI{}).Regions(), func(scwClient *scw.Client, region scw.Region) error {
		mnqAPI := mnqSDK.NewNatsAPI(scwClient)
		logging.L.Debugf("sweeper: destroying the mnq nats accounts in (%s)", region)
		listNatsAccounts, err := mnqAPI.ListNatsAccounts(
			&mnqSDK.NatsAPIListNatsAccountsRequest{
				Region: region,
			}, scw.WithAllPages())
		if err != nil {
			return fmt.Errorf("error listing nats account in (%s) in sweeper: %s", region, err)
		}

		for _, account := range listNatsAccounts.NatsAccounts {
			err := mnqAPI.DeleteNatsAccount(&mnqSDK.NatsAPIDeleteNatsAccountRequest{
				NatsAccountID: account.ID,
				Region:        region,
			})
			if err != nil {
				logging.L.Debugf("sweeper: error (%s)", err)

				return fmt.Errorf("error deleting nats account in sweeper: %s", err)
			}
		}

		return nil
	})
}
