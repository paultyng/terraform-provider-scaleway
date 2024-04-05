package mnq_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	mnqSDK "github.com/scaleway/scaleway-sdk-go/api/mnq/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/logging"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/mnq"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

func init() {
	resource.AddTestSweepers("scaleway_mnq_nats_account", &resource.Sweeper{
		Name: "scaleway_mnq_nats_account",
		F:    testSweepNatsAccount,
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

func TestAccNatsAccount_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isNatsAccountDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_mnq_nats_account main {
						name = "test-mnq-nats-account-basic"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isNatsAccountPresent(tt, "scaleway_mnq_nats_account.main"),
					resource.TestCheckResourceAttr("scaleway_mnq_nats_account.main", "name", "test-mnq-nats-account-basic"),
					resource.TestCheckResourceAttrSet("scaleway_mnq_nats_account.main", "endpoint"),
				),
			},
		},
	})
}

func isNatsAccountPresent(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		api, region, id, err := mnq.NewNatsAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = api.GetNatsAccount(&mnqSDK.NatsAPIGetNatsAccountRequest{
			NatsAccountID: id,
			Region:        region,
		})
		if err != nil {
			return err
		}

		return nil
	}
}

func isNatsAccountDestroyed(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_mnq_nats_account" {
				continue
			}

			api, region, id, err := mnq.NewNatsAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}

			err = api.DeleteNatsAccount(&mnqSDK.NatsAPIDeleteNatsAccountRequest{
				NatsAccountID: id,
				Region:        region,
			})

			if err == nil {
				return fmt.Errorf("mnq nats account (%s) still exists", rs.Primary.ID)
			}

			if !httperrors.Is404(err) {
				return err
			}
		}

		return nil
	}
}
