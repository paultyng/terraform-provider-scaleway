---
subcategory: "Instances"
page_title: "Scaleway: scaleway_instance_snapshot"
---

# scaleway_instance_snapshot

Gets information about an instance snapshot.

## Example Usage

```hcl
# Get info by snapshot name
data "scaleway_instance_snapshot" "by_name" {
  name = "my-snapshot-name"
}

# Get info by snapshot ID
data "scaleway_instance_snapshot" "by_id" {
  snapshot_id = "11111111-1111-1111-1111-111111111111"
}
```

## Argument Reference

- `name` - (Optional) The snapshot name.
  Only one of `name` and `snapshot_id` should be specified.

- `snapshot_id` - (Optional) The snapshot id.
  Only one of `name` and `snapshot_id` should be specified.

- `zone` - (Defaults to [provider](../index.md#zone) `zone`) The [zone](../guides/regions_and_zones.md#zones) in which the snapshot exists.


- `project_id` - (Defaults to [provider](../index.md#project_id) `project_id`) The ID of the project the snapshot is associated with.

## Attributes Reference

Exported attributes are the ones from `instance_snapshot` [resource](../resources/instance_snapshot.md)

