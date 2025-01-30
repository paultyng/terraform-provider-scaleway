package flexibleip_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	flexibleiptestfuncs "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/flexibleip/testfuncs"
)

func init() {
	flexibleiptestfuncs.AddTestSweepers()
}

func TestMain(m *testing.M) {
	resource.TestMain(m)
}
