package object_test

import (
	"fmt"
	"testing"

	awspolicy "github.com/hashicorp/awspolicyequivalence"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	objectchecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/object/testfuncs"
)

func TestAccDataSourceObjectBucketPolicy_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	bucketName := sdkacctest.RandomWithPrefix("tf-tests-scw-obp-data-basic")

	expectedPolicyText := `{
	"Version":"2012-10-17",
	"Id":"MyPolicy",
	"Statement": [
		{
			"Sid":"GrantToEveryone",
			"Effect":"Allow",
			"Principal":{
				"SCW":"*"
			},
			"Action":[
				"s3:ListBucket",
				"s3:GetObject"
			]
		}
  ]
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      objectchecks.IsBucketDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "scaleway_object_bucket" "main" {
						name = "%[1]s"
						region = "%[2]s"
					}

					resource "scaleway_object_bucket_policy" "main" {
						bucket = scaleway_object_bucket.main.id
						region = "%[2]s"
						policy = jsonencode(
							{
								Id = "MyPolicy"
								Statement = [
									{
										Action = [
											"s3:ListBucket",
											"s3:GetObject",
										]
										Effect = "Allow"
										Principal = {
											SCW = "*"
										}
										Resource  = [
											"${scaleway_object_bucket.main.name}",
											"${scaleway_object_bucket.main.name}/*",
										]
										Sid = "GrantToEveryone"
									},
								]
								Version = "2012-10-17"
							}
						)
					}

					data "scaleway_object_bucket_policy" "selected" {
						bucket = scaleway_object_bucket_policy.main.bucket
					}
				`, bucketName, objectTestsMainRegion),
				Check: resource.ComposeTestCheckFunc(
					objectchecks.CheckBucketExists(tt, "scaleway_object_bucket.main", true),
					resource.TestCheckResourceAttr("data.scaleway_object_bucket_policy.selected", "bucket", objectTestsMainRegion+"/"+bucketName),
					resource.TestCheckResourceAttrSet("data.scaleway_object_bucket_policy.selected", "policy"),
					testAccCheckDataSourcePolicyIsEquivalent("data.scaleway_object_bucket_policy.selected", expectedPolicyText),
				),
				ExpectNonEmptyPlan: !*acctest.UpdateCassettes,
			},
		},
	})
}

func testAccCheckDataSourcePolicyIsEquivalent(n, expectedPolicyText string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		dataSourcePolicy := ds.Primary.Attributes["policy"]

		dataSourcePolicyToCompare, err := removePolicyStatementResources(dataSourcePolicy)
		if err != nil {
			return err
		}

		equivalent, err := awspolicy.PoliciesAreEquivalent(expectedPolicyText, dataSourcePolicyToCompare)
		if err != nil {
			return fmt.Errorf("error testing policy equivalence: %s", err)
		}
		if !equivalent {
			return fmt.Errorf("non equivalent policy error:\n\nexpected: %s\n\n     got: %s",
				expectedPolicyText, dataSourcePolicyToCompare)
		}

		return nil
	}
}
