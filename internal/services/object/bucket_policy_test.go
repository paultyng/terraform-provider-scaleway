package object_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	awspolicy "github.com/hashicorp/awspolicyequivalence"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/object"
	objectchecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/object/testfuncs"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/types"
)

func TestAccObjectBucketPolicy_Basic(t *testing.T) {
	bucketName := sdkacctest.RandomWithPrefix("tf-tests-scw-obp-basic")

	tfConfig := fmt.Sprintf(`
		resource "scaleway_object_bucket" "bucket" {
			name = %[1]q
			region = %[2]q
			tags = {
				TestName = "TestAccScalewayObjectBucketPolicy_Basic"
			}
		}

		resource "scaleway_object_bucket_policy" "bucket" {
			bucket = scaleway_object_bucket.bucket.id
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
								"%[1]s",
								"%[1]s/*",
							]
							Sid = "GrantToEveryone"
						},
					]
					Version = "2012-10-17"
				}
			)
		}`, bucketName, objectTestsMainRegion)

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

	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        object.ErrorCheck(t, EndpointsID),
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      objectchecks.IsBucketDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: tfConfig,
				Check: resource.ComposeTestCheckFunc(
					objectchecks.CheckBucketExists(tt, "scaleway_object_bucket.bucket", true),
					resource.TestCheckResourceAttrPair("scaleway_object_bucket_policy.bucket", "region", "scaleway_object_bucket.bucket", "region"),
					testAccCheckBucketHasPolicy(tt, "scaleway_object_bucket.bucket", expectedPolicyText),
				),
				ExpectNonEmptyPlan: !*acctest.UpdateCassettes,
			},
			{
				ResourceName: "scaleway_object_bucket_policy.bucket",
				ImportState:  true,
			},
			{
				Config:             tfConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: !*acctest.UpdateCassettes,
			},
		},
	})
}

func TestAccObjectBucketPolicy_OtherRegionWithBucketID(t *testing.T) {
	bucketName := sdkacctest.RandomWithPrefix("tf-tests-scw-obp-with-bucket-id")

	tfConfig := fmt.Sprintf(`
		resource "scaleway_object_bucket" "bucket" {
			name = %[1]q
			region = %[2]q
			tags = {
				TestName = "TestAccScalewayObjectBucketPolicy_OtherRegionWithBucketID"
			}
		}

		resource "scaleway_object_bucket_policy" "bucket" {
			bucket = scaleway_object_bucket.bucket.id
			policy = jsonencode(
				{
					Id = "MyPolicy"
					Statement = [
						{
							Action = [
								"s3:*"
							]
							Effect = "Allow"
							Principal = {
								SCW = "*"
							}
							Resource  = [
								"%[1]s",
								"%[1]s/*",
							]
							Sid = "GrantToEveryone"
						},
					]
					Version = "2023-04-17"
				}
			)
		}`, bucketName, objectTestsSecondaryRegion)

	expectedPolicyText := `{
	"Version":"2023-04-17",
	"Id":"MyPolicy",
	"Statement": [
		{
			"Sid":"GrantToEveryone",
			"Effect":"Allow",
			"Principal":{
				"SCW":"*"
			},
			"Action":[
				"s3:*"
			]
		}
  ]
}`

	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        object.ErrorCheck(t, EndpointsID),
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      objectchecks.IsBucketDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: tfConfig,
				Check: resource.ComposeTestCheckFunc(
					objectchecks.CheckBucketExists(tt, "scaleway_object_bucket.bucket", true),
					resource.TestCheckResourceAttrPair("scaleway_object_bucket_policy.bucket", "region", "scaleway_object_bucket.bucket", "region"),
					testAccCheckBucketHasPolicy(tt, "scaleway_object_bucket.bucket", expectedPolicyText),
				),
				ExpectNonEmptyPlan: !*acctest.UpdateCassettes,
			},
			{
				ResourceName: "scaleway_object_bucket_policy.bucket",
				ImportState:  true,
			},
			{
				Config:             tfConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: !*acctest.UpdateCassettes,
			},
		},
	})
}

func TestAccObjectBucketPolicy_OtherRegionWithBucketName(t *testing.T) {
	bucketName := sdkacctest.RandomWithPrefix("tf-tests-scw-obp-with-bucket-name")

	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        object.ErrorCheck(t, EndpointsID),
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      objectchecks.IsBucketDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "scaleway_object_bucket" "bucket" {
						name = %[1]q
						region = %[2]q
						tags = {
							TestName = "TestAccScalewayObjectBucketPolicy_OtherRegionWithBucketName"
						}
					}

					resource "scaleway_object_bucket_policy" "bucket" {
						bucket = scaleway_object_bucket.bucket.name
						policy = jsonencode(
							{
								Id = "MyPolicy"
								Statement = [
									{
										Action = [
											"s3:*"
										]
										Effect = "Allow"
										Principal = {
											SCW = "*"
										}
										Resource  = [
											"%[1]s",
											"%[1]s/*",
										]
										Sid = "GrantToEveryone"
									},
								]
								Version = "2023-04-17"
							}
						)
					}`, bucketName, objectTestsSecondaryRegion),
				ExpectError: regexp.MustCompile(`error putting SCW bucket policy:.*NoSuchBucket.*The specified bucket does not exist`),
			},
		},
	})
}

func testAccCheckBucketHasPolicy(tt *acctest.TestTools, n string, expectedPolicyText string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ctx := context.Background()
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		bucketRegion := rs.Primary.Attributes["region"]
		s3Client, err := object.NewS3ClientFromMeta(ctx, tt.Meta, bucketRegion)
		if err != nil {
			return err
		}

		if rs.Primary.ID == "" {
			return errors.New("no ID is set")
		}

		bucketName := rs.Primary.Attributes["name"]
		policy, err := s3Client.GetBucketPolicy(ctx, &s3.GetBucketPolicyInput{
			Bucket: types.ExpandStringPtr(bucketName),
		})
		if err != nil {
			return fmt.Errorf("GetBucketPolicy error: %v", err)
		}

		actualPolicyText := *policy.Policy
		actualPolicyText, err = removePolicyStatementResources(actualPolicyText)
		if err != nil {
			return err
		}

		equivalent, err := awspolicy.PoliciesAreEquivalent(actualPolicyText, expectedPolicyText)
		if err != nil {
			return fmt.Errorf("error testing policy equivalence: %s", err)
		}
		if !equivalent {
			return fmt.Errorf("non equivalent policy error:\n\nexpected: %s\n\n     got: %s",
				expectedPolicyText, actualPolicyText)
		}

		return nil
	}
}

// remove the following:
//
//	policy["Statement"][i]["Resource"]
func removePolicyStatementResources(policy string) (string, error) {
	actualPolicyJSON := make(map[string]interface{})
	err := json.Unmarshal([]byte(policy), &actualPolicyJSON)
	if err != nil {
		return "", fmt.Errorf("json.Unmarshal error: %v", err)
	}

	if statement, ok := actualPolicyJSON["Statement"].([]interface{}); ok && len(statement) > 0 {
		for _, rule := range statement {
			if rule, ok := rule.(map[string]interface{}); ok {
				delete(rule, "Resource")
			}
		}
	}

	actualPolicyTextBytes, err := json.Marshal(actualPolicyJSON)
	if err != nil {
		return "", fmt.Errorf("json.Marshal error: %v", err)
	}

	return string(actualPolicyTextBytes), nil
}
