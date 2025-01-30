package vpcgw_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	vpcgwtestfuncs "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/vpcgw/testfuncs"
)

func init() {
	vpcgwtestfuncs.AddTestSweepers()
}

func TestMain(m *testing.M) {
	resource.TestMain(m)
}
