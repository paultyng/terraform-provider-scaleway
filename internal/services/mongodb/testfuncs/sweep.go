package mongodbtestfuncs

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	mongodb "github.com/scaleway/scaleway-sdk-go/api/mongodb/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/logging"
)

func AddTestSweepers() {
	resource.AddTestSweepers("scaleway_mongodb_instance", &resource.Sweeper{
		Name: "scaleway_mongodb_instance",
		F:    testSweepMongodbInstance,
	})
}

func testSweepMongodbInstance(_ string) error {
	return acctest.SweepZones(scw.AllZones, func(scwClient *scw.Client, zone scw.Zone) error {
		mongodbAPI := mongodb.NewAPI(scwClient)
		logging.L.Debugf("sweeper: destroying the mongodb instance in (%s)", zone)
		extractRegion, err := zone.Region()
		if err != nil {
			return fmt.Errorf("error extract region in (%s) in sweeper: %w", zone, err)
		}
		listInstance, err := mongodbAPI.ListInstances(&mongodb.ListInstancesRequest{
			Region: extractRegion,
		})
		if err != nil {
			return fmt.Errorf("error listing mongodb instance in (%s) in sweeper: %w", zone, err)
		}

		for _, instance := range listInstance.Instances {
			_, err := mongodbAPI.DeleteInstance(&mongodb.DeleteInstanceRequest{
				Region:     extractRegion,
				InstanceID: instance.ID,
			})
			if err != nil {
				return fmt.Errorf("error deleting mongodb instance in sweeper: %w", err)
			}
		}

		return nil
	})
}
