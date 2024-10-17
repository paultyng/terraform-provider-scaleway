---
subcategory: "Instances"
page_title: "Scaleway: scaleway_instance_private_nic"
---

# Resource: scaleway_instance_private_nic

Creates and manages Scaleway Instance Private NICs. For more information, see
[the documentation](https://www.scaleway.com/en/developers/api/instance/#path-private-nics-list-all-private-nics).

## Example Usage

### Basic

```terraform
resource "scaleway_instance_private_nic" "pnic01" {
  server_id          = "fr-par-1/11111111-1111-1111-1111-111111111111"
  private_network_id = "fr-par-1/aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
}
```

### With zone

```terraform
resource "scaleway_vpc_private_network" "pn01" {
  name   = "private_network_instance"
  region = "fr-par"
}

resource "scaleway_instance_server" "base" {
  image = "ubuntu_jammy"
  type  = "DEV1-S"
  zone  = scaleway_vpc_private_network.pn01.zone
}

resource "scaleway_instance_private_nic" "pnic01" {
  server_id          = scaleway_instance_server.base.id
  private_network_id = scaleway_vpc_private_network.pn01.id
  zone               = scaleway_vpc_private_network.pn01.zone
}
```

### With IPAM IP IDs

```terraform
resource "scaleway_vpc" "vpc01" {
  name   = "vpc_instance"
}

resource "scaleway_vpc_private_network" "pn01" {
  name   = "private_network_instance"
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

resource "scaleway_instance_server" "server01" {
  image = "ubuntu_focal"
  type  = "PLAY2-MICRO"
}

resource "scaleway_instance_private_nic" "pnic01" {
  private_network_id = scaleway_vpc_private_network.pn01.id
  server_id          = scaleway_instance_server.server01.id
  ipam_ip_ids        = [scaleway_ipam_ip.ip01.id]
}     
```

## Argument Reference

The following arguments are required:

- `server_id` - (Required) The ID of the server associated with.
- `private_network_id` - (Required) The ID of the private network attached to.
- `ipam_ip_ids` - (Optional) IPAM IDs of a pre-reserved IP addresses to assign to the Instance in the requested private network.
- `tags` - (Optional) The tags associated with the private NIC.
- `zone` - (Defaults to [provider](../index.md#zone) `zone`) The [zone](../guides/regions_and_zones.md#zones) in which the server must be created.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the private NIC.

~> **Important:** Instance private NICs' IDs are [zoned](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{zone}/{id}`, e.g. `fr-par-1/11111111-1111-1111-1111-111111111111`

- `mac_address` - The MAC address of the private NIC.
- `private_ip` - The list of private IP addresses associated with the resource.
    - `id` - The ID of the IP address resource.
    - `address` - The private IP address.

## Import

Private NICs can be imported using the `{zone}/{server_id}/{private_nic_id}`, e.g.

```bash
terraform import scaleway_instance_private_nic.pnic01 fr-par-1/11111111-1111-1111-1111-111111111111/22222222-2222-2222-2222-222222222222
```
