---
subcategory: "Instances"
page_title: "Scaleway: scaleway_instance_volume"
---

# Resource: scaleway_instance_volume

Creates and manages Scaleway Compute Instance Volumes.
For more information, see [the documentation](https://developers.scaleway.com/en/products/instance/api/#volumes-7e8a39).

## Example Usage

```terraform
resource "scaleway_instance_volume" "server_volume" {
    type       = "l_ssd"
    name       = "some-volume-name"
    size_in_gb = 20
}
```

## Argument Reference

The following arguments are supported:

- `type` - (Required) The type of the volume. The possible values are: `b_ssd` (Block SSD), `l_ssd` (Local SSD), `scratch` (Local Scratch SSD).
- `size_in_gb` - (Optional) The size of the volume. Only one of `size_in_gb` and `from_snapshot_id` should be specified.
- `from_snapshot_id` - (Optional) If set, the new volume will be created from this snapshot. Only one of `size_in_gb` and `from_snapshot_id` should be specified.
- `name` - (Optional) The name of the volume. If not provided it will be randomly generated.
- `zone` - (Defaults to [provider](../index.md#zone) `zone`) The [zone](../guides/regions_and_zones.md#zones) in which the volume should be created.
- `project_id` - (Defaults to [provider](../index.md#project_id) `project_id`) The ID of the project the volume is associated with.
- `tags` - (Optional) A list of tags to apply to the volume.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the volume.

~> **Important:** Instance volumes' IDs are [zoned](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{zone}/{id}`, e.g. `fr-par-1/11111111-1111-1111-1111-111111111111`

- `server_id` - The id of the associated server.
- `organization_id` - The organization ID the volume is associated with.

## Import

volumes can be imported using the `{zone}/{id}`, e.g.

```bash
$ terraform import scaleway_instance_volume.server_volume fr-par-1/11111111-1111-1111-1111-111111111111
```
