---
subcategory: "Instances"
page_title: "Scaleway: scaleway_instance_image"
---

# Resource: scaleway_instance_image

Creates and manages Scaleway Compute Images.
For more information, see [the documentation](https://developers.scaleway.com/en/products/instance/api/#images-41389b).

## Example Usage

### From a volume

```terraform
resource "scaleway_instance_volume" "volume" {
  type       	= "b_ssd"
  size_in_gb 	= 20
}

resource "scaleway_instance_snapshot" "volume_snapshot" {
  volume_id 	= scaleway_instance_volume.volume.id
}

resource "scaleway_instance_image" "volume_image" {
  name 		  = "image_from_volume"
  root_volume_id  = scaleway_instance_snapshot.volume_snapshot.id
}
```

### From a server

```terraform
resource "scaleway_instance_server" "server" {
  image = "ubuntu_jammy"
  type 	= "DEV1-S"
}

resource "scaleway_instance_snapshot" "server_snapshot" {
  volume_id	= scaleway_instance_server.main.root_volume.0.volume_id
}

resource "scaleway_instance_image" "server_image" {
  name            = "image_from_server"
  root_volume_id  = scaleway_instance_snapshot.server_snapshot.id
}
```

### With additional volumes

```terraform
resource "scaleway_instance_server" "server" {
  image = "ubuntu_jammy"
  type 	= "DEV1-S"
}

resource "scaleway_instance_volume" "volume" {
  type       	= "b_ssd"
  size_in_gb 	= 20
}

resource "scaleway_instance_snapshot" "volume_snapshot" {
  volume_id     = scaleway_instance_volume.volume.id
}
resource "scaleway_instance_snapshot" "server_snapshot" {
  volume_id     = scaleway_instance_server.main.root_volume.0.volume_id
}

resource "scaleway_instance_image" "image" {
  name            = "image_with_extra_volumes"
  root_volume_id  = scaleway_instance_snapshot.server_snapshot.id
  additional_volumes = [
    scaleway_instance_snapshot.volume_snapshot.id
  ]
}
```

## Argument Reference

The following arguments are supported:

- `root_volume_id` - (Required) The ID of the snapshot of the volume to be used as root in the image.
- `name` - (Optional) The name of the image. If not provided it will be randomly generated.
- `architecture` - (Optional, default `x86_64`) The architecture the image is compatible with. Possible values are: `x86_64` or `arm`.
- `additional_volume_ids` - (Optional) List of IDs of the snapshots of the additional volumes to be attached to the image.

-> **Important:** For now it is only possible to have 1 additional_volume.

- `tags` - (Optional) A list of tags to apply to the image.
- `public` - (Optional) Set to `true` if the image is public.
- `zone` - (Defaults to provider `zone`) The [zone](../guides/regions_and_zones.md#zones) in which the image should be created.
- `project_id` - (Defaults to provider `project_id`) The ID of the project the image is associated with.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the image.

~> **Important:** Instance images' IDs are [zoned](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{zone}/{id}`, e.g. `fr-par-1/11111111-1111-1111-1111-111111111111`

- `creation_date` - Date of the image creation.
- `modification_date` - Date of image latest update.
- `from_server_id` - ID of the server the image is based on (in case it is a backup).
- `state` - State of the image. Possible values are: `available`, `creating` or `error`.
- `organization_id` - The organization ID the image is associated with.
- `additional_volumes` - The description of the extra volumes attached to the image.

    -> The `additional_volumes` block contains :
    - `id` - The ID of the volume.
    - `name` - The name of the volume.
    - `export_uri` - The export URI of the volume.
    - `size` - The size of the volume.
    - `volume_type` - The type of volume, possible values are `l_ssd` and `b_ssd`.
    - `creation_date` - Date of the volume creation.
    - `modification_date` - Date of volume latest update.
    - `organization` - The organization ID the volume is associated with.
    - `project` - ID of the project the volume is associated with
    - `tags` - List of tags associated with the volume.
    - `state` - State of the volume.
    - `zone` - The [zone](../guides/regions_and_zones.md#zones) in which the volume is.
    - `server` - Description of the server containing the volume (in case the image is a backup from a server).
  
    -> The `server` block contains :
      - `id` - ID of the server containing the volume.
      - `name` - Name of the server containing the volume.

## Import

Images can be imported using the `{zone}/{id}`, e.g.

```bash
$ terraform import scaleway_instance_image.main fr-par-1/11111111-1111-1111-1111-111111111111
```
