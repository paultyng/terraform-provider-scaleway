package cockpit_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	cockpittestfuncs "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/cockpit/testfuncs"
)

func init() {
	cockpittestfuncs.AddTestSweepers()
}

func TestMain(m *testing.M) {
	resource.TestMain(m)
}
