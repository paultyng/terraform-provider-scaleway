package account_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	accounttestfuncs "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/account/testfuncs"
)

func init() {
	accounttestfuncs.AddTestSweepers()
}

func TestMain(m *testing.M) {
	resource.TestMain(m)
}
