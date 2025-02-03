package mnqtestfuncs

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	mnqSDK "github.com/scaleway/scaleway-sdk-go/api/mnq/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/api/mnq/v1beta1/sweepers"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
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
	return acctest.SweepRegions((&mnqSDK.SqsAPI{}).Regions(), sweepers.SweepSQSCredentials)
}

func testSweepSQS(_ string) error {
	return acctest.SweepRegions((&mnqSDK.SqsAPI{}).Regions(), sweepers.SweepSQS)
}

func testSweepSNSCredentials(_ string) error {
	return acctest.SweepRegions((&mnqSDK.SnsAPI{}).Regions(), sweepers.SweepSNSCredentials)
}

func testSweepSNS(_ string) error {
	return acctest.SweepRegions((&mnqSDK.SnsAPI{}).Regions(), sweepers.SweepSNS)
}

func testSweepNatsAccount(_ string) error {
	return acctest.SweepRegions((&mnqSDK.NatsAPI{}).Regions(), sweepers.SweepNatsAccount)
}
