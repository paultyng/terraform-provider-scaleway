package functiontestfuncs

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	functionSDK "github.com/scaleway/scaleway-sdk-go/api/function/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/api/function/v1beta1/sweepers"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
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
	return acctest.SweepRegions((&functionSDK.API{}).Regions(), sweepers.SweepTriggers)
}

func testSweepFunctionNamespace(_ string) error {
	return acctest.SweepRegions((&functionSDK.API{}).Regions(), sweepers.SweepNamespaces)
}

func testSweepFunction(_ string) error {
	return acctest.SweepRegions((&functionSDK.API{}).Regions(), sweepers.SweepFunctions)
}

func testSweepFunctionCron(_ string) error {
	return acctest.SweepRegions((&functionSDK.API{}).Regions(), sweepers.SweepCrons)
}
