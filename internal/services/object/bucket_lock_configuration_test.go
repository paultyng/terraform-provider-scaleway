package object_test

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/locality/regional"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/object"
	objectchecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/object/testfuncs"
)

const (
	LockResourcePrefix   = "tf-acc-test"
	lockResourceTestName = "scaleway_object_bucket_lock_configuration.test"
)

func TestAccObjectBucketLockConfiguration_Basic(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(LockResourcePrefix)
	resourceName := lockResourceTestName

	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        object.ErrorCheck(t, EndpointsID),
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckBucketLockConfigurationDestroy(tt),
			objectchecks.IsBucketDestroyed(tt),
		),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "scaleway_object_bucket" "test" {
						name = %[1]q
						region = %[2]q
						tags = {
							TestName = "TestAccSCW_LockConfig_basic"
						}
			
						object_lock_enabled = true
					}
			
					resource "scaleway_object_bucket_acl" "test" {
						bucket = scaleway_object_bucket.test.id
						acl = "public-read"
					}
			
					resource "scaleway_object_bucket_lock_configuration" "test" {
						bucket = scaleway_object_bucket.test.id
						rule {
							default_retention {
								mode = "GOVERNANCE"
								days = 1
							}
						}
					}
				`, rName, objectTestsMainRegion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketLockConfigurationExists(tt, resourceName),
					objectchecks.CheckBucketExists(tt, "scaleway_object_bucket.test", true),
					resource.TestCheckResourceAttrPair(resourceName, "bucket", "scaleway_object_bucket.test", "name"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.default_retention.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.default_retention.0.mode", "GOVERNANCE"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.default_retention.0.days", "1"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "scaleway_object_bucket" "test" {
						name = %[1]q
						region = %[2]q
						tags = {
							TestName = "TestAccSCW_LockConfig_basic"
						}
			
						object_lock_enabled = true
					}
			
					resource "scaleway_object_bucket_acl" "test" {
						bucket = scaleway_object_bucket.test.name
						acl = "public-read"
					}
			
					resource "scaleway_object_bucket_lock_configuration" "test" {
						bucket = scaleway_object_bucket.test.name
						rule {
							default_retention {
								mode = "GOVERNANCE"
								years = 1
							}
						}
					}
				`, rName, objectTestsMainRegion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketLockConfigurationExists(tt, resourceName),
					objectchecks.CheckBucketExists(tt, "scaleway_object_bucket.test", true),
					resource.TestCheckResourceAttrPair(resourceName, "bucket", "scaleway_object_bucket.test", "name"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.default_retention.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.default_retention.0.mode", "GOVERNANCE"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.default_retention.0.years", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccObjectBucketLockConfiguration_Update(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(LockResourcePrefix)
	resourceName := lockResourceTestName

	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        object.ErrorCheck(t, EndpointsID),
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckBucketLockConfigurationDestroy(tt),
			objectchecks.IsBucketDestroyed(tt),
		),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "scaleway_object_bucket" "test" {
						name = %[1]q
						region = %[2]q
						tags = {
							TestName = "TestAccSCW_LockConfig_update"
						}

						object_lock_enabled = true
					}

					resource "scaleway_object_bucket_acl" "test" {
						bucket = scaleway_object_bucket.test.id
						acl = "public-read"
					}

				  	resource "scaleway_object_bucket_lock_configuration" "test" {
						bucket = scaleway_object_bucket.test.id
						rule {
							default_retention {
								mode = "GOVERNANCE"
								days = 1
							}
						}
				  	}
				`, rName, objectTestsMainRegion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketLockConfigurationExists(tt, resourceName),
					objectchecks.CheckBucketExists(tt, "scaleway_object_bucket.test", true),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "scaleway_object_bucket" "test" {
						name = %[1]q
						region = %[2]q
						tags = {
							TestName = "TestAccSCW_LockConfig_basic"
						}

						object_lock_enabled = true
					}

					resource "scaleway_object_bucket_acl" "test" {
						bucket = scaleway_object_bucket.test.id
						acl = "public-read"
					}

				  	resource "scaleway_object_bucket_lock_configuration" "test" {
						bucket = scaleway_object_bucket.test.id
						rule {
							default_retention {
								mode = "COMPLIANCE"
								days = 2
							}
						}
				  	}
				`, rName, objectTestsMainRegion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketLockConfigurationExists(tt, resourceName),
					objectchecks.CheckBucketExists(tt, "scaleway_object_bucket.test", true),
					resource.TestCheckResourceAttrPair(resourceName, "bucket", "scaleway_object_bucket.test", "name"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.default_retention.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.default_retention.0.mode", "COMPLIANCE"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.default_retention.0.days", "2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccObjectBucketLockConfiguration_WithBucketName(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(LockResourcePrefix)
	resourceName := lockResourceTestName

	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        object.ErrorCheck(t, EndpointsID),
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckBucketLockConfigurationDestroy(tt),
			objectchecks.IsBucketDestroyed(tt),
		),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "scaleway_object_bucket" "test" {
						name = %[1]q
						region = %[2]q
						tags = {
							TestName = "TestAccSCW_LockConfig_WithBucketName"
						}

						object_lock_enabled = true
					}

					resource "scaleway_object_bucket_acl" "test" {
						bucket = scaleway_object_bucket.test.id
						acl = "public-read"
					}

					resource "scaleway_object_bucket_lock_configuration" "test" {
						bucket = scaleway_object_bucket.test.name
						rule {
							default_retention {
								mode = "GOVERNANCE"
								days = 1
							}
						}
					}
				`, rName, objectTestsMainRegion),
				ExpectError: regexp.MustCompile("NoSuchBucket: The specified bucket does not exist"),
			},
			{
				Config: fmt.Sprintf(`
					resource "scaleway_object_bucket" "test" {
						name = %[1]q
						region = %[2]q
						tags = {
							TestName = "TestAccSCW_LockConfig_WithBucketName"
						}

						object_lock_enabled = true
					}

					resource "scaleway_object_bucket_acl" "test" {
						bucket = scaleway_object_bucket.test.id
						acl = "public-read"
					}

					resource "scaleway_object_bucket_lock_configuration" "test" {
						bucket = scaleway_object_bucket.test.name
						region = %[2]q
						rule {
							default_retention {
								mode = "GOVERNANCE"
								days = 1
							}
						}
					}
				`, rName, objectTestsMainRegion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketLockConfigurationExists(tt, resourceName),
					objectchecks.CheckBucketExists(tt, "scaleway_object_bucket.test", true),
					resource.TestCheckResourceAttrPair(resourceName, "bucket", "scaleway_object_bucket.test", "name"),
				),
			},
		},
	})
}

func testAccCheckBucketLockConfigurationDestroy(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ctx := context.Background()
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "scaleway_object_bucket_lock_configuration" {
				continue
			}

			regionalID := regional.ExpandID(rs.Primary.ID)
			bucketRegion := regionalID.Region
			bucket := regionalID.ID
			conn, err := object.NewS3ClientFromMeta(ctx, tt.Meta, bucketRegion.String())
			if err != nil {
				return err
			}

			input := &s3.GetObjectLockConfigurationInput{
				Bucket: aws.String(bucket),
			}

			output, err := conn.GetObjectLockConfiguration(ctx, input)

			if object.IsS3Err(err, object.ErrCodeNoSuchBucket, "The specified bucket does not exist") {
				continue
			}

			if err != nil {
				return fmt.Errorf("error getting object bucket lock configuration (%s): %w", rs.Primary.ID, err)
			}

			if output != nil {
				return fmt.Errorf("object bucket lock configuration (%s) still exists", rs.Primary.ID)
			}
		}

		return nil
	}
}

func testAccCheckBucketLockConfigurationExists(tt *acctest.TestTools, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ctx := context.Background()
		rs := s.RootModule().Resources[resourceName]
		if rs == nil {
			return errors.New("resource not found")
		}

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource (%s) ID not set", resourceName)
		}

		regionalID := regional.ExpandID(rs.Primary.ID)
		bucketRegion := regionalID.Region
		bucket := regionalID.ID
		conn, err := object.NewS3ClientFromMeta(ctx, tt.Meta, bucketRegion.String())
		if err != nil {
			return err
		}

		input := &s3.GetObjectLockConfigurationInput{
			Bucket: aws.String(bucket),
		}

		output, err := conn.GetObjectLockConfiguration(ctx, input)
		if err != nil {
			return fmt.Errorf("error getting object bucket lock configuration (%s): %w", rs.Primary.ID, err)
		}

		if output == nil {
			return fmt.Errorf("object bucket lock configuration (%s) not found", rs.Primary.ID)
		}

		return nil
	}
}
