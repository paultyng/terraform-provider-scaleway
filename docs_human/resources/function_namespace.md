---
subcategory: "Functions"
page_title: "Scaleway: scaleway_function_namespace"
---

# Resource: scaleway_function_namespace

Creates and manages Scaleway Function Namespace.
For more information see [the documentation](https://developers.scaleway.com/en/products/functions/api/).

## Example Usage

### Basic

```terraform
resource "scaleway_function_namespace" "main" {
  name        = "main-function-namespace"
  description = "Main function namespace"
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) The unique name of the function namespace.

~> **Important** Updates to `name` will recreate the namespace.

- `description` (Optional) The description of the namespace.

- `region` - (Defaults to [provider](../index.md#region) `region`). The [region](../guides/regions_and_zones.md#regions) in which the namespace should be created.

- `project_id` - (Defaults to [provider](../index.md#project_id) `project_id`) The ID of the project the namespace is associated with.

- `environment_variables` - (Optional) The environment variables of the namespace.

- `secret_environment_variables` - (Optional) The [secret environment](https://www.scaleway.com/en/docs/compute/containers/concepts/#secrets) variables of the namespace.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the namespace

~> **Important:** Function namespaces' IDs are [regional](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{region}/{id}`, e.g. `fr-par/11111111-1111-1111-1111-111111111111`

- `organization_id` - The organization ID the namespace is associated with.
- `registry_endpoint` - The registry endpoint of the namespace.
- `registry_namespace_id` - The registry namespace ID of the namespace.


## Import

Namespaces can be imported using the `{region}/{id}`, e.g.

```bash
$ terraform import scaleway_function_namespace.main fr-par/11111111-1111-1111-1111-111111111111
```
