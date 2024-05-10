---
subcategory: "Containers"
page_title: "Scaleway: scaleway_container_namespace"
---

# Resource: scaleway_container_namespace

Creates and manages Scaleway Serverless Container Namespace.
For more information see [the documentation](https://developers.scaleway.com/en/products/containers/api/#namespaces-cdce79).

## Example Usage

### Basic

```terraform
resource "scaleway_container_namespace" "main" {
  name        = "main-container-namespace"
  description = "Main container namespace"
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) The unique name of the container namespace.

~> **Important** Updates to `name` will recreate the namespace.

- `description` (Optional) The description of the namespace.

- `region` - (Defaults to [provider](../index.md#region) `region`). The [region](../guides/regions_and_zones.md#regions) in which the namespace should be created.

- `project_id` - (Defaults to [provider](../index.md#project_id) `project_id`) The ID of the project the namespace is associated with.

- `environment_variables` - The environment variables of the namespace.

- `secret_environment_variables` - The secret environment variables of the namespace.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the namespace

~> **Important:** Container namespaces' IDs are [regional](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{region}/{id}`, e.g. `fr-par/11111111-1111-1111-1111-111111111111`

- `organization_id` - The organization ID the namespace is associated with.
- `registry_endpoint` - The registry endpoint of the namespace.
- `registry_namespace_id` - The registry namespace ID of the namespace.


## Import

Namespaces can be imported using the `{region}/{id}`, e.g.

```bash
$ terraform import scaleway_container_namespace.main fr-par/11111111-1111-1111-1111-111111111111
```
