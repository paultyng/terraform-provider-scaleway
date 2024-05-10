---
subcategory: "Object Storage"
page_title: "Scaleway: scaleway_object_bucket"
---

# Resource: scaleway_object_bucket

Creates and manages Scaleway object storage buckets.
For more information, see [the documentation](https://www.scaleway.com/en/docs/object-storage-feature/).

## Example Usage

```terraform
resource "scaleway_object_bucket" "some_bucket" {
  name = "some-unique-name"
  tags = {
    key = "value"
  }
}
```

### Creating the bucket in a specific project

```terraform
resource "scaleway_object_bucket" "some_bucket" {
  name = "some-unique-name"
  project_id = "11111111-1111-1111-1111-111111111111"
}
```

### Using object lifecycle

```terraform
resource "scaleway_object_bucket" "main"{
  name = "mybuckectid"
  region = "fr-par"
  
  # This lifecycle configuration rule will make that all objects that got a filter key that start with (path1/) be transferred
  # from their default storage class (STANDARD, ONEZONE_IA) to GLACIER after 120 days counting 
  # from their creation and then 365 days after that they will be expired and deleted.
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
  
  # This lifecycle configuration rule specifies that all objects (identified by the key name prefix (path2/) in the rule)
  # from their creation and then 50 days after that they will be expired and deleted.
  lifecycle_rule {
      id      = "id2"
      prefix  = "path2/"
      enabled = true
  
      expiration {
        days = "50"
      }
  }
  
  # This lifecycle configuration rule remove any object with (path3/) prefix that match
  # with the tags one day after creation.
  lifecycle_rule {
      id      = "id3"
      prefix  = "path3/"
      enabled = false
  
      tags = {
        "tagKey"    = "tagValue"
        "terraform" = "hashicorp"
      }
  
      expiration {
        days = "1"
      }
  }
  
  # This lifecycle configuration rule specifies a tag-based filter (tag1/value1).
  # This rule directs Scaleway S3 to transition objects S3 Glacier class soon after creation.
  # It is also disable temporaly.
  lifecycle_rule {
      id      = "id4"
      enabled = true
      
      tags = {
        "tag1"    = "value1"
      }
      
      transition {
        days          = 1
        storage_class = "GLACIER"
      }
  }
 
  # This lifecycle configuration rule specifies with the AbortIncompleteMultipartUpload action to 
  # stop incomplete multipart uploads (identified by the key name prefix (path5/) in the rule)
  # if they aren't completed within a specified number of days after initiation.
  # Note: It's not recommended using prefix/ for AbortIncompleteMultipartUpload as any incomplete multipart upload will be billed
  lifecycle_rule {
      #  prefix  = "path5/"
      enabled = true
      abort_incomplete_multipart_upload_days = 30
  }
}
```

## Argument Reference


The following arguments are supported:

* `name` - (Required) The name of the bucket.
* `tags` - (Optional) A list of tags (key / value) for the bucket.

* ~> **Important:** The Scaleway console does not support `key/value` tags yet, so only the tags' values will be displayed.
Keep in mind that if you make any change to your bucket's tags using the console, it will overwrite them with the format `value/value`.

* `acl` - (Optional)(Deprecated) The canned ACL you want to apply to the bucket.
* `region` - (Optional) The [region](https://developers.scaleway.com/en/quickstart/#region-definition) in which the bucket should be created.
* `versioning` - (Optional) A state of [versioning](https://docs.aws.amazon.com/AmazonS3/latest/dev/Versioning.html) (documented below)
* `cors_rule` - (Optional) A rule of [Cross-Origin Resource Sharing](https://docs.aws.amazon.com/AmazonS3/latest/dev/cors.html) (documented below).
* `force_destroy` - (Optional) Enable deletion of objects in bucket before destroying, locked objects or under legal hold are also deleted and **not** recoverable
* `project_id` - (Defaults to [provider](../index.md#project_id) `project_id`) The ID of the project the bucket is associated with.

The `acl` attribute is deprecated. See [scaleway_object_bucket_acl](object_bucket_acl.md) resource documentation.
Please check the [canned ACL](https://docs.aws.amazon.com/AmazonS3/latest/userguide/acl_overview.html#canned-acl) documentation for supported values.

The `CORS` object supports the following:

* `allowed_headers` (Optional) Specifies which headers are allowed.
* `allowed_methods` (Required) Specifies which methods are allowed. Can be `GET`, `PUT`, `POST`, `DELETE` or `HEAD`.
* `allowed_origins` (Required) Specifies which origins are allowed.
* `expose_headers` (Optional) Specifies expose header in the response.
* `max_age_seconds` (Optional) Specifies time in seconds that browser can cache the response for a preflight request.

The `lifecycle_rule` (Optional) object supports the following:

* `id` - (Optional) Unique identifier for the rule. Must be less than or equal to 255 characters in length.
* `prefix` - (Optional) Object key prefix identifying one or more objects to which the rule applies.
* `tags` - (Optional) Specifies object tags key and value.
* `enabled` - (Required) The element value can be either Enabled or Disabled. If a rule is disabled, Scaleway S3 doesn't perform any of the actions defined in the rule.

* `abort_incomplete_multipart_upload_days` (Optional) Specifies the number of days after initiating a multipart upload when the multipart upload must be completed.

    * ~> **Important:** It's not recommended using `prefix` for `AbortIncompleteMultipartUpload` as any incomplete multipart upload will be billed

* `expiration` - (Optional) Specifies a period in the object's expire (documented below).
* `transition` - (Optional) Specifies a period in the object's transitions (documented below).

At least one of `abort_incomplete_multipart_upload_days`, `expiration`, `transition` must be specified.

The `expiration` object supports the following

* `days` (Optional) Specifies the number of days after object creation when the specific rule action takes effect.

~> **Important:**  If versioning is enabled, this rule only deletes the current version of an object.

The `transition` object supports the following

* `days` (Optional) Specifies the number of days after object creation when the specific rule action takes effect.
* `storage_class` (Required) Specifies the Scaleway [storage class](https://www.scaleway.com/en/docs/storage/object/concepts/#storage-class) `STANDARD`, `GLACIER`, `ONEZONE_IA`  to which you want the object to transition.

~> **Important:**  `ONEZONE_IA` is only available in `fr-par` region. The storage class `GLACIER` is not available in `pl-waw` region.

The `versioning` object supports the following:

* `enabled` - (Optional) Enable versioning. Once you version-enable a bucket, it can never return to an unversioned state. You can, however, suspend versioning on that bucket.

## Attributes Reference

In addition to all arguments above, the following attribute is exported:

* `id` - The unique name of the bucket.

~> **Important:** Object buckets' IDs are [regional](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{region}/{name}`, e.g. `fr-par/bucket-name`

* `endpoint` - The endpoint URL of the bucket
* `region` - The Scaleway region this bucket resides in.

## Import

Buckets can be imported using the `{region}/{bucketName}` identifier, e.g.

```bash
$ terraform import scaleway_object_bucket.some_bucket fr-par/some-bucket
```

If you are importing a bucket from a specific project (that is not your default project), you can use the following syntax:

```bash
$ terraform import scaleway_object_bucket.some_bucket fr-par/some-bucket@11111111-1111-1111-1111-111111111111
```
