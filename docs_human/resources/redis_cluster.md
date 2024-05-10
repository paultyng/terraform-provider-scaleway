---
subcategory: "Redis"
page_title: "Scaleway: scaleway_redis_cluster"
---

# Resource: scaleway_redis_cluster

Creates and manages Scaleway Redis Clusters.
For more information, see [the documentation](https://developers.scaleway.com/en/products/redis/api/v1alpha1/).

## Example Usage

### Basic

```terraform
resource "scaleway_redis_cluster" "main" {
  name         = "test_redis_basic"
  version      = "6.2.6"
  node_type    = "RED1-MICRO"
  user_name    = "my_initial_user"
  password     = "thiZ_is_v&ry_s3cret"
  tags         = ["test", "redis"]
  cluster_size = 1
  tls_enabled  = "true"

  acl {
    ip          = "0.0.0.0/0"
    description = "Allow all"
  }
}
```

### With settings

```terraform
resource "scaleway_redis_cluster" "main" {
  name      = "test_redis_basic"
  version   = "6.2.6"
  node_type = "RED1-MICRO"
  user_name = "my_initial_user"
  password  = "thiZ_is_v&ry_s3cret"

  settings = {
    "maxclients"    = "1000"
    "tcp-keepalive" = "120"
  }
}
```

### With a private network

```terraform
resource "scaleway_vpc_private_network" "pn" {
  name = "private-network"
}

resource "scaleway_redis_cluster" "main" {
  name         = "test_redis_endpoints"
  version      = "6.2.6"
  node_type    = "RED1-MICRO"
  user_name    = "my_initial_user"
  password     = "thiZ_is_v&ry_s3cret"
  cluster_size = 1
  private_network {
    id          = "${scaleway_vpc_private_network.pn.id}"
    service_ips = [
      "10.12.1.1/20",
    ]
  }
  depends_on = [
    scaleway_vpc_private_network.pn
  ]
}
```

## Argument Reference

The following arguments are supported:

- `version` - (Required) Redis's Cluster version (e.g. `6.2.6`).

~> **Important:** Updates to `version` will migrate the Redis Cluster to the desired `version`. Keep in mind that you
cannot downgrade a Redis Cluster.

- `node_type` - (Required) The type of Redis Cluster you want to create (e.g. `RED1-M`).

~> **Important:** Updates to `node_type` will migrate the Redis Cluster to the desired `node_type`. Keep in mind that
you cannot downgrade a Redis Cluster.

- `user_name` - (Required) Identifier for the first user of the Redis Cluster.

- `password` - (Required) Password for the first user of the Redis Cluster.

- `name` - (Optional) The name of the Redis Cluster.

- `tags` - (Optional) The tags associated with the Redis Cluster.

- `zone` - (Defaults to [provider](../index.md) `zone`) The [zone](../guides/regions_and_zones.md#zones) in which the
  Redis Cluster should be created.

- `cluster_size` - (Optional) The number of nodes in the Redis Cluster.

~> **Important:** You cannot set `cluster_size` to 2, you either have to choose Standalone mode (1 node) or Cluster mode
which is minimum 3 (1 main node + 2 secondary nodes)

~> **Important:** If you are using the Cluster mode (>=3 nodes), you can set a bigger `cluster_size` than you initially
did, it will migrate the Redis Cluster but keep in mind that you cannot downgrade a Redis Cluster, so setting a smaller
`cluster_size` will destroy and recreate your Cluster.

~> **Important:** If you are using the Standalone mode (1 node), setting a bigger `cluster_size` will destroy and
recreate your Cluster as you will be switching to the Cluster mode.

- `tls_enabled` - (Defaults to false) Whether TLS is enabled or not.

  ~> The changes on `tls_enabled` will force the resource creation.

- `project_id` - (Defaults to [provider](../index.md) `project_id`) The ID of the project the Redis Cluster is
  associated with.

- `acl` - (Optional) List of acl rules, this is cluster's authorized IPs. More details on the [ACL section.](#acl)

- `settings` - (Optional) Map of settings for redis cluster. Available settings can be found by listing redis versions
  with scaleway API or CLI

- `private_network` - (Optional) Describes the private network you want to connect to your cluster. If not set, a public
  network will be provided. More details on the [Private Network section](#private-network)

### ACL

The `acl` block supports:

- `ip` - (Required) The ip range to whitelist
  in [CIDR notation](https://en.wikipedia.org/wiki/Classless_Inter-Domain_Routing#CIDR_notation)
- `description` - (Optional) A text describing this rule. Default description: `Allow IP`

  ~> The `acl` conflict with `private_network`. Only one should be specified.

### Private Network

The `private_network` block supports :

- `id` - (Required) The UUID of the Private Network resource.
- `service_ips` - (Optional) Endpoint IPv4 addresses in [CIDR notation](https://en.wikipedia.org/wiki/Classless_Inter-Domain_Routing#CIDR_notation). You must provide at least one IP per node.
  Keep in mind that in Cluster mode you cannot edit your Private Network after its creation so if you want to be able to
  scale your Cluster horizontally (adding nodes) later, you should provide more IPs than nodes.
  If not set, the IP network address within the private subnet is determined by the IP Address Management (IPAM) service.

~> The `private_network` conflicts with `acl`. Only one should be specified.

~> **Important:** The way to use private networks differs whether you are using Redis in Standalone or Cluster mode.

- Standalone mode (`cluster_size` = 1) : you can attach as many Private Networks as you want (each must be a separate
  block). If you detach your only private network, your cluster won't be reachable until you define a new Private or
  Public Network. You can modify your `private_network` and its specs, you can have both a Private and Public Network side
  by side.

- Cluster mode (`cluster_size` > 2) : you can define a single Private Network as you create your Cluster, you won't be
  able to edit or detach it afterward, unless you create another Cluster. This also means that, if you are using a static
  configuration (`service_ips`), you won't be able to scale your Cluster horizontally (add more nodes) since it would
  require updating the private network to add IPs.
  Your `service_ips` must be listed as follows:

```terraform
  service_ips = [
  "10.12.1.10/20",
  "10.12.1.11/20",
  "10.12.1.12/20",
]
```

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the Redis cluster.

~> **Important:** Redis clusters' IDs are [zoned](../guides/regions_and_zones.md#resource-ids), which means they are of
the form `{zone}/{id}`, e.g. `fr-par-1/11111111-1111-1111-1111-111111111111`

- `public_network` - (Optional) Public network details. Only one of `private_network` and `public_network` may be set.
  ~> The `public_network` block exports:
    - `id` - (Required) The UUID of the endpoint.
    - `ips` - Lis of IPv4 address of the endpoint (IP address).
    - `port` - TCP port of the endpoint.

- `private_network` - List of private networks endpoints of the Redis Cluster.
    - `endpoint_id` - The ID of the endpoint.
    - `zone` - The zone of the private network.

- `created_at` - The date and time of creation of the Redis Cluster.
- `updated_at` - The date and time of the last update of the Redis Cluster.
- `certificate` - The PEM of the certificate used by redis, only when `tls_enabled` is true

## Import

Redis Cluster can be imported using the `{zone}/{id}`, e.g.

```bash
$ terraform import scaleway_redis_cluster.main fr-par-1/11111111-1111-1111-1111-111111111111
```
