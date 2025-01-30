package object_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	objectchecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/object/testfuncs"
	"github.com/stretchr/testify/require"
)

func TestAccDataSourceObjectBucket_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	bucketName := sdkacctest.RandomWithPrefix("tf-tests-scaleway-object-bucket")
	objectBucketTestDefaultRegion, _ := tt.Meta.ScwClient().GetDefaultRegion()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      objectchecks.IsBucketDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				resource "scaleway_object_bucket" "base-01" {
					name = "%s"
					region = "%s"
					tags = {
						foo = "bar"
					}
				}

				data "scaleway_object_bucket" "by-id" {
					name = scaleway_object_bucket.base-01.id
				}
				`, bucketName, objectTestsMainRegion),
				Check: resource.ComposeTestCheckFunc(
					objectchecks.CheckBucketExists(tt, "scaleway_object_bucket.base-01", true),
					resource.TestCheckResourceAttr("data.scaleway_object_bucket.by-id", "name", bucketName),
				),
			},
			{
				Config: fmt.Sprintf(`
				resource "scaleway_object_bucket" "base-01" {
					name = "%s"
					region = "%s"
					tags = {
						foo = "bar"
					}
				}

				data "scaleway_object_bucket" "by-name" {
					name = scaleway_object_bucket.base-01.name
				}
				`, bucketName, objectBucketTestDefaultRegion),
				Check: resource.ComposeTestCheckFunc(
					objectchecks.CheckBucketExists(tt, "scaleway_object_bucket.base-01", true),
					resource.TestCheckResourceAttr("data.scaleway_object_bucket.by-name", "name", bucketName),
				),
			},
			{
				Config: fmt.Sprintf(`
				resource "scaleway_object_bucket" "base-01" {
					name = "%s"
					region = "%s"
					tags = {
						foo = "bar"
					}
				}

				data "scaleway_object_bucket" "by-name" {
					name = scaleway_object_bucket.base-01.name
				}
				`, bucketName, objectTestsMainRegion),
				Check: resource.ComposeTestCheckFunc(
					objectchecks.CheckBucketExists(tt, "scaleway_object_bucket.base-01", true),
					resource.TestCheckResourceAttr("data.scaleway_object_bucket.by-name", "name", bucketName),
				),
				ExpectError: regexp.MustCompile("failed getting Object Storage bucket"),
			},
		},
	})
}

func TestAccDataSourceObjectBucket_ProjectIDAllowed(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	bucketName := sdkacctest.RandomWithPrefix("tf-tests-scaleway-object-bucket")

	project, iamAPIKey, terminateFakeSideProject, err := acctest.CreateFakeSideProject(tt)
	require.NoError(t, err)

	ctx := context.Background()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: acctest.FakeSideProjectProviders(ctx, tt, project, iamAPIKey),
		CheckDestroy: resource.ComposeAggregateTestCheckFunc(
			func(_ *terraform.State) error {
				return terminateFakeSideProject()
			},
			objectchecks.IsBucketDestroyed(tt),
		),
		Steps: []resource.TestStep{
			// Create a bucket from the main provider into the side project and read it from the side provider
			// The side provider should only be able to read the bucket from the side project
			{
				Config: fmt.Sprintf(`
					resource "scaleway_object_bucket" "base" {
						name = "%[1]s"
						project_id = "%[2]s"
						region = "%[3]s"
					}

					data "scaleway_object_bucket" "selected" {
						name = scaleway_object_bucket.base.id
						provider = side
					}
				`,
					bucketName,
					project.ID,
					objectTestsMainRegion,
				),
				Check: resource.ComposeTestCheckFunc(
					objectchecks.CheckBucketExists(tt, "scaleway_object_bucket.base", false),
					resource.TestCheckResourceAttr("data.scaleway_object_bucket.selected", "name", bucketName),
					resource.TestCheckResourceAttr("data.scaleway_object_bucket.selected", "project_id", project.ID),
				),
			},
		},
	})
}

func TestAccDataSourceObjectBucket_ProjectIDForbidden(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	bucketName := sdkacctest.RandomWithPrefix("test-acc-scaleway-object-bucket")

	project, iamAPIKey, terminateFakeSideProject, err := acctest.CreateFakeSideProject(tt)
	require.NoError(t, err)

	ctx := context.Background()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: acctest.FakeSideProjectProviders(ctx, tt, project, iamAPIKey),
		CheckDestroy: resource.ComposeAggregateTestCheckFunc(
			func(_ *terraform.State) error {
				return terminateFakeSideProject()
			},
			objectchecks.IsBucketDestroyed(tt),
		),
		Steps: []resource.TestStep{
			// The side provider should not be able to read the bucket from the main project
			{
				Config: fmt.Sprintf(`
					resource "scaleway_object_bucket" "base" {
						name = "%[1]s"
						region = "%[3]s"
					}

					data "scaleway_object_bucket" "selected" {
						name = scaleway_object_bucket.base.id
						provider = side
					}
				`,
					bucketName,
					project.ID,
					objectTestsMainRegion,
				),
				Check:       objectchecks.CheckBucketExists(tt, "scaleway_object_bucket.base", false),
				ExpectError: regexp.MustCompile("failed getting Object Storage bucket"),
			},
		},
	})
}
