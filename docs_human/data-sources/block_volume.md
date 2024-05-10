---
subcategory: "Block"
page_title: "Scaleway: scaleway_block_volume"
---

# scaleway_block_volume

Gets information about a Block Volume.

## Example Usage

```terraform
// Get info by volume name
data "scaleway_block_volume" "my_volume" {
  name = "my-name"
}

// Get info by volume ID
data "scaleway_block_volume" "my_volume" {
  volume_id = "11111111-1111-1111-1111-111111111111"
}
```

## Argument Reference

The following arguments are supported:

- `volume_id` - (Optional) The ID of the volume. Only one of `name` and `volume_id` should be specified.
- `name` - (Optional) The name of the volume. Only one of `name` and `volume_id` should be specified.
- `zone` - (Defaults to [provider](../index.md#zone) `zone`) The [zone](../guides/regions_and_zones.md#zones) in which the volume exists.
- `project_id` - (Optional) The ID of the project the volume is associated with.

## Attributes Reference

Exported attributes are the ones from `scaleway_block_volume` [resource](../resources/block_volume.md)
