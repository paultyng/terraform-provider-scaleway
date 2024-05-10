---
subcategory: "Instances"
page_title: "Scaleway: scaleway_instance_snapshot"
---

# Resource: scaleway_instance_snapshot

Creates and manages Scaleway Compute Snapshots.
For more information,
see [the documentation](https://developers.scaleway.com/en/products/instance/api/#snapshots-756fae).

## Example Usage

```terraform
resource "scaleway_instance_snapshot" "main" {
    name       = "some-snapshot-name"
    volume_id  = "11111111-1111-1111-1111-111111111111"
}
```

### Example with Unified type snapshot

```terraform
resource "scaleway_instance_volume" "main" {
    type       = "l_ssd"
    size_in_gb = 10
}

resource "scaleway_instance_server" "main" {
    image    = "ubuntu_jammy"
    type     = "DEV1-S"
    root_volume {
        size_in_gb = 10
        volume_type = "l_ssd"
    }
    additional_volume_ids = [
        scaleway_instance_volume.main.id
    ]
}

resource "scaleway_instance_snapshot" "main" {
    volume_id = scaleway_instance_volume.main.id
    type = "unified"
    depends_on = [scaleway_instance_server.main]
}
```

### Example importing a local qcow2 file

```terraform
resource "scaleway_object_bucket" "bucket" {
  name = "snapshot-qcow-import"
}

resource "scaleway_object" "qcow" {
  bucket = scaleway_object_bucket.bucket.name
  key = "server.qcow2"
  file = "myqcow.qcow2"
}

resource "scaleway_instance_snapshot" "snapshot" {
  type = "unified"
  import {
    bucket = scaleway_object.qcow.bucket
    key = scaleway_object.qcow.key
  }
}
```

## Argument Reference

The following arguments are supported:

- `volume_id` - (Optional) The ID of the volume to take a snapshot from.
- `type` - (Optional) The snapshot's volume type.  The possible values are: `b_ssd` (Block SSD), `l_ssd` (Local SSD) and `unified`.
Updates to this field will recreate a new resource.
- `name` - (Optional) The name of the snapshot. If not provided it will be randomly generated.
- `zone` - (Defaults to [provider](../index.md#zone) `zone`) The [zone](../guides/regions_and_zones.md#zones) in which
  the snapshot should be created.
- `project_id` - (Defaults to [provider](../index.md#project_id) `project_id`) The ID of the project the snapshot is
  associated with.
- `tags` - (Optional) A list of tags to apply to the snapshot.
- `import` - (Optional) Import a snapshot from a qcow2 file located in a bucket
    - `bucket` - Bucket name containing [qcow2](https://en.wikipedia.org/wiki/Qcow) to import
    - `key` - Key of the object to import

-> **Note:** The type `unified` could be instantiated on both `l_ssd` and `b_ssd` volumes.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the snapshot.

~> **Important:** Instance snapshots' IDs are [zoned](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{zone}/{id}`, e.g. `fr-par-1/11111111-1111-1111-1111-111111111111`

- `size_in_gb` - (Optional) The size of the snapshot.
- `organization_id` - The organization ID the snapshot is associated with.
- `project_id` - The project ID the snapshot is associated with.
- `created_at` - The snapshot creation time.

## Import

Snapshots can be imported using the `{zone}/{id}`, e.g.

```bash
$ terraform import scaleway_instance_snapshot.main fr-par-1/11111111-1111-1111-1111-111111111111
```
