---
subcategory: "VPC"
page_title: "Scaleway: scaleway_vpc_private_network"
---

# Resource: scaleway_vpc_private_network

Creates and manages Scaleway VPC Private Networks.
For more information, see [the documentation](https://developers.scaleway.com/en/products/vpc/api/#private-networks-ac2df4).

## Example Usage

### Basic

```terraform
resource "scaleway_vpc_private_network" "pn_priv" {
    name = "subnet_demo"
    tags = ["demo", "terraform"]
}
```

### With subnets

```terraform
resource "scaleway_vpc_private_network" "pn_priv" {
    name = "subnet_demo"
    tags = ["demo", "terraform"]
    
    ipv4_subnet {
      subnet = "192.168.0.0/24"
    }
    ipv6_subnets {
      subnet = "fd46:78ab:30b8:177c::/64"
    }
    ipv6_subnets {
      subnet = "fd46:78ab:30b8:c7df::/64"
    }
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Optional) The name of the private network. If not provided it will be randomly generated.
- `tags` - (Optional) The tags associated with the private network.
- `project_id` - (Defaults to [provider](../index.md#project_id) `project_id`) The ID of the project the private network is associated with.
- `ipv4_subnet` - (Optional) The IPv4 subnet to associate with the private network.
    - `subnet` - (Optional) The subnet CIDR.
- `ipv6_subnets` - (Optional) The IPv6 subnets to associate with the private network.
    - `subnet` - (Optional) The subnet CIDR.
- `region` - (Defaults to [provider](../index.md#region) `region`) The [region](../guides/regions_and_zones.md#regions) of the private network.
- `vpc_id` - (Optional) The VPC in which to create the private network.
- `is_regional` - (Deprecated) The private networks are necessarily regional now.
- `zone` - (Deprecated) please use `region` instead - (Defaults to [provider](../index.md#zone) `zone`) The [zone](../guides/regions_and_zones.md#zones) in which the private network should be created.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the private network.
- `ipv4_subnet` - The IPv4 subnet associated with the private network.
    - `subnet` - The subnet CIDR.
    - `id` - The subnet ID.
    - `created_at` - The date and time of the creation of the subnet.
    - `updated_at` - The date and time of the last update of the subnet.
    - `address` - The network address of the subnet in dotted decimal notation, e.g., '192.168.0.0' for a '192.168.0.0/24' subnet.
    - `subnet_mask` - The subnet mask expressed in dotted decimal notation, e.g., '255.255.255.0' for a /24 subnet
    - `prefix_length` - The length of the network prefix, e.g., 24 for a 255.255.255.0 mask.
- `ipv6_subnets` - The IPv6 subnets associated with the private network.
    - `subnet` - The subnet CIDR.
    - `id` - The subnet ID.
    - `created_at` - The date and time of the creation of the subnet.
    - `updated_at` - The date and time of the last update of the subnet.
    - `address` - The network address of the subnet in dotted decimal notation, e.g., '192.168.0.0' for a '192.168.0.0/24' subnet.
    - `subnet_mask` - The subnet mask expressed in dotted decimal notation, e.g., '255.255.255.0' for a /24 subnet
    - `prefix_length` - The length of the network prefix, e.g., 24 for a 255.255.255.0 mask.

~> **Important:** Private networks' IDs are [regional](../guides/regions_and_zones.md#resource-ids), which means they are of the form of `{region}/{id}`, e.g. `fr-par/11111111-1111-1111-1111-111111111111

- `organization_id` - The organization ID the private network is associated with.

## Import

Private networks can be imported using the `{region}/{id}`, e.g.

```bash
$ terraform import scaleway_vpc_private_network.vpc_demo fr-par/11111111-1111-1111-1111-111111111111
```
