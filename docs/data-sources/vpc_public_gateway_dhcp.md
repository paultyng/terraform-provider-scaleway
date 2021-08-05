---
page_title: "Scaleway: scaleway_vpc_public_gateway_dhcp"
description: |-
  Get information about Scaleway VPC Public Gateway DHCP.
---

# scaleway_vpc_public_gateway_dhcp  

Gets information about a public gateway DHCP.

## Example Usage

```hcl
resource "scaleway_vpc_public_gateway_dhcp" "main" {
}

data "scaleway_vpc_public_gateway_dhcp" "dhcp_by_id" {
    dhcp_id = "${scaleway_vpc_public_gateway_dhcp.main.id}"
}
```

## Argument Reference


## Attributes Reference

`id` is set to the ID of the found public gateway DHCP config. Addition attributes are
exported.
