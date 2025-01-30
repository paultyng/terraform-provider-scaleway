package baremetal_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	baremetaltestfuncs "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/baremetal/testfuncs"
)

func init() {
	baremetaltestfuncs.AddTestSweepers()
}

func TestMain(m *testing.M) {
	resource.TestMain(m)
}
