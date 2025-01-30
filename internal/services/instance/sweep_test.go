package instance_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	instancetestfuncs "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/instance/testfuncs"
)

func init() {
	instancetestfuncs.AddTestSweepers()
}

func TestMain(m *testing.M) {
	resource.TestMain(m)
}
