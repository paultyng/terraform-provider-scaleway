package object_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	objecttestfuncs "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/object/testfuncs"
)

func init() {
	objecttestfuncs.AddTestSweepers()
}

func TestMain(m *testing.M) {
	resource.TestMain(m)
}
