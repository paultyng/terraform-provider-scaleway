package ipam_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	ipamtestfuncs "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/ipam/testfuncs"
)

func init() {
	ipamtestfuncs.AddTestSweepers()
}

func TestMain(m *testing.M) {
	resource.TestMain(m)
}
