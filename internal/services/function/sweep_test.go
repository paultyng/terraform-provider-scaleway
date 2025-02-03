package function_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	functiontestfuncs "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/function/testfuncs"
)

func init() {
	functiontestfuncs.AddTestSweepers()
}

func TestMain(m *testing.M) {
	resource.TestMain(m)
}
