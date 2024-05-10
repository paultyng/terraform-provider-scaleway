---
subcategory: "Secrets"
page_title: "Scaleway: scaleway_secret"
---

# scaleway_secret

Gets information about Scaleway Secrets.
For more information, see [the documentation](https://developers.scaleway.com/en/products/secret_manager/api/v1alpha1/).

## Examples

### Basic

```hcl
resource "scaleway_secret" "main" {
  name        = "foo"
  description = "barr"
}

// Get info by secret ID
data "scaleway_secret" "my_secret" {
  secret_id = "11111111-1111-1111-1111-111111111111"
}

// Get info by secret Name
data "scaleway_secret" "by_name" {
  name = "your_secret_name"
}
```

## Argument Reference

- `name` - (Optional) The secret name.
  Only one of `name` and `secret_id` should be specified.

- `path` - (Optional) The secret path.
  Conflicts with `secret_id`.

- `secret_id` - (Optional) The secret id.
  Only one of `name` and `secret_id` should be specified.

- `organization_id` - (Optional) The organization ID the Project is associated with.
  If no default organization_id is set, one must be set explicitly in this datasource

- `region` - (Defaults to [provider](../index.md#region) `region`) The [region](../guides/regions_and_zones.md#regions) in which the secret exists.

- `project_id` - (Optional. Defaults to [provider](../index.md#project_id) `project_id`) The ID of the
  project the secret is associated with.


## Attributes Reference

Exported attributes are the ones from `scaleway_secret` [resource](../resources/secret.md)
