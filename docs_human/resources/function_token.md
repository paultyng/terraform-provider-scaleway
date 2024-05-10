---
subcategory: "Functions"
page_title: "Scaleway: scaleway_function_token"
---

# Resource: scaleway_function_token

Creates and manages Scaleway Function Token.
For more information see [the documentation](https://developers.scaleway.com/en/products/functions/api/#tokens-26b085).

## Example Usage

### Basic

```terraform
resource scaleway_function_namespace main {
  name = "test-function-token-ns"
}

resource scaleway_function main {
  namespace_id = scaleway_function_namespace.main.id
  runtime      = "go118"
  handler      = "Handle"
  privacy      = "private"
}

// Namespace Token
resource scaleway_function_token namespace {
  namespace_id = scaleway_function_namespace.main.id
  expires_at = "2022-10-18T11:35:15+02:00"
}

// Function Token
resource scaleway_function_token function {
  function_id = scaleway_function.main.id
}
```

## Argument Reference

The following arguments are supported:

- `namespace_id` - (Required) The ID of the function namespace.

- `function_id` - (Required) The ID of the function.

~> Only one of `namespace_id` or `function_id` must be set.

- `description` (Optional) The description of the token.

- `expires_at` (Optional) The expiration date of the token.

- `region` - (Defaults to [provider](../index.md#region) `region`). The [region](../guides/regions_and_zones.md#regions) in which the namespace should be created.

~> **Important** Updates to any fields will recreate the token.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the token.

~> **Important:** Function tokens' IDs are [regional](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{region}/{id}`, e.g. `fr-par/11111111-1111-1111-1111-111111111111`

- `token` - The token.

## Import

Tokens can be imported using the `{region}/{id}`, e.g.

```bash
$ terraform import scaleway_function_token.main fr-par/11111111-1111-1111-1111-111111111111
```
