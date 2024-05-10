---
subcategory: "Instances"
page_title: "Scaleway: scaleway_instance_image"
---

# scaleway_instance_image

Gets information about an instance image.

## Example Usage

```hcl
# Get info by image name
data "scaleway_instance_image" "my_image" {
  name  = "my-image-name"
}

# Get info by image id
data "scaleway_instance_image" "my_image" {
  image_id = "11111111-1111-1111-1111-111111111111"
}
```

## Argument Reference

- `name` - (Optional) The image name. Only one of `name` and `image_id` should be specified.

- `image_id` - (Optional) The image id. Only one of `name` and `image_id` should be specified.

- `architecture` - (Optional, default `x86_64`) The architecture the image is compatible with. Possible values are: `x86_64` or `arm`.

- `latest` - (Optional, default `true`) Use the latest image ID.

- `zone` - (Defaults to [provider](../index.md#zone) `zone`) The [zone](../guides/regions_and_zones.md#zones) in which the image exists.

## Attributes Reference

In addition to all above arguments, the following attributes are exported:

- `id` - The ID of the image.

~> **Important:** Instance images' IDs are [zoned](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{zone}/{id}`, e.g. `fr-par-1/11111111-1111-1111-1111-111111111111`

- `organization_id` - The ID of the organization the image is associated with.

- `project_id` - The ID of the project the image is associated with.

- `creation_date` - Date of the image creation.

- `modification_date` - Date of image latest update.

- `public` - Set to `true` if the image is public.

- `from_server_id` - ID of the server the image if based from.

- `state` - State of the image. Possible values are: `available`, `creating` or `error`.

- `default_bootscript_id` - ID of the default bootscript for this image.

- `root_volume_id` - ID of the root volume in this image.

- `additional_volume_ids` - IDs of the additional volumes in this image.
