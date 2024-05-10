---
subcategory: "VPC"
page_title: "Scaleway: scaleway_vpc_public_gateway"
---

# Resource: scaleway_vpc_public_gateway

Creates and manages Scaleway VPC Public Gateway.
For more information, see [the documentation](https://developers.scaleway.com/en/products/vpc-gw/api/v1).

## Example Usage

```terraform
resource "scaleway_vpc_public_gateway" "main" {
    name = "public_gateway_demo"
    type = "VPC-GW-S"
    tags = ["demo", "terraform"]
}
```

## Argument Reference

The following arguments are supported:

- `type` - (Required) The gateway type.
- `name` - (Optional) The name of the public gateway. If not provided it will be randomly generated.
- `tags` - (Optional) The tags associated with the public gateway.
- `zone` - (Defaults to [provider](../index.md#zone) `zone`) The [zone](../guides/regions_and_zones.md#zones) in which the public gateway should be created.
- `project_id` - (Defaults to [provider](../index.md#project_id) `project_id`) The ID of the project the public gateway is associated with.
- `upstream_dns_servers` - (Optional) override the gateway's default recursive DNS servers, if DNS features are enabled.
- `ip_id` - (Optional) attach an existing flexible IP to the gateway
- `bastion_enabled` - (Optional) Enable SSH bastion on the gateway
- `bastion_port` - (Optional) The port on which the SSH bastion will listen.
- `enable_smtp` - (Optional) Enable SMTP on the gateway

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the public gateway.

~> **Important:** Public Gateways' IDs are [zoned](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{zone}/{id}`, e.g. `fr-par-1/11111111-1111-1111-1111-111111111111`

- `organization_id` - The organization ID the public gateway is associated with.
- `created_at` - The date and time of the creation of the public gateway.
- `updated_at` - The date and time of the last update of the public gateway.
- `status` - The status of the public gateway.

## Import

Public gateway can be imported using the `{zone}/{id}`, e.g.

```bash
$ terraform import scaleway_vpc_public_gateway.main fr-par-1/11111111-1111-1111-1111-111111111111
```
