package jobs_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	jobstestfuncs "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/jobs/testfuncs"
)

func init() {
	jobstestfuncs.AddTestSweepers()
}

func TestMain(m *testing.M) {
	resource.TestMain(m)
}
