---
subcategory: "VPC"
page_title: "Scaleway: scaleway_vpc_gateway_network"
---

# Resource: scaleway_vpc_gateway_network

Creates and manages GatewayNetworks (connections between a Public Gateway and a Private Network).

It allows the attachment of Private Networks to Public Gateways and DHCP configurations.
For more information, see [the API documentation](https://www.scaleway.com/en/developers/api/public-gateway/#step-3-attach-private-networks-to-the-vpc-public-gateway).

## Example Usage

### Create a GatewayNetwork with IPAM configuration

```terraform
resource scaleway_vpc vpc01 {
  name = "my vpc"
}

resource scaleway_vpc_private_network pn01 {
  name = "pn_test_network"
  ipv4_subnet {
    subnet = "172.16.64.0/22"
  }
  vpc_id = scaleway_vpc.vpc01.id
}

resource scaleway_vpc_public_gateway pg01 {
  name = "foobar"
  type = "VPC-GW-S"
}

resource scaleway_vpc_gateway_network main {
  gateway_id = scaleway_vpc_public_gateway.pg01.id
  private_network_id = scaleway_vpc_private_network.pn01.id
  enable_masquerade = true
  ipam_config {
    push_default_route = true
  }
}
```

### Create a GatewayNetwork with a booked IPAM IP

```terraform
resource scaleway_vpc vpc01 {
  name = "my vpc"
}

resource scaleway_vpc_private_network pn01 {
  name = "pn_test_network"
  ipv4_subnet {
    subnet = "172.16.64.0/22"
  }
  vpc_id = scaleway_vpc.vpc01.id
}

resource "scaleway_ipam_ip" "ip01" {
  address = "172.16.64.7"
  source {
    private_network_id = scaleway_vpc_private_network.pn01.id
  }
}

resource scaleway_vpc_public_gateway pg01 {
  name = "foobar"
  type = "VPC-GW-S"
}

resource scaleway_vpc_gateway_network main {
  gateway_id = scaleway_vpc_public_gateway.pg01.id
  private_network_id = scaleway_vpc_private_network.pn01.id
  enable_masquerade = true
  ipam_config {
    push_default_route = true
    ipam_ip_id = scaleway_ipam_ip.ip01.id
  }
}
```

### Create a GatewayNetwork with DHCP

```terraform
resource "scaleway_vpc_private_network" "pn01" {
  name = "pn_test_network"
}

resource "scaleway_vpc_public_gateway_ip" "gw01" {
}

resource "scaleway_vpc_public_gateway_dhcp" "dhcp01" {
  subnet = "192.168.1.0/24"
  push_default_route = true
}

resource "scaleway_vpc_public_gateway" "pg01" {
  name = "foobar"
  type = "VPC-GW-S"
  ip_id = scaleway_vpc_public_gateway_ip.gw01.id
}

resource "scaleway_vpc_gateway_network" "main" {
  gateway_id = scaleway_vpc_public_gateway.pg01.id
  private_network_id = scaleway_vpc_private_network.pn01.id
  dhcp_id = scaleway_vpc_public_gateway_dhcp.dhcp01.id
  cleanup_dhcp       = true
  enable_masquerade  = true
}
```

### Create a GatewayNetwork with a static IP address

```terraform
resource scaleway_vpc_private_network pn01 {
  name = "pn_test_network"
}

resource scaleway_vpc_public_gateway pg01 {
  name = "foobar"
  type = "VPC-GW-S"
}

resource scaleway_vpc_gateway_network main {
  gateway_id = scaleway_vpc_public_gateway.pg01.id
  private_network_id = scaleway_vpc_private_network.pn01.id
  enable_dhcp = false
  enable_masquerade = true
  static_address = "192.168.1.42/24"
}
```

## Argument Reference

The following arguments are supported:

- `gateway_id` - (Required) The ID of the Public Gateway.
- `private_network_id` - (Required) The ID of the Private Network.
- `dhcp_id` - (Required) The ID of the Public Gateway DHCP configuration. Only one of `dhcp_id`, `static_address` and `ipam_config` should be specified.
- `enable_masquerade` - (Defaults to true) Whether masquerade (dynamic NAT) should be enabled on this GatewayNetwork
- `enable_dhcp` - (Defaults to true) WWhether a DHCP configuration should be enabled on this GatewayNetwork. Requires a DHCP ID.
- `cleanup_dhcp` - (Defaults to false) Whether to remove DHCP configuration on this GatewayNetwork upon destroy. Requires DHCP ID.
- `static_address` - Enable DHCP configration on this GatewayNetwork. Only one of `dhcp_id`, `static_address` and `ipam_config` should be specified.
- `ipam_config` - Auto-configure the GatewayNetwork using Scaleway's IPAM (IP address management service). Only one of `dhcp_id`, `static_address` and `ipam_config` should be specified.
    - `push_default_route` - Defines whether to enable the default route on the GatewayNetwork.
    - `ipam_ip_id` - Use this IPAM-booked IP ID as the Gateway's IP in this Private Network.
- `zone` - (Defaults to [provider](../index.md#zone) `zone`) The [zone](../guides/regions_and_zones.md#zones) in which the gateway network should be created.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the GatewayNetwork

~> **Important:** GatewayNetwork IDs are [zoned](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{zone}/{id}`, e.g. `fr-par-1/11111111-1111-1111-1111-111111111111`

- `mac_address` - The MAC address of the GatewayNetwork.
- `created_at` - The date and time of the creation of the GatewayNetwork.
- `updated_at` - The date and time of the last update of the GatewayNetwork.
- `status` - The status of the Public Gateway's connection to the Private Network.
- `private_ip` - The list of private IP addresses associated with the resource.
    - `id` - The ID of the IP address resource.
    - `address` - The private IP address.

## Import

GatewayNetwork can be imported using `{zone}/{id}`, e.g.

```bash
terraform import scaleway_vpc_gateway_network.main fr-par-1/11111111-1111-1111-1111-111111111111
```

