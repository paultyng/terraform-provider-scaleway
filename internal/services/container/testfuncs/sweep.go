package containertestfuncs

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	containerSDK "github.com/scaleway/scaleway-sdk-go/api/container/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/logging"
)

func AddTestSweepers() {
	resource.AddTestSweepers("scaleway_container_namespace", &resource.Sweeper{
		Name:         "scaleway_container_namespace",
		F:            testSweepNamespace,
		Dependencies: []string{"scaleway_container"},
	})
	resource.AddTestSweepers("scaleway_container", &resource.Sweeper{
		Name: "scaleway_container",
		F:    testSweepContainer,
	})
	resource.AddTestSweepers("scaleway_container_trigger", &resource.Sweeper{
		Name: "scaleway_container_trigger",
		F:    testSweepTrigger,
	})
}

func testSweepTrigger(_ string) error {
	return acctest.SweepRegions((&containerSDK.API{}).Regions(), func(scwClient *scw.Client, region scw.Region) error {
		containerAPI := containerSDK.NewAPI(scwClient)
		logging.L.Debugf("sweeper: destroying the container triggers in (%s)", region)
		listTriggers, err := containerAPI.ListTriggers(
			&containerSDK.ListTriggersRequest{
				Region: region,
			}, scw.WithAllPages())
		if err != nil {
			return fmt.Errorf("error listing trigger in (%s) in sweeper: %s", region, err)
		}

		for _, trigger := range listTriggers.Triggers {
			_, err := containerAPI.DeleteTrigger(&containerSDK.DeleteTriggerRequest{
				TriggerID: trigger.ID,
				Region:    region,
			})
			if err != nil {
				logging.L.Debugf("sweeper: error (%s)", err)

				return fmt.Errorf("error deleting trigger in sweeper: %s", err)
			}
		}

		return nil
	})
}

func testSweepContainer(_ string) error {
	return acctest.SweepRegions(scw.AllRegions, func(scwClient *scw.Client, region scw.Region) error {
		containerAPI := containerSDK.NewAPI(scwClient)
		logging.L.Debugf("sweeper: destroying the container in (%s)", region)
		listNamespaces, err := containerAPI.ListContainers(
			&containerSDK.ListContainersRequest{
				Region: region,
			}, scw.WithAllPages())
		if err != nil {
			return fmt.Errorf("error listing containers in (%s) in sweeper: %s", region, err)
		}

		for _, cont := range listNamespaces.Containers {
			_, err := containerAPI.DeleteContainer(&containerSDK.DeleteContainerRequest{
				ContainerID: cont.ID,
				Region:      region,
			})
			if err != nil {
				logging.L.Debugf("sweeper: error (%s)", err)

				return fmt.Errorf("error deleting container in sweeper: %s", err)
			}
		}

		return nil
	})
}

func testSweepNamespace(_ string) error {
	return acctest.SweepRegions([]scw.Region{scw.RegionFrPar}, func(scwClient *scw.Client, region scw.Region) error {
		containerAPI := containerSDK.NewAPI(scwClient)
		logging.L.Debugf("sweeper: destroying the container namespaces in (%s)", region)
		listNamespaces, err := containerAPI.ListNamespaces(
			&containerSDK.ListNamespacesRequest{
				Region: region,
			}, scw.WithAllPages())
		if err != nil {
			return fmt.Errorf("error listing namespaces in (%s) in sweeper: %s", region, err)
		}

		for _, ns := range listNamespaces.Namespaces {
			_, err := containerAPI.DeleteNamespace(&containerSDK.DeleteNamespaceRequest{
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
