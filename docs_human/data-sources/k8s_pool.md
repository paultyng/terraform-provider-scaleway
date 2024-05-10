---
subcategory: "Kubernetes"
page_title: "Scaleway: scaleway_k8s_pool"
---

# scaleway_k8s_pool

Gets information about a Kubernetes Cluster's Pool.

## Example Usage

```hcl
# Get info by pokl name (need cluster_id)
data "scaleway_k8s_pool" "my_key" {
  name  = "my-pool-name"
  cluster_id = "11111111-1111-1111-1111-111111111111"
}

# Get info by pool id
data "scaleway_k8s_pool" "my_key" {
  pool_id = "11111111-1111-1111-1111-111111111111"
}
```

## Argument Reference

- `name` - The pool name. Only one of `name` and `pool_id` should be specified. `cluster_id` should be specified with `name`.

- `pool_id` - (Optional) The pool's ID. Only one of `name` and `pool_id` should be specified.

- `cluster_id` - (Optional) The cluster ID. Required when `name` is set.

- `region` - (Defaults to [provider](../index.md#region) `region`) The [region](../guides/regions_and_zones.md#regions) in which the pool exists.

## Attributes Reference

In addition to all above arguments, the following attributes are exported:

- `id` - The ID of the pool.

~> **Important:** Kubernetes clusters pools' IDs are [regional](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{region}/{id}`, e.g. `fr-par/11111111-1111-1111-1111-111111111111`

- `status` - The status of the pool.

- `nodes` - (List of) The nodes in the default pool.

    - `name` - The name of the node.

    - `public_ip` - The public IPv4.

    - `public_ip_v6` - The public IPv6.

    - `status` - The status of the node.

- `created_at` - The creation date of the pool.

- `updated_at` - The last update date of the pool.

- `version` - The version of the pool.

- `current_size` - The size of the pool at the time the terraform state was updated.

- `node_type` - The commercial type of the pool instances.

- `size` - The size of the pool.

- `min_size` - The minimum size of the pool, used by the autoscaling feature.

- `max_size` - The maximum size of the pool, used by the autoscaling feature.

- `tags` - The tags associated with the pool.

- `placement_group_id` - [placement group](https://developers.scaleway.com/en/products/instance/api/#placement-groups-d8f653) the nodes of the pool are attached to.

- `autoscaling` - True if the autoscaling feature is enabled for this pool.

- `autohealing` - True if the autohealing feature is enabled for this pool.

- `container_runtime` - The container runtime of the pool.
