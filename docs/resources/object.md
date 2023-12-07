---
subcategory: "Object Storage"
page_title: "Scaleway: scaleway_object"
---

# scaleway_object

Creates and manages Scaleway object storage objects.
For more information, see [the documentation](https://www.scaleway.com/en/docs/object-storage-feature/).

## Example Usage

```hcl
resource "scaleway_object_bucket" "some_bucket" {
  name = "some-unique-name"
}

resource scaleway_object "some_file" {
  bucket = scaleway_object_bucket.some_bucket.id
  key = "object_path"
  
  file = "myfile"
  hash = filemd5("myfile")
}
```

## Arguments Reference


The following arguments are supported:

* `bucket` - (Required) The name of the bucket, or its Terraform ID.
* `key` - (Required) The path of the object.
* `file` - (Optional) The name of the file to upload, defaults to an empty file. Only one of `file`, `content` or `content_base64` can be defined.
* `content` - (Optional) The content of the file to upload. Only one of `file`, `content` or `content_base64` can be defined.
* `content_base64` - (Optional) The base64-encoded content of the file to upload. Only one of `file`, `content` or `content_base64` can be defined.
* `hash` - (Optional) Hash of the file, used to trigger upload on file change
* `storage_class` - (Optional) Specifies the Scaleway [storage class](https://www.scaleway.com/en/docs/storage/object/concepts/#storage-class) `STANDARD`, `GLACIER`, `ONEZONE_IA` used to store the object.
* `visibility` - (Optional) Visibility of the object, `public-read` or `private`
* `metadata` - (Optional) Map of metadata used for the object, keys must be lowercase
* `tags` - (Optional) Map of tags
* `project_id` - (Defaults to [provider](../index.md#arguments-reference) `project_id`) The ID of the project the bucket is associated with.

## Attributes Reference

In addition to all above arguments, the following attribute is exported:

* `id` - The path of the object, including bucket name.

~> **Important:** Objects' IDs are [regional](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{region}/{bucket-name}/{key}`, e.g. `fr-par/bucket-name/object-key`

* `region` - The Scaleway region this bucket resides in.

## Import

Objects can be imported using the `{region}/{bucketName}/{objectKey}` identifier, e.g.

```bash
$ terraform import scaleway_object.some_object fr-par/some-bucket/some-file
```
