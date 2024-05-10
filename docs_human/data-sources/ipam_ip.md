---
subcategory: "IPAM"
page_title: "Scaleway: scaleway_ipam_ip"
---

# scaleway_ipam_ip

Gets information about IP managed by IPAM service. IPAM service is used for dhcp bundled in VPCs' private networks.

## Examples

### Instance Private Network IP

```hcl
# Get info by ipam ip id
data "scaleway_ipam_ip" "by_id" {
  ipam_ip_id = "11111111-1111-1111-1111-111111111111"
}

# Get Instance IP in a private network
resource "scaleway_instance_private_nic" "nic" {
  server_id = scaleway_instance_server.server.id
  private_network_id = scaleway_vpc_private_network.pn.id
}

# Find server private IPv4 using private-nic mac address
data "scaleway_ipam_ip" "by_mac" {
  mac_address = scaleway_instance_private_nic.nic.mac_address
  type = "ipv4"
}

# Find server private IPv4 using private-nic id
data "scaleway_ipam_ip" "by_id" {
  resource {
    id = scaleway_instance_private_nic.nic.id
    type = "instance_private_nic"
  }
  type = "ipv4"
}

# Find the private IPv4 using resource name
resource "scaleway_vpc_private_network" "pn" {}

resource "scaleway_rdb_instance" "main" {
  name           = "test-rdb"
  node_type      = "DB-DEV-S"
  engine         = "PostgreSQL-15"
  is_ha_cluster  = true
  disable_backup = true
  user_name      = "my_initial_user"
  password       = "thiZ_is_v&ry_s3cret"
  private_network {
    pn_id = scaleway_vpc_private_network.pn.id
  }
}

data "scaleway_ipam_ip" "by_name" {
  resource {
    name = scaleway_rdb_instance.main.name
    type = "rdb_instance"
  }
  type = "ipv4"
}

```

## Argument Reference

- `ipam_ip_id` - (Optional) The IPAM IP ID. Cannot be used with the rest of the arguments.

- `type` - (Optional) The type of IP to search for (ipv4, ipv6). Cannot be used with `ipam_ip_id`.

- `private_network_id` - (Optional) The ID of the private network the IP belong to. Cannot be used with `ipam_ip_id`.

- `resource` - (Optional) Filter by resource ID, type or name. Cannot be used with `ipam_ip_id`.
If specified, `type` is required, and at least one of `id` or `name` must be set.
    - `id` - The ID of the resource that the IP is bound to.
    - `type` - The type of the resource to get the IP from. [Documentation](https://pkg.go.dev/github.com/scaleway/scaleway-sdk-go@master/api/ipam/v1#pkg-constants) with type list.
    - `name` - The name of the resource to get the IP from.

- `mac_address` - (Optional) The Mac Address linked to the IP. Cannot be used with `ipam_ip_id`.

- `region` - (Defaults to [provider](../index.md#zone) `region`) The [region](../guides/regions_and_zones.md#regions) in which the IP exists.

- `tags` (Optional) The tags associated with the IP. Cannot be used with `ipam_ip_id`.
  As datasource only returns one IP, the search with given tags must return only one result.

- `zonal` - (Optional) Only IPs that are zonal, and in this zone, will be returned.

- `attached` - (Optional) Defines whether to filter only for IPs which are attached to a resource. Cannot be used with `ipam_ip_id`.

- `project_id` - (Defaults to [provider](../index.md#project_id) `project_id`) The ID of the project the IP is associated with.

- `organization_id` - (Defaults to [provider](../index.md#organization_id) `organization_id`) The ID of the organization the IP is in.

## Attributes Reference

In addition to all above arguments, the following attributes are exported:

- `id` - The ID of the IP in IPAM.
- `address` - The IP address.
- `address_cidr` - the IP address with a CIDR notation.
