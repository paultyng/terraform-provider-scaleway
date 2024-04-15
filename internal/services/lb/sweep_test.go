package lb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	lbtestfuncs "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/lb/testfuncs"
)

func init() {
	lbtestfuncs.AddTestSweepers()
}

func TestMain(m *testing.M) {
	resource.TestMain(m)
}
