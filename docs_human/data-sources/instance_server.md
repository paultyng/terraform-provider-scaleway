---
subcategory: "Instances"
page_title: "Scaleway: scaleway_instance_server"
---

# scaleway_instance_server

Gets information about an instance server.

## Example Usage

```hcl
# Get info by server name
data "scaleway_instance_server" "my_key" {
  name  = "my-server-name"
}

# Get info by server id
data "scaleway_instance_server" "my_key" {
  server_id = "11111111-1111-1111-1111-111111111111"
}
```

## Argument Reference

- `name` - (Optional) The server name. Only one of `name` and `server_id` should be specified.

- `server_id` - (Optional) The server id. Only one of `name` and `server_id` should be specified.

- `zone` - (Defaults to [provider](../index.md#zone) `zone`) The [zone](../guides/regions_and_zones.md#zones) in which the server exists.

- `project_id` - (Optional) The ID of the project the instance server is associated with.

## Attributes Reference

In addition to all above arguments, the following attributes are exported:

- `id` - The ID of the server.

~> **Important:** Instance servers' IDs are [zoned](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{zone}/{id}`, e.g. `fr-par-1/11111111-1111-1111-1111-111111111111`

- `type` - The commercial type of the server.
You find all the available types on the [pricing page](https://www.scaleway.com/en/pricing/).

- `image` - The UUID and the label of the base image used by the server.

- `organization_id` - The ID of the organization the server is associated with.

- `tags` - The tags associated with the server.

- `security_group_id` - The [security group](https://developers.scaleway.com/en/products/instance/api/#security-groups-8d7f89) the server is attached to.

- `placement_group_id` - The [placement group](https://developers.scaleway.com/en/products/instance/api/#placement-groups-d8f653) the server is attached to.

- `root_volume` - Root [volume](https://developers.scaleway.com/en/products/instance/api/#volumes-7e8a39) attached to the server on creation.
    - `size_in_gb` - Size of the root volume in gigabytes.
    - `delete_on_termination` - Forces deletion of the root volume on instance termination.

- `additional_volume_ids` - The [additional volumes](https://developers.scaleway.com/en/products/instance/api/#volumes-7e8a39)
attached to the server.

- `enable_ipv6` - Determines if IPv6 is enabled for the server.

- `enable_dynamic_ip` - True if dynamic IP in enable on the server.

- `routed_ip_enabled` - True if the server support routed ip only.

- `state` - The state of the server. Possible values are: `started`, `stopped` or `standby`.

- `cloud_init` - The cloud init script associated with this server.

- `user_data` - The user data associated with the server.

    - `key` - The user data key. The `cloud-init` key is reserved, please use `cloud_init` attribute instead.

    - `value` - The user data content.

- `placement_group_policy_respected` - True when the placement group policy is respected.

- `root_volume`
    - `volume_id` - The volume ID of the root volume of the server.

- `private_ip` - The Scaleway internal IP address of the server.

- `public_ip` - The public IP address of the server.

- `public_ips` - The list of public IPs of the server
    - `id` - The ID of the IP
    - `address` - The address of the IP

- `ipv6_address` - The default ipv6 address routed to the server. ( Only set when enable_ipv6 is set to true )

- `ipv6_gateway` - The ipv6 gateway address. ( Only set when enable_ipv6 is set to true )

- `ipv6_prefix_length` - The prefix length of the ipv6 subnet routed to the server. ( Only set when enable_ipv6 is set to true )
