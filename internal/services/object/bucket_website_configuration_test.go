package object_test

import (
	"fmt"
	"regexp"
	"testing"

	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/object"
	objectchecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/object/testfuncs"
)

const (
	ResourcePrefix   = "tf-acc-test"
	resourceTestName = "scaleway_object_bucket_website_configuration.test"
)

func TestAccObjectBucketWebsiteConfiguration_Basic(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(ResourcePrefix)
	resourceName := resourceTestName

	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        object.ErrorCheck(t, EndpointsID),
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			objectchecks.IsWebsiteConfigurationDestroyed(tt),
			objectchecks.IsBucketDestroyed(tt),
		),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
			  		resource "scaleway_object_bucket" "test" {
						name = %[1]q
						region = %[2]q
						acl  = "public-read"
						tags = {
							TestName = "TestAccScalewayObjectBucketWebsiteConfiguration_Basic"
						}
					}

				  	resource "scaleway_object_bucket_website_configuration" "test" {
						bucket = scaleway_object_bucket.test.id
						index_document {
						  suffix = "index.html"
						}
				  	}
				`, rName, objectTestsMainRegion),
				Check: resource.ComposeTestCheckFunc(
					objectchecks.CheckBucketExists(tt, "scaleway_object_bucket.test", true),
					objectchecks.IsWebsiteConfigurationPresent(tt, resourceName),
					resource.TestCheckResourceAttrPair(resourceName, "bucket", "scaleway_object_bucket.test", "name"),
					resource.TestCheckResourceAttr(resourceName, "index_document.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "index_document.0.suffix", "index.html"),
					resource.TestCheckResourceAttrSet(resourceName, "website_domain"),
					resource.TestCheckResourceAttrSet(resourceName, "website_endpoint"),
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

func TestAccObjectBucketWebsiteConfiguration_WithPolicy(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(ResourcePrefix)
	resourceName := resourceTestName

	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        object.ErrorCheck(t, EndpointsID),
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			objectchecks.IsWebsiteConfigurationDestroyed(tt),
			objectchecks.IsBucketDestroyed(tt),
		),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
			  		resource "scaleway_object_bucket" "test" {
						name = %[1]q
						region = %[2]q
						acl  = "public-read"
						tags = {
							TestName = "TestAccScalewayObjectBucketWebsiteConfiguration_WithPolicy"
						}
					}

					resource "scaleway_object_bucket_policy" "main" {
						bucket = scaleway_object_bucket.test.id
						policy = jsonencode(
						{
							"Version" = "2012-10-17",
							"Id" = "MyPolicy",
							"Statement" = [
							{
							   "Sid" = "GrantToEveryone",
							   "Effect" = "Allow",
							   "Principal" = "*",
							   "Action" = [
								  "s3:GetObject"
							   ],
							   "Resource":[
								  "%[1]s/*"
							   ]
							}
							]
						})
					}

				  	resource "scaleway_object_bucket_website_configuration" "test" {
						bucket = scaleway_object_bucket.test.id
						index_document {
						  suffix = "index.html"
						}
				  	}
				`, rName, objectTestsMainRegion),
				Check: resource.ComposeTestCheckFunc(
					objectchecks.CheckBucketExists(tt, "scaleway_object_bucket.test", true),
					objectchecks.IsWebsiteConfigurationPresent(tt, resourceName),
					resource.TestCheckResourceAttrPair(resourceName, "bucket", "scaleway_object_bucket.test", "name"),
					resource.TestCheckResourceAttr(resourceName, "index_document.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "index_document.0.suffix", "index.html"),
					resource.TestCheckResourceAttrSet(resourceName, "website_domain"),
					resource.TestCheckResourceAttrSet(resourceName, "website_endpoint"),
				),
				ExpectNonEmptyPlan: !*acctest.UpdateCassettes,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccObjectBucketWebsiteConfiguration_Update(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(ResourcePrefix)
	resourceName := resourceTestName

	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        object.ErrorCheck(t, EndpointsID),
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			objectchecks.IsWebsiteConfigurationDestroyed(tt),
			objectchecks.IsBucketDestroyed(tt),
		),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
			  		resource "scaleway_object_bucket" "test" {
						name = %[1]q
						region = %[2]q
						acl  = "public-read"
						tags = {
							TestName = "TestAccScalewayObjectBucketWebsiteConfiguration_Update"
						}
					}

				  	resource "scaleway_object_bucket_website_configuration" "test" {
						bucket = scaleway_object_bucket.test.id
						index_document {
						  suffix = "index.html"
						}
				  	}
				`, rName, objectTestsMainRegion),
				Check: resource.ComposeTestCheckFunc(
					objectchecks.CheckBucketExists(tt, "scaleway_object_bucket.test", true),
					objectchecks.IsWebsiteConfigurationPresent(tt, resourceName),
				),
			},
			{
				Config: fmt.Sprintf(`
			  		resource "scaleway_object_bucket" "test" {
						name = %[1]q
						region = %[2]q
						acl  = "public-read"
						tags = {
							TestName = "TestAccScalewayObjectBucketWebsiteConfiguration_Update"
						}
					}

				  	resource "scaleway_object_bucket_website_configuration" "test" {
						bucket = scaleway_object_bucket.test.id
						index_document {
						  suffix = "index.html"
						}

						error_document {
							key = "error.html"
						}
				  	}
				`, rName, objectTestsMainRegion),
				Check: resource.ComposeTestCheckFunc(
					objectchecks.CheckBucketExists(tt, "scaleway_object_bucket.test", true),
					objectchecks.IsWebsiteConfigurationPresent(tt, resourceName),
					resource.TestCheckResourceAttrPair(resourceName, "bucket", "scaleway_object_bucket.test", "name"),
					resource.TestCheckResourceAttr(resourceName, "index_document.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "index_document.0.suffix", "index.html"),
					resource.TestCheckResourceAttr(resourceName, "error_document.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "error_document.0.key", "error.html"),
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

func TestAccObjectBucketWebsiteConfiguration_WithBucketName(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(ResourcePrefix)
	resourceName := resourceTestName

	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        object.ErrorCheck(t, EndpointsID),
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			objectchecks.IsWebsiteConfigurationDestroyed(tt),
			objectchecks.IsBucketDestroyed(tt),
		),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
			  		resource "scaleway_object_bucket" "test" {
						name = %[1]q
						region = %[2]q
						acl  = "public-read"
					}

				  	resource "scaleway_object_bucket_website_configuration" "test" {
						bucket = scaleway_object_bucket.test.name
						index_document {
						  suffix = "index.html"
						}
				  	}
				`, rName, objectTestsMainRegion),
				ExpectError: regexp.MustCompile(`couldn't read bucket:.*NoSuchBucket.*`),
			},
			{
				Config: fmt.Sprintf(`
			  		resource "scaleway_object_bucket" "test" {
						name = %[1]q
						region = %[2]q
						acl  = "public-read"
					}

				  	resource "scaleway_object_bucket_website_configuration" "test" {
						bucket = scaleway_object_bucket.test.name
						region = %[2]q
						index_document {
						  suffix = "index.html"
						}
				  	}
				`, rName, objectTestsMainRegion),
				Check: resource.ComposeTestCheckFunc(
					objectchecks.CheckBucketExists(tt, "scaleway_object_bucket.test", true),
					objectchecks.IsWebsiteConfigurationPresent(tt, resourceName),
				),
			},
		},
	})
}
