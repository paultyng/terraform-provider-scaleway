---
subcategory: "Block"
page_title: "Scaleway: scaleway_block_volume"
---

# Resource: scaleway_block_volume

Creates and manages Scaleway Block Volumes.
For more information, see [the documentation](https://www.scaleway.com/en/developers/api/block/).

## Example Usage

```terraform
resource "scaleway_block_volume" "block_volume" {
    iops       = 5000
    name       = "some-volume-name"
    size_in_gb = 20
}
```

## Argument Reference

The following arguments are supported:

- `iops` - (Required) The maximum IO/s expected, must match available options.
- `name` - (Optional) The name of the volume. If not provided it will be randomly generated.
- `size_in_gb` - (Optional) The size of the volume. Only one of `size_in_gb`, and `snapshot_id` should be specified.
- `snapshot_id` - (Optional) If set, the new volume will be created from this snapshot. Only one of `size_in_gb`, `snapshot_id` should be specified.
- `tags` - (Optional) A list of tags to apply to the volume.
- `zone` - (Defaults to [provider](../index.md#zone) `zone`) The [zone](../guides/regions_and_zones.md#zones) in which the volume should be created.
- `project_id` - (Defaults to [provider](../index.md#project_id) `project_id`) The ID of the project the volume is associated with.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the volume.

~> **Important:** Block volumes' IDs are [zoned](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{zone}/{id}`, e.g. `fr-par-1/11111111-1111-1111-1111-111111111111`

- `organization_id` - The organization ID the volume is associated with.

## Import

Block Volumes can be imported using the `{zone}/{id}`, e.g.

```bash
$ terraform import scaleway_block_volume.block_volume fr-par-1/11111111-1111-1111-1111-111111111111
```
