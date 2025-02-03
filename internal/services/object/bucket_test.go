package object_test

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3Types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/locality/regional"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/object"
	objectchecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/object/testfuncs"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/transport"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/types"
)

func TestAccObjectBucket_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	testBucketACL := "private"
	testBucketUpdatedACL := "public-read"
	bucketBasic := sdkacctest.RandomWithPrefix("tf-tests-scaleway-object-bucket-basic")
	bucketMainRegion := sdkacctest.RandomWithPrefix("tf-tests-scaleway-object-bucket-main-region")
	bucketSecondary := sdkacctest.RandomWithPrefix("tf-tests-scaleway-object-bucket-secondary")
	objectBucketTestMainRegion := scw.RegionFrPar
	objectBucketTestSecondaryRegion := scw.RegionNlAms
	objectBucketTestDefaultRegion, _ := tt.Meta.ScwClient().GetDefaultRegion()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      objectchecks.IsBucketDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "scaleway_object_bucket" "base-01" {
						name = "%[1]s"
						tags = {
							foo = "bar"
						}
					}

					resource "scaleway_object_bucket" "secondary-bucket-01" {
						name = "%[2]s"
						region = "%[4]s"
						tags = {
							foo = "bar"
							baz = "qux"
						}
					}

					resource "scaleway_object_bucket" "main-bucket-01" {
						name = "%[3]s"
						region = "%[5]s"
					}
				`, bucketBasic, bucketSecondary, bucketMainRegion, objectBucketTestSecondaryRegion, objectBucketTestMainRegion),
				Check: resource.ComposeTestCheckFunc(
					objectchecks.CheckBucketExists(tt, "scaleway_object_bucket.base-01", true),
					objectchecks.CheckBucketExists(tt, "scaleway_object_bucket.secondary-bucket-01", true),
					objectchecks.CheckBucketExists(tt, "scaleway_object_bucket.main-bucket-01", true),

					resource.TestCheckResourceAttr("scaleway_object_bucket.base-01", "name", bucketBasic),
					resource.TestCheckResourceAttr("scaleway_object_bucket.base-01", "acl", testBucketACL),
					resource.TestCheckResourceAttr("scaleway_object_bucket.base-01", "tags.%", "1"),
					resource.TestCheckResourceAttr("scaleway_object_bucket.base-01", "tags.foo", "bar"),
					resource.TestCheckResourceAttr("scaleway_object_bucket.base-01", "endpoint", fmt.Sprintf("https://%s.s3.%s.scw.cloud", bucketBasic, objectBucketTestDefaultRegion)),
					resource.TestCheckResourceAttr("scaleway_object_bucket.base-01", "api_endpoint", fmt.Sprintf("https://s3.%s.scw.cloud", objectBucketTestDefaultRegion)),

					resource.TestCheckResourceAttr("scaleway_object_bucket.secondary-bucket-01", "name", bucketSecondary),
					resource.TestCheckResourceAttr("scaleway_object_bucket.secondary-bucket-01", "tags.%", "2"),
					resource.TestCheckResourceAttr("scaleway_object_bucket.secondary-bucket-01", "tags.foo", "bar"),
					resource.TestCheckResourceAttr("scaleway_object_bucket.secondary-bucket-01", "tags.baz", "qux"),
					resource.TestCheckResourceAttr("scaleway_object_bucket.secondary-bucket-01", "endpoint", fmt.Sprintf("https://%s.s3.%s.scw.cloud", bucketSecondary, objectBucketTestSecondaryRegion)),
					resource.TestCheckResourceAttr("scaleway_object_bucket.secondary-bucket-01", "api_endpoint", fmt.Sprintf("https://s3.%s.scw.cloud", objectBucketTestSecondaryRegion)),

					resource.TestCheckResourceAttr("scaleway_object_bucket.main-bucket-01", "name", bucketMainRegion),
					resource.TestCheckResourceAttr("scaleway_object_bucket.main-bucket-01", "endpoint", fmt.Sprintf("https://%s.s3.%s.scw.cloud", bucketMainRegion, objectBucketTestMainRegion)),
					resource.TestCheckResourceAttr("scaleway_object_bucket.main-bucket-01", "api_endpoint", fmt.Sprintf("https://s3.%s.scw.cloud", objectBucketTestMainRegion)),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "scaleway_object_bucket" "base-01" {
						name = "%[1]s"
						acl = "%[2]s"
					}
			
					resource "scaleway_object_bucket" "secondary-bucket-01" {
						name = "%[3]s"
						region = "%[4]s"
						tags = {
							foo = "bar"
						}
					}
				`, bucketBasic, testBucketUpdatedACL, bucketSecondary, objectTestsSecondaryRegion),
				Check: resource.ComposeTestCheckFunc(
					objectchecks.CheckBucketExists(tt, "scaleway_object_bucket.base-01", true),
					objectchecks.CheckBucketExists(tt, "scaleway_object_bucket.secondary-bucket-01", true),

					resource.TestCheckResourceAttr("scaleway_object_bucket.base-01", "name", bucketBasic),
					resource.TestCheckResourceAttr("scaleway_object_bucket.base-01", "acl", testBucketUpdatedACL),
					resource.TestCheckResourceAttr("scaleway_object_bucket.base-01", "tags.%", "0"),

					resource.TestCheckResourceAttr("scaleway_object_bucket.secondary-bucket-01", "name", bucketSecondary),
					resource.TestCheckResourceAttr("scaleway_object_bucket.secondary-bucket-01", "tags.%", "1"),
					resource.TestCheckResourceAttr("scaleway_object_bucket.secondary-bucket-01", "tags.foo", "bar"),
				),
			},
		},
	})
}

func TestAccObjectBucket_Lifecycle(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	bucketLifecycle := sdkacctest.RandomWithPrefix("tf-tests-scaleway-object-bucket-lifecycle")
	resourceNameLifecycle := "scaleway_object_bucket.main-bucket-lifecycle"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      objectchecks.IsBucketDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "scaleway_object_bucket" "main-bucket-lifecycle"{
						name = "%s"
						region = "%s"
						acl = "private"
			
						lifecycle_rule {
							id      = "id1"
							prefix  = "path1/"
							enabled = true
							expiration {
						  		days = 365
							}
							transition {
								days          = 30
								storage_class = "GLACIER"
							}
						}
					}
				`, bucketLifecycle, objectTestsMainRegion),
				Check: resource.ComposeTestCheckFunc(
					objectchecks.CheckBucketExists(tt, "scaleway_object_bucket.main-bucket-lifecycle", true),
					testAccCheckObjectBucketLifecycleConfigurationExists(tt, resourceNameLifecycle),
					resource.TestCheckResourceAttr("scaleway_object_bucket.main-bucket-lifecycle", "name", bucketLifecycle),
					resource.TestCheckResourceAttr(resourceNameLifecycle, "lifecycle_rule.0.id", "id1"),
					resource.TestCheckResourceAttr(resourceNameLifecycle, "lifecycle_rule.0.prefix", "path1/"),
					resource.TestCheckResourceAttr(resourceNameLifecycle, "lifecycle_rule.0.expiration.0.days", "365"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceNameLifecycle, "lifecycle_rule.0.transition.*", map[string]string{
						"days":          "30",
						"storage_class": "GLACIER",
					}),
				),
			},
			{
				Config: fmt.Sprintf(`
						resource "scaleway_object_bucket" "main-bucket-lifecycle"{
							name = "%s"
							region = "%s"
							acl = "private"
			
							lifecycle_rule {
								id      = "id1"
								prefix  = "path1/"
								enabled = true
								expiration {
								 	days = 365
								}
								transition {
								  	days          = 90
								  	storage_class = "ONEZONE_IA"
								}
							}
						}
					`, bucketLifecycle, objectTestsMainRegion),
				Check: resource.ComposeTestCheckFunc(
					objectchecks.CheckBucketExists(tt, "scaleway_object_bucket.main-bucket-lifecycle", true),
					testAccCheckObjectBucketLifecycleConfigurationExists(tt, resourceNameLifecycle),
					resource.TestCheckResourceAttr("scaleway_object_bucket.main-bucket-lifecycle", "name", bucketLifecycle),
					resource.TestCheckResourceAttr(resourceNameLifecycle, "lifecycle_rule.0.id", "id1"),
					resource.TestCheckResourceAttr(resourceNameLifecycle, "lifecycle_rule.0.prefix", "path1/"),
					resource.TestCheckResourceAttr(resourceNameLifecycle, "lifecycle_rule.0.expiration.0.days", "365"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceNameLifecycle, "lifecycle_rule.0.transition.*", map[string]string{
						"days":          "90",
						"storage_class": "ONEZONE_IA",
					}),
				),
			},
			{
				Config: fmt.Sprintf(`
						resource "scaleway_object_bucket" "main-bucket-lifecycle"{
							name = "%s"
							region = "%s"
							acl = "private"
			
							lifecycle_rule {
								id      = "id1"
								prefix  = "path1/"
								enabled = true
								expiration {
								  	days = 365
								}
								transition {
								  	days          = 120
								  	storage_class = "GLACIER"
								}
							}
			
							lifecycle_rule {
								id      = "id2"
								prefix  = "path2/"
								enabled = true
								expiration {
									days = "50"
								}
							}
			
							lifecycle_rule {
								id      = "id3"
								prefix  = "path3/"
								enabled = true
								tags = {
								  	"tagKey"    = "tagValue"
								  	"terraform" = "hashicorp"
								}
								expiration {
								  	days = "1"
								}
							}
			
							lifecycle_rule {
								id      = "id4"
								enabled = true
								tags = {
								  	"tagKey"    = "tagValue"
								  	"terraform" = "hashicorp"
								}
								transition {
								  	days          = 1
								  	storage_class = "GLACIER"
								}
							}
			
							lifecycle_rule {
								id      = "id5"
								enabled = true
								tags = {
								  	"tagKey" = "tagValue"
								}
								transition {
								  	days          = 1
								  	storage_class = "GLACIER"
								}
							}
						}
					`, bucketLifecycle, objectTestsMainRegion),
				Check: resource.ComposeTestCheckFunc(
					objectchecks.CheckBucketExists(tt, "scaleway_object_bucket.main-bucket-lifecycle", true),
					testAccCheckObjectBucketLifecycleConfigurationExists(tt, resourceNameLifecycle),
					resource.TestCheckResourceAttr("scaleway_object_bucket.main-bucket-lifecycle", "name", bucketLifecycle),
					resource.TestCheckResourceAttr(resourceNameLifecycle, "lifecycle_rule.0.id", "id1"),
					resource.TestCheckResourceAttr(resourceNameLifecycle, "lifecycle_rule.0.prefix", "path1/"),
					resource.TestCheckResourceAttr(resourceNameLifecycle, "lifecycle_rule.0.expiration.0.days", "365"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceNameLifecycle, "lifecycle_rule.0.transition.*", map[string]string{
						"days":          "120",
						"storage_class": "GLACIER",
					}),
					resource.TestCheckResourceAttr(resourceNameLifecycle, "lifecycle_rule.1.id", "id2"),
					resource.TestCheckResourceAttr(resourceNameLifecycle, "lifecycle_rule.1.prefix", "path2/"),
					resource.TestCheckResourceAttr(resourceNameLifecycle, "lifecycle_rule.1.expiration.0.days", "50"),
					resource.TestCheckResourceAttr(resourceNameLifecycle, "lifecycle_rule.2.id", "id3"),
					resource.TestCheckResourceAttr(resourceNameLifecycle, "lifecycle_rule.2.prefix", "path3/"),
					resource.TestCheckResourceAttr(resourceNameLifecycle, "lifecycle_rule.2.tags.tagKey", "tagValue"),
					resource.TestCheckResourceAttr(resourceNameLifecycle, "lifecycle_rule.2.tags.terraform", "hashicorp"),
					resource.TestCheckResourceAttr(resourceNameLifecycle, "lifecycle_rule.3.id", "id4"),
					resource.TestCheckResourceAttr(resourceNameLifecycle, "lifecycle_rule.3.tags.tagKey", "tagValue"),
					resource.TestCheckResourceAttr(resourceNameLifecycle, "lifecycle_rule.3.tags.terraform", "hashicorp"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceNameLifecycle, "lifecycle_rule.3.transition.*", map[string]string{
						"days":          "1",
						"storage_class": "GLACIER",
					}),
					resource.TestCheckResourceAttr(resourceNameLifecycle, "lifecycle_rule.4.id", "id5"),
					resource.TestCheckResourceAttr(resourceNameLifecycle, "lifecycle_rule.4.tags.tagKey", "tagValue"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceNameLifecycle, "lifecycle_rule.4.transition.*", map[string]string{
						"days":          "1",
						"storage_class": "GLACIER",
					}),
				),
			},
			{
				Config: fmt.Sprintf(`
						resource "scaleway_object_bucket" "main-bucket-lifecycle"{
							name = "%s"
							region = "%s"
							acl = "private"
			
							lifecycle_rule {
								id      = "id1"
								prefix  = "path1/"
								enabled = true
								abort_incomplete_multipart_upload_days = 30
							}
						}
					`, bucketLifecycle, objectTestsMainRegion),
				Check: resource.ComposeTestCheckFunc(
					objectchecks.CheckBucketExists(tt, "scaleway_object_bucket.main-bucket-lifecycle", true),
					testAccCheckObjectBucketLifecycleConfigurationExists(tt, resourceNameLifecycle),
					resource.TestCheckResourceAttr("scaleway_object_bucket.main-bucket-lifecycle", "name", bucketLifecycle),
					resource.TestCheckResourceAttr(resourceNameLifecycle, "lifecycle_rule.0.id", "id1"),
					resource.TestCheckResourceAttr(resourceNameLifecycle, "lifecycle_rule.0.prefix", "path1/"),
					resource.TestCheckResourceAttr(resourceNameLifecycle, "lifecycle_rule.0.abort_incomplete_multipart_upload_days", "30"),
				),
			},
			{
				Config: fmt.Sprintf(`
						resource "scaleway_object_bucket" "main-bucket-lifecycle"{
							name = "%s"
							region = "%s"
							acl = "private"
			
							lifecycle_rule {
								prefix  = "path1/"
								enabled = true
								abort_incomplete_multipart_upload_days = 30
							}
						}
					`, bucketLifecycle, objectTestsMainRegion),
				Check: resource.ComposeTestCheckFunc(
					objectchecks.CheckBucketExists(tt, "scaleway_object_bucket.main-bucket-lifecycle", true),
					testAccCheckObjectBucketLifecycleConfigurationExists(tt, resourceNameLifecycle),
					resource.TestCheckResourceAttrSet(resourceNameLifecycle, "lifecycle_rule.0.id"),
					resource.TestCheckResourceAttr(resourceNameLifecycle, "lifecycle_rule.0.abort_incomplete_multipart_upload_days", "30"),
				),
			},
			{
				Config: fmt.Sprintf(`
							resource "scaleway_object_bucket" "main-bucket-lifecycle"{
								name = "%s"
								region = "%s"
								acl = "private"
			
								lifecycle_rule {
									prefix  = "path1/"
									enabled = true
									tags    = {
										"deleted" = "true"
									}
									expiration {
										days = 1
									}
								}
							}
						`, bucketLifecycle, objectTestsMainRegion),
				Check: resource.ComposeTestCheckFunc(
					objectchecks.CheckBucketExists(tt, "scaleway_object_bucket.main-bucket-lifecycle", true),
					testAccCheckObjectBucketLifecycleConfigurationExists(tt, resourceNameLifecycle),
					resource.TestCheckResourceAttrSet(resourceNameLifecycle, "lifecycle_rule.0.id"),
					resource.TestCheckResourceAttr(resourceNameLifecycle, "lifecycle_rule.0.tags.deleted", "true"),
					resource.TestCheckResourceAttr(resourceNameLifecycle, "lifecycle_rule.0.prefix", "path1/"),
					resource.TestCheckResourceAttr(resourceNameLifecycle, "lifecycle_rule.0.expiration.0.days", "1"),
				),
			},
			{
				Config: fmt.Sprintf(`
						resource "scaleway_object_bucket" "main-bucket-lifecycle" {
							name                = "%s"
							region = "%s"
							object_lock_enabled = true

							lifecycle_rule {
								enabled = true
								prefix  = ""
								expiration {
									days = 2
								}
							}
			
							lifecycle_rule {
								enabled = true
								abort_incomplete_multipart_upload_days = 30
							}
						}`, bucketLifecycle, objectTestsMainRegion),
				Check: resource.ComposeTestCheckFunc(
					objectchecks.CheckBucketExists(tt, "scaleway_object_bucket.main-bucket-lifecycle", true),
					testAccCheckObjectBucketLifecycleConfigurationExists(tt, resourceNameLifecycle),
					resource.TestCheckResourceAttrSet(resourceNameLifecycle, "lifecycle_rule.0.id"),
					resource.TestCheckResourceAttrSet(resourceNameLifecycle, "lifecycle_rule.1.id"),
					resource.TestCheckResourceAttr(resourceNameLifecycle, "lifecycle_rule.0.prefix", ""),
					resource.TestCheckResourceAttr(resourceNameLifecycle, "lifecycle_rule.0.expiration.0.days", "2"),
					resource.TestCheckResourceAttr(resourceNameLifecycle, "lifecycle_rule.1.abort_incomplete_multipart_upload_days", "30"),
				),
			},
		},
	})
}

func TestAccObjectBucket_ObjectLock(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	bucketObjectLock := sdkacctest.RandomWithPrefix("tf-tests-scaleway-object-bucket-lock")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      objectchecks.IsBucketDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "scaleway_object_bucket" "object-locked-bucket"{
						name = "%s"
						region = "%s"

						object_lock_enabled = true
					}
				`, bucketObjectLock, objectTestsMainRegion),
				Check: resource.ComposeTestCheckFunc(
					objectchecks.CheckBucketExists(tt, "scaleway_object_bucket.object-locked-bucket", true),
					resource.TestCheckResourceAttr("scaleway_object_bucket.object-locked-bucket", "name", bucketObjectLock),
					resource.TestCheckResourceAttr("scaleway_object_bucket.object-locked-bucket", "object_lock_enabled", "true"),
					resource.TestCheckResourceAttr("scaleway_object_bucket.object-locked-bucket", "versioning.0.enabled", "true"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "scaleway_object_bucket" "object-locked-bucket"{
						name = "%s"
						region = "%s"

						object_lock_enabled = true

						versioning {
							enabled = true
						}
					}
				`, bucketObjectLock, objectTestsMainRegion),
				Check: resource.ComposeTestCheckFunc(
					objectchecks.CheckBucketExists(tt, "scaleway_object_bucket.object-locked-bucket", true),
					resource.TestCheckResourceAttr("scaleway_object_bucket.object-locked-bucket", "name", bucketObjectLock),
					resource.TestCheckResourceAttr("scaleway_object_bucket.object-locked-bucket", "object_lock_enabled", "true"),
					resource.TestCheckResourceAttr("scaleway_object_bucket.object-locked-bucket", "versioning.0.enabled", "true"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "scaleway_object_bucket" "object-locked-bucket"{
						name = "%s"
						region = "%s"

						object_lock_enabled = true

						versioning {
							enabled = false
						}
					}
				`, bucketObjectLock, objectTestsMainRegion),
				ExpectError: regexp.MustCompile("versioning must be enabled when object lock is enabled"),
			},
		},
	})
}

func TestAccObjectBucket_Cors_Update(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resourceName := "scaleway_object_bucket.bucket-cors-update"
	bucketName := sdkacctest.RandomWithPrefix("tf-tests-scaleway-object-bucket-cors-update")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      objectchecks.IsBucketDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "scaleway_object_bucket" "bucket-cors-update" {
						name = %[1]q
						region = %[2]q
						cors_rule {
							allowed_headers = ["*"]
							allowed_methods = ["PUT", "POST"]
							allowed_origins = ["https://www.example.com"]
							expose_headers  = ["x-amz-server-side-encryption", "ETag"]
							max_age_seconds = 3000
						}
					}`, bucketName, objectTestsMainRegion),
				Check: resource.ComposeTestCheckFunc(
					objectchecks.CheckBucketExists(tt, resourceName, true),
					testAccCheckObjectBucketCors(tt,
						resourceName,
						[]s3Types.CORSRule{
							{
								AllowedHeaders: []string{"*"},
								AllowedMethods: []string{"PUT", "POST"},
								AllowedOrigins: []string{"https://www.example.com"},
								ExposeHeaders:  []string{"x-amz-server-side-encryption", "ETag"},
								MaxAgeSeconds:  scw.Int32Ptr(3000),
							},
						},
					),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "scaleway_object_bucket" "bucket-cors-update" {
						name = %[1]q
						region = %[2]q
						cors_rule {
							allowed_headers = ["*"]
							allowed_methods = ["PUT", "POST", "GET"]
							allowed_origins = ["https://www.example.com"]
							expose_headers  = ["x-amz-server-side-encryption", "ETag"]
							max_age_seconds = 3000
						}
					}`, bucketName, objectTestsMainRegion),
				Check: resource.ComposeTestCheckFunc(
					objectchecks.CheckBucketExists(tt, resourceName, true),
					testAccCheckObjectBucketCors(tt,
						resourceName,
						[]s3Types.CORSRule{
							{
								AllowedHeaders: []string{"*"},
								AllowedMethods: []string{"PUT", "POST", "GET"},
								AllowedOrigins: []string{"https://www.example.com"},
								ExposeHeaders:  []string{"x-amz-server-side-encryption", "ETag"},
								MaxAgeSeconds:  scw.Int32Ptr(3000),
							},
						},
					),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "scaleway_object_bucket" "bucket-cors-update" {
						name = %[1]q
						region = %[2]q
						cors_rule {
							allowed_headers = ["*"]
							allowed_methods = ["PUT", "POST"]
							allowed_origins = ["https://www.example.com"]
							expose_headers  = ["x-amz-server-side-encryption", "ETag"]
							max_age_seconds = 3000
						}
					}`, bucketName, objectTestsMainRegion),
				Check: resource.ComposeTestCheckFunc(
					objectchecks.CheckBucketExists(tt, resourceName, true),
					testAccCheckObjectBucketCors(tt,
						resourceName,
						[]s3Types.CORSRule{
							{
								AllowedHeaders: []string{"*"},
								AllowedMethods: []string{"PUT", "POST"},
								AllowedOrigins: []string{"https://www.example.com"},
								ExposeHeaders:  []string{"x-amz-server-side-encryption", "ETag"},
								MaxAgeSeconds:  scw.Int32Ptr(3000),
							},
						},
					),
				),
			},
		},
	})
}

func TestAccObjectBucket_Cors_Delete(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	ctx := context.Background()

	resourceName := "scaleway_object_bucket.bucket"
	bucketName := sdkacctest.RandomWithPrefix("tf-tests-scaleway-object-bucket-cors-delete")
	deleteBucketCors := func(tt *acctest.TestTools, n string) resource.TestCheckFunc {
		return func(s *terraform.State) error {
			rs, ok := s.RootModule().Resources[n]
			if !ok {
				return fmt.Errorf("not found: %s", n)
			}
			bucketRegion := rs.Primary.Attributes["region"]
			conn, err := object.NewS3ClientFromMeta(ctx, tt.Meta, bucketRegion)
			if err != nil {
				return err
			}
			_, err = conn.DeleteBucketCors(ctx, &s3.DeleteBucketCorsInput{
				Bucket: scw.StringPtr(rs.Primary.Attributes["name"]),
			})
			if err != nil && !object.IsS3Err(err, object.ErrCodeNoSuchCORSConfiguration, "") {
				return err
			}
			return nil
		}
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      objectchecks.IsBucketDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
						resource "scaleway_object_bucket" "bucket" {
							name = %[1]q
							region = %[2]q
							cors_rule {
								allowed_headers = ["*"]
								allowed_methods = ["PUT", "POST"]
								allowed_origins = ["https://www.example.com"]
								expose_headers  = ["x-amz-server-side-encryption", "ETag"]
								max_age_seconds = 3000
							}
						}`, bucketName, objectTestsMainRegion),
				Check: resource.ComposeTestCheckFunc(
					objectchecks.CheckBucketExists(tt, resourceName, true),
					deleteBucketCors(tt, resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccObjectBucket_Cors_EmptyOrigin(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	bucketName := sdkacctest.RandomWithPrefix("tf-tests-scaleway-object-cors-empty-origin")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      objectchecks.IsBucketDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
						resource "scaleway_object_bucket" "bucket" {
							name = %[1]q
							region = %[2]q
							cors_rule {
								allowed_headers = ["*"]
								allowed_methods = ["PUT", "POST"]
								allowed_origins = [""]
								expose_headers  = ["x-amz-server-side-encryption", "ETag"]
								max_age_seconds = 3000
							}
						}`, bucketName, objectTestsMainRegion),
				ExpectError: regexp.MustCompile("error putting S3 CORS"),
			},
		},
	})
}

func testAccCheckObjectBucketCors(tt *acctest.TestTools, n string, corsRules []s3Types.CORSRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ctx := context.Background()

		rs := s.RootModule().Resources[n]
		bucketName := rs.Primary.Attributes["name"]
		bucketRegion := rs.Primary.Attributes["region"]
		s3Client, err := object.NewS3ClientFromMeta(ctx, tt.Meta, bucketRegion)
		if err != nil {
			return err
		}

		_, err = s3Client.HeadBucket(ctx, &s3.HeadBucketInput{
			Bucket: scw.StringPtr(bucketName),
		})
		if err != nil {
			return err
		}

		out, err := s3Client.GetBucketCors(ctx, &s3.GetBucketCorsInput{
			Bucket: scw.StringPtr(bucketName),
		})
		if err != nil {
			if !object.IsS3Err(err, object.ErrCodeNoSuchCORSConfiguration, "") {
				return fmt.Errorf("GetBucketCors error: %v", err)
			}
		}

		if out == nil {
			return errors.New("CORS Rules nil")
		}

		if !reflect.DeepEqual(out.CORSRules, corsRules) {
			return fmt.Errorf("bad error cors rule, expected: %v, got %v", corsRules, out.CORSRules)
		}

		return nil
	}
}

func TestAccObjectBucket_DestroyForce(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resourceName := "scaleway_object_bucket.bucket"
	bucketName := sdkacctest.RandomWithPrefix("tf-tests-scaleway-object-bucket-force")

	addObjectToBucket := func(tt *acctest.TestTools, n string) resource.TestCheckFunc {
		return func(s *terraform.State) error {
			ctx := context.Background()
			rs, ok := s.RootModule().Resources[n]
			if !ok {
				return fmt.Errorf("not found: %s", n)
			}
			bucketRegion := rs.Primary.Attributes["region"]
			conn, err := object.NewS3ClientFromMeta(ctx, tt.Meta, bucketRegion)
			if err != nil {
				return err
			}

			req := s3.PutObjectInput{
				Bucket: scw.StringPtr(rs.Primary.Attributes["name"]),
				Key:    scw.StringPtr("test-file"),
				Body:   strings.NewReader("test content"),
			}

			_, err = conn.PutObject(ctx, &req)
			if err != nil {
				return fmt.Errorf("failed to put object in test bucket: %s", err)
			}

			_, err = conn.PutObject(ctx, &s3.PutObjectInput{
				Bucket: scw.StringPtr(rs.Primary.Attributes["name"]),
				Key:    scw.StringPtr("folder/test-file-in-folder"),
				Body:   strings.NewReader("folder test content"), // Example body content
			})
			if err != nil {
				return fmt.Errorf("failed to put object in test bucket sub folder: %s", err)
			}
			return nil
		}
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      objectchecks.IsBucketDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
						resource "scaleway_object_bucket" "bucket" {
							name = %[1]q
							region = %[2]q
							force_destroy = true
							versioning {
								enabled = true
							}
						}`, bucketName, objectTestsMainRegion),
				Check: resource.ComposeTestCheckFunc(
					objectchecks.CheckBucketExists(tt, resourceName, true),
					addObjectToBucket(tt, resourceName),
				),
			},
			{
				Config: fmt.Sprintf(`
						resource "scaleway_object_bucket" "bucket" {
							name = %[1]q
							region = %[2]q
							force_destroy = true
						}`, bucketName, objectTestsMainRegion),
				Check: resource.ComposeTestCheckFunc(
					objectchecks.CheckBucketExists(tt, resourceName, true),
					addObjectToBucket(tt, resourceName),
				),
			},
		},
	})
}

func testAccCheckObjectBucketLifecycleConfigurationExists(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ctx := context.Background()
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no ID is set")
		}

		bucketRegion := rs.Primary.Attributes["region"]
		s3Client, err := object.NewS3ClientFromMeta(ctx, tt.Meta, bucketRegion)
		if err != nil {
			return err
		}

		bucketRegionalID := regional.ExpandID(rs.Primary.ID)

		input := &s3.GetBucketLifecycleConfigurationInput{
			Bucket: types.ExpandStringPtr(bucketRegionalID.ID),
		}

		_, err = s3Client.GetBucketLifecycleConfiguration(ctx, input)
		if err != nil {
			if errors.Is(err, transport.ErrRetryWhenTimeout) {
				return fmt.Errorf("object Storage Bucket Replication Configuration for bucket (%s) not found", rs.Primary.ID)
			}
			return err
		}

		return nil
	}
}
