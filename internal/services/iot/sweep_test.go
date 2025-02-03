package iot_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	iottestfuncs "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/iot/testfuncs"
)

func init() {
	iottestfuncs.AddTestSweepers()
}

func TestMain(m *testing.M) {
	resource.TestMain(m)
}
