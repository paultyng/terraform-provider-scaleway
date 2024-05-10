---
subcategory: "VPC"
page_title: "Scaleway: scaleway_vpc_public_gateway_dhcp"
---

# scaleway_vpc_public_gateway_dhcp  

Gets information about a public gateway DHCP.

## Example Usage

```hcl
resource "scaleway_vpc_public_gateway_dhcp" "main" {
  subnet = "192.168.0.0/24"
}

data "scaleway_vpc_public_gateway_dhcp" "dhcp_by_id" {
    dhcp_id = "${scaleway_vpc_public_gateway_dhcp.main.id}"
}
```

## Argument Reference


## Attributes Reference

In addition to all above arguments, the following attributes are exported:

- `id` - The ID of the public gateway DHCP config.

~> **Important:** Public gateway DHCP configs' IDs are [zoned](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{zone}/{id}`, e.g. `fr-par-1/11111111-1111-1111-1111-111111111111`

