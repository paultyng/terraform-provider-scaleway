---
subcategory: "Object Storage"
page_title: "Scaleway: scaleway_object_bucket_lock_configuration"
---

# Resource: scaleway_object_bucket_lock_configuration

Provides an Object bucket lock configuration resource.
For more information, see [Setting up object lock](https://www.scaleway.com/en/docs/storage/object/api-cli/object-lock/).

## Example Usage

### Configure an Object Lock for a new bucket

Please note that `object_lock_enabled` must be set to `true` before configuring the lock.

```terraform
resource "scaleway_object_bucket" "main" {
    name = "MyBucket"
    acl  = "public-read"

    object_lock_enabled = true
}

resource "scaleway_object_bucket_lock_configuration" "main" {
    bucket = scaleway_object_bucket.main.name

    rule {
        default_retention {
            mode = "GOVERNANCE"
            days = 1
        }
    }
}
```

### Configure an Object Lock for an existing bucket

You should [contact Scaleway support](https://console.scaleway.com/support/tickets/create) to enable object lock on an existing bucket.

## Argument Reference

The following arguments are supported:

- `bucket` - (Required, Forces new resource) The name of the bucket, or its Terraform ID.

- `rule` - (Optional) Specifies the Object Lock rule for the specified object.

    - `default_retention` - (Required) The default retention for the lock.

        - `mode` - (Required) The default Object Lock retention mode you want to apply to new objects placed in the specified bucket. Valid values are `GOVERNANCE` or `COMPLIANCE`. To learn more about the difference between these modes, see [Object Lock retention modes](https://www.scaleway.com/en/docs/storage/object/api-cli/object-lock/#retention-modes).

        - `days` - (Optional) The number of days that you want to specify for the default retention period.

        - `years` - (Optional) The number of years that you want to specify for the default retention period.

- `project_id` - (Defaults to [provider](../index.md#arguments-reference) `project_id`) The ID of the project the bucket is associated with.

~> **Important:** The `project_id` attribute has a particular behavior with s3 products because the s3 API is scoped by project.
If you are using a project different from the default one, you have to specify the `project_id` for every child resource of the bucket,
like object lock configurations. Otherwise, Terraform will try to create the child resource with the default project ID and you will get a 403 error.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the Object bucket lock configuration.

~> **Important:** Object buckets lock configurations' IDs are [regional](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{region}/{id}`, e.g. `fr-par/11111111-1111-1111-1111-111111111111`

## Import

Bucket lock configurations can be imported using the `{region}/{bucketName}` identifier, e.g.

```bash
$ terraform import scaleway_object_bucket_lock_configuration.some_bucket fr-par/some-bucket
```

~> **Important:** The `project_id` attribute has a particular behavior with s3 products because the s3 API is scoped by project.
If you are using a project different from the default one, you have to specify the project ID at the end of the import command.

```bash
$ terraform import scaleway_object_bucket_lock_configuration.some_bucket fr-par/some-bucket@xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxx
```
