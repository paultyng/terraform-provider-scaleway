package registrytestfuncs

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	registrySDK "github.com/scaleway/scaleway-sdk-go/api/registry/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/logging"
)

func AddTestSweepers() {
	resource.AddTestSweepers("scaleway_registry_namespace", &resource.Sweeper{
		Name: "scaleway_registry_namespace",
		F:    testSweepNamespace,
	})
}

func testSweepNamespace(_ string) error {
	return acctest.SweepRegions([]scw.Region{scw.RegionFrPar, scw.RegionNlAms}, func(scwClient *scw.Client, region scw.Region) error {
		registryAPI := registrySDK.NewAPI(scwClient)
		logging.L.Debugf("sweeper: destroying the registry namespaces in (%s)", region)
		listNamespaces, err := registryAPI.ListNamespaces(
			&registrySDK.ListNamespacesRequest{Region: region}, scw.WithAllPages())
		if err != nil {
			return fmt.Errorf("error listing namespaces in (%s) in sweeper: %s", region, err)
		}

		for _, ns := range listNamespaces.Namespaces {
			_, err := registryAPI.DeleteNamespace(&registrySDK.DeleteNamespaceRequest{
				NamespaceID: ns.ID,
				Region:      region,
			})
			if err != nil {
				logging.L.Debugf("sweeper: error (%s)", err)

				return fmt.Errorf("error deleting namespace in sweeper: %s", err)
			}
		}

		return nil
	})
}
