package acctest_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

func TestIsTestResource(t *testing.T) {
	assert.True(t, acctest.IsTestResource("tf_tests_mnq_sqs_queue_default_project"))
}
