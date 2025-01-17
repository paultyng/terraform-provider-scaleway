---
subcategory: "Object Storage"
page_title: "Scaleway: scaleway_object_bucket"
---

# Resource: scaleway_object_bucket

The `scaleway_object_bucket` resource allows you to create and manage buckets for [Scaleway Object storage](https://www.scaleway.com/en/docs/storage/object/).

Refer to the [dedicated documentation](https://www.scaleway.com/en/docs/storage/object/how-to/create-a-bucket/) for more information on Object Storage buckets.

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

* `tags` - (Optional) A list of tags (key/value) for the bucket.

* ~> **Important:** The Scaleway console does not support `key/value` tags yet, so only the tags' values will be displayed.
If you make any change to your bucket's tags using the console, it will overwrite them with the format `value/value`.

* `acl` - (Optional)(Deprecated) The canned ACL you want to apply to the bucket.

-> **Note:** The `acl` attribute is deprecated. See [scaleway_object_bucket_acl](object_bucket_acl.md) resource documentation. Refer to the [official canned ACL documentation](https://docs.aws.amazon.com/AmazonS3/latest/userguide/acl_overview.html#canned-acl) for more information on the different roles.

* `region` - (Optional) The [region](https://www.scaleway.com/en/developers/api/#region-definition) in which the bucket will be created.

* `versioning` - (Optional) A state of [versioning](https://www.scaleway.com/en/docs/storage/object/how-to/use-bucket-versioning/). The `versioning` object supports the following:

    * `enabled` - (Optional) Enable versioning. Once you version-enable a bucket, it can never return to an unversioned state. You can, however, suspend versioning on that bucket.

* `cors_rule` - (Optional) A rule of [Cross-Origin Resource Sharing](https://www.scaleway.com/en/docs/storage/object/api-cli/setting-cors-rules/). The `CORS` object supports the following:

    * `allowed_headers` (Optional) Specifies which headers are allowed.
    * `allowed_methods` (Required) Specifies which methods are allowed (`GET`, `PUT`, `POST`, `DELETE` or `HEAD`).
    * `allowed_origins` (Required) Specifies which origins are allowed.
    * `expose_headers` (Optional) Specifies header exposure in the response.
    * `max_age_seconds` (Optional) Specifies time in seconds that the browser can cache the response for a preflight request.

* `force_destroy` - (Optional) Boolean that, when set to true, allows the deletion of all objects (including locked objects) when the bucket is destroyed. This operation is irreversible, and the objects cannot be recovered. The default is false.

* `project_id` - (Defaults to [provider](../index.md#project_id) `project_id`) The ID of the project the bucket is associated with.

* `lifecycle_rule` (Optional) - A set of rules that defines actions applied to a group of objects. The `lifecycle_rule` object supports the following:

    * `id` - (Optional) Unique identifier for the rule. Must be less than or equal to 255 characters in length.
    * `prefix` - (Optional) Object key prefix identifying one or more objects to which the rule applies.
    * `tags` - (Optional) Specifies object tags key and value.
    * `enabled` - (Required) The element value can be either Enabled or Disabled. If a rule is disabled, Scaleway Object Storage does not perform any of the actions defined in the rule.

* `abort_incomplete_multipart_upload_days` (Optional) Specifies the number of days after initiating a multipart upload when the multipart upload must be completed.

    ~> **Important:** Avoid using `prefix` for `AbortIncompleteMultipartUpload`, as any incomplete multipart upload will be billed

* `expiration` - (Optional) Specifies a period of expiration for the object. The `expiration` object supports the following:

    * `days` (Optional) Specifies the number of days after object creation when the specific rule action takes effect.
    * `transition` - (Optional) Specifies a period in the object's transitions.

At least one of `abort_incomplete_multipart_upload_days`, `expiration`, `transition` must be specified. The `transition` object supports the following:

* `days` (Optional) Specifies the number of days after object creation when the specific rule action takes effect.

* `storage_class` (Required) Specifies the Scaleway [storage class](https://www.scaleway.com/en/docs/storage/object/concepts/#storage-class) `STANDARD`, `GLACIER`, `ONEZONE_IA`  to which you want the object to transition.

~> **Important:**  If versioning is enabled, this rule only deletes the current version of an object.

~> **Important:**  `ONEZONE_IA` is only available in `fr-par` region. The storage class `GLACIER` is not available in `pl-waw` region.

## Attributes Reference

The `scaleway_object_bucket` resource exports certain attributes once the bucket is retrieved. These attributes can be referenced in other parts of your Terraform configuration.

* `id` - The unique name of the bucket.

~> **Important:** Object bucket IDs are [regional](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{region}/{name}`, e.g. `fr-par/bucket-name`.

* `endpoint` - The endpoint URL of the bucket.

* `region` - The Scaleway [region](../guides/regions_and_zones.md) the bucket resides in.

## Import

Buckets can be imported using the `{region}/{bucketName}` identifier, as shown below:

```bash
terraform import scaleway_object_bucket.some_bucket fr-par/some-bucket
```

~> **Important:** The `project_id` attribute has a particular behavior with s3 products because the s3 API is scoped by project.
If you are using a project different from the default one, you have to specify the project ID at the end of the import command.

```bash
terraform import scaleway_object_bucket.some_bucket fr-par/some-bucket@11111111-1111-1111-1111-111111111111
```

