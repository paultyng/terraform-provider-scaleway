---
subcategory: "Databases"
page_title: "Scaleway: scaleway_rdb_read_replica"
---

# Resource: scaleway_rdb_read_replica

Creates and manages Scaleway Database read replicas.
For more information, see [the documentation](https://developers.scaleway.com/en/products/rdb/api).

## Example Usage

### Basic

```terraform
resource "scaleway_rdb_instance" "instance" {
  name           = "test-rdb-rr-update"
  node_type      = "db-dev-s"
  engine         = "PostgreSQL-14"
  is_ha_cluster  = false
  disable_backup = true
  user_name      = "my_initial_user"
  password       = "thiZ_is_v&ry_s3cret"
  tags           = ["terraform-test", "scaleway_rdb_read_replica", "minimal"]
}

resource scaleway_rdb_read_replica "replica" {
  instance_id = scaleway_rdb_instance.instance.id
  direct_access {}
}
```

### Private network with static endpoint

```terraform
resource "scaleway_rdb_instance" "instance" {
  name           = "rdb_instance"
  node_type      = "db-dev-s"
  engine         = "PostgreSQL-14"
  is_ha_cluster  = false
  disable_backup = true
  user_name      = "my_initial_user"
  password       = "thiZ_is_v&ry_s3cret"
}

resource "scaleway_vpc_private_network" "pn" {}

resource "scaleway_rdb_read_replica" "replica" {
  instance_id = scaleway_rdb_instance.instance.id
  private_network {
    private_network_id = scaleway_vpc_private_network.pn.id
    service_ip         = "192.168.1.254/24"
    # enable_ipam = false
  }
}
```

### Private network with IPAM

```terraform
resource "scaleway_rdb_instance" "instance" {
  name           = "rdb_instance"
  node_type      = "db-dev-s"
  engine         = "PostgreSQL-14"
  is_ha_cluster  = false
  disable_backup = true
  user_name      = "my_initial_user"
  password       = "thiZ_is_v&ry_s3cret"
}

resource "scaleway_vpc_private_network" "pn" {}

resource "scaleway_rdb_read_replica" "replica" {
  instance_id = scaleway_rdb_instance.instance.id
  private_network {
    private_network_id = scaleway_vpc_private_network.pn.id
    enable_ipam = true
  }
}
```

## Argument Reference

The following arguments are supported:

- `instance_id` - (Required) UUID of the rdb instance.

~> **Important:** The replica musts contains at least one of `direct_access` or `private_network`. It can contain both.

- `direct_access` - (Optional) Creates a direct access endpoint to rdb replica.

- `private_network` - (Optional) Create an endpoint in a private network.
    - `private_network_id` - (Required) UUID of the private network to be connected to the read replica.
    - `service_ip` - (Optional) The IP network address within the private subnet. This must be an IPv4 address with a CIDR notation. If not set, The IP network address within the private subnet is determined by the IP Address Management (IPAM) service.
    - `enable_ipam` - (Optional) If true, the IP network address within the private subnet is determined by the IP Address Management (IPAM) service.

- `same_zone` - (Defaults to `true`) Defines whether to create the replica in the same availability zone as the main instance nodes or not.

- `region` - (Defaults to [provider](../index.md#arguments-reference) `region`) The [region](../guides/regions_and_zones.md#regions)
  in which the Database read replica should be created.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the Database read replica.

~> **Important:** Database read replicas' IDs are [regional](../guides/regions_and_zones.md#resource-ids), which means
they are of the form `{region}/{id}`, e.g. `fr-par/11111111-1111-1111-1111-111111111111`

- `direct_access` - List of load balancer endpoints of the database read replica.
    - `endpoint_id` - The ID of the endpoint of the read replica.
    - `ip` - IPv4 address of the endpoint (IP address). Only one of ip and hostname may be set.
    - `port` - TCP port of the endpoint.
    - `name` - Name of the endpoint.
    - `hostname` - Hostname of the endpoint. Only one of ip and hostname may be set.
- `private_network` - List of private networks endpoints of the database read replica.
    - `endpoint_id` - The ID of the endpoint of the read replica.
    - `ip` - IPv4 address of the endpoint (IP address). Only one of ip and hostname may be set.
    - `port` - TCP port of the endpoint.
    - `name` - Name of the endpoint.
    - `hostname` - Hostname of the endpoint. Only one of ip and hostname may be set.
    - `enable_ipam` - Indicates whether the IP is managed by IPAM.

## Import

Database Read replica can be imported using the `{region}/{id}`, e.g.

```bash
$ terraform import scaleway_rdb_read_replica.rr fr-par/11111111-1111-1111-1111-111111111111
```
