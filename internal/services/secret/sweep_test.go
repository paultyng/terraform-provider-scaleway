package secret_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	secrettestfuncs "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/secret/testfuncs"
)

func init() {
	secrettestfuncs.AddTestSweepers()
}

func TestMain(m *testing.M) {
	resource.TestMain(m)
}
