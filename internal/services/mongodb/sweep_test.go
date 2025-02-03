package mongodb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	mongodbtestfuncs "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/mongodb/testfuncs"
)

func init() {
	mongodbtestfuncs.AddTestSweepers()
}

func TestMain(m *testing.M) {
	resource.TestMain(m)
}
