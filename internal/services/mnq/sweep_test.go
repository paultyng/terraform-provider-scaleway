package mnq_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	mnqtestfuncs "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/mnq/testfuncs"
)

func init() {
	mnqtestfuncs.AddTestSweepers()
}

func TestMain(m *testing.M) {
	resource.TestMain(m)
}
