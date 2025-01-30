package webhosting_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	webhostingtestfuncs "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/webhosting/testfuncs"
)

func init() {
	webhostingtestfuncs.AddTestSweepers()
}

func TestMain(m *testing.M) {
	resource.TestMain(m)
}
