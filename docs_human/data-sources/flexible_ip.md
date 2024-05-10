---
subcategory: "Elastic Metal"
page_title: "Scaleway: scaleway_flexible_ip"
---

# scaleway_flexible_ip

Gets information about a Flexible IP.

## Example Usage

```hcl
# Get info by IP address
data "scaleway_flexible_ip" "my_ip" {
  ip_address = "1.2.3.4"
}

# Get info by IP ID
data "scaleway_flexible_ip" "my_ip" {
  ip_id = "11111111-1111-1111-1111-111111111111"
}
```

## Argument Reference

- `ip_address` - (Optional) The IP address.
  Only one of `ip_address` and `ip_id` should be specified.

- `ip_id` - (Optional) The IP ID.
  Only one of `ip_address` and `ip_id` should be specified.

## Attributes Reference

In addition to all above arguments, the following attributes are exported:

- `id` - The ID of the flexible IP.

~> **Important:** Flexible IPs' IDs are [zoned](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{zone}/{id}`, e.g. `fr-par-1/11111111-1111-1111-1111-111111111111`

- `reverse` - The reverse domain associated with this IP.
- `server_id` - The associated server ID if any
- `organization_id` - (Defaults to [provider](../index.md#organization_id) `organization_id`) The ID of the organization the IP is in.
- `project_id` - (Defaults to [provider](../index.md#project_id) `project_id`) The ID of the project the IP is in.
