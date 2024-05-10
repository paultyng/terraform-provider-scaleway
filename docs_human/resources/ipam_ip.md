---
subcategory: "IPAM"
page_title: "Scaleway: scaleway_ipam_ip"
---

# Resource: scaleway_ipam_ip

Books and manages Scaleway IPAM IPs.

## Example Usage

### Basic

```terraform
resource "scaleway_vpc" "vpc01" {
  name = "my vpc"
}

resource "scaleway_vpc_private_network" "pn01" {
  vpc_id = scaleway_vpc.vpc01.id
  ipv4_subnet {
    subnet = "172.16.32.0/22"
  }
}

resource "scaleway_ipam_ip" "ip01" {
  source {
    private_network_id = scaleway_vpc_private_network.pn01.id
  }
}
```

### Request a specific IPv4

```terraform
resource "scaleway_vpc" "vpc01" {
  name = "my vpc"
}

resource "scaleway_vpc_private_network" "pn01" {
  vpc_id = scaleway_vpc.vpc01.id
  ipv4_subnet {
    subnet = "172.16.32.0/22"
  }
}

resource "scaleway_ipam_ip" "ip01" {
  address = "172.16.32.7"
  source {
    private_network_id = scaleway_vpc_private_network.pn01.id
  }
}
```

### Request an IPv6

```terraform
resource "scaleway_vpc" "vpc01" {
  name = "my vpc"
}

resource "scaleway_vpc_private_network" "pn01" {
  vpc_id = scaleway_vpc.vpc01.id
  ipv6_subnets {
    subnet = "fd46:78ab:30b8:177c::/64"
  }
}

resource "scaleway_ipam_ip" "ip01" {
  is_ipv6 = true
  source {
    private_network_id = scaleway_vpc_private_network.pn01.id
  }
}
```

## Argument Reference

The following arguments are supported:

- `address` - (Optional) Request a specific IP in the requested source pool.
- `tags` - (Optional) The tags associated with the IP.
- `source` - (Required) The source in which to book the IP.
    - `zonal` - The zone the IP lives in if the IP is a public zoned one
    - `private_network_id` - The private network the IP lives in if the IP is a private IP.
    - `subnet_id` - The private network subnet the IP lives in if the IP is a private IP in a private network.
- `is_ipv6` - (Optional) Defines whether to request an IPv6 instead of an IPv4.
- `region` - (Defaults to [provider](../index.md#region) `region`) The [region](../guides/regions_and_zones.md#regions) of the IP.
- `project_id` - (Defaults to [provider](../index.md#project_id) `project_id`) The ID of the project the IP is associated with.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the IP in IPAM.
- `resource` - The IP resource.
    - `id` - The ID of the resource that the IP is bound to.
    - `type` - The type of resource the IP is attached to.
    - `name` - The name of the resource the IP is attached to.
    - `mac_address` - The MAC Address of the resource the IP is attached to.
- `reverses` - The reverses DNS for this IP.
    - `hostname` The reverse domain name.
    - `address` The IP corresponding to the hostname.
- `created_at` - Date and time of IP's creation (RFC 3339 format).
- `updated_at` - Date and time of IP's last update (RFC 3339 format).
- `zone` - The zone of the IP.

~> **Important:** IPAM IPs' IDs are [regional](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{region}/{id}`, e.g. `fr-par/11111111-1111-1111-1111-111111111111

## Import

IPAM IPs can be imported using the `{region}/{id}`, e.g.

```bash
$ terraform import scaleway_ipam_ip.ip_demo fr-par/11111111-1111-1111-1111-111111111111
```
