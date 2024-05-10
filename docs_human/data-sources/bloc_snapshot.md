---
subcategory: "Block"
page_title: "Scaleway: scaleway_block_snapshot"
---

# scaleway_block_snapshot

Gets information about a Block Snapshot.

## Example Usage

```terraform
// Get info by snapshot name
data "scaleway_block_snapshot" "my_snapshot" {
  name = "my-name"
}

// Get info by snapshot name and volume id
data "scaleway_block_snapshot" "my_snapshot" {
  name = "my-name"
  volume_id = "11111111-1111-1111-1111-111111111111"
}

// Get info by snapshot ID
data "scaleway_block_snapshot" "my_snapshot" {
  snapshot_id = "11111111-1111-1111-1111-111111111111"
}
```

## Argument Reference

The following arguments are supported:

- `snapshot_id` - (Optional) The ID of the snapshot. Only one of `name` and `snapshot_id` should be specified.
- `name` - (Optional) The name of the snapshot. Only one of `name` and `snapshot_id` should be specified.
- `volume_id` - (Optional) The ID of the volume from which the snapshot has been created.
- `zone` - (Defaults to [provider](../index.md#zone) `zone`) The [zone](../guides/regions_and_zones.md#zones) in which the snapshot exists.
- `project_id` - (Optional) The ID of the project the snapshot is associated with.

## Attributes Reference

Exported attributes are the ones from `scaleway_block_snapshot` [resource](../resources/block_snapshot.md)
