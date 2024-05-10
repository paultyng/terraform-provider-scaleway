---
subcategory: "IPAM"
page_title: "Scaleway: scaleway_ipam_ips"
---

# scaleway_ipam_ips

Gets information about multiple IPs managed by IPAM service.

## Examples

### By tag

```terraform
data "scaleway_ipam_ips" "by_tag" {
  tags = ["tag"]
}
```

### By type and resource

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

resource "scaleway_redis_cluster" "redis01" {
  name         = "my_redis_cluster"
  version      = "7.0.5"
  node_type    = "RED1-XS"
  user_name    = "my_initial_user"
  password     = "thiZ_is_v&ry_s3cret"
  cluster_size = 3
  private_network {
    id = scaleway_vpc_private_network.pn01.id
  }
}

data "scaleway_ipam_ips" "by_type_and_resource" {
  type = "ipv4"
  resource {
    id   = scaleway_redis_cluster.redis01.id
    type = "redis_cluster"
  }
}
```

## Argument Reference

- `type` - (Optional) The type of IP used as filter (ipv4, ipv6).

- `private_network_id` - (Optional) The ID of the private network used as filter.

- `resource` - (Optional) Filter by resource ID, type or name.
    - `id` - The ID of the resource that the IP is bound to.
    - `type` - The type of the resource to get the IP from. [Documentation](https://pkg.go.dev/github.com/scaleway/scaleway-sdk-go@master/api/ipam/v1#pkg-constants) with type list.
    - `name` - The name of the resource to get the IP from.

- `mac_address` - (Optional) The Mac Address used as filter.

- `tags` (Optional) The tags used as filter.

- `attached` - (Optional) Defines whether to filter only for IPs which are attached to a resource.

- `zonal` - (Optional) Only IPs that are zonal, and in this zone, will be returned.

- `region` - (Optional) The region used as filter.

- `project_id` - (Optional) The ID of the project used as filter.

- `organization_id` - (Optional) The ID of the organization used as filter.

## Attributes Reference

In addition to all above arguments, the following attributes are exported:

- `id` - The region of the IPS
- `ips` - List of found IPs
    - `id` - The ID of the IP.

    ~> **Important:** IPAM IPs' IDs are [regional](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{region}/{id}`, e.g. `fr-par/11111111-1111-1111-1111-111111111111`

    - `address` - The Scaleway internal IP address of the server.
    - `resource` - The list of public IPs of the server.
        - `id` - The ID of the resource.
        - `type` - The type of resource.
        - `mac_address` - The mac address.
        - `name` - The name of the resource.
    - `tags` - The tags associated with the IP.
    - `created_at` - The date and time of the creation of the IP.
    - `updated_at` - The date and time of the last update of the IP.
    - `zone` - The [zone](../guides/regions_and_zones.md#zones) in which the IP is.
    - `region` - The [region](../guides/regions_and_zones.md#regions) in which the IP is.
    - `project_id` - The ID of the project the server is associated with.
  
