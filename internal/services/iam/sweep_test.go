package iam_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	iamtestfuncs "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/iam/testfuncs"
)

func init() {
	iamtestfuncs.AddTestSweepers()
}

func TestMain(m *testing.M) {
	resource.TestMain(m)
}
