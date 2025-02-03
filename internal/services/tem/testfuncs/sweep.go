package temtestfuncs

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	temSDK "github.com/scaleway/scaleway-sdk-go/api/tem/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/logging"
)

func AddTestSweepers() {
	resource.AddTestSweepers("scaleway_tem_domain", &resource.Sweeper{
		Name: "scaleway_tem_domain",
		F:    testSweepDomain,
	})
}

func testSweepDomain(_ string) error {
	return acctest.SweepRegions([]scw.Region{scw.RegionFrPar, scw.RegionNlAms}, func(scwClient *scw.Client, region scw.Region) error {
		temAPI := temSDK.NewAPI(scwClient)
		logging.L.Debugf("sweeper: revoking the tem domains in (%s)", region)

		listDomains, err := temAPI.ListDomains(&temSDK.ListDomainsRequest{Region: region}, scw.WithAllPages())
		if err != nil {
			return fmt.Errorf("error listing domains in (%s) in sweeper: %s", region, err)
		}

		for _, ns := range listDomains.Domains {
			if ns.Name == "test.scaleway-terraform.com" {
				logging.L.Debugf("sweeper: skipping deletion of domain %s", ns.Name)
				continue
			}
			_, err := temAPI.RevokeDomain(&temSDK.RevokeDomainRequest{
				DomainID: ns.ID,
				Region:   region,
			})
			if err != nil {
				logging.L.Debugf("sweeper: error (%s)", err)

				return fmt.Errorf("error revoking domain in sweeper: %s", err)
			}
		}

		return nil
	})
}
