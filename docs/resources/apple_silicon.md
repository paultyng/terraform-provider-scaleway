---
subcategory: "Apple Silicon"
page_title: "Scaleway: scaleway_apple_silicon"
---

# Resource: scaleway_apple_silicon

Creates and manages Scaleway Apple silicon. For more information,
see [the documentation](https://www.scaleway.com/en/developers/api/apple-silicon/).

## Example Usage

### Basic

```terraform
resource scaleway_apple_silicon_server server {
    name = "test-m1"
    type = "M1-M"
}
```

## Argument Reference

The following arguments are supported:

- `type` - (Required) The commercial type of the server. You find all the available types on
  the [pricing page](https://www.scaleway.com/en/pricing/apple-silicon/). Updates to this field will recreate a new
  resource.

- `name` - (Optional) The name of the server.

- `zone` - (Defaults to [provider](../index.md#zone) `zone`) The [zone](../guides/regions_and_zones.md#zones) in which
  the server should be created.

- `project_id` - (Defaults to [provider](../index.md#project_id) `project_id`) The ID of the project the server is
  associated with.
- `enable_vpc` - (Optional, Default: false): Enables the VPC option when set to true.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the server.

~> **Important:** Apple Silicon servers' IDs are [zoned](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{zone}/{id}`, e.g. `fr-par-1/11111111-1111-1111-1111-111111111111`

- `state` - The state of the server.
- `ip` - IPv4 address of the server (IPv4 address).
- `vnc_url` - URL of the VNC.
- `created_at` - The date and time of the creation of the Apple Silicon server.
- `updated_at` - The date and time of the last update of the Apple Silicon server.
- `deleted_at` - The minimal date and time on which you can delete this server due to Apple licence.
- `organization_id` - The organization ID the server is associated with.
- `vpc_status` - The current status of the VPC option.

## Import

Instance servers can be imported using the `{zone}/{id}`, e.g.

```bash
terraform import scaleway_apple_silicon_server.main fr-par-1/11111111-1111-1111-1111-111111111111
```
