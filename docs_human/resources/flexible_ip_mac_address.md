---
subcategory: "Elastic Metal"
page_title: "Scaleway: scaleway_flexible_ip_mac_address"
---

# Resource: scaleway_flexible_ip_mac_address

Creates and manages Scaleway Flexible IP Mac Addresses.
For more information, see [the documentation](https://developers.scaleway.com/en/products/flexible-ip/api).

## Example Usage

### Basic

```terraform
resource "scaleway_flexible_ip" "main" {}

resource "scaleway_flexible_ip_mac_address" "main" {
  flexible_ip_id = scaleway_flexible_ip.main.id
  type = "kvm"
}
```

### Duplicate on many other flexible IPs

```terraform
data "scaleway_baremetal_offer" "my_offer" {
  name = "EM-B112X-SSD"
}

resource "scaleway_baremetal_server" "base" {
  name 			             = "TestAccScalewayBaremetalServer_WithoutInstallConfig"
  offer     				 = data.scaleway_baremetal_offer.my_offer.offer_id
  install_config_afterward   = true
}

resource "scaleway_flexible_ip" "ip01" {
  server_id = scaleway_baremetal_server.base.id
}

resource "scaleway_flexible_ip" "ip02" {
  server_id = scaleway_baremetal_server.base.id
}

resource "scaleway_flexible_ip" "ip03" {
  server_id = scaleway_baremetal_server.base.id
}

resource "scaleway_flexible_ip_mac_address" "main" {
  flexible_ip_id = scaleway_flexible_ip.ip01.id
  type = "kvm"
  flexible_ip_ids_to_duplicate = [
    scaleway_flexible_ip.ip02.id,
    scaleway_flexible_ip.ip03.id
  ]
}
```

## Argument Reference

The following arguments are supported:

- `flexible_ip_id`: (Required) The ID of the flexible IP for which to generate a virtual MAC.
- `type`: (Required) The type of the virtual MAC.
- `flexible_ip_ids_to_duplicate` - (Optional) The IDs of the flexible IPs on which to duplicate the virtual MAC.
~> **Important:** The flexible IPs need to be attached to the same server for the operation to work.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the Flexible IP Mac Address

~> **Important:** Flexible IP Mac Addresses' IDs are [zoned](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{zone}/{id}`, e.g. `fr-par-1/11111111-1111-1111-1111-111111111111`

- `address` -  The Virtual MAC address.
- `zone` - The zone of the Virtual Mac Address.
- `status` - The Virtual MAC status.
- `created_at` - The date at which the Virtual Mac Address was created (RFC 3339 format).
- `updated_at` - The date at which the Virtual Mac Address was last updated (RFC 3339 format).

## Import

Flexible IP Mac Addresses can be imported using the `{zone}/{id}`, e.g.

```bash
$ terraform import scaleway_flexible_ip_mac_address.main fr-par-1/11111111-1111-1111-1111-111111111111
```
