---
subcategory: "Databases"
page_title: "Scaleway: scaleway_documentdb_read_replica"
---

# Resource: scaleway_documentdb_read_replica

Creates and manages Scaleway DocumentDB Database read replicas.
For more information, see [the documentation](https://www.scaleway.com/en/developers/api/document_db/).

## Example Usage

### Basic

```terraform
resource scaleway_documentdb_read_replica "replica" {
  instance_id     = "11111111-1111-1111-1111-111111111111"
  direct_access {}
}
```

### Private network

```terraform
resource "scaleway_vpc_private_network" "pn" {}

resource "scaleway_documentdb_read_replica" "replica" {
  instance_id = scaleway_rdb_instance.instance.id
  private_network {
    private_network_id = scaleway_vpc_private_network.pn.id
    service_ip         = "192.168.1.254/24" // omit this attribute if private IP is determined by the IP Address Management (IPAM)
  }
}
```

## Argument Reference

The following arguments are supported:

- `instance_id` - (Required) UUID of the documentdb instance.

~> **Important:** The replica musts contains at least one of `direct_access` or `private_network`. It can contain both.

- `direct_access` - (Optional) Creates a direct access endpoint to documentdb replica.

- `private_network` - (Optional) Create an endpoint in a private network.
    - `private_network_id` - (Required) UUID of the private network to be connected to the read replica.
    - `service_ip` - (Optional) The IP network address within the private subnet. This must be an IPv4 address with a
      CIDR notation. The IP network address within the private subnet is determined by the IP Address Management (IPAM)
      service if not set.

- `region` - (Defaults to [provider](../index.md#region) `region`) The [region](../guides/regions_and_zones.md#regions)
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

## Import

Database Read replica can be imported using the `{region}/{id}`, e.g.

```bash
$ terraform import scaleway_documentdb_read_replica.rr fr-par/11111111-1111-1111-1111-111111111111
```
