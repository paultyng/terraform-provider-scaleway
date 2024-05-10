---
subcategory: "Instances"
page_title: "Scaleway: scaleway_instance_servers"
---

# scaleway_instance_servers

Gets information about multiple instance servers.

## Examples

### Basic

```hcl
# Find servers by tag
data "scaleway_instance_servers" "my_key" {
  tags  = ["tag"]
}

# Find servers by name and zone
data "scaleway_instance_servers" "my_key" {
  name = "myserver"
  zone = "fr-par-2"
}
```

## Argument Reference

- `name` - (Optional) The server name used as filter. Servers with a name like it are listed.

- `tags` - (Optional) List of tags used as filter. Servers with these exact tags are listed.

- `zone` - (Defaults to [provider](../index.md#zone) `zone`) The [zone](../guides/regions_and_zones.md#zones) in which servers exist.

## Attributes Reference

In addition to all above arguments, the following attributes are exported:

- `id` - The zone of the servers

- `servers` - List of found servers
    - `id` - The ID of the server.

        ~> **Important:** Instance servers' IDs are [zoned](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{zone}/{id}`, e.g. `fr-par-1/11111111-1111-1111-1111-111111111111`

    - `tags` - The tags associated with the server.
    - `public_ip` - The public IP address of the server.
    - `private_ip` - The Scaleway internal IP address of the server.
    - `public_ips` - The list of public IPs of the server
        - `id` - The ID of the IP
        - `address` - The address of the IP
    - `prefix` - The public IP prefix of the server.
    - `state` - The state of the server. Possible values are: `started`, `stopped` or `standby`.
    - `zone` - The [zone](../guides/regions_and_zones.md#zones) in which the server is.
    - `name` - The name of the server.
    - `boot_type` - The boot Type of the server. Possible values are: `local`, `bootscript` or `rescue`.
    - `bootscript_id` - The ID of the bootscript.
    - `type` - The commercial type of the server.
    - `security_group_id` - The [security group](https://developers.scaleway.com/en/products/instance/api/#security-groups-8d7f89) the server is attached to.
    - `enable_ipv6` - Determines if IPv6 is enabled for the server.
    - `ipv6_address` - The default ipv6 address routed to the server. ( Only set when enable_ipv6 is set to true )
    - `ipv6_gateway` - The ipv6 gateway address. ( Only set when enable_ipv6 is set to true )
    - `ipv6_prefix_length` - The prefix length of the ipv6 subnet routed to the server. ( Only set when enable_ipv6 is set to true )
    - `routed_ip_enabled` - True if the server support routed ip only.
    - `enable_dynamic_ip` - If true a dynamic IP will be attached to the server.
    - `image` - The UUID or the label of the base image used by the server.
    - `placement_group_id` - The [placement group](https://developers.scaleway.com/en/products/instance/api/#placement-groups-d8f653) the server is attached to.
    - `organization_id` - The organization ID the server is associated with.
    - `project_id` - The ID of the project the server is associated with.
  

