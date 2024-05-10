---
subcategory: "Object Storage"
page_title: "Scaleway: scaleway_object_bucket_policy"
---

# scaleway_object_bucket_policy

Gets information about the Bucket's policy.
For more information, see [the documentation](https://www.scaleway.com/en/docs/object-storage-feature/).

## Example Usage

```hcl
data "scaleway_object_bucket_policy" "main" {
    bucket = "bucket.test.com"
}
```

## Argument Reference

- `bucket` - (Required) The bucket name.
- `region` - (Defaults to [provider](../index.md#arguments-reference) `region`) The [region](../guides/regions_and_zones.md#zones) in which the Object Storage exists.
- `project_id` - (Defaults to [provider](../index.md#arguments-reference) `project_id`) The ID of the project the bucket is associated with.

~> **Important:** The `project_id` attribute has a particular behavior with s3 products because the s3 API is scoped by project.
If you are using a project different from the default one, you have to specify the `project_id` for every child resource of the bucket,
like bucket policies. Otherwise, Terraform will try to create the child resource with the default project ID and you will get a 403 error.


## Attributes Reference

In addition to all above arguments, the following attribute is exported:

* `policy` - The bucket's policy in JSON format.
