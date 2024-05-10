---
subcategory: "Elastic Metal"
page_title: "Scaleway: scaleway_flexible_ips"
---

# scaleway_flexible_ips

Gets information about multiple Flexible IPs.

## Example Usage

```hcl
# Find ips that share the same tags
data "scaleway_flexible_ips" "fips_by_tags" {
  tags = [ "a tag" ]
}

# Find ips that share the same Server ID
data "scaleway_baremetal_offer" "my_offer" {
  name = "EM-B112X-SSD"
}

resource "scaleway_baremetal_server" "base" {
  name  = "MyServer"
  offer = data.scaleway_baremetal_offer.my_offer.offer_id
  install_config_afterward  = true
}

resource "scaleway_flexible_ip" "first" {
  server_id = scaleway_baremetal_server.base.id
  tags = [ "foo", "first" ]
}

resource "scaleway_flexible_ip" "second" {
  server_id = scaleway_baremetal_server.base.id
  tags = [ "foo", "second" ]
}

data "scaleway_flexible_ips" "fips_by_server_id" {
  server_ids = [scaleway_baremetal_server.base.id]
  depends_on = [scaleway_flexible_ip.first, scaleway_flexible_ip.second]
}
```

## Argument Reference

- `server_ids` - (Optional)  List of server IDs used as filter. IPs with these exact server IDs are listed.

- `tags` - (Optional)  List of tags used as filter. IPs with these exact tags are listed.

- `zone` - (Defaults to [provider](../index.md#zone) `zone`) The [zone](../guides/regions_and_zones.md#zones) in which IPs exist.

## Attributes Reference

In addition to all above arguments, the following attributes are exported:

- `ips` - List of found flexible IPS
    - `id` - The associated flexible IP ID.
    - `description` - The description of the flexible IP.
    - `tags` - The list of tags which are attached to the flexible IP.
    - `status` - The status of the flexible IP.
    - `mac_address` - The MAC address of the server associated with this flexible IP.
        - `id` - The MAC address ID.
        - `mac_address` - The MAC address of the Virtual MAC.
        - `mac_type` - The type of virtual MAC.
        - `status` - The status of virtual MAC.
        - `created_at` - The date on which the virtual MAC was created (RFC 3339 format).
        - `updated_at` - The date on which the virtual MAC was last updated (RFC 3339 format).
        - `zone` - (Defaults to [provider](../index.md#zone) `zone`) The [zone](../guides/regions_and_zones.md#zones) in which the MAC address exist.
    - `created_at` - The date on which the flexible IP was created (RFC 3339 format).
    - `updated_at` - The date on which the flexible IP was last updated (RFC 3339 format).
    - `reverse` - The reverse domain associated with this IP.
    - `server_id` - The associated server ID if any.
    - `organization_id` - (Defaults to [provider](../index.md#organization_id) `organization_id`) The ID of the organization the IP is in.
    - `project_id` - (Defaults to [provider](../index.md#project_id) `project_id`) The ID of the project the IP is in.
