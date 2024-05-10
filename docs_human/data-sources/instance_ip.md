---
subcategory: "Instances"
page_title: "Scaleway: scaleway_instance_ip"
---

# scaleway_instance_ip

Gets information about an instance IP.

## Example Usage

```hcl
# Get info by IP address
data "scaleway_instance_ip" "my_ip" {
  address = "0.0.0.0"
}

# Get info by ID
data "scaleway_instance_ip" "my_ip" {
  id = "fr-par-1/11111111-1111-1111-1111-111111111111"
}
```

## Argument Reference

- `address` - (Optional) The IPv4 address to retrieve
  Only one of `address` and `id` should be specified.

- `id` - (Optional) The ID of the IP address to retrieve
  Only one of `address` and `id` should be specified.

- `zone` - (Defaults to [provider](../index.md#zone) `zone`) The [zone](../guides/regions_and_zones.md#zones) in which the IP should be reserved.

- `project_id` - (Defaults to [provider](../index.md#project_id) `project_id`) The ID of the project the IP is associated with.

## Attributes Reference

In addition to all above arguments, the following attributes are exported:

- `id` - The ID of the IP.

~> **Important:** Instance IPs' IDs are [zoned](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{zone}/{id}`, e.g. `fr-par-1/11111111-1111-1111-1111-111111111111`

- `type` - The type of the IP
- `address` - The IP address.
- `prefix` - The IP Prefix.
- `reverse` - The reverse dns attached to this IP
- `organization_id` - The organization ID the IP is associated with.
