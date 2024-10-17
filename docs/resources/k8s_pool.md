---
subcategory: "Kubernetes"
page_title: "Scaleway: scaleway_k8s_pool"
---

# Resource: scaleway_k8s_pool

Creates and manages Scaleway Kubernetes cluster pools. For more information, see [the documentation](https://www.scaleway.com/en/developers/api/kubernetes/).

## Example Usage

### Basic

```terraform
resource "scaleway_k8s_cluster" "jack" {
  name    = "jack"
  version = "1.24.3"
  cni     = "cilium"
}

resource "scaleway_k8s_pool" "bill" {
  cluster_id         = scaleway_k8s_cluster.jack.id
  name               = "bill"
  node_type          = "DEV1-M"
  size               = 3
  min_size           = 0
  max_size           = 10
  autoscaling        = true
  autohealing        = true
  container_runtime  = "containerd"
  placement_group_id = "1267e3fd-a51c-49ed-ad12-857092ee3a3d"
}
```

## Argument Reference

The following arguments are supported:

- `cluster_id` - (Required) The ID of the Kubernetes cluster on which this pool will be created.

- `name` - (Required) The name for the pool.
~> **Important:** Updates to this field will recreate a new resource.

- `node_type` - (Required) The commercial type of the pool instances. Instances with insufficient memory are not eligible (DEV1-S, PLAY2-PICO, STARDUST). `external` is a special node type used to provision from other Cloud providers.

~> **Important:** Updates to this field will recreate a new resource.

- `size` - (Required) The size of the pool.
~> **Important:** This field will only be used at creation if autoscaling is enabled.

- `min_size` - (Defaults to `1`) The minimum size of the pool, used by the autoscaling feature.

- `max_size` - (Defaults to `size`) The maximum size of the pool, used by the autoscaling feature.

- `tags` - (Optional) The tags associated with the pool.
  > Note: As mentionned in [this document](https://github.com/scaleway/scaleway-cloud-controller-manager/blob/master/docs/tags.md#taints), taints of a pool's nodes are applied using tags. (Example: "taint=taintName=taineValue:Effect")

- `placement_group_id` - (Optional) The [placement group](https://www.scaleway.com/en/developers/api/instance/#path-placement-groups-create-a-placement-group) the nodes of the pool will be attached to.
~> **Important:** Updates to this field will recreate a new resource.

- `autoscaling` - (Defaults to `false`) Enables the autoscaling feature for this pool.
~> **Important:** When enabled, an update of the `size` will not be taken into account.

- `autohealing` - (Defaults to `false`) Enables the autohealing feature for this pool.

- `container_runtime` - (Defaults to `containerd`) The container runtime of the pool.
~> **Important:** Updates to this field will recreate a new resource.

- `kubelet_args` - (Optional) The Kubelet arguments to be used by this pool

- `upgrade_policy` - (Optional) The Pool upgrade policy

    - `max_surge` - (Defaults to `0`) The maximum number of nodes to be created during the upgrade

    - `max_unavailable` - (Defaults to `1`) The maximum number of nodes that can be not ready at the same time

- `root_volume_type` - (Optional) System volume type of the nodes composing the pool

- `root_volume_size_in_gb` - (Optional) The size of the system volume of the nodes in gigabyte

- `zone` - (Defaults to [provider](../index.md#zone) `zone`) The [zone](../guides/regions_and_zones.md#regions) in which the pool should be created.
~> **Important:** Updates to this field will recreate a new resource.

- `region` - (Defaults to [provider](../index.md#region) `region`) The [region](../guides/regions_and_zones.md#regions) in which the pool should be created.

- `wait_for_pool_ready` - (Defaults to `false`) Whether to wait for the pool to be ready.

- `public_ip_disabled` - (Defaults to `false`) Defines if the public IP should be removed from Nodes. To use this feature, your Cluster must have an attached [Private Network](vpc_private_network.md) set up with a [Public Gateway](vpc_public_gateway.md).
~> **Important:** Updates to this field will recreate a new resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the pool.

~> **Important:** Kubernetes clusters pools' IDs are [regional](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{region}/{id}`, e.g. `fr-par/11111111-1111-1111-1111-111111111111`

- `status` - The status of the pool.
- `nodes` - (List of) The nodes in the default pool.
    - `name` - The name of the node.
    - `public_ip` - The public IPv4. (Deprecated, Please use the official Kubernetes provider and the kubernetes_nodes data source)
    - `public_ip_v6` - The public IPv6. (Deprecated, Please use the official Kubernetes provider and the kubernetes_nodes data source)
    - `status` - The status of the node.
- `created_at` - The creation date of the pool.
- `updated_at` - The last update date of the pool.
- `version` - The version of the pool.
- `current_size` - The size of the pool at the time the terraform state was updated.
- `private_ip` - The list of private IP addresses associated with the resource.
    - `id` - The ID of the IP address resource.
    - `address` - The private IP address.

## Zone

The option `zone` indicate where you the resource of your pool should be created, and it could be different from `region`

Please note that a pool belongs to only one cluster, in the same region.`region`.

## Placement Group

If you are working with cluster type `multicloud` please set the `zone` where your placement group is e.g:

```terraform
resource "scaleway_instance_placement_group" "placement_group" { 
  name        = "pool-placement-group"
  policy_type = "max_availability"
  policy_mode = "optional"
  zone        = "nl-ams-1"
}

resource "scaleway_k8s_pool" "pool" {
  name               = "placement_group"
  cluster_id         = scaleway_k8s_cluster.cluster.id
  node_type          = "gp1_xs"
  placement_group_id = scaleway_instance_placement_group.placement_group.id
  size               = 1
  region             = scaleway_k8s_cluster.cluster.region
  zone               = scaleway_instance_placement_group.placement_group.zone
}

resource "scaleway_k8s_cluster" "cluster" {
  name     = "placement_group"
  cni      = "kilo"
  version  = "%s"
  tags     = [ "terraform-test", "scaleway_k8s_cluster", "placement_group" ]
  region   = "fr-par"
  type     = "multicloud"
}
```

## Import

Kubernetes pools can be imported using the `{region}/{id}`, e.g.

```bash
terraform import scaleway_k8s_pool.mypool fr-par/11111111-1111-1111-1111-111111111111
```

## Changing the node-type of a pool

As your needs evolve, you can migrate your workflow from one pool to another.
Pools have a unique name, and they also have an immutable node type.
Just changing the pool node type will recreate a new pool which could lead to service disruption.
To migrate your application with as little downtime as possible we recommend using the following workflow:

### General workflow to upgrade a pool

- Create a new pool with a different name and the type you target.
- Use [`kubectl drain`](https://kubernetes.io/docs/reference/generated/kubectl/kubectl-commands#drain) on nodes composing your old pool to drain the remaining workflows of this pool.
  Normally it should transfer your workflows to the new pool. Check out the official documentation about [how to safely drain your nodes](https://kubernetes.io/docs/tasks/administer-cluster/safely-drain-node/).
- Delete the old pool from your terraform configuration.

### Using a composite name to force creation of a new pool when a variable updates

If you want to have a new pool created when a variable changes, you can use a name derived from node type such as:

```terraform
resource "scaleway_k8s_pool" "kubernetes_cluster_workers_1" {
  cluster_id    = scaleway_k8s_cluster.kubernetes_cluster.id
  name          = "${var.kubernetes_cluster_id}_${var.node_type}_1"
  node_type     = "${var.node_type}"

  # use Scaleway built-in cluster autoscaler
  autoscaling         = true
  autohealing         = true
  size                = "5"
  min_size            = "5"
  max_size            = "10"
  wait_for_pool_ready = true
}
```

Thanks to [@deimosfr](https://github.com/deimosfr) for the contribution.
