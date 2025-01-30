package tem_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	temtestfuncs "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/tem/testfuncs"
)

func init() {
	temtestfuncs.AddTestSweepers()
}

func TestMain(m *testing.M) {
	resource.TestMain(m)
}
