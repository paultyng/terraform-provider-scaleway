package functiontestfuncs

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	functionSDK "github.com/scaleway/scaleway-sdk-go/api/function/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/logging"
)

func AddTestSweepers() {
	resource.AddTestSweepers("scaleway_function_cron", &resource.Sweeper{
		Name: "scaleway_function_cron",
		F:    testSweepFunctionCron,
	})
	resource.AddTestSweepers("scaleway_function", &resource.Sweeper{
		Name: "scaleway_function",
		F:    testSweepFunction,
	})
	resource.AddTestSweepers("scaleway_function_namespace", &resource.Sweeper{
		Name: "scaleway_function_namespace",
		F:    testSweepFunctionNamespace,
	})
	resource.AddTestSweepers("scaleway_function_trigger", &resource.Sweeper{
		Name: "scaleway_function_trigger",
		F:    testSweepFunctionTrigger,
	})
}

func testSweepFunctionTrigger(_ string) error {
	return acctest.SweepRegions((&functionSDK.API{}).Regions(), func(scwClient *scw.Client, region scw.Region) error {
		functionAPI := functionSDK.NewAPI(scwClient)
		logging.L.Debugf("sweeper: destroying the function triggers in (%s)", region)
		listTriggers, err := functionAPI.ListTriggers(
			&functionSDK.ListTriggersRequest{
				Region: region,
			}, scw.WithAllPages())
		if err != nil {
			return fmt.Errorf("error listing trigger in (%s) in sweeper: %s", region, err)
		}

		for _, trigger := range listTriggers.Triggers {
			_, err := functionAPI.DeleteTrigger(&functionSDK.DeleteTriggerRequest{
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

func testSweepFunctionNamespace(_ string) error {
	return acctest.SweepRegions([]scw.Region{scw.RegionFrPar}, func(scwClient *scw.Client, region scw.Region) error {
		functionAPI := functionSDK.NewAPI(scwClient)
		logging.L.Debugf("sweeper: destroying the function namespaces in (%s)", region)
		listNamespaces, err := functionAPI.ListNamespaces(
			&functionSDK.ListNamespacesRequest{
				Region: region,
			}, scw.WithAllPages())
		if err != nil {
			return fmt.Errorf("error listing namespaces in (%s) in sweeper: %s", region, err)
		}

		for _, ns := range listNamespaces.Namespaces {
			_, err := functionAPI.DeleteNamespace(&functionSDK.DeleteNamespaceRequest{
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

func testSweepFunction(_ string) error {
	return acctest.SweepRegions([]scw.Region{scw.RegionFrPar}, func(scwClient *scw.Client, region scw.Region) error {
		functionAPI := functionSDK.NewAPI(scwClient)
		logging.L.Debugf("sweeper: destroying the function in (%s)", region)
		listFunctions, err := functionAPI.ListFunctions(
			&functionSDK.ListFunctionsRequest{
				Region: region,
			}, scw.WithAllPages())
		if err != nil {
			return fmt.Errorf("error listing functions in (%s) in sweeper: %s", region, err)
		}

		for _, f := range listFunctions.Functions {
			_, err := functionAPI.DeleteFunction(&functionSDK.DeleteFunctionRequest{
				FunctionID: f.ID,
				Region:     region,
			})
			if err != nil && !httperrors.Is404(err) {
				logging.L.Debugf("sweeper: error (%s)", err)

				return fmt.Errorf("error deleting functions in sweeper: %s", err)
			}
		}

		return nil
	})
}

func testSweepFunctionCron(_ string) error {
	return acctest.SweepRegions([]scw.Region{scw.RegionFrPar}, func(scwClient *scw.Client, region scw.Region) error {
		functionAPI := functionSDK.NewAPI(scwClient)
		logging.L.Debugf("sweeper: destroying the function cron in (%s)", region)
		listCron, err := functionAPI.ListCrons(
			&functionSDK.ListCronsRequest{
				Region: region,
			}, scw.WithAllPages())
		if err != nil {
			return fmt.Errorf("error listing cron in (%s) in sweeper: %s", region, err)
		}

		for _, cron := range listCron.Crons {
			_, err := functionAPI.DeleteCron(&functionSDK.DeleteCronRequest{
				CronID: cron.ID,
				Region: region,
			})
			if err != nil {
				logging.L.Debugf("sweeper: error (%s)", err)

				return fmt.Errorf("error deleting cron in sweeper: %s", err)
			}
		}

		return nil
	})
}
