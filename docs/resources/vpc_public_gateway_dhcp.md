---
page_title: "Scaleway: scaleway_vpc_public_gateway_dhcp"
description: |-
  Manages Scaleway VPC Public Gateways IP.
---

# scaleway_vpc_public_gateway_dhcp

Creates and manages Scaleway VPC Public Gateway DHCP.
For more information, see [the documentation](https://developers.scaleway.com/en/products/vpc-gw/api/).

## Example

```hcl
resource "scaleway_vpc_public_gateway_dhcp" "main" {
    subnet = "192.168.1.0/24"
}
```

## Arguments Reference

The following arguments are supported:

- `zone` - (Defaults to [provider](../index.md#zone) `zone`) The [zone](../guides/regions_and_zones.md#zones) in which the public gateway DHCP config should be created.
- `project_id` - (Defaults to [provider](../index.md#project_id) `project_id`) The ID of the project the public gateway DHCP config is associated with.

## Attributes Reference

In addition to all above arguments, the following attributes are exported:

- `id` - The ID of the public gateway DHCP config.
- `organization_id` - The organization ID the public gateway DHCP config is associated with.
- `created_at` - The date and time of the creation of the public gateway DHCP config.
- `updated_at` - The date and time of the last update of the public gateway DHCP config.

## Import

Public gateway DHCP config can be imported using the `{zone}/{id}`, e.g.

```bash
$ terraform import scaleway_vpc_public_gateway_dhcp.main fr-par-1/11111111-1111-1111-1111-111111111111
```
