package k8s_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	k8stestfuncs "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/k8s/testfuncs"
)

func init() {
	k8stestfuncs.AddTestSweepers()
}

func TestMain(m *testing.M) {
	resource.TestMain(m)
}
