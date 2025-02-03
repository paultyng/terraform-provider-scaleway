package container_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	containertestfuncs "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/container/testfuncs"
)

func init() {
	containertestfuncs.AddTestSweepers()
}

func TestMain(m *testing.M) {
	resource.TestMain(m)
}
