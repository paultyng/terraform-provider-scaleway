package objecttestfuncs

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/locality/regional"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/object"
)

func CheckBucketExists(tt *acctest.TestTools, n string, shouldBeAllowed bool) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		ctx := context.Background()
		rs := state.RootModule().Resources[n]
		if rs == nil {
			return errors.New("resource not found")
		}
		bucketName := rs.Primary.Attributes["name"]
		bucketRegion := rs.Primary.Attributes["region"]

		s3Client, err := object.NewS3ClientFromMeta(ctx, tt.Meta, bucketRegion)
		if err != nil {
			return err
		}

		if rs.Primary.ID == "" {
			return errors.New("no ID is set")
		}

		_, err = s3Client.HeadBucket(ctx, &s3.HeadBucketInput{
			Bucket: scw.StringPtr(bucketName),
		})
		if err != nil {
			if !shouldBeAllowed && object.IsS3Err(err, object.ErrCodeForbidden, object.ErrCodeForbidden) {
				return nil
			}
			if errors.As(err, new(*types.NoSuchBucket)) {
				return errors.New("s3 bucket not found")
			}
			return err
		}
		return nil
	}
}

func IsBucketDestroyed(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		ctx := context.Background()
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway" {
				continue
			}

			regionalID := regional.ExpandID(rs.Primary.ID)
			bucketRegion := regionalID.Region.String()
			bucketName := regionalID.ID

			s3Client, err := object.NewS3ClientFromMeta(ctx, tt.Meta, bucketRegion)
			if err != nil {
				return err
			}

			_, err = s3Client.ListObjects(ctx, &s3.ListObjectsInput{
				Bucket: &bucketName,
			})
			if err != nil {
				if errors.As(err, new(*types.NoSuchBucket)) {
					// Bucket doesn't exist
					continue
				}
				return fmt.Errorf("couldn't get bucket to verify if it still exists: %s", err)
			}

			return errors.New("bucket should be deleted")
		}
		return nil
	}
}

func IsObjectDestroyed(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		ctx := context.Background()
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway" {
				continue
			}

			regionalID := regional.ExpandID(rs.Primary.Attributes["bucket"])
			bucketRegion := regionalID.Region.String()
			bucketName := regionalID.ID
			key := rs.Primary.Attributes["key"]

			s3Client, err := object.NewS3ClientFromMeta(ctx, tt.Meta, bucketRegion)
			if err != nil {
				return err
			}

			_, err = s3Client.GetObject(ctx, &s3.GetObjectInput{
				Bucket: scw.StringPtr(bucketName),
				Key:    scw.StringPtr(key),
			})
			if err != nil {
				if object.IsS3Err(err, object.ErrCodeNoSuchBucket, "The specified bucket does not exist") {
					continue
				}
				return fmt.Errorf("couldn't get object to verify if it still exists: %s", err)
			}

			return errors.New("object should be deleted")
		}
		return nil
	}
}

func IsWebsiteConfigurationDestroyed(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ctx := context.Background()
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "scaleway_object_bucket_website_configuration" {
				continue
			}

			regionalID := regional.ExpandID(rs.Primary.ID)
			bucket := regionalID.ID
			bucketRegion := regionalID.Region

			conn, err := object.NewS3ClientFromMeta(ctx, tt.Meta, bucketRegion.String())
			if err != nil {
				return err
			}

			input := &s3.GetBucketWebsiteInput{
				Bucket: aws.String(bucket),
			}

			output, err := conn.GetBucketWebsite(ctx, input)
			if object.IsS3Err(err, object.ErrCodeNoSuchBucket, "The specified bucket does not exist") {
				continue
			}

			if err != nil {
				return fmt.Errorf("error getting object bucket website configuration (%s): %w", rs.Primary.ID, err)
			}

			if output != nil {
				return fmt.Errorf("object bucket website configuration (%s) still exists", rs.Primary.ID)
			}
		}

		return nil
	}
}

func IsWebsiteConfigurationPresent(tt *acctest.TestTools, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs := s.RootModule().Resources[resourceName]
		if rs == nil {
			return errors.New("resource not found")
		}

		regionalID := regional.ExpandID(rs.Primary.ID)
		bucket := regionalID.ID
		bucketRegion := regionalID.Region
		ctx := context.Background()

		conn, err := object.NewS3ClientFromMeta(ctx, tt.Meta, bucketRegion.String())
		if err != nil {
			return err
		}

		input := &s3.GetBucketWebsiteInput{
			Bucket: aws.String(bucket),
		}

		output, err := conn.GetBucketWebsite(ctx, input)
		if err != nil {
			return fmt.Errorf("error getting object bucket website configuration (%s): %w", rs.Primary.ID, err)
		}

		if output == nil {
			return fmt.Errorf("object bucket website configuration (%s) not found", rs.Primary.ID)
		}

		return nil
	}
}
