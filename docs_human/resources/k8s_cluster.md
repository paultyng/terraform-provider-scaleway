---
subcategory: "Kubernetes"
page_title: "Scaleway: scaleway_k8s_cluster"
---

# Resource: scaleway_k8s_cluster

Creates and manages Scaleway Kubernetes clusters. For more information, see [the documentation](https://developers.scaleway.com/en/products/k8s/api/).

## Example Usage

### Basic

```terraform
resource "scaleway_vpc_private_network" "hedy" {}

resource "scaleway_k8s_cluster" "jack" {
  name    = "jack"
  version = "1.24.3"
  cni     = "cilium"
  private_network_id = scaleway_vpc_private_network.hedy.id
  delete_additional_resources = false
}

resource "scaleway_k8s_pool" "john" {
  cluster_id = scaleway_k8s_cluster.jack.id
  name       = "john"
  node_type  = "DEV1-M"
  size       = 1
}
```

### Multicloud

```terraform
resource "scaleway_k8s_cluster" "henry" {
  name = "henry"
  type = "multicloud"
  version = "1.24.3"
  cni     = "kilo"
  delete_additional_resources = false
}

resource "scaleway_k8s_pool" "friend_from_outer_space" {
  cluster_id = scaleway_k8s_cluster.henry.id
  name = "henry_friend"
  node_type = "external"
  size = 0
  min_size = 0
}
```

For a detailed example of how to add or run Elastic Metal servers instead of instances on your cluster, please refer to [this guide](../guides/multicloud_cluster_with_baremetal_servers.md).

### With additional configuration

```terraform
resource "scaleway_vpc_private_network" "hedy" {}

resource "scaleway_k8s_cluster" "john" {
  name             = "john"
  description      = "my awesome cluster"
  version          = "1.24.3"
  cni              = "calico"
  tags             = ["i'm an awesome tag", "yay"]
  private_network_id = scaleway_vpc_private_network.hedy.id
  delete_additional_resources = false

  autoscaler_config {
    disable_scale_down              = false
    scale_down_delay_after_add      = "5m"
    estimator                       = "binpacking"
    expander                        = "random"
    ignore_daemonsets_utilization   = true
    balance_similar_node_groups     = true
    expendable_pods_priority_cutoff = -5
  }
}

resource "scaleway_k8s_pool" "john" {
  cluster_id  = scaleway_k8s_cluster.john.id
  name        = "john"
  node_type   = "DEV1-M"
  size        = 3
  autoscaling = true
  autohealing = true
  min_size    = 1
  max_size    = 5
}
```

### With the kubernetes provider

```terraform
resource "scaleway_vpc_private_network" "hedy" {}

resource "scaleway_k8s_cluster" "joy" {
  name    = "joy"
  version = "1.24.3"
  cni     = "flannel"
  private_network_id = scaleway_vpc_private_network.hedy.id
  delete_additional_resources = false
}

resource "scaleway_k8s_pool" "john" {
  cluster_id = scaleway_k8s_cluster.joy.id
  name       = "john"
  node_type  = "DEV1-M"
  size       = 1
}

resource "null_resource" "kubeconfig" {
  depends_on = [scaleway_k8s_pool.john] # at least one pool here
  triggers = {
    host                   = scaleway_k8s_cluster.joy.kubeconfig[0].host
    token                  = scaleway_k8s_cluster.joy.kubeconfig[0].token
    cluster_ca_certificate = scaleway_k8s_cluster.joy.kubeconfig[0].cluster_ca_certificate
  }
}

provider "kubernetes" {
  host  = null_resource.kubeconfig.triggers.host
  token = null_resource.kubeconfig.triggers.token
  cluster_ca_certificate = base64decode(
    null_resource.kubeconfig.triggers.cluster_ca_certificate
  )
}
```

The `null_resource` is needed because when the cluster is created, it's status is `pool_required`, but the kubeconfig can already be downloaded.
It leads the `kubernetes` provider to start creating its objects, but the DNS entry for the Kubernetes master is not yet ready, that's why it's needed to wait for at least a pool.

### With the Helm provider

```terraform
resource "scaleway_vpc_private_network" "hedy" {}

resource "scaleway_k8s_cluster" "joy" {
  name    = "joy"
  version = "1.24.3"
  cni     = "flannel"
  delete_additional_resources = false
  private_network_id = scaleway_vpc_private_network.hedy.id
}

resource "scaleway_k8s_pool" "john" {
  cluster_id = scaleway_k8s_cluster.joy.id
  name       = "john"
  node_type  = "DEV1-M"
  size       = 1
}

resource "null_resource" "kubeconfig" {
  depends_on = [scaleway_k8s_pool.john] # at least one pool here
  triggers = {
    host                   = scaleway_k8s_cluster.joy.kubeconfig[0].host
    token                  = scaleway_k8s_cluster.joy.kubeconfig[0].token
    cluster_ca_certificate = scaleway_k8s_cluster.joy.kubeconfig[0].cluster_ca_certificate
  }
}

provider "helm" {
  kubernetes {
    host = null_resource.kubeconfig.triggers.host
    token = null_resource.kubeconfig.triggers.token
    cluster_ca_certificate = base64decode(
    null_resource.kubeconfig.triggers.cluster_ca_certificate
    )
  }
}

resource "scaleway_lb_ip" "nginx_ip" {
  zone       = "fr-par-1"
  project_id = scaleway_k8s_cluster.joy.project_id
}

resource "helm_release" "nginx_ingress" {
  name      = "nginx-ingress"
  namespace = "kube-system"

  repository = "https://kubernetes.github.io/ingress-nginx"
  chart = "ingress-nginx"

  set {
    name = "controller.service.loadBalancerIP"
    value = scaleway_lb_ip.nginx_ip.ip_address
  }

  // enable proxy protocol to get client ip addr instead of loadbalancer one
  set {
    name = "controller.config.use-proxy-protocol"
    value = "true"
  }
  set {
    name = "controller.service.annotations.service\\.beta\\.kubernetes\\.io/scw-loadbalancer-proxy-protocol-v2"
    value = "true"
  }

  // indicates in which zone to create the loadbalancer
  set {
    name = "controller.service.annotations.service\\.beta\\.kubernetes\\.io/scw-loadbalancer-zone"
    value = scaleway_lb_ip.nginx_ip.zone
  }

  // enable to avoid node forwarding
  set {
    name = "controller.service.externalTrafficPolicy"
    value = "Local"
  }

  // enable this annotation to use cert-manager
  //set {
  //  name  = "controller.service.annotations.service\\.beta\\.kubernetes\\.io/scw-loadbalancer-use-hostname"
  //  value = "true"
  //}
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) The name for the Kubernetes cluster.

- `type` - (Optional) The type of Kubernetes cluster. Possible values are:

    - for mutualized clusters: `kapsule` or `multicloud`

    - for dedicated Kapsule clusters: `kapsule-dedicated-4`, `kapsule-dedicated-8` or `kapsule-dedicated-16`.

    - for dedicated Kosmos clusters: `multicloud-dedicated-4`, `multicloud-dedicated-8` or `multicloud-dedicated-16`.

- `description` - (Optional) A description for the Kubernetes cluster.

- `version` - (Required) The version of the Kubernetes cluster.

- `cni` - (Required) The Container Network Interface (CNI) for the Kubernetes cluster.
~> **Important:** Updates to this field will recreate a new resource.

- `delete_additional_resources` - (Required) Delete additional resources like block volumes, load-balancers and the cluster's private network (if empty) that were created in Kubernetes on cluster deletion.
~> **Important:** Setting this field to `true` means that you will lose all your cluster data and network configuration when you delete your cluster.
If you prefer keeping it, you should instead set it as `false`.

- `private_network_id` - (Required) The ID of the private network of the cluster.

~> **Important:** Changes to this field will recreate a new resource.

~> **Important:** Private Networks are now mandatory with Kapsule Clusters. If you have a legacy cluster (no `private_network_id` set),
you can still set it now. In this case it will not destroy and recreate your cluster but migrate it to the Private Network.

- `tags` - (Optional) The tags associated with the Kubernetes cluster.

- `autoscaler_config` - (Optional) The configuration options for the [Kubernetes cluster autoscaler](https://github.com/kubernetes/autoscaler/tree/master/cluster-autoscaler).

    - `disable_scale_down` - (Defaults to `false`) Disables the scale down feature of the autoscaler.

    - `scale_down_delay_after_add` - (Defaults to `10m`) How long after scale up that scale down evaluation resumes.

    - `scale_down_unneeded_time` - (Default to `10m`) How long a node should be unneeded before it is eligible for scale down.

    - `estimator` - (Defaults to `binpacking`) Type of resource estimator to be used in scale up.

    - `expander` - (Default to `random`) Type of node group expander to be used in scale up.

    - `ignore_daemonsets_utilization` - (Defaults to `false`) Ignore DaemonSet pods when calculating resource utilization for scaling down.

    - `balance_similar_node_groups` - (Defaults to `false`) Detect similar node groups and balance the number of nodes between them.

    - `expendable_pods_priority_cutoff` - (Defaults to `-10`) Pods with priority below cutoff will be expendable. They can be killed without any consideration during scale down and they don't cause scale up. Pods with null priority (PodPriority disabled) are non expendable.

    - `scale_down_utilization_threshold` - (Defaults to `0.5`) Node utilization level, defined as sum of requested resources divided by capacity, below which a node can be considered for scale down

    - `max_graceful_termination_sec` - (Defaults to `600`) Maximum number of seconds the cluster autoscaler waits for pod termination when trying to scale down a node

- `auto_upgrade` - (Optional) The auto upgrade configuration.

    - `enable` - (Optional) Set to `true` to enable Kubernetes patch version auto upgrades.
~> **Important:** When enabling auto upgrades, the `version` field take a minor version like x.y (ie 1.18).

    - `maintenance_window_start_hour` - (Optional) The start hour (UTC) of the 2-hour auto upgrade maintenance window (0 to 23).

    - `maintenance_window_day` - (Optional) The day of the auto upgrade maintenance window (`monday` to `sunday`, or `any`).

- `feature_gates` - (Optional) The list of [feature gates](https://kubernetes.io/docs/reference/command-line-tools-reference/feature-gates/) to enable on the cluster.

- `admission_plugins` - (Optional) The list of [admission plugins](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/) to enable on the cluster.

- `apiserver_cert_sans` - (Optional) Additional Subject Alternative Names for the Kubernetes API server certificate

- `open_id_connect_config` - (Optional) The OpenID Connect configuration of the cluster

    - `issuer_url` - (Required) URL of the provider which allows the API server to discover public signing keys

    - `client_id` - (Required) A client id that all tokens must be issued for

    - `username_claim` - (Optional) JWT claim to use as the user name

    - `username_prefix` - (Optional) Prefix prepended to username

    - `groups_claim` - (Optional) JWT claim to use as the user's group

    - `groups_prefix` - (Optional) Prefix prepended to group claims

    - `required_claim` - (Optional) Multiple key=value pairs that describes a required claim in the ID Token

- `default_pool` - (Deprecated) See below.

- `region` - (Defaults to [provider](../index.md#arguments-reference) `region`) The [region](../guides/regions_and_zones.md#regions) in which the cluster should be created.

- `project_id` - (Defaults to [provider](../index.md#arguments-reference) `project_id`) The ID of the project the cluster is associated with.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

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
- `upgrade_available` - Set to `true` if a newer Kubernetes version is available.
- `organization_id` - The organization ID the cluster is associated with.

## Import

Kubernetes clusters can be imported using the `{region}/{id}`, e.g.

```bash
$ terraform import scaleway_k8s_cluster.mycluster fr-par/11111111-1111-1111-1111-111111111111
```

## Deprecation of default_pool

`default_pool` is deprecated in favour the `scaleway_k8s_pool` resource. Here is a migration example.

Before:

```terraform
resource "scaleway_k8s_cluster" "jack" {
  name    = "jack"
  version = "1.18.0"
  cni     = "cilium"

  default_pool {
    node_type = "DEV1-M"
    size      = 1
  }
}
```

After:

```terraform
resource "scaleway_k8s_cluster" "jack" {
  name    = "jack"
  version = "1.18.0"
  cni     = "cilium"
}

resource "scaleway_k8s_pool" "default" {
  cluster_id = scaleway_k8s_cluster.jack.id
  name       = "default"
  node_type  = "DEV1-M"
  size       = 1
}
```

Once you have moved all the `default_pool` into their own object, you will need to import them. If your pool had the ID 11111111-1111-1111-1111-111111111111 in the `fr-par` region, you can import it by typing:

```bash
$ terraform import scaleway_k8s_pool.default fr-par/11111111-1111-1111-1111-111111111111
```

Then you will only need to type `terraform apply` to have a smooth migration.
