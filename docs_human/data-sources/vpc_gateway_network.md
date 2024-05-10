---
subcategory: "VPC"
page_title: "Scaleway: scaleway_vpc_gateway_network"
---

# scaleway_vpc_gateway_network

Gets information about a gateway network.

## Example Usage

```hcl
resource "scaleway_vpc_gateway_network" "main" {
  gateway_id = scaleway_vpc_public_gateway.pg01.id
  private_network_id = scaleway_vpc_private_network.pn01.id
  dhcp_id = scaleway_vpc_public_gateway_dhcp.dhcp01.id
  cleanup_dhcp       = true
  enable_masquerade  = true
}

data scaleway_vpc_gateway_network by_id {
    gateway_network_id = scaleway_vpc_gateway_network.main.id
}

data scaleway_vpc_gateway_network by_gateway_and_pn {
    gateway_id = scaleway_vpc_public_gateway.pg01.id
    private_network_id = scaleway_vpc_private_network.pn01.id
}
```

## Argument Reference

* `gateway_network_id` - (Optional) ID of the gateway network.

~> Only one of `gateway_network_id` or filters should be specified. You can use all the filters you want.

* `gateway_id` - (Optional) ID of the public gateway the gateway network is linked to
* `private_network_id` - (Optional) ID of the private network the gateway network is linked to
* `enable_masquerade` - (Optional) If masquerade is enabled on requested network
* `dhcp_id` - (Optional) ID of the public gateway DHCP config

## Attributes Reference

Exported attributes are the attributes of the [resource](../resources/vpc_gateway_network.md)
