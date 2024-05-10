---
subcategory: "Containers"
page_title: "Scaleway: scaleway_container_token"
---

# Resource: scaleway_container_token

Creates and manages Scaleway Container Token.
For more information see [the documentation](https://developers.scaleway.com/en/products/containers/api/#tokens-26b085).

## Example Usage

### Basic

```terraform
resource scaleway_container_namespace main {
  name = "test-container-token-ns"
}

resource scaleway_container main {
  namespace_id = scaleway_container_namespace.main.id
}

// Namespace Token
resource scaleway_container_token namespace {
  namespace_id = scaleway_container_namespace.main.id
  expires_at = "2022-10-18T11:35:15+02:00"
}

// Container Token
resource scaleway_container_token container {
  container_id = scaleway_container.main.id
}
```

## Argument Reference

The following arguments are supported:

- `namespace_id` - (Required) The ID of the container namespace.

- `container_id` - (Required) The ID of the container.

~> Only one of `namespace_id` or `container_id` must be set.

- `description` (Optional) The description of the token.

- `expires_at` (Optional) The expiration date of the token.

- `region` - (Defaults to [provider](../index.md#region) `region`). The [region](../guides/regions_and_zones.md#regions) in which the namespace should be created.

~> **Important** Updates to any fields will recreate the token.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the token.

~> **Important:** Container tokens' IDs are [regional](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{region}/{id}`, e.g. `fr-par/11111111-1111-1111-1111-111111111111`

- `token` - The token.

## Import

Tokens can be imported using the `{region}/{id}`, e.g.

```bash
$ terraform import scaleway_container_token.main fr-par/11111111-1111-1111-1111-111111111111
```
