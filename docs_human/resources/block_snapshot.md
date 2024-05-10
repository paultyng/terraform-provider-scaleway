---
subcategory: "Block"
page_title: "Scaleway: scaleway_block_snapshot"
---

# Resource: scaleway_block_snapshot

Creates and manages Scaleway Block Snapshots.
For more information, see [the documentation](https://www.scaleway.com/en/developers/api/block/).

## Example Usage

```terraform
resource "scaleway_block_snapshot" "block_snapshot" {
    name       = "some-snapshot-name"
    volume_id  = "11111111-1111-1111-1111-111111111111"
}
```

## Argument Reference

The following arguments are supported:

- `volume_id` - (Optional) The ID of the volume to take a snapshot from.
- `name` - (Optional) The name of the snapshot. If not provided it will be randomly generated.
- `zone` - (Defaults to [provider](../index.md#zone) `zone`) The [zone](../guides/regions_and_zones.md#zones) in which the snapshot should be created.
- `project_id` - (Defaults to [provider](../index.md#project_id) `project_id`) The ID of the project the snapshot is associated with.
- `tags` - (Optional) A list of tags to apply to the snapshot.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the snapshot.

~> **Important:** Block snapshots' IDs are [zoned](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{zone}/{id}`, e.g. `fr-par-1/11111111-1111-1111-1111-111111111111`

## Import

Block Snapshots can be imported using the `{zone}/{id}`, e.g.

```bash
$ terraform import scaleway_block_snapshot.main fr-par-1/11111111-1111-1111-1111-111111111111
```
