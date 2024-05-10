---
subcategory: "Container Registry"
page_title: "Scaleway: scaleway_registry_namespace"
---

# Resource: scaleway_registry_namespace

Creates and manages Scaleway Container Registry.
For more information see [the documentation](https://developers.scaleway.com/en/products/registry/api/).

## Example Usage

### Basic

```terraform
resource "scaleway_registry_namespace" "main" {
  name        = "main-cr"
  description = "Main container registry"
  is_public   = false
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) The unique name of the namespace.

~> **Important** Updates to `name` will recreate the namespace.

- `description` (Optional) The description of the namespace.

- `is_public` (Defaults to `false`) Whether the images stored in the namespace should be downloadable publicly (docker pull).

- `region` - (Defaults to [provider](../index.md#region) `region`). The [region](../guides/regions_and_zones.md#regions) in which the namespace should be created.

- `project_id` - (Defaults to [provider](../index.md#project_id) `project_id`) The ID of the project the namespace is associated with.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the namespace

~> **Important:** Registry namespaces' IDs are [regional](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{region}/{id}`, e.g. `fr-par/11111111-1111-1111-1111-111111111111`

- `endpoint` - Endpoint reachable by Docker.
- `organization_id` - The organization ID the namespace is associated with.

## Import

Namespaces can be imported using the `{region}/{id}`, e.g.

```bash
$ terraform import scaleway_registry_namespace.main fr-par/11111111-1111-1111-1111-111111111111
```
