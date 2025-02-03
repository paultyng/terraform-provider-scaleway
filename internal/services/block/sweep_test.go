package block_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	blocktestfuncs "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/block/testfuncs"
)

func init() {
	blocktestfuncs.AddTestSweepers()
}

func TestMain(m *testing.M) {
	resource.TestMain(m)
}
