---
subcategory: "Instances"
page_title: "Scaleway: scaleway_instance_volume"
---

# scaleway_instance_volume

Gets information about an instance volume.

## Example Usage

```hcl
# Get info by volume name
data "scaleway_instance_volume" "my_volume" {
  name = "my-volume-name"
}

# Get info by volume ID
data "scaleway_instance_volume" "my_volume" {
  volume_id = "11111111-1111-1111-1111-111111111111"
}
```

## Argument Reference

- `name` - (Optional) The volume name.
  Only one of `name` and `volume_id` should be specified.

- `volume_id` - (Optional) The volume id.
  Only one of `name` and `volume_id` should be specified.

- `zone` - (Defaults to [provider](../index.md#zone) `zone`) The [zone](../guides/regions_and_zones.md#zones) in which the volume exists.

- `project_id` - (Optional) The ID of the project the volume is associated with.

## Attributes Reference

In addition to all above arguments, the following attributes are exported:

- `id` - The ID of the volume.

~> **Important:** Instance volumes' IDs are [zoned](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{zone}/{id}`, e.g. `fr-par-1/11111111-1111-1111-1111-111111111111`

- `volume_type` - The type of the volume.
  `l_ssd` for local SSD, `b_ssd` for block storage SSD.

- `creation_date` - Volume creation date.

- `modification_date` - Volume last modification date.

- `state` - State of the volume. Possible values are `available`, `snapshotting` and `error`.
  The default value is available.

- `size` - The volumes disk size (in bytes).

- `server` - Information about the server attached to the volume.

- `organization_id` - The ID of the organization the volume is associated with.
