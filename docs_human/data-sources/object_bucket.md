---
subcategory: "Object Storage"
page_title: "Scaleway: scaleway_object_bucket"
---

# scaleway_object_bucket

Gets information about the Bucket.
For more information, see [the documentation](https://www.scaleway.com/en/docs/object-storage-feature/).

## Example Usage

```hcl
resource "scaleway_object_bucket" "main" {
    name = "bucket.test.com"
    tags = {
        foo = "bar"
    }
}

data "scaleway_object_bucket" "selected" {
  name = scaleway_object_bucket.main.id
}
```


### Fetching the bucket from a specific project

```hcl
data "scaleway_object_bucket" "selected" {
  name = "bucket.test.com"
  project_id = "11111111-1111-1111-1111-111111111111"
}
```

## Argument Reference

- `name` - (Required) The bucket name, or its terraform's ID (`{region}/{name}`)
- `object_lock_enabled` - (Optional) Enable object lock on the bucket. Defaults to `false`. Updating this field will force creating a new bucket.
- `region` - (Defaults to [provider](../index.md#region) `region`) The [region](../guides/regions_and_zones.md#zones) in which the bucket exists.
- `project_id` - (Defaults to [provider](../index.md#project_id) `project_id`) The ID of the project the bucket is associated with.


## Attributes Reference

In addition to all above arguments, the following attribute is exported:

* `id` - The unique name of the bucket.

~> **Important:** Object buckets' IDs are [regional](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{region}/{name}`, e.g. `fr-par/bucket-name`

* `endpoint` - The endpoint URL of the bucket