---
subcategory: "Secrets"
page_title: "Scaleway: scaleway_secret"
---

# Resource: scaleway_secret

Creates and manages Scaleway Secrets.
For more information, see [the documentation](https://developers.scaleway.com/en/products/secret_manager/api/v1alpha1/).

## Example Usage

### Basic

```terraform
resource "scaleway_secret" "main" {
  name        = "foo"
  description = "barr"
  tags        = ["foo", "terraform"]
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) Name of the secret (e.g. `my-secret`).
- `path` - (Optional) Path of the secret, defaults to `/`.
- `description` - (Optional) Description of the secret (e.g. `my-new-description`).
- `tags` - (Optional) Tags of the secret (e.g. `["tag", "secret"]`).
- `region` - (Defaults to [provider](../index.md#region) `region`) The [region](../guides/regions_and_zones.md#regions)
  in which the resource exists.
- `project_id` - (Optional) The project ID containing is the secret.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `version_count` - The number of versions for this Secret.
- `status` - The status of the Secret.
- `created_at` - Date and time of secret's creation (RFC 3339 format).
- `updated_at` - Date and time of secret's last update (RFC 3339 format).

## Import

The Secret can be imported using the `{region}/{id}`, e.g.

```bash
$ terraform import scaleway_secret.main fr-par/11111111-1111-1111-1111-111111111111
```
