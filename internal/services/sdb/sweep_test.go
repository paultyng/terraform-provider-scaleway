package sdb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	sdbtestfuncs "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/sdb/testfuncs"
)

func init() {
	sdbtestfuncs.AddTestSweepers()
}

func TestMain(m *testing.M) {
	resource.TestMain(m)
}
