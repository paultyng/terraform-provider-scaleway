package redis_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	redistestfuncs "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/redis/testfuncs"
)

func init() {
	redistestfuncs.AddTestSweepers()
}

func TestMain(m *testing.M) {
	resource.TestMain(m)
}
