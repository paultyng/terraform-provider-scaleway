---
subcategory: "Redis"
page_title: "Scaleway: scaleway_redis_instance"
---

# scaleway_redis_cluster

Gets information about a Redis cluster. For further information check our [api documentation](https://developers.scaleway.com/en/products/redis/api/v1alpha1/#clusters-a85816)

## Example Usage

```hcl
# Get info by name
data "scaleway_redis_cluster" "my_cluster" {
  name = "foobar"
}

# Get info by cluster ID
data "scaleway_redis_cluster" "my_cluster" {
  cluster_id = "11111111-1111-1111-1111-111111111111"
}
```

## Argument Reference

- `name` - (Optional) The name of the Redis cluster.
  Only one of the `name` and `cluster_id` should be specified.

- `cluster_id` - (Optional) The Redis cluster ID.
  Only one of the `name` and `cluster_id` should be specified.

- `zone` - (Default to [provider](../index.md) `region`) The [zone](../guides/regions_and_zones.md#zones) in which the server exists.

- `project_id` - (Optional) The ID of the project the Redis cluster is associated with.

## Attributes Reference

In addition to all above arguments, the following attributes are exported:

- `id` - The ID of the Redis cluster.

~> **Important:** Redis clusters' IDs are [zoned](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{zone}/{id}`, e.g. `fr-par-1/11111111-1111-1111-1111-111111111111`
