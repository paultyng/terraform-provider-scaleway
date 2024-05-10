---
subcategory: "Kubernetes"
page_title: "Scaleway: scaleway_k8s_version"
---

# scaleway_k8s_version

Gets information about a Kubernetes version.
For more information, see [the documentation](https://developers.scaleway.com/en/products/k8s/api).

You can also use the [scaleway-cli](https://github.com/scaleway/scaleway-cli) with `scw k8s version list` to list all available versions.

## Example Usage

### Use the latest version

```hcl
data "scaleway_k8s_version" "latest" {
  name = "latest"
}
```

### Use a specific version

```hcl
data "scaleway_k8s_version" "by_name" {
  name = "1.26.0"
}
```

## Argument Reference

- `name` - (Required) The name of the Kubernetes version.
- `region` - (Defaults to [provider](../index.md) `region`) The [region](../guides/regions_and_zones.md#regions) in which the version exists.

## Attributes Reference

In addition to all above arguments, the following attributes are exported:

- `id` - The ID of the version.

~> **Important:** Kubernetes versions' IDs are [regional](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{region}/{name}`, e.g. `fr-par/1.1.1`

- `available_cnis` - The list of supported Container Network Interface (CNI) plugins for this version.
- `available_container_runtimes` - The list of supported container runtimes for this version.
- `available_feature_gates` - The list of supported feature gates for this version.