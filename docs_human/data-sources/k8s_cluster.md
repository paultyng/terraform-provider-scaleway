---
subcategory: "Kubernetes"
page_title: "Scaleway: scaleway_k8s_cluster"
---

# scaleway_k8s_cluster

Gets information about a Kubernetes Cluster.

## Example Usage

```hcl
# Get info by cluster name
data "scaleway_k8s_cluster" "my_key" {
  name  = "my-cluster-name"
}

# Get info by cluster id
data "scaleway_k8s_cluster" "my_key" {
  cluster_id = "11111111-1111-1111-1111-111111111111"
}
```

## Argument Reference

- `name` - (Optional) The cluster name. Only one of `name` and `cluster_id` should be specified.

- `cluster_id` - (Optional) The cluster ID. Only one of `name` and `cluster_id` should be specified.

- `region` - (Defaults to [provider](../index.md#region) `region`) The [region](../guides/regions_and_zones.md#regions) in which the cluster exists.

- `project_id` - (Optional) The ID of the project the cluster is associated with.

## Attributes Reference

In addition to all above arguments, the following attributes are exported:

- `id` - The ID of the cluster.

~> **Important:** Kubernetes clusters' IDs are [regional](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{region}/{id}`, e.g. `fr-par/11111111-1111-1111-1111-111111111111`

- `created_at` - The creation date of the cluster.

- `updated_at` - The last update date of the cluster.

- `apiserver_url` - The URL of the Kubernetes API server.

- `wildcard_dns` - The DNS wildcard that points to all ready nodes.

- `kubeconfig`

    - `config_file` - The raw kubeconfig file.

    - `host` - The URL of the Kubernetes API server.

    - `cluster_ca_certificate` - The CA certificate of the Kubernetes API server.

    - `token` - The token to connect to the Kubernetes API server.

- `status` - The status of the Kubernetes cluster.

- `type` - The type of the Kubernetes cluster.

- `upgrade_available` - True if a newer Kubernetes version is available.

- `description` - A description for the Kubernetes cluster.

- `version` - The version of the Kubernetes cluster.

- `cni` - The Container Network Interface (CNI) for the Kubernetes cluster.

- `tags` - The tags associated with the Kubernetes cluster.

- `autoscaler_config` - The configuration options for the [Kubernetes cluster autoscaler](https://github.com/kubernetes/autoscaler/tree/master/cluster-autoscaler).

    - `disable_scale_down` - True if the scale down feature of the autoscaler is disabled.

    - `scale_down_delay_after_add` - The duration after scale up that scale down evaluation resumes.

    - `scale_down_unneeded_time` - The duration a node should be unneeded before it is eligible for scale down.

    - `estimator` - The type of resource estimator used in scale up.

    - `expander` - The type of node group expander be used in scale up.

    - `ignore_daemonsets_utilization` - True if ignoring DaemonSet pods when calculating resource utilization for scaling down is enabled.

    - `balance_similar_node_groups` - True if detecting similar node groups and balance the number of nodes between them is enabled.

    - `expendable_pods_priority_cutoff` - Pods with priority below cutoff will be expendable. They can be killed without any consideration during scale down and they don't cause scale up. Pods with null priority (PodPriority disabled) are non expendable.

- `auto_upgrade` - The auto upgrade configuration.

    - `enable` - True if Kubernetes patch version auto upgrades is enabled.

    - `maintenance_window_start_hour` - The start hour (UTC) of the 2-hour auto upgrade maintenance window (0 to 23).

    - `maintenance_window_day` - The day of the auto upgrade maintenance window (`monday` to `sunday`, or `any`).

- `feature_gates` - The list of [feature gates](https://kubernetes.io/docs/reference/command-line-tools-reference/feature-gates/) enabled on the cluster.

- `admission_plugins` - The list of [admission plugins](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/) enabled on the cluster.

- `private_network_id` - The ID of the private network of the cluster.

- `region` - The [region](../guides/regions_and_zones.md#regions) in which the cluster is.

- `organization_id` - The ID of the organization the cluster is associated with.
