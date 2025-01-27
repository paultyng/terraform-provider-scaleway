---
subcategory: "VPC"
page_title: "Scaleway: scaleway_vpc_public_gateway_dhcp"
---

# Resource: scaleway_vpc_public_gateway_dhcp



Creates and manages Scaleway VPC Public Gateway DHCP configurations.
For more information, see [the documentation](https://www.scaleway.com/en/developers/api/public-gateway/#dhcp-c05544).

## Example Usage

```terraform
resource "scaleway_vpc_public_gateway_dhcp" "main" {
    subnet = "192.168.1.0/24"
}
```

## Argument Reference

The following arguments are supported:

- `zone` - (Defaults to [provider](../index.md#zone) `zone`) The [zone](../guides/regions_and_zones.md#zones) in which the Public Gateway DHCP configuration should be created.
- `project_id` - (Defaults to [provider](../index.md#project_id) `project_id`) The ID of the Project the Public Gateway DHCP configuration is associated with.
- `subnet` - (Required) The subnet to associate with the Public Gateway DHCP configuration.
- `address` - (Optional) The IP address of the DHCP server. This will be the gateway's address in the Private Network.
- `pool_low` - (Optional) Low IP (included) of the dynamic address pool. Defaults to the second address of the subnet.
- `pool_high` - (Optional) High IP (excluded) of the dynamic address pool. Defaults to the last address of the subnet.
- `enable_dynamic` - (Optional) Whether to enable dynamic pooling of IPs. By turning the dynamic pool off, only pre-existing DHCP reservations will be handed out. Defaults to `true`.
- `valid_lifetime` - (Optional) How long, in seconds, DHCP entries will be valid. Defaults to 1h (3600s).
- `renew_timer` - (Optional) After how long, in seconds, a renewal will be attempted. Must be 30s lower than `rebind_timer`. Defaults to 50m (3000s).
- `rebind_timer` - (Optional) After how long, in seconds, a DHCP client will query for a new lease if previous renews fail. Must be 30s lower than `valid_lifetime`. Defaults to 51m (3060s).
- `push_default_route` - (Optional) Whether the gateway should push a default route to DHCP clients or only hand out IPs. Defaults to `true`.

~> **Warning**: If you need to setup a default route, it's recommended to use the [`scaleway_vpc_gateway_network`](vpc_gateway_network.md#create-a-gatewaynetwork-with-ipam-configuration) resource instead.

- `push_dns_server` - (Optional) Whether the gateway should push custom DNS servers to clients. This allows for instance hostname -> IP resolution. Defaults to `true`.
- `dns_servers_override` - (Optional) Override the DNS server list pushed to DHCP clients, instead of the gateway itself.
- `dns_search` - (Optional) Additional DNS search paths
- `dns_local_name` - (Optional) TLD given to hostnames in the Private Network. Allowed characters are `a-z0-9-.`. Defaults to the slugified Private Network name if created along a GatewayNetwork, or else to `priv`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the public gateway DHCP config.

~> **Important:** Public Gateway DHCP configuration IDs are [zoned](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{zone}/{id}`, e.g. `fr-par-1/11111111-1111-1111-1111-111111111111`

- `organization_id` - The Organization ID the Public Gateway DHCP config is associated with.
- `created_at` - The date and time of the creation of the Public Gateway DHCP configuration.
- `updated_at` - The date and time of the last update of the Public Gateway DHCP configuration.

## Import

Public Gateway DHCP configuration can be imported using `{zone}/{id}`, e.g.

```bash
terraform import scaleway_vpc_public_gateway_dhcp.main fr-par-1/11111111-1111-1111-1111-111111111111
```
