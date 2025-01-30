package registry_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	registrytestfuncs "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/registry/testfuncs"
)

func init() {
	registrytestfuncs.AddTestSweepers()
}

func TestMain(m *testing.M) {
	resource.TestMain(m)
}
